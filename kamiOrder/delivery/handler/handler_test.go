package handler_test

import (
	"encoding/json"
	"fmt"
	"kami/domain"
	"kami/domain/mocks"
	kamiOrderHandler "kami/kamiOrder/delivery/handler"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type KamiOrderHandlerSuite struct {
	suite.Suite
	kamiOrderUsecase *mocks.KamiOrderUsecase
	kamiUserUsecase  *mocks.KamiUserUsercase
	engine           *gin.Engine
	router           *gin.RouterGroup
	w                *httptest.ResponseRecorder
}

func TestStart(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	suite.Run(t, &KamiOrderHandlerSuite{})
}

func mockmiddleware(ctx *gin.Context) {
	ctx.Set("id", uint(1))
	ctx.Set("phone", "0900123123")
}

func (s *KamiOrderHandlerSuite) SetupTest() {
	s.kamiOrderUsecase = new(mocks.KamiOrderUsecase)
	s.kamiUserUsecase = new(mocks.KamiUserUsercase)
	s.engine = gin.Default()
	s.router = s.engine.Group("/api/order", mockmiddleware)
	s.w = httptest.NewRecorder()
}

func (s *KamiOrderHandlerSuite) TestOrderRegister() {
	mockKamiUser := &domain.KamiUser{
		Model:    gorm.Model{ID: 1},
		Email:    "test@test.com",
		Phone:    "0900123123",
		Birthday: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
	}

	mockKamiOrder := &domain.KamiOrder{
		Model:            gorm.Model{ID: 1},
		OrderId:          "123-123",
		Restaurant:       "ABC",
		Status:           "Delivered",
		BillingStatus:    "Payable",
		OrderDeliveredAt: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
		OwnerPhone:       "0900123123",
	}

	mockOrderInput := &domain.OrderInput{
		Prefix:           "123",
		Suffix:           "123",
		OrderDeliveredAt: "2000-01-01",
	}
	mockOrderInputMarshal, _ := json.Marshal(mockOrderInput)
	body := strings.NewReader(string(mockOrderInputMarshal))

	s.kamiUserUsecase.On("GetKamiUser", mock.Anything, mock.MatchedBy(func(user *domain.KamiUser) bool { return user.Phone == mockKamiUser.Phone })).
		Return(mockKamiUser, nil).Once()
	s.kamiUserUsecase.On("UpdateKamiUser", mock.Anything, mock.MatchedBy(func(user *domain.KamiUser) bool { return user == mockKamiUser })).
		Return(nil).Once()
	s.kamiOrderUsecase.On("CheckOrderInput", mock.Anything, mock.MatchedBy(func(input *domain.OrderInput) bool {
		return fmt.Sprintf("%s-%s", input.Prefix, input.Suffix) == mockKamiOrder.OrderId
	})).
		Return(mockKamiOrder, nil).Once()
	s.kamiOrderUsecase.On("RegisterOrder", mock.Anything, mockKamiOrder, mockKamiUser).
		Return(nil).Once()

	handler := kamiOrderHandler.KamiOrderHandler{
		KamiOrderUsecase: s.kamiOrderUsecase,
		KamiUserUsercase: s.kamiUserUsecase,
	}
	s.router.PUT("/register", handler.OrderRegister)
	req, _ := http.NewRequest("PUT", "/api/order/register", body)
	s.engine.ServeHTTP(s.w, req)

	assert.Equal(s.Suite.T(), http.StatusOK, s.w.Code)
}

func (s *KamiOrderHandlerSuite) TestGetOrders() {
	mockKamiOrders := []*domain.KamiOrder{
		{
			Model:            gorm.Model{ID: 1},
			OrderId:          "123-123",
			Restaurant:       "ABC",
			Status:           "Delivered",
			BillingStatus:    "Payable",
			OrderDeliveredAt: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
			OwnerPhone:       "0900123123",
		},
	}
	mockKamiOrdersMarshal, _ := json.Marshal(mockKamiOrders)

	s.kamiOrderUsecase.On("GetOrderListByDate", mock.Anything, mock.MatchedBy(func(order *domain.KamiOrder) bool { return order.OwnerPhone == mockKamiOrders[0].OwnerPhone }), "2000-01", "2000-02").
		Return(mockKamiOrders, nil).Once()

	handler := kamiOrderHandler.KamiOrderHandler{
		KamiOrderUsecase: s.kamiOrderUsecase,
		KamiUserUsercase: s.kamiUserUsecase,
	}
	s.router.GET("/getOrders", handler.GetOrders)
	req, _ := http.NewRequest("GET", "/api/order/getOrders?startDate=2000-01&endDate=2000-02", nil)
	s.engine.ServeHTTP(s.w, req)

	assert.Equal(s.Suite.T(), http.StatusOK, s.w.Code)
	assert.Equal(s.Suite.T(), string(mockKamiOrdersMarshal), s.w.Body.String())
}
