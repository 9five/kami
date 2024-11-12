package handler

import (
	"context"
	"errors"
	"fmt"
	"kami/domain"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LoginHandler struct {
	KamiUserUsecase      domain.KamiUserUsercase
	TwilioServiceUsecase domain.TwilioServiceUsecase
	EncryptionUsecase    domain.EncryptionUsecase
}

func NewLoginHandler(router *gin.RouterGroup, kamiUserUsecase domain.KamiUserUsercase, twilioServiceUsecase domain.TwilioServiceUsecase, encryptionUsecase domain.EncryptionUsecase) {
	handler := &LoginHandler{
		KamiUserUsecase:      kamiUserUsecase,
		TwilioServiceUsecase: twilioServiceUsecase,
		EncryptionUsecase:    encryptionUsecase,
	}

	router.POST("/enterPhone", handler.EnterPhone)
	router.POST("/verificationCheck", handler.VerificationCheck)
	router.POST("/enterPassword", handler.EnterPassword)
	router.POST("/forgotPassword", handler.ForgotPassword)
}

// @Summary     取得用戶電話號碼並寄出驗證碼
// @Description 用戶輸入電話號碼
// @Tags        login
// @Accept      json
// @Produce     json
// @Param       phone   query  string  true  "user's phone number"
// @Success     200  {object}  map[string]interface{}
// @Failure     400  {object}  map[string]interface{}
// @Router      /api/login/enterPhone [post]
func (l *LoginHandler) EnterPhone(ctx *gin.Context) {
	phone := ctx.Query("phone")
	if phone == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "please enter correct phone number",
		})
		return
	}

	status := true
	if user, err := l.KamiUserUsecase.GetKamiUser(ctx, &domain.KamiUser{Phone: phone}); err != nil && err != gorm.ErrRecordNotFound {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	} else if err == gorm.ErrRecordNotFound || user.Password == "" {
		if err := l.KamiUserUsecase.CheckKamiUserLog(ctx, &domain.KamiUserLog{Phone: phone}); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := l.TwilioServiceUsecase.SendVerificationSMS(phone); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		status = false
	}

	userLog, err := l.KamiUserUsecase.GetKamiUserLog(ctx, &domain.KamiUserLog{Phone: phone})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	token, err := l.EncryptionUsecase.DesEncrypt(ctx, fmt.Sprintf("%s-%d", phone, userLog.Model.ID))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token":  token,
		"status": status,
	})
}

// @Summary     驗證並回傳token
// @Description 用戶輸入驗證碼
// @Tags        login
// @Accept      json
// @Produce     json
// @Param       verification   body  domain.VerificationInput  true  "login token & 6 digit verification code"
// @Success     200  {object}  map[string]interface{}
// @Failure     400  {object}  map[string]interface{}
// @Router      /api/login/verificationCheck [post]
func (l *LoginHandler) VerificationCheck(ctx *gin.Context) {
	var input domain.VerificationInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	phone, _, err := l.loginTokenDecode(ctx, input.Token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := l.TwilioServiceUsecase.VerificationCheck(phone, input.Code); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if input.ForgotPw {
		if user, err := l.KamiUserUsecase.GetKamiUser(ctx, &domain.KamiUser{Phone: phone}); err != nil {
			user.Password = ""
			if err = l.KamiUserUsecase.UpdateKamiUser(ctx, user); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}
		}
	}

	ctx.JSON(http.StatusOK, nil)
}

// @Summary     驗證密碼並回傳token
// @Description 用戶輸入密碼
// @Tags        login
// @Accept      json
// @Produce     json
// @Param       token   query  string  true  "login token"
// @Param       password   query  string  true  "user's password"
// @Success     200  {object}  map[string]interface{}
// @Failure     400  {object}  map[string]interface{}
// @Router      /api/login/enterPassword [post]
func (l *LoginHandler) EnterPassword(ctx *gin.Context) {
	token := ctx.Query("token")
	password := ctx.Query("password")
	if token == "" || password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "please enter correct token or password",
		})
		return
	}

	phone, _, err := l.loginTokenDecode(ctx, token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	pwHash, err := l.EncryptionUsecase.HashPassword(ctx, password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := l.KamiUserUsecase.LoginKamiUser(ctx, &domain.KamiUser{Phone: phone, Password: password})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if !l.EncryptionUsecase.CheckPwHash(ctx, password, pwHash) {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "password incorrect",
		})
		return

	}

	userLog, err := l.KamiUserUsecase.GetKamiUserLog(ctx, &domain.KamiUserLog{Phone: phone})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userLog.AuthFrequency = 0
	if err = l.KamiUserUsecase.UpdateKamiUserLog(ctx, userLog); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	jwtToken, err := l.KamiUserUsecase.GenerateToken(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": jwtToken,
		"kami_user": domain.KamiUserOutput{
			Phone:  user.Phone,
			Email:  user.Email,
			Status: user.Status,
		},
	})
}

// @Summary     取得用戶電話號碼並寄出驗證碼
// @Description 用戶輸入電話號碼
// @Tags        login
// @Accept      json
// @Produce     json
// @Param       token   query  string  true  "login token"
// @Success     200  {object}  map[string]interface{}
// @Failure     400  {object}  map[string]interface{}
// @Router      /api/login/forgotPassword [post]
func (l *LoginHandler) ForgotPassword(ctx *gin.Context) {
	token := ctx.Query("token")
	if token == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "please enter correct phone number",
		})
		return
	}

	phone, _, err := l.loginTokenDecode(ctx, token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if _, err := l.KamiUserUsecase.GetKamiUser(ctx, &domain.KamiUser{Phone: phone}); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := l.KamiUserUsecase.CheckKamiUserLog(ctx, &domain.KamiUserLog{Phone: phone}); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := l.TwilioServiceUsecase.SendVerificationSMS(phone); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

func (l *LoginHandler) loginTokenDecode(ctx context.Context, token string) (string, string, error) {
	tokenText, err := l.EncryptionUsecase.DesDecrypt(ctx, token)
	if err != nil {
		return "", "", err
	}

	tokenSlice := strings.Split(tokenText, "-")
	if len(tokenSlice) != 2 {
		return "", "", errors.New("token incorrect")
	}

	return tokenSlice[0], tokenSlice[1], nil
}
