package dto

type User_RegisterRequest struct {
	ID       int64  `json:"id"`
	Email    string `json:"email,omitempty" binding:"required,email"`
	Password string `json:"password,omitempty" binding:"required,min=4,max=16"`
	Role     string `json:"role,omitempty"`
}

type User_LoginRequest struct {
	ID       int64  `json:"id"`
	Email    string `json:"email,omitempty" binding:"required,email"`
	Password string `json:"password,omitempty" binding:"required"`
}
