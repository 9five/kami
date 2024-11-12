package usecase

import (
	"errors"
	"fmt"
	"kami/domain"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type middlewareUsecase struct {
	jwtKey []byte
}

func NewMiddlewareUsecase(jwtKey []byte) domain.MiddlewareUsecase {
	return &middlewareUsecase{
		jwtKey: jwtKey,
	}
}

func (m *middlewareUsecase) VerifyToken(ctx *gin.Context) {
	token, ok := getToken(ctx)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "token required",
		})
		return
	}
	tokenInfo, err := validateToken(m.jwtKey, token)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.Set("id", tokenInfo.ID)
	ctx.Set("phone", tokenInfo.Phone)
	ctx.Writer.Header().Set("Authorization", "Bearer "+token)
	ctx.Next()
}

func getToken(ctx *gin.Context) (string, bool) {
	authValue := ctx.GetHeader("Authorization")
	arr := strings.Split(authValue, " ")
	if len(arr) != 2 {
		return "", false
	}
	authType := strings.Trim(arr[0], "\n\r\t")
	if strings.ToLower(authType) != strings.ToLower("Bearer") {
		return "", false
	}
	return strings.Trim(arr[1], "\n\t\r"), true
}

func validateToken(jwtKey []byte, tokenString string) (*domain.TokenInfo, error) {
	var result domain.TokenInfo
	var claims domain.PrivateClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})
	if err != nil {
		return &result, err
	}
	if !token.Valid {
		return &result, errors.New("invalid token")
	}
	expiryDate := time.Unix(claims.ExpiresAt, 0)
	if time.Now().After(expiryDate) {
		return &result, errors.New("token expired")
	}
	result = domain.TokenInfo{
		ID:    claims.UID,
		Phone: claims.Subject,
	}
	return &result, nil
}
