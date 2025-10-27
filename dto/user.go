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

type EnableSMSDTO struct {
	PhoneNumber		string `json:"phone_number" binding:"required"`
}

type UpdatePhoneNumberDTO struct {
	PhoneNumber		string `json:"phone_number" binding:"required"`
}

type RequestSMSDTO struct {
	PhoneNumber		string `json:"phone_number" binding:"required"`
}

type VerifySMSDTO struct {
	Code 			string `json:"code" binding:"required"`
}
