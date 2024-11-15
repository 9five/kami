// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	context "context"
	domain "kami/domain"

	mock "github.com/stretchr/testify/mock"
)

// KamiOrderRepository is an autogenerated mock type for the KamiOrderRepository type
type KamiOrderRepository struct {
	mock.Mock
}

// Get provides a mock function with given fields: ctx, order, optsWhere
func (_m *KamiOrderRepository) Get(ctx context.Context, order *domain.KamiOrder, optsWhere ...map[string]interface{}) (*domain.KamiOrder, error) {
	_va := make([]interface{}, len(optsWhere))
	for _i := range optsWhere {
		_va[_i] = optsWhere[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, order)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *domain.KamiOrder
	if rf, ok := ret.Get(0).(func(context.Context, *domain.KamiOrder, ...map[string]interface{}) *domain.KamiOrder); ok {
		r0 = rf(ctx, order, optsWhere...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.KamiOrder)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *domain.KamiOrder, ...map[string]interface{}) error); ok {
		r1 = rf(ctx, order, optsWhere...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Gets provides a mock function with given fields: ctx, order, optsWhere
func (_m *KamiOrderRepository) Gets(ctx context.Context, order *domain.KamiOrder, optsWhere ...map[string]interface{}) ([]*domain.KamiOrder, error) {
	_va := make([]interface{}, len(optsWhere))
	for _i := range optsWhere {
		_va[_i] = optsWhere[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, order)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []*domain.KamiOrder
	if rf, ok := ret.Get(0).(func(context.Context, *domain.KamiOrder, ...map[string]interface{}) []*domain.KamiOrder); ok {
		r0 = rf(ctx, order, optsWhere...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.KamiOrder)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *domain.KamiOrder, ...map[string]interface{}) error); ok {
		r1 = rf(ctx, order, optsWhere...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Store provides a mock function with given fields: ctx, order
func (_m *KamiOrderRepository) Store(ctx context.Context, order *domain.KamiOrder) error {
	ret := _m.Called(ctx, order)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.KamiOrder) error); ok {
		r0 = rf(ctx, order)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: ctx, order
func (_m *KamiOrderRepository) Update(ctx context.Context, order *domain.KamiOrder) error {
	ret := _m.Called(ctx, order)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.KamiOrder) error); ok {
		r0 = rf(ctx, order)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewKamiOrderRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewKamiOrderRepository creates a new instance of KamiOrderRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewKamiOrderRepository(t mockConstructorTestingTNewKamiOrderRepository) *KamiOrderRepository {
	mock := &KamiOrderRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
