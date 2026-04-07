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
		// public auth routes
		public.POST("/register", userController.CreateUser)
		public.GET("/verify-email", userController.VerifyEmail)
		public.POST("/login", userController.LoginUser)
		public.POST("/forgot-password", userController.ForgotPassword)
		public.POST("/reset-password", userController.ResetPassword)
	}

	// protected routes - require JWT
	protected := router.Group("/api/users")
	protected.Use(middleware.JWTAuthMiddleware())
	{
		// email-based 2FA routes
		protected.POST("/enable-email-2fa", userController.EnableEmail2FA)
		protected.POST("/verify-email-2fa", userController.VerifyEmailOTP)

		// totp-based 2FA routes
		protected.POST("/enable-totp", userController.EnableTOTP)
		protected.POST("/verify-totp", userController.VerifyTOTP)

		// disable 2FA
		protected.POST("/disable-2fa", userController.Disable2FA)
	}
}
