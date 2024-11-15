// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	context "context"

	cdp "github.com/chromedp/cdproto/cdp"

	domain "kami/domain"

	mock "github.com/stretchr/testify/mock"
)

// KamiOrderUsecase is an autogenerated mock type for the KamiOrderUsecase type
type KamiOrderUsecase struct {
	mock.Mock
}

// BatchStore provides a mock function with given fields: ctx, platform, orders
func (_m *KamiOrderUsecase) BatchStore(ctx context.Context, platform string, orders []*domain.KamiOrder) error {
	ret := _m.Called(ctx, platform, orders)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, []*domain.KamiOrder) error); ok {
		r0 = rf(ctx, platform, orders)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CheckOrderInput provides a mock function with given fields: ctx, orderInput
func (_m *KamiOrderUsecase) CheckOrderInput(ctx context.Context, orderInput *domain.OrderInput) (*domain.KamiOrder, error) {
	ret := _m.Called(ctx, orderInput)

	var r0 *domain.KamiOrder
	if rf, ok := ret.Get(0).(func(context.Context, *domain.OrderInput) *domain.KamiOrder); ok {
		r0 = rf(ctx, orderInput)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.KamiOrder)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *domain.OrderInput) error); ok {
		r1 = rf(ctx, orderInput)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMoreOrderDetail provides a mock function with given fields: ctx, orders
func (_m *KamiOrderUsecase) GetMoreOrderDetail(ctx context.Context, orders []*domain.KamiOrder) ([]*domain.KamiOrder, error) {
	ret := _m.Called(ctx, orders)

	var r0 []*domain.KamiOrder
	if rf, ok := ret.Get(0).(func(context.Context, []*domain.KamiOrder) []*domain.KamiOrder); ok {
		r0 = rf(ctx, orders)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.KamiOrder)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []*domain.KamiOrder) error); ok {
		r1 = rf(ctx, orders)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOrderList provides a mock function with given fields: ctx, order
func (_m *KamiOrderUsecase) GetOrderList(ctx context.Context, order *domain.KamiOrder) ([]*domain.KamiOrder, error) {
	ret := _m.Called(ctx, order)

	var r0 []*domain.KamiOrder
	if rf, ok := ret.Get(0).(func(context.Context, *domain.KamiOrder) []*domain.KamiOrder); ok {
		r0 = rf(ctx, order)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.KamiOrder)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *domain.KamiOrder) error); ok {
		r1 = rf(ctx, order)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOrderListByDate provides a mock function with given fields: ctx, order, startDate, endDate
func (_m *KamiOrderUsecase) GetOrderListByDate(ctx context.Context, order *domain.KamiOrder, startDate string, endDate string) ([]*domain.KamiOrder, error) {
	ret := _m.Called(ctx, order, startDate, endDate)

	var r0 []*domain.KamiOrder
	if rf, ok := ret.Get(0).(func(context.Context, *domain.KamiOrder, string, string) []*domain.KamiOrder); ok {
		r0 = rf(ctx, order, startDate, endDate)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.KamiOrder)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *domain.KamiOrder, string, string) error); ok {
		r1 = rf(ctx, order, startDate, endDate)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OrderCrawler provides a mock function with given fields: ctx, nodes
func (_m *KamiOrderUsecase) OrderCrawler(ctx context.Context, nodes []*cdp.Node) ([]*domain.KamiOrder, error) {
	ret := _m.Called(ctx, nodes)

	var r0 []*domain.KamiOrder
	if rf, ok := ret.Get(0).(func(context.Context, []*cdp.Node) []*domain.KamiOrder); ok {
		r0 = rf(ctx, nodes)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.KamiOrder)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []*cdp.Node) error); ok {
		r1 = rf(ctx, nodes)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RegisterOrder provides a mock function with given fields: ctx, order, user
func (_m *KamiOrderUsecase) RegisterOrder(ctx context.Context, order *domain.KamiOrder, user *domain.KamiUser) error {
	ret := _m.Called(ctx, order, user)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.KamiOrder, *domain.KamiUser) error); ok {
		r0 = rf(ctx, order, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewKamiOrderUsecase interface {
	mock.TestingT
	Cleanup(func())
}

// NewKamiOrderUsecase creates a new instance of KamiOrderUsecase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewKamiOrderUsecase(t mockConstructorTestingTNewKamiOrderUsecase) *KamiOrderUsecase {
	mock := &KamiOrderUsecase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
