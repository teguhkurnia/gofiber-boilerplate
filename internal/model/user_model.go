package model

type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Token    string `json:"token,omitempty"`
}

type RegisterUserRequest struct {
	Username  string `json:"username" validate:"required,alphanum,min=5"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	Country   string `json:"country" validate:"required"`
	IpAddress string `json:"-"`
}

type LoginUserRequest struct {
	Identifier string `json:"identifier" validate:"required"`
	Password   string `json:"password" validate:"required"`
	IpAddress  string `json:"-"`
}

type AuthorizeUserRequest struct {
	Token string
}
