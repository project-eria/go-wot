// Code generated by mockery v2.37.0. DO NOT EDIT.

package mocks

import (
	consumer "github.com/project-eria/go-wot/consumer"
	mock "github.com/stretchr/testify/mock"

	thing "github.com/project-eria/go-wot/thing"
)

// ConsumedThing is an autogenerated mock type for the ConsumedThing type
type ConsumedThing struct {
	mock.Mock
}

// GetThingDescription provides a mock function with given fields:
func (_m *ConsumedThing) GetThingDescription() *thing.Thing {
	ret := _m.Called()

	var r0 *thing.Thing
	if rf, ok := ret.Get(0).(func() *thing.Thing); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*thing.Thing)
		}
	}

	return r0
}

// InvokeAction provides a mock function with given fields: _a0, _a1
func (_m *ConsumedThing) InvokeAction(_a0 string, _a1 interface{}) (interface{}, error) {
	ret := _m.Called(_a0, _a1)

	var r0 interface{}
	var r1 error
	if rf, ok := ret.Get(0).(func(string, interface{}) (interface{}, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(string, interface{}) interface{}); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	if rf, ok := ret.Get(1).(func(string, interface{}) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ObserveProperty provides a mock function with given fields: _a0, _a1
func (_m *ConsumedThing) ObserveProperty(_a0 string, _a1 consumer.Listener) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, consumer.Listener) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ReadAllProperties provides a mock function with given fields:
func (_m *ConsumedThing) ReadAllProperties() {
	_m.Called()
}

// ReadMultipleProperties provides a mock function with given fields:
func (_m *ConsumedThing) ReadMultipleProperties() {
	_m.Called()
}

// ReadProperty provides a mock function with given fields: _a0
func (_m *ConsumedThing) ReadProperty(_a0 string) (interface{}, error) {
	ret := _m.Called(_a0)

	var r0 interface{}
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (interface{}, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(string) interface{}); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SubscribeEvent provides a mock function with given fields: _a0, _a1
func (_m *ConsumedThing) SubscribeEvent(_a0 string, _a1 consumer.Listener) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, consumer.Listener) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WriteMultipleProperties provides a mock function with given fields:
func (_m *ConsumedThing) WriteMultipleProperties() {
	_m.Called()
}

// WriteProperty provides a mock function with given fields: _a0, _a1
func (_m *ConsumedThing) WriteProperty(_a0 string, _a1 interface{}) (interface{}, error) {
	ret := _m.Called(_a0, _a1)

	var r0 interface{}
	var r1 error
	if rf, ok := ret.Get(0).(func(string, interface{}) (interface{}, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(string, interface{}) interface{}); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	if rf, ok := ret.Get(1).(func(string, interface{}) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewConsumedThing creates a new instance of ConsumedThing. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewConsumedThing(t interface {
	mock.TestingT
	Cleanup(func())
}) *ConsumedThing {
	mock := &ConsumedThing{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
