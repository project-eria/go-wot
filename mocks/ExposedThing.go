// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	producer "github.com/project-eria/go-wot/producer"
	mock "github.com/stretchr/testify/mock"

	thing "github.com/project-eria/go-wot/thing"
)

// ExposedThing is an autogenerated mock type for the ExposedThing type
type ExposedThing struct {
	mock.Mock
}

// Destroy provides a mock function with given fields:
func (_m *ExposedThing) Destroy() {
	_m.Called()
}

// EmitEvent provides a mock function with given fields: _a0, _a1
func (_m *ExposedThing) EmitEvent(_a0 string, _a1 map[string]string) error {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for EmitEvent")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, map[string]string) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// EmitPropertyChange provides a mock function with given fields: _a0, _a1, _a2
func (_m *ExposedThing) EmitPropertyChange(_a0 string, _a1 interface{}, _a2 map[string]string) error {
	ret := _m.Called(_a0, _a1, _a2)

	if len(ret) == 0 {
		panic("no return value specified for EmitPropertyChange")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, interface{}, map[string]string) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Expose provides a mock function with given fields:
func (_m *ExposedThing) Expose() {
	_m.Called()
}

// ExposedAction provides a mock function with given fields: _a0
func (_m *ExposedThing) ExposedAction(_a0 string) (producer.ExposedAction, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for ExposedAction")
	}

	var r0 producer.ExposedAction
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (producer.ExposedAction, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(string) producer.ExposedAction); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(producer.ExposedAction)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ExposedEvent provides a mock function with given fields: _a0
func (_m *ExposedThing) ExposedEvent(_a0 string) (producer.ExposedEvent, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for ExposedEvent")
	}

	var r0 producer.ExposedEvent
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (producer.ExposedEvent, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(string) producer.ExposedEvent); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(producer.ExposedEvent)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ExposedProperty provides a mock function with given fields: _a0
func (_m *ExposedThing) ExposedProperty(_a0 string) (producer.ExposedProperty, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for ExposedProperty")
	}

	var r0 producer.ExposedProperty
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (producer.ExposedProperty, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(string) producer.ExposedProperty); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(producer.ExposedProperty)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetEventChannel provides a mock function with given fields:
func (_m *ExposedThing) GetEventChannel() <-chan producer.Event {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetEventChannel")
	}

	var r0 <-chan producer.Event
	if rf, ok := ret.Get(0).(func() <-chan producer.Event); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan producer.Event)
		}
	}

	return r0
}

// GetPropertyChangeChannel provides a mock function with given fields:
func (_m *ExposedThing) GetPropertyChangeChannel() <-chan producer.PropertyChange {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetPropertyChangeChannel")
	}

	var r0 <-chan producer.PropertyChange
	if rf, ok := ret.Get(0).(func() <-chan producer.PropertyChange); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan producer.PropertyChange)
		}
	}

	return r0
}

// GetThingDescription provides a mock function with given fields:
func (_m *ExposedThing) GetThingDescription() *thing.Thing {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetThingDescription")
	}

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

// Ref provides a mock function with given fields:
func (_m *ExposedThing) Ref() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Ref")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// SetActionHandler provides a mock function with given fields: _a0, _a1
func (_m *ExposedThing) SetActionHandler(_a0 string, _a1 producer.ActionHandler) error {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for SetActionHandler")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, producer.ActionHandler) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetEventHandler provides a mock function with given fields: _a0, _a1
func (_m *ExposedThing) SetEventHandler(_a0 string, _a1 producer.EventListenerHandler) error {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for SetEventHandler")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, producer.EventListenerHandler) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetEventSubscribeHandler provides a mock function with given fields: _a0, _a1
func (_m *ExposedThing) SetEventSubscribeHandler(_a0 string, _a1 producer.EventSubscriptionHandler) error {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for SetEventSubscribeHandler")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, producer.EventSubscriptionHandler) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetEventUnsubscribeHandler provides a mock function with given fields: _a0
func (_m *ExposedThing) SetEventUnsubscribeHandler(_a0 string) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for SetEventUnsubscribeHandler")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetObserverSelectorHandler provides a mock function with given fields: _a0, _a1
func (_m *ExposedThing) SetObserverSelectorHandler(_a0 string, _a1 producer.ObserverSelectorHandler) error {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for SetObserverSelectorHandler")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, producer.ObserverSelectorHandler) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetPropertyObserveHandler provides a mock function with given fields: _a0, _a1
func (_m *ExposedThing) SetPropertyObserveHandler(_a0 string, _a1 producer.PropertyObserveHandler) error {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for SetPropertyObserveHandler")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, producer.PropertyObserveHandler) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetPropertyReadHandler provides a mock function with given fields: _a0, _a1
func (_m *ExposedThing) SetPropertyReadHandler(_a0 string, _a1 producer.PropertyReadHandler) error {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for SetPropertyReadHandler")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, producer.PropertyReadHandler) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetPropertyUnobserveHandler provides a mock function with given fields: _a0
func (_m *ExposedThing) SetPropertyUnobserveHandler(_a0 string) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for SetPropertyUnobserveHandler")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetPropertyWriteHandler provides a mock function with given fields: _a0, _a1
func (_m *ExposedThing) SetPropertyWriteHandler(_a0 string, _a1 producer.PropertyWriteHandler) error {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for SetPropertyWriteHandler")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, producer.PropertyWriteHandler) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TD provides a mock function with given fields:
func (_m *ExposedThing) TD() *thing.Thing {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for TD")
	}

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

// NewExposedThing creates a new instance of ExposedThing. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewExposedThing(t interface {
	mock.TestingT
	Cleanup(func())
}) *ExposedThing {
	mock := &ExposedThing{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
