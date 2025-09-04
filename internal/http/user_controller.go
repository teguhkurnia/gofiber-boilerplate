package http

import (
	"gofiber-boilerplate/internal/model"
	"gofiber-boilerplate/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	Log     *logrus.Logger
	UseCase *usecase.UserUseCase
}

func NewUserController(log *logrus.Logger,
	useCase *usecase.UserUseCase) *UserController {
	return &UserController{
		Log:     log,
		UseCase: useCase,
	}
}

func (c *UserController) Register(ctx *fiber.Ctx) error {
	request := new(model.RegisterUserRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}

	request.IpAddress = ctx.IP()
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

	request.IpAddress = ctx.IP()
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
