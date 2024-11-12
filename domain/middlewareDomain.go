package domain

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type Middleware struct {
}

type MiddlewareRepository interface {
}

type MiddlewareUsecase interface {
	VerifyToken(ctx *gin.Context)
}

type TokenInfo struct {
	ID    uint   `json:"id"`
	Phone string `json:"phone"`
}

type PrivateClaims struct {
	jwt.StandardClaims
	UID uint `json:"uId"`
}
