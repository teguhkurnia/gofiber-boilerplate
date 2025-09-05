package config

import (
	"gofiber-boilerplate/internal/http"
	"gofiber-boilerplate/internal/http/middleware"
	"gofiber-boilerplate/internal/http/route"
	"gofiber-boilerplate/internal/repository"
	"gofiber-boilerplate/internal/usecase"
	"gofiber-boilerplate/internal/util"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	App             *fiber.App
	DB              *gorm.DB
	Redis           *redis.Client
	Log             *logrus.Logger
	Validate        *validator.Validate
	Config          *viper.Viper
	TokenUtil       *util.TokenUtil
	RateLimiterUtil *util.RateLimiterUtil
}

func Bootstrap(config *BootstrapConfig) {
	// setup repositories
	userRepo := repository.NewUserRepository(config.Log)

	// setup usecases
	userUseCase := usecase.NewUserUseCase(
		config.DB,
		config.TokenUtil,
		config.Log,
		config.Validate,
		userRepo,
	)

	// setup controllers
	userController := http.NewUserController(
		config.Log,
		userUseCase,
		config.RateLimiterUtil,
	)

	// setup routes
	authMiddleware := middleware.NewAuthMiddleware(userUseCase)
	routeConfig := &route.RouteConfig{
		App:            config.App,
		UserController: userController,
		AuthMiddleware: authMiddleware,
	}

	routeConfig.Setup()
}
