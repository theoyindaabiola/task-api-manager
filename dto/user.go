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

type Toggle2FARequest struct {
	Enabled 		bool `json:"enabled"`
}

type Toggle2FAResponse struct {
    UserID     		string `json:"user_id"`
    Enabled2FA 		bool   `json:"enabled_2fa"`
    Message    		string `json:"message"`
	Token			string `json:"token"`
}

type RequestOtpDTO struct {
	Email 			string `json:"email" binding:"required"`
}

type VerifyOtpDTO struct {
	OTP				string `json:"otp" binding:"required"`
}

type UserResponseDTO struct {
	IsOtpVerified	bool `json:"is_otp_verified"`
}
