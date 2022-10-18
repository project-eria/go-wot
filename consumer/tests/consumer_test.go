package test

import (
	"testing"

	"github.com/project-eria/go-wot/dataSchema"
	"github.com/project-eria/go-wot/interaction"
	"github.com/project-eria/go-wot/securityScheme"
	"github.com/project-eria/go-wot/thing"
	"github.com/stretchr/testify/assert"
)

func Test_GetThingDescription(t *testing.T) {
	td := thing.Thing{
		ID:          "urn:dev:ops:my-actuator-1234",
		AtContext:   "https://www.w3.org/2022/wot/td/v1.1",
		Title:       "Actuator1 Example",
		Description: "An actuator 1st example",
		SecurityDefinitions: map[string]securityScheme.SecurityScheme{
			"no_sec": map[string]interface{}{
				"scheme": "nosec",
			},
			// "no_sec": securityScheme.Security{
			// 	Scheme: "nosec",
			// },
		},
		Security: []string{"no_sec"},
		Properties: map[string]*interaction.Property{
			"boolR": {
				Observable: false,
				Interaction: interaction.Interaction{
					Key: "",
					// AtType: ([]string) <nil>,
					Title: "R bool",
					// Titles: (map[string]string) <nil>,
					Description: "Readable only/Not Observable boolean",
					// Descriptions: (map[string]string) <nil>,
					Forms: []interaction.Form{
						{
							Href:        "http://127.0.0.1:8888/boolR",
							ContentType: "application/json",
							Op:          []string{"readproperty"},
						},
					},
				},
				Data: dataSchema.Data{
					// Const: (interface {}) <nil>,
					Default: false,
					Unit:    "",
					// +    Enum: ([]interface {}) <nil>,
					ReadOnly:         true,
					WriteOnly:        false,
					Format:           "",
					ContentEncoding:  "",
					ContentMediaType: "",
					Type:             "boolean",
					// +    DataSchema: (dataSchema.DataSchema) <nil>
				},
			},
			"boolW": {
				Observable: false,
				Interaction: interaction.Interaction{
					Key: "",
					// AtType: ([]string) <nil>,
					Title: "W bool",
					// Titles: (map[string]string) <nil>,
					Description: "Writable only/Not Observable boolean",
					// Descriptions: (map[string]string) <nil>,
					Forms: []interaction.Form{
						{
							Href:        "http://127.0.0.1:8888/boolW",
							ContentType: "application/json",
							Op:          []string{"writeproperty"},
						},
					},
				},
				Data: dataSchema.Data{
					// Const: (interface {}) <nil>,
					Default: false,
					Unit:    "",
					// +    Enum: ([]interface {}) <nil>,
					ReadOnly:         false,
					WriteOnly:        true,
					Format:           "",
					ContentEncoding:  "",
					ContentMediaType: "",
					Type:             "boolean",
					// +    DataSchema: (dataSchema.DataSchema) <nil>
				},
			},
			"boolRW": {
				Observable: false,
				Interaction: interaction.Interaction{
					Key: "",
					// AtType: ([]string) <nil>,
					Title: "RW bool",
					// Titles: (map[string]string) <nil>,
					Description: "Readable/Writable/Not Observable boolean",
					// Descriptions: (map[string]string) <nil>,
					Forms: []interaction.Form{
						{
							Href:        "http://127.0.0.1:8888/boolRW",
							ContentType: "application/json",
							Op:          []string{"writeproperty", "readproperty"},
						},
					},
				},
				Data: dataSchema.Data{
					// Const: (interface {}) <nil>,
					Default: false,
					Unit:    "",
					// +    Enum: ([]interface {}) <nil>,
					ReadOnly:         false,
					WriteOnly:        false,
					Format:           "",
					ContentEncoding:  "",
					ContentMediaType: "",
					Type:             "boolean",
					// +    DataSchema: (dataSchema.DataSchema) <nil>
				},
			},
		},
		Actions: map[string]*interaction.Action{
			"a": {
				//    Input: (*dataSchema.Data)(<nil>),
				//    Output: (*dataSchema.Data)(<nil>),
				Interaction: interaction.Interaction{
					Key: "",
					//    AtType: ([]string) <nil>,
					Title: "No Input, No Output",
					//    Titles: (map[string]string) <nil>,
					Description: (""),
					//    Descriptions: (map[string]string) <nil>,
					Forms: []interaction.Form{
						{
							Href:        "http://127.0.0.1:8888/a",
							ContentType: "application/json",
							Op: []string{
								"invokeaction",
							},
						},
					},
				},
			},
			"b": {
				Input: &dataSchema.Data{
					// Const: (interface {}) <nil>,
					Default: "",
					Unit:    "",
					// Enum: ([]interface {}) <nil>,
					ReadOnly:         false,
					WriteOnly:        false,
					Format:           "",
					ContentEncoding:  "",
					ContentMediaType: "",
					Type:             "string",
					// DataSchema: (dataSchema.DataSchema) <nil>
				},
				//    Output: (*dataSchema.Data)(<nil>),
				Interaction: interaction.Interaction{
					Key: "",
					//    AtType: ([]string) <nil>,
					Title: "String Input, No Output",
					//    Titles: (map[string]string) <nil>,
					Description: (""),
					//    Descriptions: (map[string]string) <nil>,
					Forms: []interaction.Form{
						{
							Href:        "http://127.0.0.1:8888/b",
							ContentType: "application/json",
							Op: []string{
								"invokeaction",
							},
						},
					},
				},
			},
			"c": {
				Input: &dataSchema.Data{
					// Const: (interface {}) <nil>,
					Default: "",
					Unit:    "",
					// Enum: ([]interface {}) <nil>,
					ReadOnly:         false,
					WriteOnly:        false,
					Format:           "",
					ContentEncoding:  "",
					ContentMediaType: "",
					Type:             "string",
					// DataSchema: (dataSchema.DataSchema) <nil>
				},
				Output: &dataSchema.Data{
					// Const: (interface {}) <nil>,
					Default: "",
					Unit:    "",
					// Enum: ([]interface {}) <nil>,
					ReadOnly:         false,
					WriteOnly:        false,
					Format:           "",
					ContentEncoding:  "",
					ContentMediaType: "",
					Type:             "string",
					// DataSchema: (dataSchema.DataSchema) <nil>
				},
				Interaction: interaction.Interaction{
					Key: "",
					//    AtType: ([]string) <nil>,
					Title: "String Input, String Output",
					//    Titles: (map[string]string) <nil>,
					Description: (""),
					//    Descriptions: (map[string]string) <nil>,
					Forms: []interaction.Form{
						{
							Href:        "http://127.0.0.1:8888/c",
							ContentType: "application/json",
							Op: []string{
								"invokeaction",
							},
						},
					},
				},
			},
		},
	}

	tdGet := consumedThing.GetThingDescription()
	assert.Equal(t, &td, tdGet, "they should be equal")
}
