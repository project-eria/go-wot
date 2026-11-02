package protocolHttp_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/project-eria/go-wot/interaction"
	"github.com/project-eria/go-wot/mocks"
	"github.com/project-eria/go-wot/protocolHttp"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type HttpClientTestSuite struct {
	client              *protocolHttp.HttpClient
	httpClientProcessor *mocks.HttpClientProcessor
	suite.Suite
}

func Test_HttpClientTestSuite(t *testing.T) {
	suite.Run(t, &HttpClientTestSuite{})
}

func (ts *HttpClientTestSuite) SetupTest() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	httpClientProcessor := &mocks.HttpClientProcessor{}
	ts.client = &protocolHttp.HttpClient{
		Client:  httpClientProcessor,
		Schemes: []string{"http"},
	}
	ts.httpClientProcessor = httpClientProcessor
}

func (ts *HttpClientTestSuite) Test_ReadResource_WithoutURIVars() {
	form := &interaction.Form{
		Href: "http://127.0.0.1:8888/test",
	}
	// build response JSON
	json := `true`
	// create a new reader with that JSON
	r := io.NopCloser(bytes.NewReader([]byte(json)))
	ts.httpClientProcessor.On("Do", mock.MatchedBy(func(req *http.Request) bool { return req.URL.String() == "http://127.0.0.1:8888/test" })).Return(&http.Response{
		StatusCode: 200,
		Body:       r,
	}, nil)
	result, _, err := ts.client.ReadResource(form, nil)
	ts.NoError(err, "should not return error")
	ts.Equal(result.(bool), true, "they should be equal")
	ts.httpClientProcessor.AssertExpectations(ts.T())
}

func (ts *HttpClientTestSuite) Test_ReadResource_WithURIVars() {
	form := &interaction.Form{
		Href: "http://127.0.0.1:8888/test/{var}",
	}
	// build response JSON
	json := `true`
	// create a new reader with that JSON
	r := io.NopCloser(bytes.NewReader([]byte(json)))
	ts.httpClientProcessor.On("Do", mock.MatchedBy(func(req *http.Request) bool { return req.URL.String() == "http://127.0.0.1:8888/test/1" })).Return(&http.Response{
		StatusCode: 200,
		Body:       r,
	}, nil)
	result, _, err := ts.client.ReadResource(form, map[string]interface{}{"var": "1"})
	ts.NoError(err, "should not return error")
	ts.Equal(result.(bool), true, "they should be equal")
	ts.httpClientProcessor.AssertExpectations(ts.T())
}

func (ts *HttpClientTestSuite) Test_ReadResource_Token() {
	ts.client.Token = func() string { return "tok" }
	ts.httpClientProcessor.On("Do", mock.MatchedBy(func(r *http.Request) bool {
		return r.Header.Get("Authorization") == "Bearer tok"
	})).Return(&http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewBufferString(`1`))}, nil)
	_, _, err := ts.client.ReadResource(&interaction.Form{Href: "http://127.0.0.1:8888/test"}, nil)
	ts.NoError(err, "the service token is sent")
}
