package test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ReadProperty_boolR(t *testing.T) {
	// build response JSON
	json := `false`
	// create a new reader with that JSON
	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))
	GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	result, err := consumedThing.ReadProperty("boolR")
	assert.NoError(t, err, "should not return error")
	assert.Equal(t, result.(bool), false, "they should be equal")
}

func Test_ReadProperty_NotFound(t *testing.T) {
	_, err := consumedThing.ReadProperty("x")
	assert.Error(t, err, "should return error")
	assert.EqualError(t, err, "property x not found", "they should be equal")
}

func Test_ReadProperty_WriteOnly(t *testing.T) {
	_, err := consumedThing.ReadProperty("boolW")
	assert.Error(t, err, "should return error")
	assert.EqualError(t, err, "property boolW not readable", "they should be equal")
}

func Test_WriteProperty_NotFound(t *testing.T) {
	err := consumedThing.WriteProperty("x", true)
	assert.Error(t, err, "should return error")
	assert.EqualError(t, err, "property x not found", "they should be equal")
}

func Test_WriteProperty_ReadOnly(t *testing.T) {
	err := consumedThing.WriteProperty("boolR", true)
	assert.Error(t, err, "should return error")
	assert.EqualError(t, err, "property boolR not writable", "they should be equal")
}

func Test_WriteProperty_boolRW(t *testing.T) {
	// build response JSON
	json := `{"ok":true}`
	// create a new reader with that JSON
	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))
	GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	err := consumedThing.WriteProperty("boolRW", true)
	assert.NoError(t, err, "should not return error")
}
