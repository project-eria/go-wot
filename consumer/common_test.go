package consumer_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/project-eria/go-wot/consumer"
	"github.com/project-eria/go-wot/mocks"
	"github.com/project-eria/go-wot/protocolHttp"
	"github.com/project-eria/go-wot/thing"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

const jsonTD string = `{
	"id": "urn:dev:ops:my-actuator-1234",
	"@context": "https://www.w3.org/2022/wot/td/v1.1",
	"title": "Actuator1 Example",
	"description": "An actuator 1st example",
	"securityDefinitions": {
		"no_sec": {
			"scheme": "nosec"
		}
	},
	"security": "no_sec",
	"properties": {
		"boolR": {
			"observable": false,
			"title": "R bool",
			"description": "Readable only/Not Observable boolean",
			"forms": [{
				"href": "http://127.0.0.1:8888/boolR",
				"contentType": "application/json",
				"op": [
				  	"readproperty"
				]
			}],
			"default": false,
			"readOnly": true,
			"writeOnly": false,
			"type": "boolean"
		},
		"boolRW": {
			"observable": false,
			"title": "RW bool",
			"description": "Readable/Writable/Not Observable boolean",
			"forms": [{
				"href": "http://127.0.0.1:8888/boolRW",
				"contentType": "application/json",
				"op": [
					"writeproperty",
					"readproperty"
				]
			}],
			"default": false,
			"readOnly": false,
			"writeOnly": false,
			"type": "boolean"
		},
		"boolW": {
			"observable": false,
			"title": "W bool",
			"description": "Writable only/Not Observable boolean",
			"forms": [
				{
				"href": "http://127.0.0.1:8888/boolW",
				"contentType": "application/json",
				"op": [
					"writeproperty"
				]
				}
			],
			"default": false,
			"readOnly": false,
			"writeOnly": true,
			"type": "boolean"
		}
	},
	"actions": {
		"a": {
			"title": "No Input, No Output",
			"forms": [{
				"href": "http://127.0.0.1:8888/a",
				"contentType": "application/json",
				"op": [
					"invokeaction"
				]}
			]
		},
		"b": {
			"input": {
				"default": "",
				"readOnly": false,
				"writeOnly": false,
				"type": "string"
			},
			"title": "String Input, No Output",
			"forms": [{
				"href": "http://127.0.0.1:8888/b",
				"contentType": "application/json",
				"op": [
					"invokeaction"
				]
			}]
		},
		"c": {
			"input": {
				"default": "",
				"readOnly": false,
				"writeOnly": false,
				"type": "string"
			},
			"output": {
				"default": "",
				"readOnly": false,
				"writeOnly": false,
				"type": "string"
			},
			"title": "String Input, String Output",
			"forms": [{
				"href": "http://127.0.0.1:8888/c",
				"contentType": "application/json",
				"op": [
					"invokeaction"
				]
			}]
		}
	}
}`

//   "boolRWO": {
// 	"observable": true,
// 	"title": "RWO bool",
// 	"description": "Readable/Writable/Observable boolean",
// 	"forms": [
// 	  {
// 		"href": "http://127.0.0.1:8888/boolRWO",
// 		"contentType": "application/json",
// 		"op": [
// 		  "writeproperty",
// 		  "readproperty"
// 		]
// 	  },
// 	  {
// 		"href": "ws://127.0.0.1:8888/boolRWO",
// 		"contentType": "application/json",
// 		"op": [
// 		  "observeproperty",
// 		  "unobserveproperty"
// 		]
// 	  }
// 	],
// 	"default": false,
// 	"readOnly": false,
// 	"writeOnly": false,
// 	"type": "boolean"
//   },
// },

// func TestMain(m *testing.M) {
// 	//	zerolog.SetGlobalLevel(zerolog.Disabled)
// 	var td thing.Thing
// 	if err := json.Unmarshal([]byte(jsonTD), &td); err != nil {
// 		println(err.Error())
// 		os.Exit(1)
// 	}
// 	myConsumer := consumer.New()
// 	httpClientProcessor = &mocks.HttpClientProcessor{}
// 	mockClient := &protocolHttp.HttpClient{
// 		Client:  httpClientProcessor,
// 		Schemes: []string{"http"},
// 	}
// 	myConsumer.AddClient(mockClient)
// 	consumedThing = myConsumer.Consume(&td)
// 	exitVal := m.Run()
// 	os.Exit(exitVal)
// }

type ConsumerTestSuite struct {
	consumedThing       consumer.ConsumedThing
	httpClientProcessor *mocks.HttpClientProcessor
	suite.Suite
}

func Test_ConsumerTestSuite(t *testing.T) {
	suite.Run(t, &ConsumerTestSuite{})
}

func (ts *ConsumerTestSuite) SetupTest() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	var td thing.Thing
	if err := json.Unmarshal([]byte(jsonTD), &td); err != nil {
		println(err.Error())
		os.Exit(1)
	}
	myConsumer := consumer.New()
	httpClientProcessor := &mocks.HttpClientProcessor{}
	mockClient := &protocolHttp.HttpClient{
		Client:  httpClientProcessor,
		Schemes: []string{"http"},
	}
	myConsumer.AddClient(mockClient)
	ts.consumedThing = myConsumer.Consume(&td)
	ts.httpClientProcessor = httpClientProcessor
}
