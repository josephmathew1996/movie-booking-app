// Code generated by mockery v2.32.4. DO NOT EDIT.

package mocks

import (
	utils "movie-booking-app/users-service/pkg/utils"

	mock "github.com/stretchr/testify/mock"
)

// IAuthMiddleware is an autogenerated mock type for the IAuthMiddleware type
type IAuthMiddleware struct {
	mock.Mock
}

// Authenticate provides a mock function with given fields: _a0
func (_m *IAuthMiddleware) Authenticate(_a0 string) (*utils.Claims, error) {
	ret := _m.Called(_a0)

	var r0 *utils.Claims
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*utils.Claims, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(string) *utils.Claims); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*utils.Claims)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewIAuthMiddleware creates a new instance of IAuthMiddleware. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIAuthMiddleware(t interface {
	mock.TestingT
	Cleanup(func())
}) *IAuthMiddleware {
	mock := &IAuthMiddleware{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
