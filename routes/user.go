package routes

import (
	"taskapi/controllers"
	"taskapi/middleware"

	"github.com/gin-gonic/gin"
)

type UserRoutes struct {
	UserController *controllers.UserController
}

func RegisterUserRoutes(router *gin.Engine, userController *controllers.UserController) {
	public := router.Group("/api/users")
	{
		public.POST("/register", userController.CreateUser)
		public.GET("/verify-email", userController.VerifyEmail)
		public.POST("/login", userController.LoginUser)
		public.POST("/forgot-password", userController.ForgotPassword)
		public.POST("/reset-password", userController.ResetPassword)
	}

	// require JWT
	protected := router.Group("/api/users")
	protected.Use(middleware.JWTAuthMiddleware())
	{
		protected.POST("/enable-otp", userController.EnableEmail2FA)
		protected.POST("/verify-otp", userController.VerifyEmailOTP)
		protected.POST("/disable-otp", userController.DisableEmail2FA)
	}
}
