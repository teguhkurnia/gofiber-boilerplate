package http

import (
	"gofiber-boilerplate/internal/model"
	"gofiber-boilerplate/internal/usecase"
	"gofiber-boilerplate/internal/util"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	Log             *logrus.Logger
	UseCase         *usecase.UserUseCase
	RateLimiterUtil *util.RateLimiterUtil
}

func NewUserController(log *logrus.Logger,
	useCase *usecase.UserUseCase,
	rateLimiterUtil *util.RateLimiterUtil,
) *UserController {
	return &UserController{
		Log:             log,
		UseCase:         useCase,
		RateLimiterUtil: rateLimiterUtil,
	}
}

func (c *UserController) Register(ctx *fiber.Ctx) error {
	request := new(model.RegisterUserRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}

	reponse, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to create user: %v", err)
		return err
	}

	return ctx.Status(fiber.StatusCreated).
		JSON(model.WebResponse[*model.UserResponse]{
			Data: reponse,
		})
}

func (c *UserController) Login(ctx *fiber.Ctx) error {
	request := new(model.LoginUserRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}

	err := c.RateLimiterUtil.IsAllowed(ctx, "login:"+ctx.IP(), nil)
	if err != nil {
		c.Log.Warnf("Rate limit exceeded: %v", err)
		return err
	}

	reponse, err := c.UseCase.Login(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to login user: %v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).
		JSON(model.WebResponse[*model.UserResponse]{
			Data: reponse,
		})
}

func (c *UserController) Verify(ctx *fiber.Ctx) error {
	code := ctx.Params("code")
	if code == "" {
		c.Log.Warnf("Missing verification code")
		return fiber.ErrBadRequest
	}

	if err := c.UseCase.Verify(ctx.UserContext(), code); err != nil {
		c.Log.Warnf("Failed to verify user: %v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).
		JSON(model.WebResponse[any]{
			Data: "User verified successfully",
		})
}
