package handler

import (
	"kami/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type KamiOrderHandler struct {
	KamiOrderUsecase domain.KamiOrderUsecase
	KamiUserUsercase domain.KamiUserUsercase
}

func NewKamiOrderHandler(router *gin.RouterGroup, kamiOrderUsecase domain.KamiOrderUsecase, kamiUserUsecase domain.KamiUserUsercase) {
	handler := &KamiOrderHandler{
		KamiOrderUsecase: kamiOrderUsecase,
		KamiUserUsercase: kamiUserUsecase,
	}

	router.PUT("/register", handler.OrderRegister)
	router.GET("/getOrders", handler.GetOrders)
}

// @Summary     找到訂單並更新訂單擁有人
// @Description 用戶輸入訂單資料
// @Tags        order
// @Accept      json
// @Produce     json
// @Param       order  body  domain.OrderInput  true  "order info"
// @Success     200  {object}  map[string]interface{}
// @Failure     400  {object}  map[string]interface{}
// @Router      /api/order/register [put]
func (k *KamiOrderHandler) OrderRegister(ctx *gin.Context) {
	phone := ctx.GetString("phone")
	var input domain.OrderInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := k.KamiUserUsercase.GetKamiUser(ctx, &domain.KamiUser{Phone: phone})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	order, err := k.KamiOrderUsecase.CheckOrderInput(ctx, &input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := k.KamiOrderUsecase.RegisterOrder(ctx, order, user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := k.KamiUserUsercase.UpdateKamiUser(ctx, user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

// @Summary     取得用戶order list
// @Description 取得用戶order list
// @Tags        order
// @Produce     json
// @param startDate query string true "開始時間 etc. 2006-01"
// @param endDate query string true "結束時間 etc. 2006-01"
// @Success     200  {object}  []domain.KamiOrder
// @Failure     400  {object}  map[string]interface{}
// @Router      /api/order/getOrders [get]
func (k *KamiOrderHandler) GetOrders(ctx *gin.Context) {
	phone := ctx.MustGet("phone").(string)
	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")
	orders, err := k.KamiOrderUsecase.GetOrderListByDate(ctx, &domain.KamiOrder{OwnerPhone: phone}, startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, orders)
}
