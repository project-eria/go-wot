package test

import (
	"net/http"
	"testing"

	"github.com/project-eria/go-wot/dataSchema"
	"github.com/project-eria/go-wot/interaction"
	"github.com/project-eria/go-wot/producer"
)

// TODO test all errors /Users/cedric/Documents/PROJETS/eria/v3/src/go-wot/protocolHttp/httpServerPut.go
func TestPropertyPUT(t *testing.T) {
	mything := getThing()
	booleanData, _ := dataSchema.NewBoolean()
	propertyRWO := interaction.NewProperty(
		"boolRWO",
		"RWO bool",
		"Readable/Writable/Observable boolean",
		false,
		false,
		true,
		nil,
		booleanData,
	)
	mything.AddProperty(propertyRWO)

	httpServer, exposedThing := getProducer(mything)

	exposedThing.SetPropertyWriteHandler("boolRWO", func(t producer.ExposedThing, name string, value interface{}, parameters map[string]interface{}) error {
		return nil
	})

	e := JSONTester(t, httpServer)

	e.PUT("/boolRWO").WithHeader("Content-type", "application/json").Expect().
		Status(http.StatusBadRequest).JSON().Object().
		HasValue("error", "No data provided").
		HasValue("type", "DataError")

	e.PUT("/boolRWO").WithHeader("Content-type", "application/json").
		WithJSON(true).Expect().
		Status(http.StatusOK).JSON().Object().
		HasValue("ok", true)
}
