package routes

import (
	"taskapi/controllers"
	"taskapi/middleware"
	"taskapi/services"

	"github.com/gin-gonic/gin"
)

type UserRoutes struct {
	UserController *controllers.UserController
}

func RegisterUserRoutes(router *gin.Engine, userController *controllers.UserController, userService *services.UserService) {
	// group the routes
	userRoutes := router.Group("/api/users")
	{
		userRoutes.POST("/register", userController.CreateUser)
		userRoutes.GET("/verify-email", userController.VerifyEmail)
		userRoutes.POST("/login", userController.LoginUser)
		userRoutes.POST("/forgot-password", userController.ForgotPassword)
		userRoutes.POST("/reset-password", userController.ResetPassword)

		// 2FA routes allowed in middleware
		userRoutes.POST("/enable-sms", middleware.JWTAuthMiddleware(userService), userController.EnableSMS)
		userRoutes.POST("/request-sms", middleware.JWTAuthMiddleware(userService), userController.RequestSMS)
		userRoutes.POST("/verify-sms", middleware.JWTAuthMiddleware(userService), userController.VerifySMS)
		userRoutes.POST("/disable-sms", middleware.JWTAuthMiddleware(userService), userController.DisableSMS)
	}
}
