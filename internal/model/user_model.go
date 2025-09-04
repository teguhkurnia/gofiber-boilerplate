package model

type UserResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Token string `json:"token,omitempty"`
}

type RegisterUserRequest struct {
	Name     string `json:"name" validate:"required,alphanum,min=5"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthorizeUserRequest struct {
	Token string
}
