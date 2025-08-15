package services

import (
	"fmt"
	"taskapi/dao"    // needs to interact with it
	"taskapi/models" // needs the model for the db
	"taskapi/utils"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"

	// "github.com/google/uuid"
	"errors"
	"os"
	"time"
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

	/** 
		if valid, create a token for the user
		MapClaim is flexible
	**/
	claims := jwt.MapClaims{
		"user_id": user.ID, // from the db bcos it is reliable
		"exp": time.Now().Add(time.Hour * 48).Unix(), // set the time
		"issuer": "task-api-manager",
	}
	// get the token signed
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
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
