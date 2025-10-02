package controllers

import (
	"net/http" // allows hhtp requests
	"taskapi/dto"
	"taskapi/models"
	"taskapi/services"
	"taskapi/utils"

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

	c.JSON(http.StatusCreated, "User successfully created, kindly verify your email address")
}

func (tc *UserController) VerifyEmail(c *gin.Context) { 
	verificationToken := c.Query("code")

	if err := tc.UserService.VerificationService(verificationToken); err != nil {
		c.JSON(400, gin.H{"error": "Invalid verification token"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User successfully logged in", "token": token})
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

func (tc *UserController) EnableSMS(c *gin.Context) {
	var payload dto.EnableSMSDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Validate Nigerian phone number
    if !utils.IsValidPhoneNumber(payload.PhoneNumber) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Nigerian phone number"})
        return
    }
	
	// get userID from claims
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDVal.(string)

	// call service to generate secret + save in DB
	if err := tc.UserService.EnableSMS(userID, payload.PhoneNumber); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "SMS 2FA enabled. Check your phone for the OTP and verify",
	})	
}

func (tc *UserController) UpdatePhoneNumber(c *gin.Context) {
	var payload dto.UpdatePhoneNumberDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Validate Nigerian phone number
    if !utils.IsValidPhoneNumber(payload.PhoneNumber) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Nigerian phone number"})
        return
    }

	// get userID from claims
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDVal.(string)

	// call service to generate secret to update DB
	if err := tc.UserService.UpdatePhoneNumber(userID, payload.PhoneNumber); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Phone number updated. OTP sent for verification. SMS 2FA will be active after verification.",
	})	
}

func (tc *UserController) RequestSMS(c *gin.Context) {
    // get userID from claims
	userIDVal, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    userID := userIDVal.(string)

    err := tc.UserService.RequestSMS(userID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "A new SMS OTP has been sent to your phone.",
    })
}

func (tc *UserController) VerifySMS(c *gin.Context) {
	// get userID from claims
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDVal.(string)

	var payload dto.VerifySMSDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	token, err := tc.UserService.VerifySMS(userID, payload.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "SMS OTP verified successfully",
		"token":   token,
	})
}

func (tc *UserController) DisableSMS(c *gin.Context) {
	// get userID from claims
    userIDVal, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    userID := userIDVal.(string)

    token, err := tc.UserService.DisableSMS(userID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
		"message": "2FA has been disabled successfully, update header with the new token",
		"token": token,
	})
}
