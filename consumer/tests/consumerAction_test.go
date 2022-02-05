package test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Action_NotFound(t *testing.T) {
	_, err := consumedThing.InvokeAction("x", nil)
	assert.Error(t, err, "should return error")
	assert.EqualError(t, err, "action x not found", "they should be equal")
}

func Test_Action_NoInput_NoOutput(t *testing.T) {
	// build response JSON
	result := map[string]interface{}{"ok": true}
	jsonBytes, _ := json.Marshal(result)
	// create a new reader with that JSON
	r := ioutil.NopCloser(bytes.NewReader([]byte(jsonBytes)))
	GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	data, err := consumedThing.InvokeAction("a", nil)
	assert.NoError(t, err, "should not return error")
	assert.Equal(t, data, result)
}

func Test_Action_StringInput_NoOutput(t *testing.T) {
	// build response JSON
	result := map[string]interface{}{"ok": true}
	jsonBytes, _ := json.Marshal(result)
	// create a new reader with that JSON
	r := ioutil.NopCloser(bytes.NewReader([]byte(jsonBytes)))
	GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	data, err := consumedThing.InvokeAction("b", "test")
	assert.NoError(t, err, "should not return error")
	assert.Equal(t, data, result)
}

func Test_Action_StringInput_NoOutput_Missing_data(t *testing.T) {
	// build response JSON
	result := map[string]interface{}{
		"error": "incorrect input value: missing value",
		"type":  "DataError",
	}
	jsonBytes, _ := json.Marshal(result)
	// create a new reader with that JSON
	r := ioutil.NopCloser(bytes.NewReader([]byte(jsonBytes)))
	GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 400,
			Body:       r,
		}, nil
	}

	_, err := consumedThing.InvokeAction("b", nil)
	assert.Error(t, err, "should return error")
	assert.EqualError(t, err, "incorrect input value: missing value", "they should be equal")
}

func Test_Action_StringInput_NoOutput_Incorrect_Type(t *testing.T) {
	// build response JSON
	result := map[string]interface{}{
		"error": "incorrect input value: incorrect string value type",
		"type":  "DataError",
	}
	jsonBytes, _ := json.Marshal(result)
	// create a new reader with that JSON
	r := ioutil.NopCloser(bytes.NewReader([]byte(jsonBytes)))
	GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 400,
			Body:       r,
		}, nil
	}

	_, err := consumedThing.InvokeAction("b", true)
	assert.Error(t, err, "should return error")
	assert.EqualError(t, err, "incorrect input value: incorrect string value type", "they should be equal")
}

func Test_Action_StringInput_StringOutput(t *testing.T) {
	// build response JSON
	result := "test"
	jsonBytes, _ := json.Marshal(result)
	// create a new reader with that JSON
	r := ioutil.NopCloser(bytes.NewReader([]byte(jsonBytes)))
	GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	data, err := consumedThing.InvokeAction("c", "test")
	assert.NoError(t, err, "should not return error")
	assert.Equal(t, data, result)
}
