package test

import (
	"testing"

	"github.com/project-eria/go-wot/producer"
)

func Test_Root(t *testing.T) {
	mything := createThing()
	myProducer := producer.New("127.0.0.1", 8888, false)
	myProducer.Produce(mything)
	myProducer.Expose()
	client := createClient()

	checkGet(t, client, "", 200, "{\"id\":\"urn:dev:ops:my-actuator-1234\",\"@context\":\"http://www.w3.org/ns/td\",\"title\":\"Actuator1 Example\",\"description\":\"An actuator 1st example\",\"securityDefinitions\":{\"no_sec\":{\"scheme\":\"nosec\"}},\"security\":\"no_sec\"}")

	checkGet(t, client, "/x", 404, "{\"error\":\"Property not found\",\"type\":\"NotFoundError\"}")

	myProducer.Close()
}
