// Code generated by mockery v2.39.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// SimpleType is an autogenerated mock type for the SimpleType type
type SimpleType struct {
	mock.Mock
}

// NewSimpleType creates a new instance of SimpleType. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSimpleType(t interface {
	mock.TestingT
	Cleanup(func())
}) *SimpleType {
	mock := &SimpleType{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
