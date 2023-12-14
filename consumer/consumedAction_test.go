package consumer_test

import (
	"errors"

	"github.com/stretchr/testify/mock"
)

func (ts *ConsumerTestSuite) Test_Action_NotFound() {
	_, err := ts.consumedThing.InvokeAction("x", nil, nil)
	ts.Error(err, "should return error")
	ts.EqualError(err, "action x not found", "they should be equal")
	ts.client.AssertNotCalled(ts.T(), "InvokeResource")
}

func (ts *ConsumerTestSuite) Test_Action_NoInput_NoOutput() {
	ts.client.On("InvokeResource", mock.Anything, map[string]interface{}{}, nil).Return(true, "http://", nil)
	result, err := ts.consumedThing.InvokeAction("a", nil, nil)
	ts.NoError(err, "should not return error")
	ts.Equal(result.(bool), true, "they should be equal")
	ts.client.AssertExpectations(ts.T())
}

func (ts *ConsumerTestSuite) Test_Action_StringInput_NoOutput() {
	ts.client.On("InvokeResource", mock.Anything, map[string]interface{}{}, "test").Return(true, "http://", nil)
	result, err := ts.consumedThing.InvokeAction("b", nil, "test")
	ts.NoError(err, "should not return error")
	ts.Equal(result.(bool), true, "they should be equal")
	ts.client.AssertExpectations(ts.T())
}

func (ts *ConsumerTestSuite) Test_Action_StringInput_NoOutput_Missing_data() {
	ts.client.On("InvokeResource", mock.Anything, map[string]interface{}{}, nil).Return(false, "http://", errors.New("incorrect input value: missing value"))
	_, err := ts.consumedThing.InvokeAction("b", nil, nil)
	ts.Error(err, "should return error")
	ts.EqualError(err, "incorrect input value: missing value", "they should be equal")
	ts.client.AssertExpectations(ts.T())
}

func (ts *ConsumerTestSuite) Test_Action_StringInput_NoOutput_Incorrect_Type() {
	ts.client.On("InvokeResource", mock.Anything, map[string]interface{}{}, true).Return(false, "http://", errors.New("incorrect input value: incorrect string value type"))
	_, err := ts.consumedThing.InvokeAction("b", nil, true)
	ts.Error(err, "should return error")
	ts.EqualError(err, "incorrect input value: incorrect string value type", "they should be equal")
	ts.client.AssertExpectations(ts.T())
}

func (ts *ConsumerTestSuite) Test_Action_StringInput_StringOutput() {
	ts.client.On("InvokeResource", mock.Anything, map[string]interface{}{}, "test").Return(true, "http://", nil)
	result, err := ts.consumedThing.InvokeAction("c", nil, "test")
	ts.NoError(err, "should not return error")
	ts.Equal(result.(bool), true, "they should be equal")
	ts.client.AssertExpectations(ts.T())

}
