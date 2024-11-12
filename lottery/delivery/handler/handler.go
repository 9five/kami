package handler

import (
	"fmt"
	"kami/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LotteryHandler struct {
	PrizePoolUsecase domain.PrizePoolUsecase
	PrizeCardUsecase domain.PrizeCardUsecase
	KamiUserUsercase domain.KamiUserUsercase
}

func NewLotteryHandler(router *gin.RouterGroup, prizePoolUsecase domain.PrizePoolUsecase, prizeCardUsecase domain.PrizeCardUsecase, kamiUserUsecase domain.KamiUserUsercase) {
	handler := &LotteryHandler{
		PrizePoolUsecase: prizePoolUsecase,
		PrizeCardUsecase: prizeCardUsecase,
		KamiUserUsercase: kamiUserUsecase,
	}

	router.GET("/prizePool", handler.GetPrizePools)
	router.GET("/collection", handler.GetCollection)
	router.GET("/collection/detail", handler.GetCollectionDetail)
	router.POST("/draw", handler.DrawCard)
}

// @Summary     提取現有抽獎池
// @Description 提取現有抽獎池
// @Tags        lottery
// @Accept      json
// @Produce     json
// @Success     200  {object}  []domain.PrizePool
// @Failure     400  {object}  map[string]interface{}
// @Router      /api/lottery/prizePool [get]
func (l *LotteryHandler) GetPrizePools(ctx *gin.Context) {
	prizePoolList, err := l.PrizePoolUsecase.GetPrizePoolList(ctx, &domain.PrizePool{})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, prizePoolList)
}

// @Summary     取得user的collection
// @Description 取得user的collection
// @Tags        lottery
// @Accept      json
// @Produce     json
// @Param       pid   query  int  true  "prize pool id"
// @Success     200  {object}  domain.PrizeCardCollection
// @Failure     400  {object}  map[string]interface{}
// @Router      /api/lottery/collection [get]
func (l *LotteryHandler) GetCollection(ctx *gin.Context) {
	userId := ctx.MustGet("id").(uint)
	pid := ctx.Query("pid")
	var pidInt int64
	if pid != "" {
		cache, err := strconv.ParseInt(pid, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		pidInt = cache
	}

	prizePoolList, err := l.PrizePoolUsecase.GetPrizePoolList(ctx, &domain.PrizePool{Model: gorm.Model{ID: uint(pidInt)}})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var collection domain.PrizeCardCollection
	for _, prizePool := range prizePoolList {
		if c, err := l.PrizeCardUsecase.GetPrizeCardCollection(ctx, userId, prizePool); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		} else {
			collection.PoolName = c.PoolName
			collection.Total += c.Total
			for _, cc := range c.Cards {
				cc.Style = prizePool.Style
				collection.Cards = append(collection.Cards, cc)
			}
		}
	}

	ctx.JSON(http.StatusOK, collection)
}

// @Summary     取得user的collection detail
// @Description 取得user的collection detail
// @Tags        lottery
// @Accept      json
// @Produce     json
// @Param       cid   query  int  true  "prize card id"
// @Success     200  {object}  map[string]interface{}
// @Failure     400  {object}  map[string]interface{}
// @Router      /api/lottery/collection/detail [get]
func (l *LotteryHandler) GetCollectionDetail(ctx *gin.Context) {
	userId := ctx.MustGet("id").(uint)
	cid := ctx.Query("cid")
	if cid == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "please enter card id",
		})
		return
	}

	cidInt, err := strconv.ParseInt(cid, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	serialNumber, coupon, card, err := l.PrizeCardUsecase.GetPrizeCardCollectionDetail(ctx, &domain.UserPrizeCard{UserId: userId, CardId: uint(cidInt)})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"serial_number": serialNumber,
		"coupon":        coupon,
		"card":          card,
	})
}

// @Summary     抽獎
// @Description 抽獎
// @Tags        lottery
// @Accept      json
// @Produce     json
// @Param       pid   query  int  true  "prize pool id"
// @Success     200  {object}  domain.PrizeCardOutput
// @Failure     400  {object}  map[string]interface{}
// @Router      /api/lottery/draw [post]
func (l *LotteryHandler) DrawCard(ctx *gin.Context) {
	userId := ctx.MustGet("id").(uint)
	pid := ctx.Query("pid")
	pidInt, err := strconv.ParseInt(pid, 10, 64)
	user, err := l.KamiUserUsercase.GetKamiUser(ctx, &domain.KamiUser{Model: gorm.Model{ID: userId}})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	prizePool, err := l.PrizePoolUsecase.GetPrizePool(ctx, &domain.PrizePool{Model: gorm.Model{ID: uint(pidInt)}})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	prizeCardList, err := l.PrizeCardUsecase.GetPrizeCardList(ctx, &domain.PrizeCard{PoolId: uint(pidInt)})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	prizeCardList, err = l.PrizeCardUsecase.GetWeightedRandomList(ctx, userId, prizeCardList)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	fmt.Printf("first: %+v\n", user)
	if err := l.PrizePoolUsecase.SubtractUserPoints(ctx, user, prizePool); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	fmt.Printf("second: %+v\n", user)
	prize, err := l.PrizeCardUsecase.Draw(ctx, userId, prizeCardList)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := l.KamiUserUsercase.UpdateKamiUser(ctx, user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, prize)
}
