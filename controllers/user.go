package controllers

import (
	"log"
	"net/http" // allows http requests
	"taskapi/dto"
	"taskapi/models"
	"taskapi/services"

	"github.com/gin-gonic/gin" // web framework
)

/**
	service instance needed to send the payload to it and the functions
	has to be connected to this instance
**/

type UserController struct {
	UserService *services.UserService
}

func NewUserController(service *services.UserService) *UserController {
	return &UserController{UserService: service}
}

// context gets into the body of your package. c is an instance of pointing to the gin.Context
func (tc *UserController) CreateUser(c *gin.Context) {
	// dto instance for validation
	var payload dto.RegisterUserDTO // user payload placeholder
	
	// read the request JSON data and convert it to user of struct User
	if err := c.ShouldBindJSON(&payload); err != nil {
		// if error, return error using http in JSON format using the gin context
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) 
		return
	}

	// create models.user instance for the business layer and then db because db needs full disclosure
	// this is coming from the user payload
	user := models.User {
		Username: payload.Username,
		Email: payload.Email,
		Password: payload.Password,
	}

	// creates the database and return error in JSON format
	// tc.UserService connects to the services and call the CreateUser()
	if err := tc.UserService.RegisterUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User successfully created, kindly verify your email address"})
}

func (tc *UserController) VerifyEmail(c *gin.Context) { 
	verificationToken := c.Query("code")

	if err := tc.UserService.VerificationService(verificationToken); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid verification token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email successfully verified"})
}

// context gets into the body of your package. c is an instance of pointing to the gin.Context
func (tc *UserController) LoginUser(c *gin.Context) { 
	// dto instance for validation
	var payload dto.LoginUserDTO // user payload placeholder

	// read the request JSON data and convert it to user of struct User
	if err := c.ShouldBindJSON(&payload); err != nil {
		// if error, return error using http in JSON format using the gin context
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) 
		return
	}

	// model instance of payload that the business layer need to send to the database
	user := models.User {
		Username: payload.Username,
		Password: payload.Password,
	}

	// creates the database and return error in JSON format
	// tc.UserService connects to the services and call the LoginUser()
	token, err := tc.UserService.LoginUser(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User successfully logged in", "token": token})
}

func (tc *UserController) ForgotPassword(c *gin.Context) {
	var payload dto.ForgotPasswordDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := tc.UserService.ForgotPasswordService(payload.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "A reset link has been sent to the email"})
}

func (tc *UserController) ResetPassword(c *gin.Context) {
	var payload dto.ResetPasswordDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := tc.UserService.ResetPasswordService(payload.Token, payload.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password successfully reset"})
}

// 2FA / TOTP Routes

func (tc *UserController) EnableTOTP(c *gin.Context) {
	// get userID from claims
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID := userIDVal.(string)

	// call service to generate secret + save in DB
	qrCodeURL, err := tc.UserService.EnableTOTP(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"qr_code_url": qrCodeURL,
		"message": "TOTP enabled. Scan the QR code in your authenticator app, and complete login with 6-digit code.",
	})	
}

func (tc *UserController) EnableEmail2FA(c *gin.Context) {
	var payload dto.RequestOtpDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
	}

	_, err := tc.UserService.EnableEmail2FA(payload.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "OTP request successful. Please check your email"})
}

func (tc *UserController) VerifyEmailOTP(c *gin.Context) {
	userID := c.GetString("user_id") // take userId from JWT
	log.Println("UserID:", userID)

	var payload dto.VerifyOtpDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
	}

	token, err := tc.UserService.VerifyEmailOTP(userID, payload.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Email OTP successfully verified",
		"token": token,
	})
}

func (tc *UserController) VerifyTOTP(c *gin.Context) {
	// get userID from claims
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	
	userID := userIDVal.(string)

	var payload dto.VerifyTOTPDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	token, err := tc.UserService.VerifyTOTP(userID, payload.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Authenticator app verified successfully",
		"token":   token,
	})
}

func (tc *UserController) Disable2FA(c *gin.Context) {
	// get userID from claims
    userIDVal, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    userID := userIDVal.(string)

    token, err := tc.UserService.Disable2FA(userID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
		"message": "2FA has been disabled successfully, update header with new token",
		"token": token,
	})
}
