// Code generated by mockery v2.37.0. DO NOT EDIT.

package mocks

import (
	consumer "github.com/project-eria/go-wot/consumer"
	interaction "github.com/project-eria/go-wot/interaction"

	mock "github.com/stretchr/testify/mock"
)

// ProtocolClient is an autogenerated mock type for the ProtocolClient type
type ProtocolClient struct {
	mock.Mock
}

// GetSchemes provides a mock function with given fields:
func (_m *ProtocolClient) GetSchemes() []string {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// InvokeResource provides a mock function with given fields: _a0, _a1
func (_m *ProtocolClient) InvokeResource(_a0 *interaction.Form, _a1 interface{}) (interface{}, error) {
	ret := _m.Called(_a0, _a1)

	var r0 interface{}
	var r1 error
	if rf, ok := ret.Get(0).(func(*interaction.Form, interface{}) (interface{}, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(*interaction.Form, interface{}) interface{}); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	if rf, ok := ret.Get(1).(func(*interaction.Form, interface{}) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReadResource provides a mock function with given fields: _a0
func (_m *ProtocolClient) ReadResource(_a0 *interaction.Form) (interface{}, error) {
	ret := _m.Called(_a0)

	var r0 interface{}
	var r1 error
	if rf, ok := ret.Get(0).(func(*interaction.Form) (interface{}, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(*interaction.Form) interface{}); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	if rf, ok := ret.Get(1).(func(*interaction.Form) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Stop provides a mock function with given fields:
func (_m *ProtocolClient) Stop() {
	_m.Called()
}

// SubscribeResource provides a mock function with given fields: _a0, _a1, _a2
func (_m *ProtocolClient) SubscribeResource(_a0 *interaction.Form, _a1 *consumer.Subscription, _a2 consumer.Listener) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(*interaction.Form, *consumer.Subscription, consumer.Listener) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WriteResource provides a mock function with given fields: _a0, _a1
func (_m *ProtocolClient) WriteResource(_a0 *interaction.Form, _a1 interface{}) (interface{}, error) {
	ret := _m.Called(_a0, _a1)

	var r0 interface{}
	var r1 error
	if rf, ok := ret.Get(0).(func(*interaction.Form, interface{}) (interface{}, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(*interaction.Form, interface{}) interface{}); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	if rf, ok := ret.Get(1).(func(*interaction.Form, interface{}) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewProtocolClient creates a new instance of ProtocolClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewProtocolClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *ProtocolClient {
	mock := &ProtocolClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
