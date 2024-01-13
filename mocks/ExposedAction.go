// Code generated by mockery v2.39.1. DO NOT EDIT.

package mocks

import (
	dataSchema "github.com/project-eria/go-wot/dataSchema"
	mock "github.com/stretchr/testify/mock"

	producer "github.com/project-eria/go-wot/producer"
)

// ExposedAction is an autogenerated mock type for the ExposedAction type
type ExposedAction struct {
	mock.Mock
}

// CheckUriVariables provides a mock function with given fields: _a0
func (_m *ExposedAction) CheckUriVariables(_a0 map[string]string) (map[string]interface{}, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for CheckUriVariables")
	}

	var r0 map[string]interface{}
	var r1 error
	if rf, ok := ret.Get(0).(func(map[string]string) (map[string]interface{}, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(map[string]string) map[string]interface{}); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	if rf, ok := ret.Get(1).(func(map[string]string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Input provides a mock function with given fields:
func (_m *ExposedAction) Input() *dataSchema.Data {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Input")
	}

	var r0 *dataSchema.Data
	if rf, ok := ret.Get(0).(func() *dataSchema.Data); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dataSchema.Data)
		}
	}

	return r0
}

// Output provides a mock function with given fields:
func (_m *ExposedAction) Output() *dataSchema.Data {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Output")
	}

	var r0 *dataSchema.Data
	if rf, ok := ret.Get(0).(func() *dataSchema.Data); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dataSchema.Data)
		}
	}

	return r0
}

// Run provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *ExposedAction) Run(_a0 producer.ExposedThing, _a1 string, _a2 interface{}, _a3 map[string]string) (interface{}, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	if len(ret) == 0 {
		panic("no return value specified for Run")
	}

	var r0 interface{}
	var r1 error
	if rf, ok := ret.Get(0).(func(producer.ExposedThing, string, interface{}, map[string]string) (interface{}, error)); ok {
		return rf(_a0, _a1, _a2, _a3)
	}
	if rf, ok := ret.Get(0).(func(producer.ExposedThing, string, interface{}, map[string]string) interface{}); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	if rf, ok := ret.Get(1).(func(producer.ExposedThing, string, interface{}, map[string]string) error); ok {
		r1 = rf(_a0, _a1, _a2, _a3)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetHandler provides a mock function with given fields: _a0
func (_m *ExposedAction) SetHandler(_a0 producer.ActionHandler) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for SetHandler")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(producer.ActionHandler) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewExposedAction creates a new instance of ExposedAction. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewExposedAction(t interface {
	mock.TestingT
	Cleanup(func())
}) *ExposedAction {
	mock := &ExposedAction{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
