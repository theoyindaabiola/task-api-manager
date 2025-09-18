package routes

import (
	"taskapi/controllers"
	"github.com/gin-gonic/gin"
)

type UserRoutes struct {
	UserController *controllers.UserController
}

func RegisterUserRoutes(router *gin.Engine, userController *controllers.UserController) {
	// group the routes
	userRoutes := router.Group("/api/users")
	{
		userRoutes.POST("/register", userController.CreateUser)
		userRoutes.GET("/verify-email", userController.VerifyEmail)
		userRoutes.POST("/login", userController.LoginUser)
		userRoutes.POST("/forgot-password", userController.ForgotPassword)
		userRoutes.POST("/reset-password", userController.ResetPassword)
		userRoutes.POST("/request-otp", userController.RequestOTP)
		userRoutes.POST("/verify-otp", userController.VerifyOTP)
	}
}
