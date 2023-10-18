package test

import (
	"net/http"
	"testing"

	"github.com/project-eria/go-wot/interaction"
)

// TODO test all errors /Users/cedric/Documents/PROJETS/eria/v3/src/go-wot/protocolHttp/httpServerPost.go
func TestActionPOST(t *testing.T) {
	mything := getThing()
	aAction := interaction.NewAction(
		"a",
		"Action title",
		"Action description",
		nil,
		nil,
	)
	mything.AddAction(aAction)

	httpServer, _ := getProducer(mything)
	e := HTMLTester(t, httpServer)

	e.GET("/a").Expect().
		Status(http.StatusMethodNotAllowed)

	e = JSONTester(t, httpServer)
	e.POST("/a").WithHeader("Content-type", "application/json").Expect().
		Status(http.StatusNotImplemented).JSON().Object().
		HasValue("error", "Not Implemented").
		HasValue("type", "NotSupportedError")

}
