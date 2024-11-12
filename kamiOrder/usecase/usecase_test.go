package usecase_test

import (
	"context"
	"kami/domain"
	"kami/domain/mocks"
	"testing"
	"time"

	ucase "kami/kamiOrder/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type KamiOrderUcaseSuite struct {
	suite.Suite
	repo *mocks.KamiOrderRepository
}

func TestStart(t *testing.T) {
	suite.Run(t, &KamiOrderUcaseSuite{})
}

func (s *KamiOrderUcaseSuite) SetupTest() {
	s.repo = new(mocks.KamiOrderRepository)
}

func (s *KamiOrderUcaseSuite) TestBatchStore_Success() {
	mockKamiOrder := &domain.KamiOrder{
		Model:            gorm.Model{ID: 1},
		OrderId:          "123-123",
		Restaurant:       "ABC",
		Status:           "Delivered",
		BillingStatus:    "Payable",
		OrderDeliveredAt: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
		OwnerPhone:       "0900123123",
	}

	s.repo.On("Store", mock.Anything, mock.MatchedBy(func(order *domain.KamiOrder) bool { return order == mockKamiOrder })).
		Return(nil).Once()

	u := ucase.NewKamiOrderUsecase(s.repo)
	err := u.BatchStore(context.TODO(), "testingPlatform", []*domain.KamiOrder{mockKamiOrder})
	assert.NoError(s.Suite.T(), err)
}

func (s *KamiOrderUcaseSuite) TestCheckOrderInput_Success() {
	mockKamiOrder := &domain.KamiOrder{
		Model:            gorm.Model{ID: 1},
		OrderId:          "123-123",
		Restaurant:       "ABC",
		Status:           "Delivered",
		BillingStatus:    "Payable",
		OrderDeliveredAt: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
		OwnerPhone:       "0900123123",
	}

	s.repo.On("Get", mock.Anything, mock.MatchedBy(func(order *domain.KamiOrder) bool { return order.OrderId == mockKamiOrder.OrderId })).
		Return(mockKamiOrder, nil).Once()

	u := ucase.NewKamiOrderUsecase(s.repo)
	res, err := u.CheckOrderInput(context.TODO(), &domain.OrderInput{Prefix: "123", Suffix: "123", OrderDeliveredAt: "2000-01-01 01:01"})
	assert.NoError(s.Suite.T(), err)
	assert.Equal(s.Suite.T(), mockKamiOrder, res)
}

func (s *KamiOrderUcaseSuite) TestRegisterOrder_Success() {
	mockKamiOrder := &domain.KamiOrder{
		Model:            gorm.Model{ID: 1},
		OrderId:          "123-123",
		Restaurant:       "ABC",
		Status:           "Delivered",
		BillingStatus:    "Payable",
		OrderDeliveredAt: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
		OwnerPhone:       "0900123123",
	}

	s.repo.On("Update", mock.Anything, mock.MatchedBy(func(order *domain.KamiOrder) bool { return order.Model.ID == mockKamiOrder.Model.ID })).
		Return(nil).Once()

	u := ucase.NewKamiOrderUsecase(s.repo)
	err := u.RegisterOrder(context.TODO(), mockKamiOrder, &domain.KamiUser{Phone: "0900123123"})
	assert.NoError(s.Suite.T(), err)
}

func (s *KamiOrderUcaseSuite) TestGetOrderList_Success() {
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

	s.repo.On("Gets", mock.Anything, mock.MatchedBy(func(order *domain.KamiOrder) bool { return order.Model.ID == mockKamiOrders[0].Model.ID })).
		Return(mockKamiOrders, nil).Once()

	u := ucase.NewKamiOrderUsecase(s.repo)
	res, err := u.GetOrderList(context.TODO(), mockKamiOrders[0])
	assert.NoError(s.Suite.T(), err)
	assert.Equal(s.Suite.T(), mockKamiOrders, res)
}

func (s *KamiOrderUcaseSuite) TestGetOrderListByDate_Success() {
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

	s.repo.On("Gets", mock.Anything, mock.MatchedBy(func(order *domain.KamiOrder) bool { return order.Model.ID == mockKamiOrders[0].Model.ID }), map[string]interface{}{
		"order_delivered_at >= ?": time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
		"order_delivered_at <= ?": time.Date(2000, time.March, 1, 0, 0, 0, 0, time.UTC),
	}).
		Return(mockKamiOrders, nil).Once()

	u := ucase.NewKamiOrderUsecase(s.repo)
	res, err := u.GetOrderListByDate(context.TODO(), mockKamiOrders[0], "2000-01", "2000-02")
	assert.NoError(s.Suite.T(), err)
	assert.Equal(s.Suite.T(), mockKamiOrders, res)
}
