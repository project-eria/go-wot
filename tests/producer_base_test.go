package test

import (
	"net/http"
	"testing"

	"github.com/project-eria/go-wot/securityScheme"
)

//	{
//	    "@context": "",
//	    "id": "urn:dev:ops:my-actuator-1234",
//	    "title": "Actuator1 Example",
//	    "description": "An actuator 1st example",
//	    "version": {
//	        "instance": "0.0.0-dev"
//	    },
//	    "securityDefinitions": {},
//	    "security": []
//	}
func TestBase(t *testing.T) {
	mything := getThing()
	httpServer, _ := getProducer(mything)
	e := HTMLTester(t, httpServer)

	e.GET("/i-dont-exist").Expect().
		Status(http.StatusNotFound).
		Body().IsEqual("Cannot GET /i-dont-exist")

	e = JSONTester(t, httpServer)
	obj := e.GET("/").Expect().
		Status(http.StatusOK).JSON().Object()
	obj.HasValue("id", "urn:dev:ops:my-actuator-1234")
	obj.HasValue("title", "Example")
	obj.HasValue("description", "An example")
	obj.Value("version").Object().HasValue("instance", "0.0.0-dev")
	obj.Value("@context").String().IsEmpty()
	obj.Value("securityDefinitions").Object().IsEmpty()
	obj.Value("security").Array().IsEmpty()
}

//	{
//		"@context": [
//		  "https://www.w3.org/2022/wot/td/v1.1",
//		  {
//			"schema": "https://schema.org/"
//		  }
//		],
//		"description": "An example",
//		"id": "urn:dev:ops:my-actuator-1234",
//		"security": [],
//		"securityDefinitions": {},
//		"title": "Example",
//		"version": {
//		  "instance": "0.0.0-dev",
//		  "schema:softwareVersion": "1.1.1"
//		}
//	}
func TestBaseContext(t *testing.T) {
	mything := getThing()
	mything.AddContext("schema", "https://schema.org/")
	mything.AddVersion("schema:softwareVersion", "1.1.1")
	httpServer, _ := getProducer(mything)
	e := JSONTester(t, httpServer)

	obj := e.GET("/").Expect().
		Status(http.StatusOK).JSON().Object()
	obj.Value("@context").Array().Value(0).String().IsEqual("https://www.w3.org/2022/wot/td/v1.1")
	obj.Value("@context").Array().Value(1).Object().HasValue("schema", "https://schema.org/")
	obj.Value("version").Object().HasValue("schema:softwareVersion", "1.1.1")
}

//	{
//		"@context": "",
//		"description": "An example",
//		"id": "urn:dev:ops:my-actuator-1234",
//		"security": "no_sec",
//		"securityDefinitions": {
//		  "no_sec": {
//			"scheme": "nosec"
//		  }
//		},
//		"title": "Example",
//		"version": {
//		  "instance": "0.0.0-dev"
//		}
//	}
func TestBaseSecurity(t *testing.T) {
	mything := getThing()
	// Add Security
	noSecurityScheme := securityScheme.NewNoSecurity()
	mything.AddSecurity("no_sec", noSecurityScheme)
	httpServer, _ := getProducer(mything)

	e := JSONTester(t, httpServer)

	obj := e.GET("/").Expect().
		Status(http.StatusOK).JSON().Object()
	obj.Value("securityDefinitions").Object().Value("no_sec").Object().HasValue("scheme", "nosec")
	obj.HasValue("security", "no_sec")
}
