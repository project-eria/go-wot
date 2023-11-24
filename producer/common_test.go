package producer_test

import (
	"net/http"
	"os"
	"sync"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/project-eria/go-wot/producer"
	"github.com/project-eria/go-wot/protocolHttp"
	"github.com/project-eria/go-wot/protocolWebSocket"
	"github.com/project-eria/go-wot/securityScheme"
	"github.com/project-eria/go-wot/thing"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

// var base_url = "http://127.0.0.1:8888"

type ProducerTestSuite struct {
	myThing *thing.Thing
	suite.Suite
}

func Test_ProducerTestSuite(t *testing.T) {
	suite.Run(t, &ProducerTestSuite{})
}

func (ts *ProducerTestSuite) SetupSuite() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
}

func (ts *ProducerTestSuite) SetupTest() {
	ts.myThing = createThing()
}

func getProducer(ts *ProducerTestSuite, mything *thing.Thing) (*httpexpect.Expect, producer.ExposedThing) {
	var wait sync.WaitGroup
	myProducer := producer.New(&wait)
	exposedThing := myProducer.Produce("", mything)
	httpServer := protocolHttp.NewServer(":8888", "127.0.0.1", "My App", "My App v0.0.0")
	myProducer.AddServer(httpServer)
	wsServer := protocolWebSocket.NewServer(httpServer)
	myProducer.AddServer(wsServer)
	myProducer.Expose()
	expect := JSONTester(ts, httpServer)
	return expect, exposedThing
}

func createThing() *thing.Thing {
	mything, err := thing.New(
		"dev:ops:my-actuator-1234",
		"0.0.0-dev",
		"Example",
		"A 1st example",
		[]string{},
	)

	if err != nil {
		os.Exit(1)
	}
	// Add Security
	noSecurityScheme := securityScheme.NewNoSecurity()
	mything.AddSecurity("no_sec", noSecurityScheme)

	return mything
}

func fastHTTPTester(ts *ProducerTestSuite, httpServer *protocolHttp.HttpServer) *httpexpect.Expect {
	return httpexpect.WithConfig(httpexpect.Config{
		// Pass requests directly to FastHTTPHandler.
		Client: &http.Client{
			Transport: httpexpect.NewFastBinder(httpServer.Handler()),
			Jar:       httpexpect.NewCookieJar(),
		},
		// Report errors using testify.
		Reporter: httpexpect.NewAssertReporter(ts.T()),
	})
}

func JSONTester(ts *ProducerTestSuite, httpServer *protocolHttp.HttpServer) *httpexpect.Expect {
	e := fastHTTPTester(ts, httpServer)
	// every response should have this header
	return e.Matcher(func(resp *httpexpect.Response) {
		resp.Header("Content-type").IsEqual("application/json")
	})
}
