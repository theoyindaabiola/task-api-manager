package services

import (
	"fmt"
	"taskapi/dao" // needs to interact with it
	"taskapi/models" // needs the model for the db
	"taskapi/utils"

	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserDAO *dao.UserDAO
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func NewUserService(dao *dao.UserDAO) *UserService {
	// return an instance of the NewTaskDAO instance of the dao we are passing in
	return &UserService{UserDAO: dao}
}

func (s *UserService) RegisterUser(user *models.User) error {
	// hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	// the hash string returned is the same as the user.Password
	user.Password = string(hash)
	verificationCode := utils.GenerateVerificationCode()
	user.VerificationToken = &verificationCode
	user.Enabled2FA = false
	user.ResetToken = nil
	if err := s.UserDAO.CreateUserDB(user); err != nil {
		return err
	}
	// send the verification email 
	if err := utils.PublishMessage(
		os.Getenv("EMAIL_VERIFICATION_QUEUE"),
		user.Email,
		verificationCode,
		"verification",
		"", 
		"",
		"",
		); err != nil {
		return err
	}
	return nil
}

// verify user
func (s *UserService) VerificationService(verificationToken string) error {
	user, err := s.UserDAO.GetUserVerification(verificationToken)
	if err != nil {
		return err
	}

	user.Verified = true
	// clear after use
	user.VerificationToken = nil
	if err := s.UserDAO.Update(user); err != nil {
		return err
	}
	return nil
}

func (s *UserService) LoginUser(payload *models.User) (string, error) {
	// check if user exist
	user, err := s.UserDAO.GetUserDB(payload.Username)
	if err != nil {
		return "", err // "" is a token for the user
	}

	if !user.Verified {
		return "", errors.New("user not verified")
	}

	// compare the hashed passwords if user exists
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	otpVerified := false
	// otp should be true for user login, either verified or not
	if user.Enabled2FA {
		// last OTP record to confirm if verified
		latestOtp, err := s.UserDAO.GetLatestOTPByUser(user.ID.String())
		if err == nil && latestOtp.OtpVerified {
			otpVerified = true
		}
	} else {
    	otpVerified = true // no 2FA, no restriction
	}
    return GenerateJWTToken(user, otpVerified)
}

func (s *UserService) ForgotPasswordService(email string) error {
	user, err := s.UserDAO.GetUserByEmail(email)
    if err != nil {
        return nil
    }
	resetToken := utils.GenerateVerificationCode()
	expiresAt := time.Now().Add(time.Hour)
	user.ResetToken = &resetToken // generated
	user.ResetTokenExpiresAt = &expiresAt
	if err := s.UserDAO.Update(user); err != nil {
        return err
    }
	return utils.PublishMessage(
		os.Getenv("EMAIL_RESET_QUEUE"),
		user.Email,
		resetToken,
		"reset",
		"", 
		"",
		"",
	)
}

func (s *UserService) ResetPasswordService(token, password string) error {
	user, err := s.UserDAO.GetUserResetToken(token)
	if err != nil || user == nil || 
		user.ResetTokenExpiresAt == nil || 
		user.ResetTokenExpiresAt.Before(time.Now()) {
			return fmt.Errorf("invalid or expired token")
    }
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
	user.Password = string(hash)
	user.ResetToken = nil
	user.ResetTokenExpiresAt = nil
	return s.UserDAO.Update(user)
}

func (s *UserService) GetUserByID(id string) (*models.User, error) {
    user, err := s.UserDAO.GetUserByIdDB(id)
    if err != nil {
        return nil, err
    }
    return user, nil
}

func (s *UserService) EnableEmail2FA(email string) (*models.OtpVerification, error) {
	// get user by email
	user, err := s.UserDAO.GetUserByEmail(email)
    if err != nil {
        return nil, err
    }

	_ = s.UserDAO.InvalidateOldOTPs(user.ID.String())

	// generate otp code and expiring time
	otpCode := utils.Generate2FACode()
	expireAt := time.Now().Add(5 * time.Minute)

	// set up the user&otp model
	otp := &models.OtpVerification{
		UserID:    user.ID,
		OtpCode:   otpCode,
		ExpiresAt: expireAt,
		OtpVerified:  false, // always false at new otp creation
	}

	// store the otp in the database
	if err := s.UserDAO.CreateOTP(otp); err != nil {
		return nil, fmt.Errorf("failed to create OTP: %w", err)
	}

	// send the otp via user's email
	if err := utils.PublishMessage(
		os.Getenv("EMAIL_OTP_QUEUE"),
		user.Email,
		otpCode,
		"otp",
		"",
		"",
		"",
	); err != nil {
    	return otp, fmt.Errorf("OTP saved but not sent: %w", err)
	}

	return otp, nil
}

func (s *UserService) VerifyEmailOTP(userID, code string) (string, error) {
	// check DB for the otp code
	otp, err := s.UserDAO.GetOTPByCodeAndUser(userID, code)
	if err != nil {
		return "", errors.New("invalid OTP")
	}

	// check expiry
	if time.Now().After(otp.ExpiresAt) {
		return "", fmt.Errorf("OTP expired")
	}

	// check if already used
	if otp.OtpVerified {
		return "", fmt.Errorf("OTP already used")
	}

	// mark OTP as verified
	otp.OtpVerified = true
	if err := s.UserDAO.DB.Save(otp).Error; err != nil {
		return "", err
	}

	// update 2FA flag
	if err := s.UserDAO.Update2FA(userID, true); err != nil {
		return "", err
	}

	// reload fresh user
	user, err := s.UserDAO.GetUserByIdDB(userID)
	if err != nil {
		return "", fmt.Errorf("user not found")
	}

	// generate JWT reflecting current 2FA state
	signedToken, err := GenerateJWTToken(user, user.Enabled2FA) 
	if err != nil { 
		return "", err 
	}

	return signedToken, nil

}

func (s *UserService) DisableEmail2FA(email string) (*models.User, error) {
    // get user by email
    user, err := s.UserDAO.GetUserByEmail(email)
    if err != nil {
        return nil, err
    }

    // disable flags
    user.Enabled2FA = false
    user.IsOtpVerified = false

    // update user
    if err := s.UserDAO.Update(user); err != nil {
        return nil, err
    }

    // invalidate old OTPs
    _ = s.UserDAO.InvalidateOldOTPs(user.ID.String())

    return user, nil
}

func GenerateJWTToken(user *models.User, otpVerified bool) (string, error) {
    // set short expiry if OTP not verified yet
    var expiry time.Duration
    if otpVerified {
        expiry = time.Hour * 48 // normal long session
    } else {
        expiry = time.Minute * 10 // short token for OTP pending
    }

    claims := jwt.MapClaims{
        "user_id":         user.ID,
        "enabled_2fa":     user.Enabled2FA,
        "is_otp_verified": otpVerified,
        "exp":             time.Now().Add(expiry).Unix(),
        "issuer":          "task-api-manager",
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}


// to decode, verify the signature, check the expiration time and extract the user details
func ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// check if the token is signed with the HMAC signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid Token")
		}
		return jwtSecret, nil
	})				
}
