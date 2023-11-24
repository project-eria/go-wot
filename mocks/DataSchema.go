// Code generated by mockery v2.37.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// DataSchema is an autogenerated mock type for the DataSchema type
type DataSchema struct {
	mock.Mock
}

// Check provides a mock function with given fields: _a0
func (_m *DataSchema) Check(_a0 interface{}) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewDataSchema creates a new instance of DataSchema. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDataSchema(t interface {
	mock.TestingT
	Cleanup(func())
}) *DataSchema {
	mock := &DataSchema{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
