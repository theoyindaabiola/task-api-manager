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
	"github.com/pquerna/otp/totp"
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

	if !user.Enabled2FA {
		// no 2fa, issue long-lived token
		return GenerateJWTToken(user, true)
	}

	return GenerateJWTToken(user, false)
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

	if err := s.UserDAO.DB.Create(otp).Error; err != nil {
		return nil, err
	}

	// send otp via email queue
	if err := utils.PublishMessage(
		os.Getenv("EMAIL_OTP_QUEUE"),
		user.Email,
		otpCode,
		"otp",
		"",
		"",
		"",
	); err != nil {
		return nil, err
	}
	return otp, nil
}

func (s *UserService) VerifyEmailOTP(userID, code string) (string, error) {
	user, err := s.UserDAO.GetUserByIdDB(userID)
    if err != nil {
        return "", err
    }

	// check DB for the otp code
	otp, err := s.UserDAO.GetOTPByCodeAndUser(userID, code)
	if err != nil {
		return "", errors.New("invalid OTP")
	}

	// check if OTP is expired
    if otp.ExpiresAt.Before(time.Now()) {
        return "", errors.New("OTP expired")
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

    // update user flag for Email OTP verification
	user.Enabled2FA = true
    user.Is2FAVerified = true
    user.TwoFactorType = "email"
	
    if err := s.UserDAO.Update(user); err != nil {
        return "", err
    }

	// generate JWT reflecting current 2FA state
	return GenerateJWTToken(user, user.Is2FAVerified)
}

// enable TOTP generates then store a new TOTP secret for the user
func (s *UserService) EnableTOTP(userID string) (string, error) {
    user, err := s.UserDAO.GetUserByIdDB(userID)
	if err != nil {
        return "", err
    }

	// block enabling again if already active
	if user.Enabled2FA && user.TOTPSecret != nil {
		return "", fmt.Errorf("2FA is already enabled for this account")

	}

	// generate secret
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer: "TaskAPI",
		AccountName: user.Email,
	})
	if err != nil {
        return "", err
    }

	// store secret in DB
	secret := key.Secret()
	user.TOTPSecret = &secret

    if err := s.UserDAO.Update(user); err != nil {
        return "", err
    }

    // return URL: to be turned to qr code
    return key.URL(), nil
}

// VerifyTOTP to checks the provided OTP against the stored secret
func (s *UserService) VerifyTOTP(userID, code string) (string, error) {
    user, err := s.UserDAO.GetUserByIdDB(userID)
    if err != nil {
        return "", err
    }

	if user.TOTPSecret == nil {
		return "", fmt.Errorf("TOTP is not enabled")
	}
	valid := totp.Validate(code, *user.TOTPSecret)

	if !valid {
		return "", fmt.Errorf("invalid code")
	}

	user.Enabled2FA = true
	if err := s.UserDAO.Update(user); err != nil {
		return "", err
	}
		
	// update user flag for 2FA type
    user.TwoFactorType = "totp"
    if err := s.UserDAO.Update(user); err != nil {
        return "", err
    }

	return GenerateJWTToken(user, true)
}

func (s *UserService) Disable2FA(userID string) (string, error) {
    user, err := s.UserDAO.GetUserByIdDB(userID)
    if err != nil {
        return "", err
    }

	// security: ensure user has already verified 2FA
    if !user.Is2FAVerified {
        return "", fmt.Errorf("2FA verification required before disabling")
    }

	// update the fields to disable 2fa
	user.Enabled2FA = false
	user.Is2FAVerified = false
	user.TwoFactorType = ""
	user.TOTPSecret = nil // only for TOTP

	// update user
    if err := s.UserDAO.Update(user); err != nil {
        return "", err
    }

	 // invalidate old OTPs
    _ = s.UserDAO.InvalidateOldOTPs(user.ID.String())

    // issue a fresh token without 2FA requirement
    return GenerateJWTToken(user, true)
}

func GenerateJWTToken(user *models.User, is2FAVerified bool) (string, error) {
	// set short expiry if OTP not verified yet
	expiry := time.Hour * 48
	if !is2FAVerified {
		expiry = time.Minute * 10
	}

    claims := jwt.MapClaims{
        "user_id":          user.ID.String(),
        "enabled_2fa":      user.Enabled2FA,
        "is_2fa_verified": 	is2FAVerified,
        "exp":              time.Now().Add(expiry).Unix(),
		"issuer":           "task-api-manager",
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
