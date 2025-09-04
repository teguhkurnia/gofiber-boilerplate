package usecase

import (
	"context"
	"gofiber-boilerplate/internal/entity"
	"gofiber-boilerplate/internal/model"
	"gofiber-boilerplate/internal/model/converter"
	"gofiber-boilerplate/internal/repository"
	"gofiber-boilerplate/internal/util"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUseCase struct {
	DB             *gorm.DB
	Redis          *redis.Client
	Log            *logrus.Logger
	Validate       *validator.Validate
	UserRepository *repository.UserRepository
	TokenUtil      *util.TokenUtil
}

func NewUserUseCase(db *gorm.DB, tokenUtil *util.TokenUtil, log *logrus.Logger,
	validate *validator.Validate, userRepo *repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		DB:             db,
		Log:            log,
		Validate:       validate,
		UserRepository: userRepo,
		TokenUtil:      tokenUtil,
	}
}

func (c *UserUseCase) Create(ctx context.Context, request *model.RegisterUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body: %v", err)
		return nil, fiber.ErrBadRequest
	}

	total, err := c.UserRepository.CountByEmail(tx, request.Email)
	if err != nil {
		c.Log.Warnf("Failed to count user: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	if total > 0 {
		c.Log.Warnf("Email already exists: %s", request.Email)
		return nil, fiber.ErrConflict
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Log.Warnf("Failed to hash password: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	now := time.Now()
	verificationCode := uuid.New().String()
	user := &entity.User{
		Name:             request.Name,
		Email:            request.Email,
		Password:         string(password),
		VerificationCode: &verificationCode,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if err := c.UserRepository.Create(tx, user); err != nil {
		c.Log.Warnf("Failed to create user: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	token, err := c.TokenUtil.CreateJWTToken(ctx, user.ID)
	if err != nil {
		c.Log.Warnf("Failed to create JWT token: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.UserToResponse(user, token), nil
}

func (c *UserUseCase) Login(ctx context.Context,
	request *model.LoginUserRequest) (*model.UserResponse, error) {

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body: %v", err)
		return nil, fiber.ErrBadRequest
	}

	user, err := c.UserRepository.FindByEmail(c.DB, request.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.Log.Warnf("User not found: %s", request.Email)
			return nil, fiber.ErrUnauthorized
		}
		c.Log.Warnf("Failed to find user: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		c.Log.Warnf("Invalid password for user: %s", request.Email)
		return nil, fiber.ErrUnauthorized
	}

	token, err := c.TokenUtil.CreateJWTToken(ctx, user.ID)
	if err != nil {
		c.Log.Warnf("Failed to create JWT token: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.UserToResponse(user, token), nil
}

func (c *UserUseCase) Verify(ctx context.Context, code string) error {
	user, err := c.UserRepository.FindByVerificationCode(c.DB, code)
	if err != nil {
		c.Log.Warnf("Failed to find user by verification code: %v", err)
		return fiber.ErrNotFound
	}

	// Mark user as verified
	now := time.Now()
	user.VerifiedAt = &now
	user.VerificationCode = nil
	if err := c.UserRepository.Update(c.DB, user); err != nil {
		c.Log.Warnf("Failed to update user: %v", err)
		return fiber.ErrInternalServerError
	}

	return nil
}

func (c *UserUseCase) Authorize(ctx context.Context, token string) (*model.Auth, error) {
	auth, err := c.TokenUtil.ParseJWTToken(ctx, token)
	if err != nil {
		c.Log.Warnf("Failed to parse JWT token: %v", err)
		return nil, fiber.ErrUnauthorized
	}

	return auth, nil
}
