package producer_test

import (
	"net/http"

	"github.com/project-eria/go-wot/dataSchema"
	"github.com/project-eria/go-wot/interaction"
	"github.com/project-eria/go-wot/producer"
)

func (ts *ProducerTestSuite) Test_RWOBoolProperty() {
	booleanData := dataSchema.NewBoolean(false)
	propertyRWO := interaction.NewProperty(
		"boolRWO",
		"RWO bool",
		"Readable/Writable/Observable boolean",
		false,
		false,
		true,
		map[string]dataSchema.Data{},
		booleanData,
	)
	ts.myThing.AddProperty(propertyRWO)

	expect, _ := getProducer(ts, ts.myThing)

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
	// "properties": {
	properties := obj.Value("properties").Object()
	property := properties.Value("boolRWO").Object()
	property.HasValue("title", "RWO bool")
	property.HasValue("description", "Readable/Writable/Observable boolean")
	property.HasValue("readOnly", false)
	property.HasValue("writeOnly", false)
	property.HasValue("observable", true)
	property.HasValue("type", "boolean")
	property.HasValue("default", false)
	property.Value("forms").Array().Length().IsEqual(2)
	form1 := property.Value("forms").Array().Value(0).Object()
	form1.HasValue("href", "http://127.0.0.1/boolRWO")
	form1.HasValue("contentType", "application/json")
	form1.Value("op").Array().IsEqualUnordered([]string{"writeproperty", "readproperty"})
	form2 := property.Value("forms").Array().Value(1).Object()
	form2.HasValue("href", "ws://127.0.0.1/boolRWO")
	form2.HasValue("contentType", "application/json")
	form2.Value("op").Array().IsEqualUnordered([]string{"observeproperty", "unobserveproperty"})

	// 		"observable": true,
	// 		"title": "RWO bool",
	// 		"description": "Readable/Writable/Observable boolean",
	// 		"forms": [
	// 			{
	// 				"href": "http://127.0.0.1:8888/boolRWO",
	// 				"contentType": "application/json",
	// 				"op": [
	// 					"writeproperty",
	// 					"readproperty"
	// 				]
	// 			},
	// 			{
	// 				"href": "ws://127.0.0.1:8888/boolRWO",
	// 				"contentType": "application/json",
	// 				"op": [
	// 					"observeproperty",
	// 					"unobserveproperty"
	// 				]
	// 			}
	// 		],
	// 		"default": false,
	// 		"readOnly": false,
	// 		"writeOnly": false,
	// 		"type": "boolean"
	// 	}
	// },
}

func (ts *ProducerTestSuite) Test_RWOBoolPropertyRead() {
	booleanData := dataSchema.NewBoolean(false)
	propertyRWO := interaction.NewProperty(
		"boolRWO",
		"RWO bool",
		"Readable/Writable/Observable boolean",
		false,
		false,
		true,
		map[string]dataSchema.Data{},
		booleanData,
	)
	ts.myThing.AddProperty(propertyRWO)
	expect, exposedThing := getProducer(ts, ts.myThing)

	exposedThing.SetPropertyReadHandler("boolRWO", func(t producer.ExposedThing, name string, parameters map[string]interface{}) (interface{}, error) {
		return false, nil
	})

	expect.GET("/boolRWO").Expect().
		Status(http.StatusOK).JSON().Boolean().IsFalse()
}

func (ts *ProducerTestSuite) Test_RWOBoolPropertyWrite() {
	booleanData := dataSchema.NewBoolean(false)
	propertyRWO := interaction.NewProperty(
		"boolRWO",
		"RWO bool",
		"Readable/Writable/Observable boolean",
		false,
		false,
		true,
		map[string]dataSchema.Data{},
		booleanData,
	)
	ts.myThing.AddProperty(propertyRWO)

	expect, exposedThing := getProducer(ts, ts.myThing)

	exposedThing.SetPropertyReadHandler("boolRWO", func(t producer.ExposedThing, name string, parameters map[string]interface{}) (interface{}, error) {
		return true, nil
	})
	exposedThing.SetPropertyWriteHandler("boolRWO", func(t producer.ExposedThing, name string, value interface{}, parameters map[string]interface{}) error {
		return nil
	})

	expect.PUT("/boolRWO").WithHeader("Content-type", "application/json").
		Expect().
		Status(http.StatusBadRequest).JSON().Object().HasValue("error", "No data provided").HasValue("type", "DataError")

	expect.PUT("/boolRWO").WithHeader("Content-type", "application/json").
		WithJSON(true).
		Expect().
		Status(http.StatusOK).JSON().Object().HasValue("ok", true)

	expect.GET("/boolRWO").Expect().Status(http.StatusOK).JSON().Boolean().IsTrue()
}

func (ts *ProducerTestSuite) Test_RWBoolProperty() {
	booleanData := dataSchema.NewBoolean(false)
	propertyRW := interaction.NewProperty(
		"boolRW",
		"RW bool",
		"Readable/Writable/Not Observable boolean",
		false,
		false,
		false,
		map[string]dataSchema.Data{},
		booleanData,
	)
	ts.myThing.AddProperty(propertyRW)

	expect, _ := getProducer(ts, ts.myThing)

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
	//	"version": {
	//		"instance": "0.0.0-dev"
	//	},
	obj.Value("version").Object().HasValue("instance", "0.0.0-dev")
	//	"securityDefinitions": {
	//		"no_sec": {
	//			"scheme": "nosec"
	//		}
	//	},
	obj.Value("securityDefinitions").Object().HasValue("no_sec", map[string]interface{}{"scheme": "nosec"})
	// "security": "no_sec"
	obj.Value("security").String().IsEqual("no_sec")
	// "properties": {
	properties := obj.Value("properties").Object()
	property := properties.Value("boolRW").Object()
	property.HasValue("title", "RW bool")
	property.HasValue("description", "Readable/Writable/Not Observable boolean")
	property.HasValue("readOnly", false)
	property.HasValue("writeOnly", false)
	property.HasValue("observable", false)
	property.HasValue("type", "boolean")
	property.HasValue("default", false)
	property.Value("forms").Array().Length().IsEqual(1)
	form1 := property.Value("forms").Array().Value(0).Object()
	form1.HasValue("href", "http://127.0.0.1/boolRW")
	form1.HasValue("contentType", "application/json")
	form1.Value("op").Array().IsEqualUnordered([]string{"writeproperty", "readproperty"})
}

func (ts *ProducerTestSuite) Test_RWBoolPropertyRead() {
	booleanData := dataSchema.NewBoolean(false)
	propertyRW := interaction.NewProperty(
		"boolRW",
		"RW bool",
		"Readable/Writable/Not Observable boolean",
		false,
		false,
		false,
		map[string]dataSchema.Data{},
		booleanData,
	)
	ts.myThing.AddProperty(propertyRW)

	expect, exposedThing := getProducer(ts, ts.myThing)

	exposedThing.SetPropertyReadHandler("boolRW", func(t producer.ExposedThing, name string, parameters map[string]interface{}) (interface{}, error) {
		return false, nil
	})

	expect.GET("/boolRW").Expect().Status(http.StatusOK).JSON().Boolean().IsFalse()
}

func (ts *ProducerTestSuite) Test_RWBoolPropertyWrite() {
	booleanData := dataSchema.NewBoolean(false)
	propertyRW := interaction.NewProperty(
		"boolRW",
		"RW bool",
		"Readable/Writable/Not Observable boolean",
		false,
		false,
		false,
		map[string]dataSchema.Data{},
		booleanData,
	)
	ts.myThing.AddProperty(propertyRW)

	expect, exposedThing := getProducer(ts, ts.myThing)

	exposedThing.SetPropertyReadHandler("boolRW", func(t producer.ExposedThing, name string, parameters map[string]interface{}) (interface{}, error) {
		return true, nil
	})
	exposedThing.SetPropertyWriteHandler("boolRW", func(t producer.ExposedThing, name string, value interface{}, parameters map[string]interface{}) error {
		return nil
	})

	expect.PUT("/boolRW").WithHeader("Content-type", "application/json").
		Expect().
		Status(http.StatusBadRequest).JSON().Object().HasValue("error", "No data provided").HasValue("type", "DataError")

	expect.PUT("/boolRW").WithHeader("Content-type", "application/json").
		WithJSON(true).
		Expect().
		Status(http.StatusOK).JSON().Object().HasValue("ok", true)

	expect.GET("/boolRW").Expect().Status(http.StatusOK).JSON().Boolean().IsTrue()
}

func (ts *ProducerTestSuite) Test_RBoolProperty() {
	booleanData := dataSchema.NewBoolean(false)
	propertyR := interaction.NewProperty(
		"boolR",
		"R bool",
		"Readable only/Not Observable boolean",
		true,
		false,
		false,
		map[string]dataSchema.Data{},
		booleanData,
	)
	ts.myThing.AddProperty(propertyR)

	expect, _ := getProducer(ts, ts.myThing)

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
	//	"version": {
	//		"instance": "0.0.0-dev"
	//	},
	obj.Value("version").Object().HasValue("instance", "0.0.0-dev")
	//	"securityDefinitions": {
	//		"no_sec": {
	//			"scheme": "nosec"
	//		}
	//	},
	obj.Value("securityDefinitions").Object().HasValue("no_sec", map[string]interface{}{"scheme": "nosec"})
	// "security": "no_sec"
	obj.Value("security").String().IsEqual("no_sec")
	// "properties": {
	properties := obj.Value("properties").Object()
	property := properties.Value("boolR").Object()
	property.HasValue("title", "R bool")
	property.HasValue("description", "Readable only/Not Observable boolean")
	property.HasValue("readOnly", true)
	property.HasValue("writeOnly", false)
	property.HasValue("observable", false)
	property.HasValue("type", "boolean")
	property.HasValue("default", false)
	property.Value("forms").Array().Length().IsEqual(1)
	form1 := property.Value("forms").Array().Value(0).Object()
	form1.HasValue("href", "http://127.0.0.1/boolR")
	form1.HasValue("contentType", "application/json")
	form1.Value("op").Array().IsEqualUnordered([]string{"readproperty"})
}

func (ts *ProducerTestSuite) Test_RBoolPropertyRead() {
	booleanData := dataSchema.NewBoolean(false)
	propertyR := interaction.NewProperty(
		"boolR",
		"R bool",
		"Readable only/Not Observable boolean",
		true,
		false,
		false,
		map[string]dataSchema.Data{},
		booleanData,
	)
	ts.myThing.AddProperty(propertyR)

	expect, exposedThing := getProducer(ts, ts.myThing)

	exposedThing.SetPropertyReadHandler("boolR", func(t producer.ExposedThing, name string, parameters map[string]interface{}) (interface{}, error) {
		return false, nil
	})

	expect.GET("/boolR").Expect().Status(http.StatusOK).JSON().Boolean().IsFalse()
}

func (ts *ProducerTestSuite) Test_RBoolPropertyWrite() {
	booleanData := dataSchema.NewBoolean(false)
	propertyR := interaction.NewProperty(
		"boolR",
		"R bool",
		"Readable only/Not Observable boolean",
		true,
		false,
		false,
		map[string]dataSchema.Data{},
		booleanData,
	)
	ts.myThing.AddProperty(propertyR)

	expect, _ := getProducer(ts, ts.myThing)

	expect.PUT("/boolR").WithHeader("Content-type", "application/json").
		WithJSON(true).Expect().Status(http.StatusUnauthorized).JSON().Object().HasValue("error", "Read Only property").HasValue("type", "NotAllowedError")
}

func (ts *ProducerTestSuite) Test_WBoolProperty() {
	booleanData := dataSchema.NewBoolean(false)
	propertyW := interaction.NewProperty(
		"boolW",
		"W bool",
		"Writable only/Not Observable boolean",
		false,
		true,
		false,
		map[string]dataSchema.Data{},
		booleanData,
	)
	ts.myThing.AddProperty(propertyW)

	expect, _ := getProducer(ts, ts.myThing)

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
	//	"version": {
	//		"instance": "0.0.0-dev"
	//	},
	obj.Value("version").Object().HasValue("instance", "0.0.0-dev")
	//	"securityDefinitions": {
	//		"no_sec": {
	//			"scheme": "nosec"
	//		}
	//	},
	obj.Value("securityDefinitions").Object().HasValue("no_sec", map[string]interface{}{"scheme": "nosec"})
	// "security": "no_sec"
	obj.Value("security").String().IsEqual("no_sec")
	// "properties": {
	properties := obj.Value("properties").Object()
	property := properties.Value("boolW").Object()
	property.HasValue("title", "W bool")
	property.HasValue("description", "Writable only/Not Observable boolean")
	property.HasValue("readOnly", false)
	property.HasValue("writeOnly", true)
	property.HasValue("observable", false)
	property.HasValue("type", "boolean")
	property.HasValue("default", false)
	property.Value("forms").Array().Length().IsEqual(1)
	form1 := property.Value("forms").Array().Value(0).Object()
	form1.HasValue("href", "http://127.0.0.1/boolW")
	form1.HasValue("contentType", "application/json")
	form1.Value("op").Array().IsEqualUnordered([]string{"writeproperty"})
}

func (ts *ProducerTestSuite) Test_WBoolPropertyRead() {
	booleanData := dataSchema.NewBoolean(false)
	propertyW := interaction.NewProperty(
		"boolW",
		"W bool",
		"Writable only/Not Observable boolean",
		false,
		true,
		false,
		map[string]dataSchema.Data{},
		booleanData,
	)
	ts.myThing.AddProperty(propertyW)

	expect, _ := getProducer(ts, ts.myThing)

	expect.GET("/boolW").Expect().Status(http.StatusUnauthorized).JSON().Object().HasValue("error", "Write Only property").HasValue("type", "NotAllowedError")
}

func (ts *ProducerTestSuite) Test_WBoolPropertyWrite() {
	booleanData := dataSchema.NewBoolean(false)
	propertyW := interaction.NewProperty(
		"boolW",
		"W bool",
		"Writable only/Not Observable boolean",
		false,
		true,
		false,
		map[string]dataSchema.Data{},
		booleanData,
	)
	ts.myThing.AddProperty(propertyW)

	expect, exposedThing := getProducer(ts, ts.myThing)

	exposedThing.SetPropertyWriteHandler("boolW", func(t producer.ExposedThing, name string, value interface{}, parameters map[string]interface{}) error {
		return nil
	})

	expect.PUT("/boolW").WithHeader("Content-type", "application/json").
		Expect().
		Status(http.StatusBadRequest).JSON().Object().HasValue("error", "No data provided").HasValue("type", "DataError")

	expect.PUT("/boolW").WithHeader("Content-type", "application/json").
		WithJSON(true).
		Expect().
		Status(http.StatusOK).JSON().Object().HasValue("ok", true)
}

func (ts *ProducerTestSuite) Test_URIVariablesProperty() {
	booleanData := dataSchema.NewBoolean(false)
	propertyRWO := interaction.NewProperty(
		"uriVars",
		"URI Variables",
		"With URI Variables",
		false,
		false,
		true,
		map[string]dataSchema.Data{
			"var1": {
				Default: "",
				Type:    "string",
			},
		},
		booleanData,
	)
	ts.myThing.AddProperty(propertyRWO)

	expect, _ := getProducer(ts, ts.myThing)

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
	// "properties": {
	properties := obj.Value("properties").Object()
	property := properties.Value("uriVars").Object()
	property.HasValue("title", "URI Variables")
	property.HasValue("description", "With URI Variables")
	property.HasValue("readOnly", false)
	property.HasValue("writeOnly", false)
	property.HasValue("observable", true)
	property.HasValue("type", "boolean")
	property.HasValue("default", false)
	property.Value("forms").Array().Length().IsEqual(2)
	form1 := property.Value("forms").Array().Value(0).Object()
	form1.HasValue("href", "http://127.0.0.1/uriVars/{var1}")
	form1.HasValue("contentType", "application/json")
	form1.Value("op").Array().IsEqualUnordered([]string{"writeproperty", "readproperty"})
	form2 := property.Value("forms").Array().Value(1).Object()
	form2.HasValue("href", "ws://127.0.0.1/uriVars/{var1}")
	form2.HasValue("contentType", "application/json")
	form2.Value("op").Array().IsEqualUnordered([]string{"observeproperty", "unobserveproperty"})
	uriVariables := property.Value("uriVariables").Object()
	uriVariables.HasValue("var1", map[string]interface{}{
		"default": "",
		"type":    "string",
	})
}
