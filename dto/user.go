package dto

type RegisterUserDTO struct {
	Username 	string `Json:"username" binding:"required"`
	Email 		string `Json:"email" binding:"required"`
	Password	string `json:"password" binding:"required"`
}

type LoginUserDTO struct {
	Username 	string `Json:"username" binding:"required"`
	Password	string `json:"password" binding:"required"`
}

type ForgotPasswordDTO struct {
	Email 		string `Json:"email" binding:"required"`
}

type ResetPasswordDTO struct {
	Token 		string `Json:"token" binding:"required"`
	Password	string `json:"password" binding:"required"`
}