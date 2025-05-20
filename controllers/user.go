package controllers

import (
	"net/http" // allows hhtp requests
	"taskapi/dto"
	"taskapi/models"
	"taskapi/services"

	// "github.com/google/uuid"
	"github.com/gin-gonic/gin" // web framework
)

/**
	we need the service instance to be able to send the payload to it and the functions
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
	var payload dto.RegisterUserDTO // placeholder to hold the user/payload to be proccessed
	
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
		Email: payload.Email, // commented out fpr error testing
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
	var payload dto.LoginUserDTO // placeholder to hold the user/payload to be proccessed

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
	// tc.UserService connects to the services and call the CreateUser()
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


