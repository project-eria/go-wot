package test

import (
	"net/http"
	"testing"

	"github.com/project-eria/go-wot/dataSchema"
	"github.com/project-eria/go-wot/interaction"
	"github.com/project-eria/go-wot/producer"
)

// TODO test all errors /Users/cedric/Documents/PROJETS/eria/v3/src/go-wot/protocolHttp/httpServerGet.go

func TestPropertyGET(t *testing.T) {
	mything := getThing()
	booleanData := dataSchema.NewBoolean(false)
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

	exposedThing.SetPropertyReadHandler("boolRWO", func(t producer.ExposedThing, name string, parameters map[string]interface{}) (interface{}, error) {
		return true, nil
	})

	e := JSONTester(t, httpServer)

	e.GET("/boolRWO").Expect().
		Status(http.StatusOK).JSON().Boolean().IsTrue()
}
