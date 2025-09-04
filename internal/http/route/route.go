package route

import (
	"gofiber-boilerplate/internal/http"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App            *fiber.App
	UserController *http.UserController
	AuthMiddleware fiber.Handler
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupAuthRoute() {
	c.App.Use(c.AuthMiddleware)

	c.App.Get("api/v1/profile", func(ctx *fiber.Ctx) error {
		user := ctx.Locals("user")
		return ctx.JSON(fiber.Map{
			"status": "success",
			"user":   user,
		})
	})
}

func (c *RouteConfig) SetupGuestRoute() {
	c.App.Post("api/v1/login", c.UserController.Login)
	c.App.Post("api/v1/register", c.UserController.Register)
}
