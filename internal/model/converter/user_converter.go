package converter

import (
	"gofiber-boilerplate/internal/entity"
	"gofiber-boilerplate/internal/model"
)

func UserToResponse(user *entity.User, token string) *model.UserResponse {
	return &model.UserResponse{
		ID:    uint(user.ID),
		Email: user.Email,
		Token: token,
	}
}
