package consumer_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/stretchr/testify/mock"
)

func (ts *ConsumerTestSuite) Test_Action_NotFound() {
	_, err := ts.consumedThing.InvokeAction("x", nil)
	ts.Error(err, "should return error")
	ts.EqualError(err, "action x not found", "they should be equal")
}

func (ts *ConsumerTestSuite) Test_Action_NoInput_NoOutput() {
	// build response JSON
	result := map[string]interface{}{"ok": true}
	jsonBytes, _ := json.Marshal(result)
	// create a new reader with that JSON
	r := ioutil.NopCloser(bytes.NewReader([]byte(jsonBytes)))
	ts.httpClientProcessor.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: 200,
		Body:       r,
	}, nil)

	data, err := ts.consumedThing.InvokeAction("a", nil)
	ts.NoError(err, "should not return error")
	ts.Equal(data, result)
}

func (ts *ConsumerTestSuite) Test_Action_StringInput_NoOutput() {
	// build response JSON
	result := map[string]interface{}{"ok": true}
	jsonBytes, _ := json.Marshal(result)
	// create a new reader with that JSON
	r := ioutil.NopCloser(bytes.NewReader([]byte(jsonBytes)))
	ts.httpClientProcessor.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: 200,
		Body:       r,
	}, nil)

	data, err := ts.consumedThing.InvokeAction("b", "test")
	ts.NoError(err, "should not return error")
	ts.Equal(data, result)
}

func (ts *ConsumerTestSuite) Test_Action_StringInput_NoOutput_Missing_data() {
	// build response JSON
	result := map[string]interface{}{
		"error": "incorrect input value: missing value",
		"type":  "DataError",
	}
	jsonBytes, _ := json.Marshal(result)
	// create a new reader with that JSON
	r := ioutil.NopCloser(bytes.NewReader([]byte(jsonBytes)))
	ts.httpClientProcessor.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: 400,
		Body:       r,
	}, nil)

	_, err := ts.consumedThing.InvokeAction("b", nil)
	ts.Error(err, "should return error")
	ts.EqualError(err, "incorrect input value: missing value", "they should be equal")
}

func (ts *ConsumerTestSuite) Test_Action_StringInput_NoOutput_Incorrect_Type() {
	// build response JSON
	result := map[string]interface{}{
		"error": "incorrect input value: incorrect string value type",
		"type":  "DataError",
	}
	jsonBytes, _ := json.Marshal(result)
	// create a new reader with that JSON
	r := ioutil.NopCloser(bytes.NewReader([]byte(jsonBytes)))
	ts.httpClientProcessor.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: 400,
		Body:       r,
	}, nil)

	_, err := ts.consumedThing.InvokeAction("b", true)
	ts.Error(err, "should return error")
	ts.EqualError(err, "incorrect input value: incorrect string value type", "they should be equal")
}

func (ts *ConsumerTestSuite) Test_Action_StringInput_StringOutput() {
	// build response JSON
	result := "test"
	jsonBytes, _ := json.Marshal(result)
	// create a new reader with that JSON
	r := ioutil.NopCloser(bytes.NewReader([]byte(jsonBytes)))
	ts.httpClientProcessor.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: 200,
		Body:       r,
	}, nil)

	data, err := ts.consumedThing.InvokeAction("c", "test")
	ts.NoError(err, "should not return error")
	ts.Equal(data, result)
}
