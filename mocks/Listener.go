// Code generated by mockery v2.39.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Listener is an autogenerated mock type for the Listener type
type Listener struct {
	mock.Mock
}

// Execute provides a mock function with given fields: value, err
func (_m *Listener) Execute(value interface{}, err error) {
	_m.Called(value, err)
}

// NewListener creates a new instance of Listener. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewListener(t interface {
	mock.TestingT
	Cleanup(func())
}) *Listener {
	mock := &Listener{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
