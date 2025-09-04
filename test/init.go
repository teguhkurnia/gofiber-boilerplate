package test

import (
	"gofiber-boilerplate/internal/config"
	"gofiber-boilerplate/internal/util"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var app *fiber.App

var db *gorm.DB

var viperConfig *viper.Viper

var log *logrus.Logger

var validate *validator.Validate

func init() {
	viperConfig = config.NewViper()
	log = config.NewLogger(viperConfig)
	validate = config.NewValidator()
	app = config.NewFiber(viperConfig)
	db = config.NewDatabase(viperConfig, log, true)
	redis := config.NewRedis(viperConfig)
	tokenUtil := util.NewTokenUtil(redis, viperConfig, log)

	config.Bootstrap(&config.BootstrapConfig{
		DB:        db,
		App:       app,
		Log:       log,
		Validate:  validate,
		Config:    viperConfig,
		Redis:     redis,
		TokenUtil: tokenUtil,
	})
}
