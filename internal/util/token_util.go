package util

import (
	"context"
	"errors"
	"gofiber-boilerplate/internal/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type TokenUtil struct {
	Redis  *redis.Client
	config *viper.Viper
	logger *logrus.Logger
}

func NewTokenUtil(redis *redis.Client, config *viper.Viper, logger *logrus.Logger) *TokenUtil {
	return &TokenUtil{
		Redis:  redis,
		config: config,
		logger: logger,
	}
}

func (t *TokenUtil) ParseJWTToken(context context.Context, token string) (*model.Auth, error) {
	const BEARER_SCHEMA = "Bearer "
	if len(token) > len(BEARER_SCHEMA) && token[:len(BEARER_SCHEMA)] == BEARER_SCHEMA {
		token = token[len(BEARER_SCHEMA):]
	}
	secretKey := t.config.GetString("jwt.secret")

	t.logger.Infof("Parsing token: %s", token)
	t.logger.Infof("Using secret key: %s", secretKey)

	parsed, err := jwt.ParseWithClaims(token, &model.Auth{}, func(token *jwt.Token) (any, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		t.logger.Warnf("Failed to parse token: %v", err)
		return nil, err
	}

	exist, err := t.Redis.Exists(context, token).Result()
	if err != nil {
		t.logger.Errorf("Failed to check token in Redis: %v", err)
		return nil, err
	}

	if exist == 0 {
		return nil, errors.New("token not found or expired")
	}

	claims, ok := parsed.Claims.(*model.Auth)
	if !ok || !parsed.Valid {
		t.logger.Warn("Invalid token claims")
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}

func (t *TokenUtil) CreateJWTToken(context context.Context, uid uint64) (string, error) {
	secretKey := t.config.GetString("jwt.secret")
	ttl := t.config.GetInt("jwt.ttl")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &model.Auth{
		UserID: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(ttl) * time.Minute)),
		},
	})
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		t.logger.Errorf("Failed to sign token: %v", err)
		return "", err
	}

	_, err = t.Redis.SetEx(context, signedToken, uid, time.Duration(ttl)*time.Minute).Result()
	if err != nil {
		t.logger.Errorf("Failed to store token in Redis: %v", err)
		return "", err
	}

	return signedToken, nil
}
