package test

import (
	"net/http"
	"testing"

	"github.com/project-eria/go-wot/dataSchema"
	"github.com/project-eria/go-wot/interaction"
)

func TestActionTD(t *testing.T) {
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
	e := JSONTester(t, httpServer)

	obj := e.GET("/").Expect().
		Status(http.StatusOK).JSON().Object()
	action := obj.Value("actions").Object().Value("a").Object().
		HasValue("title", "Action title").
		HasValue("description", "Action description")
	forms := action.Value("forms").Array()
	forms.Length().IsEqual(1)
	forms.Value(0).Object().
		HasValue("contentType", "application/json").
		HasValue("href", "http://127.0.0.1/a").
		HasValue("htv:methodName", "POST").
		Value("op").Array().ContainsOnly("invokeaction")
}

func TestActionInputOutputString(t *testing.T) {
	mything := getThing()
	stringInput, _ := dataSchema.NewString("", nil, nil, "")
	stringOutput, _ := dataSchema.NewString("", nil, nil, "")
	aAction := interaction.NewAction(
		"a",
		"Action title",
		"Action description",
		&stringInput,
		&stringOutput,
	)
	mything.AddAction(aAction)

	httpServer, _ := getProducer(mything)
	e := JSONTester(t, httpServer)

	obj := e.GET("/").Expect().
		Status(http.StatusOK).JSON().Object()
	action := obj.Value("actions").Object().Value("a").Object()
	action.Value("input").Object().
		HasValue("default", "").
		HasValue("type", "string")
	action.Value("output").Object().
		HasValue("default", "").
		HasValue("type", "string")
}
