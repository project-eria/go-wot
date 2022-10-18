package test

import (
	"testing"

	"github.com/project-eria/go-wot/dataSchema"
	"github.com/project-eria/go-wot/interaction"
	"github.com/project-eria/go-wot/producer"
)

func Test_RWOBoolProperty(t *testing.T) {
	mything := createThing()
	booleanData := dataSchema.NewBoolean(false)
	propertyRWO := interaction.NewProperty(
		"boolRWO",
		"RWO bool",
		"Readable/Writable/Observable boolean",
		false,
		false,
		true,
		booleanData,
	)
	mything.AddProperty(&propertyRWO)
	myProducer := producer.New("127.0.0.1", 8888, false)
	myProducer.Produce(mything)
	myProducer.Expose()
	client := createClient()

	checkGet(t, client, "", 200, "{\"id\":\"urn:dev:ops:my-actuator-1234\",\"@context\":\"https://www.w3.org/2022/wot/td/v1.1\",\"title\":\"Actuator1 Example\",\"description\":\"An actuator 1st example\",\"properties\":{\"boolRWO\":{\"observable\":true,\"title\":\"RWO bool\",\"description\":\"Readable/Writable/Observable boolean\",\"forms\":[{\"href\":\"http://127.0.0.1:8888/boolRWO\",\"contentType\":\"application/json\",\"op\":[\"writeproperty\",\"readproperty\"]},{\"href\":\"ws://127.0.0.1:8888/boolRWO\",\"contentType\":\"application/json\",\"op\":[\"observeproperty\",\"unobserveproperty\"]}],\"default\":false,\"readOnly\":false,\"writeOnly\":false,\"type\":\"boolean\"}},\"securityDefinitions\":{\"no_sec\":{\"scheme\":\"nosec\"}},\"security\":\"no_sec\"}")

	myProducer.Close()
}

func Test_RWOBoolPropertyRead(t *testing.T) {
	mything := createThing()
	booleanData := dataSchema.NewBoolean(false)
	propertyRWO := interaction.NewProperty(
		"boolRWO",
		"RWO bool",
		"Readable/Writable/Observable boolean",
		false,
		false,
		true,
		booleanData,
	)
	mything.AddProperty(&propertyRWO)
	myProducer := producer.New("127.0.0.1", 8888, false)
	myProducer.Produce(mything)
	myProducer.Expose()
	client := createClient()

	checkGet(t, client, "/boolRWO", 200, "false")

	myProducer.Close()
}

func Test_RWOBoolPropertyWrite(t *testing.T) {
	mything := createThing()
	booleanData := dataSchema.NewBoolean(false)
	propertyRWO := interaction.NewProperty(
		"boolRWO",
		"RWO bool",
		"Readable/Writable/Observable boolean",
		false,
		false,
		true,
		booleanData,
	)
	mything.AddProperty(&propertyRWO)
	myProducer := producer.New("127.0.0.1", 8888, false)
	myProducer.Produce(mything)
	myProducer.Expose()
	client := createClient()

	checkPut(t, client, "/boolRWO", "", 400, "{\"error\":\"No data provided\",\"type\":\"DataError\"}")

	checkPut(t, client, "/boolRWO", "true", 200, "{\"ok\":true}")

	checkGet(t, client, "/boolRWO", 200, "true")
	myProducer.Close()
}

func Test_RWBoolProperty(t *testing.T) {
	mything := createThing()
	booleanData := dataSchema.NewBoolean(false)
	propertyRW := interaction.NewProperty(
		"boolRW",
		"RW bool",
		"Readable/Writable/Not Observable boolean",
		false,
		false,
		false,
		booleanData,
	)
	mything.AddProperty(&propertyRW)
	myProducer := producer.New("127.0.0.1", 8888, false)
	myProducer.Produce(mything)
	myProducer.Expose()
	client := createClient()

	checkGet(t, client, "", 200, "{\"id\":\"urn:dev:ops:my-actuator-1234\",\"@context\":\"https://www.w3.org/2022/wot/td/v1.1\",\"title\":\"Actuator1 Example\",\"description\":\"An actuator 1st example\",\"properties\":{\"boolRW\":{\"observable\":false,\"title\":\"RW bool\",\"description\":\"Readable/Writable/Not Observable boolean\",\"forms\":[{\"href\":\"http://127.0.0.1:8888/boolRW\",\"contentType\":\"application/json\",\"op\":[\"writeproperty\",\"readproperty\"]}],\"default\":false,\"readOnly\":false,\"writeOnly\":false,\"type\":\"boolean\"}},\"securityDefinitions\":{\"no_sec\":{\"scheme\":\"nosec\"}},\"security\":\"no_sec\"}")

	myProducer.Close()
}

func Test_RWBoolPropertyRead(t *testing.T) {
	mything := createThing()
	booleanData := dataSchema.NewBoolean(false)
	propertyRW := interaction.NewProperty(
		"boolRW",
		"RW bool",
		"Readable/Writable/Not Observable boolean",
		false,
		false,
		false,
		booleanData,
	)
	mything.AddProperty(&propertyRW)
	myProducer := producer.New("127.0.0.1", 8888, false)
	myProducer.Produce(mything)
	myProducer.Expose()
	client := createClient()

	checkGet(t, client, "/boolRW", 200, "false")

	myProducer.Close()
}

func Test_RWBoolPropertyWrite(t *testing.T) {
	mything := createThing()
	booleanData := dataSchema.NewBoolean(false)
	propertyRW := interaction.NewProperty(
		"boolRW",
		"RW bool",
		"Readable/Writable/Not Observable boolean",
		false,
		false,
		false,
		booleanData,
	)
	mything.AddProperty(&propertyRW)
	myProducer := producer.New("127.0.0.1", 8888, false)
	myProducer.Produce(mything)
	myProducer.Expose()
	client := createClient()

	checkPut(t, client, "/boolRW", "", 400, "{\"error\":\"No data provided\",\"type\":\"DataError\"}")

	checkPut(t, client, "/boolRW", "true", 200, "{\"ok\":true}")

	checkGet(t, client, "/boolRW", 200, "true")
	myProducer.Close()
}

func Test_RBoolProperty(t *testing.T) {
	mything := createThing()
	booleanData := dataSchema.NewBoolean(false)
	propertyR := interaction.NewProperty(
		"boolR",
		"R bool",
		"Readable only/Not Observable boolean",
		true,
		false,
		false,
		booleanData,
	)
	mything.AddProperty(&propertyR)
	myProducer := producer.New("127.0.0.1", 8888, false)
	myProducer.Produce(mything)
	myProducer.Expose()
	client := createClient()

	checkGet(t, client, "", 200, "{\"id\":\"urn:dev:ops:my-actuator-1234\",\"@context\":\"https://www.w3.org/2022/wot/td/v1.1\",\"title\":\"Actuator1 Example\",\"description\":\"An actuator 1st example\",\"properties\":{\"boolR\":{\"observable\":false,\"title\":\"R bool\",\"description\":\"Readable only/Not Observable boolean\",\"forms\":[{\"href\":\"http://127.0.0.1:8888/boolR\",\"contentType\":\"application/json\",\"op\":[\"readproperty\"]}],\"default\":false,\"readOnly\":true,\"writeOnly\":false,\"type\":\"boolean\"}},\"securityDefinitions\":{\"no_sec\":{\"scheme\":\"nosec\"}},\"security\":\"no_sec\"}")

	myProducer.Close()
}

func Test_RBoolPropertyRead(t *testing.T) {
	mything := createThing()
	booleanData := dataSchema.NewBoolean(false)
	propertyR := interaction.NewProperty(
		"boolR",
		"R bool",
		"Readable only/Not Observable boolean",
		true,
		false,
		false,
		booleanData,
	)
	mything.AddProperty(&propertyR)
	myProducer := producer.New("127.0.0.1", 8888, false)
	myProducer.Produce(mything)
	myProducer.Expose()
	client := createClient()

	checkGet(t, client, "/boolR", 200, "false")

	myProducer.Close()
}

func Test_RBoolPropertyWrite(t *testing.T) {
	mything := createThing()
	booleanData := dataSchema.NewBoolean(false)
	propertyR := interaction.NewProperty(
		"boolR",
		"R bool",
		"Readable only/Not Observable boolean",
		true,
		false,
		false,
		booleanData,
	)
	mything.AddProperty(&propertyR)
	myProducer := producer.New("127.0.0.1", 8888, false)
	myProducer.Produce(mything)
	myProducer.Expose()
	client := createClient()

	checkPut(t, client, "/boolR", "true", 401, "{\"error\":\"Read Only property\",\"type\":\"NotAllowedError\"}")

	myProducer.Close()
}

func Test_WBoolProperty(t *testing.T) {
	mything := createThing()
	booleanData := dataSchema.NewBoolean(false)
	propertyW := interaction.NewProperty(
		"boolW",
		"W bool",
		"Writable only/Not Observable boolean",
		false,
		true,
		false,
		booleanData,
	)
	mything.AddProperty(&propertyW)
	myProducer := producer.New("127.0.0.1", 8888, false)
	myProducer.Produce(mything)
	myProducer.Expose()
	client := createClient()

	checkGet(t, client, "", 200, "{\"id\":\"urn:dev:ops:my-actuator-1234\",\"@context\":\"https://www.w3.org/2022/wot/td/v1.1\",\"title\":\"Actuator1 Example\",\"description\":\"An actuator 1st example\",\"properties\":{\"boolW\":{\"observable\":false,\"title\":\"W bool\",\"description\":\"Writable only/Not Observable boolean\",\"forms\":[{\"href\":\"http://127.0.0.1:8888/boolW\",\"contentType\":\"application/json\",\"op\":[\"writeproperty\"]}],\"default\":false,\"readOnly\":false,\"writeOnly\":true,\"type\":\"boolean\"}},\"securityDefinitions\":{\"no_sec\":{\"scheme\":\"nosec\"}},\"security\":\"no_sec\"}")

	myProducer.Close()
}

func Test_WBoolPropertyRead(t *testing.T) {
	mything := createThing()
	booleanData := dataSchema.NewBoolean(false)
	propertyW := interaction.NewProperty(
		"boolW",
		"W bool",
		"Writable only/Not Observable boolean",
		false,
		true,
		false,
		booleanData,
	)
	mything.AddProperty(&propertyW)
	myProducer := producer.New("127.0.0.1", 8888, false)
	myProducer.Produce(mything)
	myProducer.Expose()
	client := createClient()

	checkGet(t, client, "/boolW", 401, "{\"error\":\"Write Only property\",\"type\":\"NotAllowedError\"}")

	myProducer.Close()
}

func Test_WBoolPropertyWrite(t *testing.T) {
	mything := createThing()
	booleanData := dataSchema.NewBoolean(false)
	propertyW := interaction.NewProperty(
		"boolW",
		"W bool",
		"Writable only/Not Observable boolean",
		false,
		true,
		false,
		booleanData,
	)
	mything.AddProperty(&propertyW)
	myProducer := producer.New("127.0.0.1", 8888, false)
	myProducer.Produce(mything)
	myProducer.Expose()
	client := createClient()

	checkPut(t, client, "/boolW", "", 400, "{\"error\":\"No data provided\",\"type\":\"DataError\"}")

	checkPut(t, client, "/boolW", "true", 200, "{\"ok\":true}")

	myProducer.Close()
}
