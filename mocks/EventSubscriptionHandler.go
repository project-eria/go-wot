// Code generated by mockery v2.39.1. DO NOT EDIT.

package mocks

import (
	producer "github.com/project-eria/go-wot/producer"
	mock "github.com/stretchr/testify/mock"
)

// EventSubscriptionHandler is an autogenerated mock type for the EventSubscriptionHandler type
type EventSubscriptionHandler struct {
	mock.Mock
}

// Execute provides a mock function with given fields: _a0, _a1, _a2
func (_m *EventSubscriptionHandler) Execute(_a0 producer.ExposedThing, _a1 string, _a2 map[string]string) (interface{}, error) {
	ret := _m.Called(_a0, _a1, _a2)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 interface{}
	var r1 error
	if rf, ok := ret.Get(0).(func(producer.ExposedThing, string, map[string]string) (interface{}, error)); ok {
		return rf(_a0, _a1, _a2)
	}
	if rf, ok := ret.Get(0).(func(producer.ExposedThing, string, map[string]string) interface{}); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	if rf, ok := ret.Get(1).(func(producer.ExposedThing, string, map[string]string) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewEventSubscriptionHandler creates a new instance of EventSubscriptionHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewEventSubscriptionHandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *EventSubscriptionHandler {
	mock := &EventSubscriptionHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
