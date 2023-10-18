package test

import (
	"net/http"
	"sync"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/project-eria/go-wot/producer"
	"github.com/project-eria/go-wot/protocolHttp"
	"github.com/project-eria/go-wot/protocolWebSocket"
	"github.com/project-eria/go-wot/thing"
	"github.com/rs/zerolog"
)

func getProducer(mything *thing.Thing) (*protocolHttp.HttpServer, *producer.ExposedThing) {
	zerolog.SetGlobalLevel(zerolog.Disabled)

	var wait sync.WaitGroup
	myProducer := producer.New(&wait)
	exposedThing := myProducer.Produce("", mything)
	httpServer := protocolHttp.NewServer(":8888", "127.0.0.1", "", "")
	myProducer.AddServer(httpServer)
	wsServer := protocolWebSocket.NewServer(httpServer)
	myProducer.AddServer(wsServer)
	myProducer.Expose()
	return httpServer, exposedThing
}

func fastHTTPTester(t *testing.T, httpServer *protocolHttp.HttpServer) *httpexpect.Expect {
	return httpexpect.WithConfig(httpexpect.Config{
		// Pass requests directly to FastHTTPHandler.
		Client: &http.Client{
			Transport: httpexpect.NewFastBinder(httpServer.Handler()),
			Jar:       httpexpect.NewCookieJar(),
		},
		// Report errors using testify.
		Reporter: httpexpect.NewAssertReporter(t),
	})
}

func HTMLTester(t *testing.T, httpServer *protocolHttp.HttpServer) *httpexpect.Expect {
	return fastHTTPTester(t, httpServer)
}

func JSONTester(t *testing.T, httpServer *protocolHttp.HttpServer) *httpexpect.Expect {
	e := fastHTTPTester(t, httpServer)
	// every response should have this header
	return e.Matcher(func(resp *httpexpect.Response) {
		resp.Header("Content-type").IsEqual("application/json")
	})
}

func getThing() *thing.Thing {
	mything, _ := thing.New(
		"dev:ops:my-actuator-1234",
		"0.0.0-dev",
		"Example",
		"An example",
		[]string{},
	)
	return mything
}
