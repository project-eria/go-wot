// Code generated by mockery v2.37.0. DO NOT EDIT.

package mocks

import (
	producer "github.com/project-eria/go-wot/producer"
	mock "github.com/stretchr/testify/mock"
)

// ExposedEvent is an autogenerated mock type for the ExposedEvent type
type ExposedEvent struct {
	mock.Mock
}

// CheckUriVariables provides a mock function with given fields: _a0
func (_m *ExposedEvent) CheckUriVariables(_a0 map[string]string) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(map[string]string) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetEventHandler provides a mock function with given fields:
func (_m *ExposedEvent) GetEventHandler() producer.EventListenerHandler {
	ret := _m.Called()

	var r0 producer.EventListenerHandler
	if rf, ok := ret.Get(0).(func() producer.EventListenerHandler); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(producer.EventListenerHandler)
		}
	}

	return r0
}

// GetListenerSelectorHandler provides a mock function with given fields:
func (_m *ExposedEvent) GetListenerSelectorHandler() producer.ListenerSelectorHandler {
	ret := _m.Called()

	var r0 producer.ListenerSelectorHandler
	if rf, ok := ret.Get(0).(func() producer.ListenerSelectorHandler); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(producer.ListenerSelectorHandler)
		}
	}

	return r0
}

// SetEventHandler provides a mock function with given fields: _a0
func (_m *ExposedEvent) SetEventHandler(_a0 producer.EventListenerHandler) {
	_m.Called(_a0)
}

// SetListenerSelectorHandler provides a mock function with given fields: _a0
func (_m *ExposedEvent) SetListenerSelectorHandler(_a0 producer.ListenerSelectorHandler) {
	_m.Called(_a0)
}

// SetSubscribeHandler provides a mock function with given fields: _a0
func (_m *ExposedEvent) SetSubscribeHandler(_a0 producer.EventSubscriptionHandler) {
	_m.Called(_a0)
}

// SetUnSubscribeHandler provides a mock function with given fields:
func (_m *ExposedEvent) SetUnSubscribeHandler() {
	_m.Called()
}

// NewExposedEvent creates a new instance of ExposedEvent. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewExposedEvent(t interface {
	mock.TestingT
	Cleanup(func())
}) *ExposedEvent {
	mock := &ExposedEvent{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
