package dto

// this is an abstract/interface of the real model validates payloads
type Task struct {
	Title 		string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Completed 	string `json:"completed"`
}
