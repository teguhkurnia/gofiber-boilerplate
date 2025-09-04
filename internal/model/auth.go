package model

import "github.com/golang-jwt/jwt/v5"

type Auth struct {
	UserID uint64
	jwt.RegisteredClaims
}
