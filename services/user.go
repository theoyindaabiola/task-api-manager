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
	if err := s.UserDAO.DB.Save(user).Error; err != nil {
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

	// if 2FA enabled but phone not verified (new or changed number)
	if !user.IsSmsVerified {
		// send a new OTP to verify the updated number
		if err := s.IssueSMSOTP(user); err != nil {
    		return "", fmt.Errorf("failed to issue OTP: %w", err)
		}
		// issue only short-lived token until verification
		return GenerateJWTToken(user, false)
	}
	
	// At Enabled2FA = true and IsSmsVerified = true, issue long-lived token 
	return GenerateJWTToken(user, true)
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

// enable TOTP generates then store a new TOTP secret for the user
func (s *UserService) EnableSMS(userID, phoneNumber string) error {
	// get user from DB
	user, err := s.UserDAO.GetUserByIdDB(userID)
	if err != nil {
		return err
	}

	// block enabling again if already active
	if user.Enabled2FA {
		return fmt.Errorf("SMS 2FA already enabled, verify or request a new OTP")
	}

	// update phone number and enable flag
	user.PhoneNumber = &phoneNumber
	user.Enabled2FA = true
	user.IsSmsVerified = false

	if err := s.UserDAO.Update(user); err != nil {
		return err
	}

	// after enabling, issue the OTP
	if err := s.IssueSMSOTP(user); err != nil {
		return fmt.Errorf("failed to issue SMS OTP: %w", err)
	}

	return nil
}

// updatePhoneNumber changes the user’s phone and reissues OTP verification.
func (s *UserService) UpdatePhoneNumber(userID, newPhoneNumber string) error {
	user, err := s.UserDAO.GetUserByIdDB(userID)
	if err != nil {
		return err
	}

	if !user.Enabled2FA {
		return fmt.Errorf("SMS 2FA is not enabled")
	}

	if user.PhoneNumber != nil && *user.PhoneNumber == newPhoneNumber {
		return fmt.Errorf("new phone number cannot be the same as the existing one")
	}

	// update phone number and reset verification
	user.PhoneNumber = &newPhoneNumber
	user.IsSmsVerified = false

	if err := s.UserDAO.Update(user); err != nil {
		return err
	}

	// send a new OTP to verify the updated number
	if err := s.IssueSMSOTP(user); err != nil {
		return fmt.Errorf("failed to issue SMS OTP: %w", err)
	}

	return nil
}


// this function below can take any phone number if the user's number (registered or from enable2fa function) has issue getting the otp code, but this is very risky
func (s *UserService) RequestSMS(userID string) error {
    user, err := s.UserDAO.GetUserByIdDB(userID)
    if err != nil {
        return err
    }
    if !user.Enabled2FA {
        return fmt.Errorf("SMS 2FA not enabled")
    }
    if user.PhoneNumber == nil || *user.PhoneNumber == "" {
        return fmt.Errorf("no phone number on record")
    }
    return s.IssueSMSOTP(user) // same number, re-send OTP
}

// VerifyTOTP to checks the provided OTP against the stored secret
func (s *UserService) VerifySMS(userID, code string) (string, error) {
    user, err := s.UserDAO.GetUserByIdDB(userID)
    if err != nil {
        return "", err
    }

	if !user.Enabled2FA || user.SmsOTP == nil || user.SmsOTPExpiresAt == nil {
        return "", fmt.Errorf("SMS 2FA not enabled or no OTP set")
    }
    if user.SmsOTPExpiresAt.Before(time.Now()) {
        return "", fmt.Errorf("SMS OTP expired")
    }
    if *user.SmsOTP != code {
        return "", fmt.Errorf("invalid SMS OTP")
    }
    user.IsSmsVerified = true
    user.SmsOTP = nil
    user.SmsOTPExpiresAt = nil
    if err := s.UserDAO.Update(user); err != nil {
        return "", err
    }

    // generate fresh long session JWT
    return GenerateJWTToken(user, true)
}

func (s *UserService) DisableSMS(userID string) (string, error) {
    user, err := s.UserDAO.GetUserByIdDB(userID)
    if err != nil {
        return "", err
    }

	// update the fields to disable 2fa
    user.Enabled2FA = false
    user.SmsOTP = nil
    user.IsSmsVerified = false
	user.SmsOTPExpiresAt = nil

    if err := s.UserDAO.Update(user); err != nil {
        return "", err
    }

    // issue a fresh token without 2FA requirement
    return GenerateJWTToken(user, true)
}

func GenerateJWTToken(user *models.User, totpVerified bool) (string, error) {
	// set short expiry if TOTP not verified yet
    var expiry time.Duration
	if totpVerified {
		expiry = time.Hour * 48 // 2 days for full session
	} else {
		expiry = time.Minute * 10 // short session until TOTP is verified
	}

    claims := jwt.MapClaims{
        "user_id":          user.ID,
        "enabled_2fa":      user.Enabled2FA,
        "is_totp_verified": totpVerified,
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

func (s *UserService) GenerateOTP() (string, time.Time) {
    otp := utils.Generate2FACode()
    expiresAt := time.Now().Add(5 * time.Minute)
    return otp, expiresAt
}

func (s *UserService) IssueSMSOTP(user *models.User) error {
	otp, expiresAt := s.GenerateOTP()

	user.SmsOTP = &otp
	user.SmsOTPExpiresAt = &expiresAt

	if err := s.UserDAO.Update(user); err != nil {
		return err
	}

	// Send OTP via SMS queue
	return utils.PublishMessage(
		os.Getenv("SMS_OTP_QUEUE"),
		*user.PhoneNumber,
		otp,
		"sms_otp",
		"",
		"",
		"",
	)
}
