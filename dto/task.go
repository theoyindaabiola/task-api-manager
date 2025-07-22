package dto

// this is an abstract/interface of the real model validates payloads
type Task struct {
	Title 		string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Completed 	string `json:"completed"`
}

type TaskDelegationInput struct {
	DelegateeID     string `json:"delegatee_id" binding:"required"`
	Permission 		string `json:"permission" binding:"required"` // expects "R", or "U" , "O" in this context
}

