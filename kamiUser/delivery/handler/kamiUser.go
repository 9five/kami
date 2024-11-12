package handler

import (
	"kami/domain"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type KamiUserHandler struct {
	KamiUserUsecase domain.KamiUserUsercase
}

func NewKamiUserHandler(router *gin.RouterGroup, kamiUserUsecase domain.KamiUserUsercase) {
	handler := &KamiUserHandler{
		KamiUserUsecase: kamiUserUsecase,
	}

	router.GET("/status", handler.GetCurrentUser)
	router.PUT("/updateInfo", handler.UpdateCurrentUserInfo)
}

// @Summary     取得當前用戶資料
// @Description 取得當前用戶資料
// @Tags        user
// @Accept      json
// @Produce     json
// @Success     200  {object}  domain.KamiUserOutput
// @Failure     400  {object}  map[string]interface{}
// @Router      /api/user/status [get]
func (k *KamiUserHandler) GetCurrentUser(ctx *gin.Context) {
	userId := ctx.MustGet("id").(uint)
	user, err := k.KamiUserUsecase.GetKamiUser(ctx, &domain.KamiUser{Model: gorm.Model{ID: userId}})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, &domain.KamiUserOutput{
		Phone:    user.Phone,
		Email:    user.Email,
		Status:   user.Status,
		Points:   user.Points,
		Gender:   user.Gender,
		Birthday: user.Birthday.Format("2006-01-02"),
		Name:     user.Name,
		Career:   user.Career,
	})
}

// @Summary     更新用戶資料
// @Description 更新用戶資料
// @Tags        user
// @Accept      json
// @Produce     json
// @Param       user  body  domain.KamiUserInput  true  "user info"
// @Success     200  {object}  domain.KamiUserOutput
// @Failure     400  {object}  map[string]interface{}
// @Router      /api/user/updateInfo [put]
func (k *KamiUserHandler) UpdateCurrentUserInfo(ctx *gin.Context) {
	userId := ctx.MustGet("id").(uint)

	var input domain.KamiUserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := k.KamiUserUsecase.GetKamiUser(ctx, &domain.KamiUser{Model: gorm.Model{ID: userId}})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := k.KamiUserUsecase.UpdateUserInfo(ctx, user, &input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, user)
}
