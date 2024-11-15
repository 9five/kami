// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	gin "github.com/gin-gonic/gin"
	mock "github.com/stretchr/testify/mock"
)

// MiddlewareUsecase is an autogenerated mock type for the MiddlewareUsecase type
type MiddlewareUsecase struct {
	mock.Mock
}

// VerifyToken provides a mock function with given fields: ctx
func (_m *MiddlewareUsecase) VerifyToken(ctx *gin.Context) {
	_m.Called(ctx)
}

type mockConstructorTestingTNewMiddlewareUsecase interface {
	mock.TestingT
	Cleanup(func())
}

// NewMiddlewareUsecase creates a new instance of MiddlewareUsecase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMiddlewareUsecase(t mockConstructorTestingTNewMiddlewareUsecase) *MiddlewareUsecase {
	mock := &MiddlewareUsecase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
