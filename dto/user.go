package dto

type RegisterUserDTO struct {
	Username 		string `json:"username" binding:"required"`
	Email 			string `json:"email" binding:"required"`
	Password		string `json:"password" binding:"required"`
}

type LoginUserDTO struct {
	Username 		string `json:"username" binding:"required"`
	Password		string `json:"password" binding:"required"`
}

type ForgotPasswordDTO struct {
	Email 			string `json:"email" binding:"required"`
}

type ResetPasswordDTO struct {
	Token 			string `json:"token" binding:"required"`
	Password		string `json:"password" binding:"required"`
}

type VerifyTOTPDTO struct {
	Code 			string `json:"code" binding:"required"`
}

type RequestOtpDTO struct {
	Email 			string `json:"email" binding:"required"`
}

type VerifyOtpDTO struct {
	OTP				string `json:"otp" binding:"required"`
}

type UserResponseDTO struct {
	EmailOtpVerified	bool `json:"email_otp_verified"`
	TOTPVerified		bool `json:"totp_verified"`
}

