package test

import (
	"net/http"
	"testing"

	"github.com/project-eria/go-wot/dataSchema"
	"github.com/project-eria/go-wot/interaction"
	"github.com/project-eria/go-wot/thing"
)

func TestPropertyTD(t *testing.T) {
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

	httpServer, _ := getProducer(mything)
	e := JSONTester(t, httpServer)

	obj := e.GET("/").Expect().
		Status(http.StatusOK).JSON().Object()
	obj.Value("properties").Object().Value("boolRWO").Object().
		HasValue("default", false).
		HasValue("type", "boolean")
	// TODO prop.HasValue("title", "RWO bool")
	// TODO prop.HasValue("description", "Readable/Writable/Observable boolean")
}

func TestPropertyGeneral(t *testing.T) {
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

	httpServer, _ := getProducer(mything)
	e := JSONTester(t, httpServer)

	e.GET("/").Expect().
		Status(http.StatusOK).JSON().Object()

	e.PUT("/").Expect().Status(http.StatusBadRequest)
	e.PUT("/boolRWO").Expect().Status(http.StatusBadRequest)
	// POST is for actions
	e.POST("/").Expect().Status(http.StatusBadRequest)
	e.POST("/boolRWO").Expect().Status(http.StatusBadRequest)

	e = HTMLTester(t, httpServer)
	e.PUT("/").WithHeader("Content-type", "application/json").
		Expect().Status(http.StatusMethodNotAllowed)
	e.POST("/").WithHeader("Content-type", "application/json").
		Expect().Status(http.StatusMethodNotAllowed)
	e.POST("/boolRWO").WithHeader("Content-type", "application/json").
		Expect().Status(http.StatusMethodNotAllowed)
}

func TestPropertyRWO(t *testing.T) {
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

	httpServer, _ := getProducer(mything)
	e := JSONTester(t, httpServer)

	obj := e.GET("/").Expect().
		Status(http.StatusOK).JSON().Object()
	prop := obj.Value("properties").Object().Value("boolRWO").Object().
		HasValue("observable", true).
		HasValue("readOnly", false).
		HasValue("writeOnly", false)
	forms := prop.Value("forms").Array()
	forms.Length().IsEqual(2)
	forms.Value(0).Object().
		HasValue("contentType", "application/json").
		HasValue("href", "http://127.0.0.1/boolRWO").
		Value("op").Array().ContainsOnly("readproperty", "writeproperty")
	//	.HasValue("htv:methodName", "GET")
	forms.Value(1).Object().
		HasValue("contentType", "application/json").
		HasValue("href", "ws://127.0.0.1/boolRWO").
		Value("op").Array().ContainsOnly("observeproperty", "unobserveproperty")
	//	.HasValue("htv:methodName", "GET")

	e.GET("/boolRWO").Expect().
		Status(http.StatusNotImplemented).JSON().Object().
		HasValue("error", "Not Implemented").
		HasValue("type", "NotSupportedError")

	e.PUT("/boolRWO").WithHeader("Content-type", "application/json").
		Expect().Status(http.StatusNotImplemented).
		JSON().Object().
		HasValue("error", "Not Implemented").
		HasValue("type", "NotSupportedError")
}

func TestPropertyReadOnly(t *testing.T) {
	mything, _ := thing.New(
		"dev:ops:my-actuator-1234",
		"0.0.0-dev",
		"Example",
		"An example",
		[]string{},
	)
	booleanData := dataSchema.NewBoolean(false)
	propertyR := interaction.NewProperty(
		"boolR",
		"R bool",
		"Readable only/Not Observable boolean",
		true,
		false,
		false,
		nil,
		booleanData,
	)
	mything.AddProperty(propertyR)

	httpServer, _ := getProducer(mything)
	e := JSONTester(t, httpServer)

	obj := e.GET("/").Expect().
		Status(http.StatusOK).JSON().Object()
	prop := obj.Value("properties").Object().Value("boolR").Object().
		HasValue("observable", false).
		HasValue("readOnly", true).
		HasValue("writeOnly", false)
	forms := prop.Value("forms").Array()
	forms.Length().IsEqual(1)
	forms.Value(0).Object().
		HasValue("contentType", "application/json").
		HasValue("href", "http://127.0.0.1/boolR").
		Value("op").Array().ContainsOnly("readproperty")
	//	form.HasValue("htv:methodName", "GET")

	e.GET("/boolR").Expect().
		Status(http.StatusNotImplemented).JSON().Object().
		HasValue("error", "Not Implemented").
		HasValue("type", "NotSupportedError")

	e.PUT("/boolR").WithHeader("Content-type", "application/json").
		Expect().Status(http.StatusUnauthorized).
		JSON().Object().
		HasValue("error", "Read Only property").
		HasValue("type", "NotAllowedError")
}

func TestPropertyNotObservable(t *testing.T) {
	mything, _ := thing.New(
		"dev:ops:my-actuator-1234",
		"0.0.0-dev",
		"Example",
		"An example",
		[]string{},
	)
	booleanData := dataSchema.NewBoolean(false)
	propertyRW := interaction.NewProperty(
		"boolRW",
		"RW bool",
		"Readable/Writable/Not Observable boolean",
		false,
		false,
		false,
		nil,
		booleanData,
	)
	mything.AddProperty(propertyRW)

	httpServer, _ := getProducer(mything)
	e := JSONTester(t, httpServer)

	obj := e.GET("/").Expect().
		Status(http.StatusOK).JSON().Object()
	prop := obj.Value("properties").Object().Value("boolRW").Object().
		HasValue("observable", false).
		HasValue("readOnly", false).
		HasValue("writeOnly", false)
	forms := prop.Value("forms").Array()
	forms.Length().IsEqual(1)
	forms.Value(0).Object().
		HasValue("contentType", "application/json").
		HasValue("href", "http://127.0.0.1/boolRW").
		Value("op").Array().ContainsOnly("readproperty", "writeproperty")
	//	form.HasValue("htv:methodName", "GET")

	// TODO Move to GET test
	e.GET("/boolRW").Expect().
		Status(http.StatusNotImplemented).JSON().Object().
		HasValue("error", "Not Implemented").
		HasValue("type", "NotSupportedError")

	e.PUT("/boolRW").WithHeader("Content-type", "application/json").
		Expect().Status(http.StatusNotImplemented).
		JSON().Object().
		HasValue("error", "Not Implemented").
		HasValue("type", "NotSupportedError")
}

func TestPropertyWriteOnly(t *testing.T) {
	mything, _ := thing.New(
		"dev:ops:my-actuator-1234",
		"0.0.0-dev",
		"Example",
		"An example",
		[]string{},
	)
	booleanData := dataSchema.NewBoolean(false)
	propertyW := interaction.NewProperty(
		"boolW",
		"W bool",
		"Writable only/Not Observable boolean",
		false,
		true,
		false,
		nil,
		booleanData,
	)
	mything.AddProperty(propertyW)

	httpServer, _ := getProducer(mything)
	e := JSONTester(t, httpServer)

	obj := e.GET("/").Expect().
		Status(http.StatusOK).JSON().Object()
	prop := obj.Value("properties").Object().Value("boolW").Object().
		HasValue("observable", false).
		HasValue("readOnly", false).
		HasValue("writeOnly", true)
	forms := prop.Value("forms").Array()
	forms.Length().IsEqual(1)
	forms.Value(0).Object().
		HasValue("contentType", "application/json").
		HasValue("href", "http://127.0.0.1/boolW").
		Value("op").Array().ContainsOnly("writeproperty")
	//	form.HasValue("htv:methodName", "GET")

	// TODO Move to PUT test

	e.GET("/boolW").Expect().
		Status(http.StatusUnauthorized).JSON().Object().
		HasValue("error", "Write Only property").
		HasValue("type", "NotAllowedError")

	e.PUT("/boolW").WithHeader("Content-type", "application/json").
		Expect().Status(http.StatusNotImplemented).
		JSON().Object().
		HasValue("error", "Not Implemented").
		HasValue("type", "NotSupportedError")
}
