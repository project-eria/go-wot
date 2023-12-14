package consumer_test

import (
	"github.com/stretchr/testify/mock"
)

func (ts *ConsumerTestSuite) Test_ReadProperty_boolR() {
	ts.client.On("ReadResource", mock.Anything, mock.Anything).Return(true, "http://", nil)
	result, err := ts.consumedThing.ReadProperty("boolR", nil)
	ts.NoError(err, "should not return error")
	ts.Equal(result.(bool), true, "they should be equal")
	ts.client.AssertExpectations(ts.T())
}

func (ts *ConsumerTestSuite) Test_ReadProperty_NotFound() {
	_, err := ts.consumedThing.ReadProperty("x", nil)
	ts.Error(err, "should return error")
	ts.EqualError(err, "property x not found", "they should be equal")
	ts.client.AssertNotCalled(ts.T(), "ReadResource")
}

func (ts *ConsumerTestSuite) Test_ReadProperty_WriteOnly() {
	_, err := ts.consumedThing.ReadProperty("boolW", nil)
	ts.Error(err, "should return error")
	ts.EqualError(err, "property boolW not readable", "they should be equal")
	ts.client.AssertNotCalled(ts.T(), "ReadResource")
}

func (ts *ConsumerTestSuite) Test_WriteProperty_NotFound() {
	_, err := ts.consumedThing.WriteProperty("x", nil, true)
	ts.Error(err, "should return error")
	ts.EqualError(err, "property x not found", "they should be equal")
	ts.client.AssertNotCalled(ts.T(), "ReadResource")
}

func (ts *ConsumerTestSuite) Test_WriteProperty_ReadOnly() {
	_, err := ts.consumedThing.WriteProperty("boolR", nil, true)
	ts.Error(err, "should return error")
	ts.EqualError(err, "property boolR not writable", "they should be equal")
	ts.client.AssertNotCalled(ts.T(), "ReadResource")
}

func (ts *ConsumerTestSuite) Test_WriteProperty_boolRW() {
	ts.client.On("WriteResource", mock.Anything, map[string]interface{}{}, true).Return(true, "http://", nil)
	result, err := ts.consumedThing.WriteProperty("boolRW", nil, true)
	ts.NoError(err, "should not return error")
	ts.Equal(result.(bool), true, "they should be equal")
	ts.client.AssertExpectations(ts.T())
}

func (ts *ConsumerTestSuite) Test_ReadProperty_WithURIVars() {
	ts.client.On("ReadResource", mock.Anything, map[string]interface{}{"var1": "1", "var2": "2"}).Return(true, "http://", nil)

	result, err := ts.consumedThing.ReadProperty("uriVars", map[string]interface{}{"var1": "1", "var2": "2"})
	ts.NoError(err, "should not return error")
	ts.Equal(result.(bool), true, "they should be equal")
	ts.client.AssertExpectations(ts.T())
}

func (ts *ConsumerTestSuite) Test_ReadProperty_WithURIVarsNil() {
	_, err := ts.consumedThing.ReadProperty("uriVars", nil)
	ts.Error(err, "should return error")
	ts.EqualError(err, "uri variables not set", "they should be equal")
	ts.client.AssertNotCalled(ts.T(), "ReadResource")
}

func (ts *ConsumerTestSuite) Test_ReadProperty_WithURIVarsEmpty() {
	_, err := ts.consumedThing.ReadProperty("uriVars", map[string]interface{}{})
	ts.Error(err, "should return error")
	ts.EqualError(err, "uri variables not set", "they should be equal")
	ts.client.AssertNotCalled(ts.T(), "ReadResource")
}

func (ts *ConsumerTestSuite) Test_ReadProperty_WithURIVarsMissing() {
	_, err := ts.consumedThing.ReadProperty("uriVars", map[string]interface{}{"var2": "2"})
	ts.Error(err, "should return error")
	ts.EqualError(err, "uri variable var1 not set", "they should be equal")
	ts.client.AssertNotCalled(ts.T(), "ReadResource")

}

func (ts *ConsumerTestSuite) Test_ReadProperty_WithURIVarsMissingWithDefault() {
	ts.client.On("ReadResource", mock.Anything, map[string]interface{}{"var1": "1", "var2": "test"}).Return(true, "http://", nil)
	result, err := ts.consumedThing.ReadProperty("uriVars", map[string]interface{}{"var1": "1"})
	ts.NoError(err, "should not return error")
	ts.Equal(result.(bool), true, "they should be equal")
	ts.client.AssertExpectations(ts.T())
}

func (ts *ConsumerTestSuite) Test_ReadProperty_WithURIVarsExtra() {
	ts.client.On("ReadResource", mock.Anything, map[string]interface{}{"var1": "1", "var2": "test"}).Return(true, "http://", nil)

	result, err := ts.consumedThing.ReadProperty("uriVars", map[string]interface{}{"var1": "1", "var3": "3"})
	ts.NoError(err, "should not return error")
	ts.Equal(result.(bool), true, "they should be equal")
	ts.client.AssertExpectations(ts.T())
}
