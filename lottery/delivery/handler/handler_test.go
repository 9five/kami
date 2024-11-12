package handler_test

import (
	"encoding/json"
	"kami/domain"
	"kami/domain/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	lotteryHandler "kami/lottery/delivery/handler"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type LotteryHandlerSuite struct {
	suite.Suite
	prizePoolUsecase *mocks.PrizePoolUsecase
	prizeCardUsecase *mocks.PrizeCardUsecase
	kamiUserUsercase *mocks.KamiUserUsercase
	engine           *gin.Engine
	router           *gin.RouterGroup
	w                *httptest.ResponseRecorder
}

func TestStart(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	suite.Run(t, &LotteryHandlerSuite{})
}

func mockmiddleware(ctx *gin.Context) {
	ctx.Set("id", uint(1))
	ctx.Set("phone", "0900123123")
}

func (s *LotteryHandlerSuite) SetupTest() {
	s.prizePoolUsecase = new(mocks.PrizePoolUsecase)
	s.prizeCardUsecase = new(mocks.PrizeCardUsecase)
	s.kamiUserUsercase = new(mocks.KamiUserUsercase)
	s.engine = gin.Default()
	s.router = s.engine.Group("/api/lottery", mockmiddleware)
	s.w = httptest.NewRecorder()
}

func (s *LotteryHandlerSuite) TestGetPrizePools() {
	mockPrizePools := []*domain.PrizePool{
		{
			Model:  gorm.Model{ID: 1},
			Owner:  "ABC",
			Name:   "testingPool1",
			Points: 1,
		},
		{
			Model:  gorm.Model{ID: 2},
			Owner:  "ABC",
			Name:   "testingPool2",
			Points: 1,
		},
	}
	mockPrizePoolsMarshal, _ := json.Marshal(mockPrizePools)

	s.prizePoolUsecase.On("GetPrizePoolList", mock.Anything, &domain.PrizePool{}).
		Return(mockPrizePools, nil).Once()

	handler := lotteryHandler.LotteryHandler{
		PrizePoolUsecase: s.prizePoolUsecase,
		PrizeCardUsecase: s.prizeCardUsecase,
		KamiUserUsercase: s.kamiUserUsercase,
	}
	s.router.GET("/prizePool", handler.GetPrizePools)
	req, _ := http.NewRequest("GET", "/api/lottery/prizePool", nil)
	s.engine.ServeHTTP(s.w, req)

	assert.Equal(s.Suite.T(), http.StatusOK, s.w.Code)
	assert.Equal(s.Suite.T(), string(mockPrizePoolsMarshal), s.w.Body.String())
}

func (s *LotteryHandlerSuite) TestGetCollection() {
	mockPrizePools := []*domain.PrizePool{
		{
			Model:  gorm.Model{ID: 1},
			Owner:  "ABC",
			Name:   "testingPool1",
			Points: 1,
		},
	}

	mockPrizeCards := []*domain.PrizeCard{
		{
			Model:       gorm.Model{ID: 1},
			PoolId:      1,
			Title:       "testingCard",
			Description: "testingCardDescription",
			Probability: 100,
		},
	}

	mockUserPrizeCards := []*domain.UserPrizeCard{
		{
			ID:           1,
			UserId:       1,
			CardId:       1,
			SerialNumber: "a1b2c3",
		},
	}

	mockPrizeCardOutputs := []domain.PrizeCardOutput{
		{
			ID:           mockPrizeCards[0].Model.ID,
			Title:        mockPrizeCards[0].Title,
			Description:  mockPrizeCards[0].Description,
			SerialNumber: mockUserPrizeCards[0].SerialNumber,
		},
	}

	mockPrizeCardCollection := &domain.PrizeCardCollection{
		PoolName: mockPrizePools[0].Name,
		Cards:    mockPrizeCardOutputs,
		Total:    int64(len(mockPrizeCards)),
	}
	mockPrizeCardCollectionMarshal, _ := json.Marshal(mockPrizeCardCollection)

	s.prizePoolUsecase.On("GetPrizePoolList", mock.Anything, mock.MatchedBy(func(prizePool *domain.PrizePool) bool { return prizePool.Model.ID == mockPrizePools[0].Model.ID })).
		Return(mockPrizePools, nil).Once()
	s.prizeCardUsecase.On("GetPrizeCardCollection", mock.Anything, uint(1), mockPrizePools[0]).
		Return(mockPrizeCardCollection, nil).Once()

	handler := lotteryHandler.LotteryHandler{
		PrizePoolUsecase: s.prizePoolUsecase,
		PrizeCardUsecase: s.prizeCardUsecase,
		KamiUserUsercase: s.kamiUserUsercase,
	}
	s.router.GET("/collection", handler.GetCollection)
	req, _ := http.NewRequest("GET", "/api/lottery/collection?pid=1", nil)
	s.engine.ServeHTTP(s.w, req)

	assert.Equal(s.Suite.T(), http.StatusOK, s.w.Code)
	assert.Equal(s.Suite.T(), string(mockPrizeCardCollectionMarshal), s.w.Body.String())
}

func (s *LotteryHandlerSuite) TestGetCollectionDetail() {
	mockPrizeCard := &domain.PrizeCard{
		Model:       gorm.Model{ID: 1},
		PoolId:      1,
		Title:       "testingCard",
		Description: "testingCardDescription",
		Probability: 100,
	}

	mockUserPrizeCard := &domain.UserPrizeCard{
		ID:           2,
		UserId:       1,
		CardId:       1,
		SerialNumber: "112",
	}

	mockCoupon := &domain.Coupon{
		ID:         1,
		CardId:     1,
		UserId:     1,
		Code:       "123",
		ExpiryDate: time.Now().AddDate(0, 0, 1),
	}

	mockResult := gin.H{
		"serial_number": mockUserPrizeCard.SerialNumber,
		"coupon":        mockCoupon,
		"card":          mockPrizeCard,
	}
	mockResultMarshal, _ := json.Marshal(mockResult)

	s.prizeCardUsecase.On("GetPrizeCardCollectionDetail", mock.Anything, mock.MatchedBy(func(userPrizeCard *domain.UserPrizeCard) bool {
		return userPrizeCard.UserId == mockUserPrizeCard.UserId && userPrizeCard.CardId == mockUserPrizeCard.CardId
	})).
		Return(mockUserPrizeCard.SerialNumber, mockCoupon, mockPrizeCard, nil).Once()

	handler := lotteryHandler.LotteryHandler{
		PrizePoolUsecase: s.prizePoolUsecase,
		PrizeCardUsecase: s.prizeCardUsecase,
		KamiUserUsercase: s.kamiUserUsercase,
	}
	s.router.GET("/collection/detail", handler.GetCollectionDetail)
	req, _ := http.NewRequest("GET", "/api/lottery/collection/detail?cid=1", nil)
	s.engine.ServeHTTP(s.w, req)

	assert.Equal(s.Suite.T(), http.StatusOK, s.w.Code)
	assert.Equal(s.Suite.T(), string(mockResultMarshal), s.w.Body.String())
}

func (s *LotteryHandlerSuite) TestDrawCard() {
	mockKamiUser := &domain.KamiUser{
		Model:    gorm.Model{ID: 1},
		Email:    "test@test.com",
		Phone:    "0900123123",
		Password: "123",
		Points:   10,
	}

	mockPrizePool := &domain.PrizePool{
		Model:  gorm.Model{ID: 1},
		Owner:  "ABC",
		Name:   "testingPool",
		Points: 1,
	}

	mockPrizeCards := []*domain.PrizeCard{
		{
			Model:       gorm.Model{ID: 1},
			PoolId:      1,
			Title:       "testingCard",
			Description: "testingCardDescription",
			Probability: 100,
		},
	}

	mockUserPrizeCard := &domain.UserPrizeCard{
		ID:           2,
		UserId:       1,
		CardId:       1,
		SerialNumber: "112",
	}

	mockPrizeCardOutput := &domain.PrizeCardOutput{
		ID:           mockPrizeCards[0].Model.ID,
		Title:        mockPrizeCards[0].Title,
		Description:  mockPrizeCards[0].Description,
		SerialNumber: mockUserPrizeCard.SerialNumber,
	}
	mockPrizeCardOutputMarshal, _ := json.Marshal(mockPrizeCardOutput)

	s.kamiUserUsercase.On("GetKamiUser", mock.Anything, mock.MatchedBy(func(user *domain.KamiUser) bool { return user.Model.ID == mockKamiUser.Model.ID })).
		Return(mockKamiUser, nil).Once()
	s.kamiUserUsercase.On("UpdateKamiUser", mock.Anything, mock.MatchedBy(func(user *domain.KamiUser) bool { return user.Model.ID == mockKamiUser.Model.ID })).
		Return(nil).Once()
	s.prizePoolUsecase.On("GetPrizePool", mock.Anything, mock.MatchedBy(func(prizePool *domain.PrizePool) bool { return prizePool.Model.ID == mockPrizePool.Model.ID })).
		Return(mockPrizePool, nil).Once()
	s.prizePoolUsecase.On("SubtractUserPoints", mock.Anything, mockKamiUser, mockPrizePool).
		Return(nil).Once()
	s.prizeCardUsecase.On("GetPrizeCardList", mock.Anything, mock.MatchedBy(func(prizeCard *domain.PrizeCard) bool { return prizeCard.PoolId == mockPrizeCards[0].PoolId })).
		Return(mockPrizeCards, nil).Once()
	s.prizeCardUsecase.On("GetWeightedRandomList", mock.Anything, uint(1), mock.MatchedBy(func(prizeCardList []*domain.PrizeCard) bool {
		return prizeCardList[0].Model.ID == mockPrizeCards[0].Model.ID
	})).
		Return(mockPrizeCards, nil).Once()
	s.prizeCardUsecase.On("Draw", mock.Anything, uint(1), mockPrizeCards).
		Return(mockPrizeCardOutput, nil).Once()

	handler := lotteryHandler.LotteryHandler{
		PrizePoolUsecase: s.prizePoolUsecase,
		PrizeCardUsecase: s.prizeCardUsecase,
		KamiUserUsercase: s.kamiUserUsercase,
	}
	s.router.POST("/draw", handler.DrawCard)
	req, _ := http.NewRequest("POST", "/api/lottery/draw?pid=1", nil)
	s.engine.ServeHTTP(s.w, req)

	assert.Equal(s.Suite.T(), http.StatusOK, s.w.Code)
	assert.Equal(s.Suite.T(), string(mockPrizeCardOutputMarshal), s.w.Body.String())
}
