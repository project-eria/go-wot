package consumer_test

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/stretchr/testify/mock"
)

func (ts *ConsumerTestSuite) Test_ReadProperty_boolR() {
	// build response JSON
	json := `false`
	// create a new reader with that JSON
	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))
	ts.httpClientProcessor.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: 200,
		Body:       r,
	}, nil)

	result, err := ts.consumedThing.ReadProperty("boolR")
	ts.NoError(err, "should not return error")
	ts.Equal(result.(bool), false, "they should be equal")
}

func (ts *ConsumerTestSuite) Test_ReadProperty_NotFound() {
	_, err := ts.consumedThing.ReadProperty("x")
	ts.Error(err, "should return error")
	ts.EqualError(err, "property x not found", "they should be equal")
}

func (ts *ConsumerTestSuite) Test_ReadProperty_WriteOnly() {
	_, err := ts.consumedThing.ReadProperty("boolW")
	ts.Error(err, "should return error")
	ts.EqualError(err, "property boolW not readable", "they should be equal")
}

func (ts *ConsumerTestSuite) Test_WriteProperty_NotFound() {
	_, err := ts.consumedThing.WriteProperty("x", true)
	ts.Error(err, "should return error")
	ts.EqualError(err, "property x not found", "they should be equal")
}

func (ts *ConsumerTestSuite) Test_WriteProperty_ReadOnly() {
	_, err := ts.consumedThing.WriteProperty("boolR", true)
	ts.Error(err, "should return error")
	ts.EqualError(err, "property boolR not writable", "they should be equal")
}

func (ts *ConsumerTestSuite) Test_WriteProperty_boolRW() {
	// build response JSON
	json := `{"ok":true}`
	// create a new reader with that JSON
	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))
	ts.httpClientProcessor.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: 200,
		Body:       r,
	}, nil)

	_, err := ts.consumedThing.WriteProperty("boolRW", true)
	ts.NoError(err, "should not return error")
}
