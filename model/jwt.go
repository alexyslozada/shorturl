package model

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type JWTCustomClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.StandardClaims
}
