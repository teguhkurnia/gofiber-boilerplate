package middleware

import (
	"gofiber-boilerplate/internal/model"
	"gofiber-boilerplate/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

func NewAuthMiddleware(userUseCase *usecase.UserUseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := &model.AuthorizeUserRequest{
			Token: c.Get("Authorization", "NOT_FOUND"),
		}

		auth, err := userUseCase.Authorize(c.Context(), request.Token)
		if err != nil {
			return fiber.ErrUnauthorized
		}

		c.Locals("user", auth)
		return c.Next()
	}
}
