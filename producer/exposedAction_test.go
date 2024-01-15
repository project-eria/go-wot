package producer_test

import (
	"net/http"

	"github.com/project-eria/go-wot/dataSchema"
	"github.com/project-eria/go-wot/interaction"
)

func (ts *ProducerTestSuite) Test_NoInputNoOutputAction() {
	iAction := interaction.NewAction(
		"i",
		"No Input, No Output",
		"",
		nil,
		nil,
	)
	ts.myThing.AddAction(iAction)
	expect, _ := getProducer(ts, ts.myThing)

	//	expect.GET("/").Expect().Body().IsEqual("")

	obj := expect.GET("/").Expect().
		Status(http.StatusOK).JSON().Object()
	// "id": "urn:dev:ops:my-actuator-1234",
	obj.HasValue("id", "urn:dev:ops:my-actuator-1234")
	// "@context": "https://www.w3.org/2022/wot/td/v1.1",
	obj.Value("@context").String().IsEqual("https://www.w3.org/2022/wot/td/v1.1")
	// "title": "Example",
	obj.HasValue("title", "Example")
	// "description": "A 1st example",
	obj.HasValue("description", "A 1st example")
	// "version": {
	// 	"instance": "0.0.0-dev"
	// },
	obj.Value("version").Object().HasValue("instance", "0.0.0-dev")
	// "securityDefinitions": {
	// 	"no_sec": {
	// 		"scheme": "nosec"
	// 	}
	// },
	obj.Value("securityDefinitions").Object().HasValue("no_sec", map[string]interface{}{"scheme": "nosec"})
	// "security": "no_sec"
	obj.Value("security").String().IsEqual("no_sec")
	// "actions": {
	actions := obj.Value("actions").Object()
	action := actions.Value("i").Object()
	action.HasValue("title", "No Input, No Output")
	action.Value("forms").Array().Length().IsEqual(1)
	form1 := action.Value("forms").Array().Value(0).Object()
	form1.HasValue("href", "http://127.0.0.1/i")
	form1.HasValue("contentType", "application/json")
	form1.Value("op").Array().IsEqualUnordered([]string{"invokeaction"})
}

func (ts *ProducerTestSuite) Test_NoInputNoOutputActionInvokeNoHandler() {
	iAction := interaction.NewAction(
		"i",
		"No Input, No Output",
		"",
		nil,
		nil,
	)
	ts.myThing.AddAction(iAction)
	expect, _ := getProducer(ts, ts.myThing)

	expect.POST("/i").WithHeader("Content-type", "application/json").
		Expect().
		Status(http.StatusNotImplemented).JSON().Object().HasValue("error", "No handler function for the action").HasValue("type", "NotSupportedError")
}

func (ts *ProducerTestSuite) Test_NoInputNoOutputActionInvoke() {
	iAction := interaction.NewAction(
		"i",
		"No Input, No Output",
		"",
		nil,
		nil,
	)
	ts.myThing.AddAction(iAction)
	expect, exposedThing := getProducer(ts, ts.myThing)

	exposedThing.SetActionHandler("i", handlerA)
	// expect.POST("/i").WithHeader("Content-type", "application/json").Expect().Body().IsEqual("")
	expect.POST("/i").WithHeader("Content-type", "application/json").
		Expect().
		Status(http.StatusOK).JSON().Object().HasValue("ok", true)
}

func (ts *ProducerTestSuite) Test_StringInputNoOutputActionInvokeNoData() {
	stringInput, _ := dataSchema.NewString("", nil, nil, "")
	iAction := interaction.NewAction(
		"i",
		"String Input, No Output",
		"",
		&stringInput,
		nil,
	)
	ts.myThing.AddAction(iAction)
	expect, exposedThing := getProducer(ts, ts.myThing)

	exposedThing.SetActionHandler("i", handlerB)

	expect.POST("/i").WithHeader("Content-type", "application/json").
		Expect().
		Status(http.StatusBadRequest).JSON().Object().HasValue("error", "input value is required").HasValue("type", "DataError")
}

func (ts *ProducerTestSuite) Test_StringInputNoOutputActionInvoke() {
	stringInput, _ := dataSchema.NewString("", nil, nil, "")
	iAction := interaction.NewAction(
		"i",
		"String Input, No Output",
		"",
		&stringInput,
		nil,
	)
	ts.myThing.AddAction(iAction)
	expect, exposedThing := getProducer(ts, ts.myThing)
	exposedThing.SetActionHandler("i", handlerB)

	expect.POST("/i").WithHeader("Content-type", "application/json").
		WithJSON("\"text\"").
		Expect().
		Status(http.StatusOK).JSON().Object().HasValue("ok", true)
}

func (ts *ProducerTestSuite) Test_StringInputStringOutputActionInvoke() {
	stringInput, _ := dataSchema.NewString("", nil, nil, "")
	stringOutput, _ := dataSchema.NewString("", nil, nil, "")
	iAction := interaction.NewAction(
		"i",
		"String Input, No Output",
		"",
		&stringInput,
		&stringOutput,
	)
	ts.myThing.AddAction(iAction)
	expect, exposedThing := getProducer(ts, ts.myThing)

	exposedThing.SetActionHandler("i", handlerC)

	expect.POST("/i").WithHeader("Content-type", "application/json").
		WithJSON("text").
		Expect().
		Status(http.StatusOK).JSON().String().IsEqual("text")

}

func handlerA(interface{}, map[string]interface{}) (interface{}, error) {
	return nil, nil
}

func handlerB(interface{}, map[string]interface{}) (interface{}, error) {
	return nil, nil
}

func handlerC(value interface{}, _ map[string]interface{}) (interface{}, error) {
	return value, nil
}
