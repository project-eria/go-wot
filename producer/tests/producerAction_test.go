package test

import (
	"testing"

	"github.com/project-eria/go-wot/dataSchema"
	"github.com/project-eria/go-wot/interaction"
	"github.com/project-eria/go-wot/producer"
)

func Test_NoInputNoOutputAction(t *testing.T) {
	mything := createThing()
	iAction := interaction.NewAction(
		"i",
		"No Input, No Output",
		"",
		nil,
		nil,
	)
	mything.AddAction(iAction)
	myProducer := producer.New("127.0.0.1", 8888, false)
	myProducer.Produce(mything)
	myProducer.Expose()
	client := createClient()

	checkGet(t, client, "", 200, "{\"id\":\"urn:dev:ops:my-actuator-1234\",\"@context\":\"https://www.w3.org/2022/wot/td/v1.1\",\"title\":\"Actuator1 Example\",\"description\":\"An actuator 1st example\",\"actions\":{\"i\":{\"title\":\"No Input, No Output\",\"forms\":[{\"href\":\"http://127.0.0.1:8888/i\",\"contentType\":\"application/json\",\"op\":[\"invokeaction\"]}]}},\"securityDefinitions\":{\"no_sec\":{\"scheme\":\"nosec\"}},\"security\":\"no_sec\"}")

	myProducer.Close()
}

func Test_NoInputNoOutputActionInvokeNoHandler(t *testing.T) {
	mything := createThing()
	iAction := interaction.NewAction(
		"i",
		"No Input, No Output",
		"",
		nil,
		nil,
	)
	mything.AddAction(iAction)
	myProducer := producer.New("127.0.0.1", 8888, false)
	myProducer.Produce(mything)
	myProducer.Expose()
	client := createClient()

	checkPost(t, client, "/i", "", 400, "{\"error\":\"Not Implemented\",\"type\":\"NotSupportedError\"}")

	myProducer.Close()
}

func Test_NoInputNoOutputActionInvoke(t *testing.T) {
	mything := createThing()
	iAction := interaction.NewAction(
		"i",
		"No Input, No Output",
		"",
		nil,
		nil,
	)
	mything.AddAction(iAction)
	myProducer := producer.New("127.0.0.1", 8888, false)
	exposedThing := myProducer.Produce(mything)
	myProducer.Expose()
	exposedThing.SetActionHandler("i", handlerA)
	client := createClient()

	checkPost(t, client, "/i", "", 200, "{\"ok\":true}")

	myProducer.Close()
}

func Test_StringInputNoOutputActionInvokeNoData(t *testing.T) {
	mything := createThing()
	stringInput := dataSchema.NewString("")
	iAction := interaction.NewAction(
		"i",
		"String Input, No Output",
		"",
		&stringInput,
		nil,
	)
	mything.AddAction(iAction)
	myProducer := producer.New("127.0.0.1", 8888, false)
	exposedThing := myProducer.Produce(mything)
	myProducer.Expose()
	exposedThing.SetActionHandler("i", handlerB)
	client := createClient()

	checkPost(t, client, "/i", "", 400, "{\"error\":\"incorrect input value: missing value\",\"type\":\"DataError\"}")

	myProducer.Close()
}

func Test_StringInputNoOutputActionInvoke(t *testing.T) {
	mything := createThing()
	stringInput := dataSchema.NewString("")
	iAction := interaction.NewAction(
		"i",
		"String Input, No Output",
		"",
		&stringInput,
		nil,
	)
	mything.AddAction(iAction)
	myProducer := producer.New("127.0.0.1", 8888, false)
	exposedThing := myProducer.Produce(mything)
	myProducer.Expose()
	exposedThing.SetActionHandler("i", handlerB)
	client := createClient()

	checkPost(t, client, "/i", "\"text\"", 200, "{\"ok\":true}")

	myProducer.Close()
}

func Test_StringInputStringOutputActionInvoke(t *testing.T) {
	mything := createThing()
	stringInput := dataSchema.NewString("")
	stringOutput := dataSchema.NewString("")
	iAction := interaction.NewAction(
		"i",
		"String Input, No Output",
		"",
		&stringInput,
		&stringOutput,
	)
	mything.AddAction(iAction)
	myProducer := producer.New("127.0.0.1", 8888, false)
	exposedThing := myProducer.Produce(mything)
	myProducer.Expose()
	exposedThing.SetActionHandler("i", handlerC)
	client := createClient()

	checkPost(t, client, "/i", "\"text\"", 200, "text")

	myProducer.Close()
}

func handlerA(interface{}) (interface{}, error) {
	return nil, nil
}

func handlerB(value interface{}) (interface{}, error) {
	return nil, nil
}

func handlerC(value interface{}) (interface{}, error) {
	return value, nil
}
