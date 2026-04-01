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
		// 2FA routes allowed in middleware

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
		protected.POST("/enable-email-2fa", userController.EnableEmail2FA)
		protected.POST("/verify-email-2fa", userController.VerifyEmailOTP)
		protected.POST("/disable-email-2fa", userController.DisableEmail2FA)
		// TOTP
		protected.POST("/enable-totp", userController.EnableTOTP)
		protected.POST("/verify-totp", userController.VerifyTOTP)
		protected.POST("/disable-totp", userController.DisableTOTP)
	}
}
