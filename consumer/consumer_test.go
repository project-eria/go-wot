package consumer_test

import (
	"github.com/project-eria/go-wot/dataSchema"
	"github.com/project-eria/go-wot/interaction"
	"github.com/project-eria/go-wot/securityScheme"
	"github.com/project-eria/go-wot/thing"
)

func (ts *ConsumerTestSuite) Test_GetThingDescription() {
	td := thing.Thing{
		ID:          "urn:dev:ops:my-actuator-1234",
		AtContext:   map[string]string{"": "https://www.w3.org/2022/wot/td/v1.1"},
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
				ReadOnly:   true,
				Interaction: interaction.Interaction{
					Key: "",
					// AtType: ([]string) <nil>,
					Title: "R bool",
					// Titles: (map[string]string) <nil>,
					Description: "Readable only/Not Observable boolean",
					// Descriptions: (map[string]string) <nil>,
					Forms: []*interaction.Form{
						{
							Href:        "http://127.0.0.1:8888/boolR",
							ContentType: "application/json",
							Op:          []string{"readproperty"},
							Supplement:  map[string]interface{}{},
						},
					},
				},
				Data: dataSchema.Data{
					// Const: (interface {}) <nil>,
					Default: false,
					Unit:    "",
					// +    Enum: ([]interface {}) <nil>,
					// Format:     "",
					Type:       "boolean",
					DataSchema: dataSchema.Boolean{},
				},
			},
			"boolW": {
				Observable: false,
				WriteOnly:  true,
				Interaction: interaction.Interaction{
					Key: "",
					// AtType: ([]string) <nil>,
					Title: "W bool",
					// Titles: (map[string]string) <nil>,
					Description: "Writable only/Not Observable boolean",
					// Descriptions: (map[string]string) <nil>,
					Forms: []*interaction.Form{
						{
							Href:        "http://127.0.0.1:8888/boolW",
							ContentType: "application/json",
							Op:          []string{"writeproperty"},
							Supplement:  map[string]interface{}{},
						},
					},
				},
				Data: dataSchema.Data{
					// Const: (interface {}) <nil>,
					Default: false,
					Unit:    "",
					// +    Enum: ([]interface {}) <nil>,
					// Format:     "",
					Type:       "boolean",
					DataSchema: dataSchema.Boolean{},
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
					Forms: []*interaction.Form{
						{
							Href:        "http://127.0.0.1:8888/boolRW",
							ContentType: "application/json",
							Op:          []string{"writeproperty", "readproperty"},
							Supplement:  map[string]interface{}{},
						},
					},
				},
				Data: dataSchema.Data{
					// Const: (interface {}) <nil>,
					Default: false,
					Unit:    "",
					// +    Enum: ([]interface {}) <nil>,
					// Format:     "",
					Type:       "boolean",
					DataSchema: dataSchema.Boolean{},
				},
			},
			"uriVars": {
				Observable: false,
				Interaction: interaction.Interaction{
					Key: "",
					// AtType: ([]string) <nil>,
					Title: "URI Vars",
					// Titles: (map[string]string) <nil>,
					Description: "With URI Vars",
					// Descriptions: (map[string]string) <nil>,
					Forms: []*interaction.Form{
						{
							Href:        "http://127.0.0.1:8888/uriVars/{var1}/{var2}",
							ContentType: "application/json",
							Op:          []string{"writeproperty", "readproperty"},
							Supplement:  map[string]interface{}{},
						},
					},
					UriVariables: map[string]dataSchema.Data{
						"var1": {
							Default:    "",
							Type:       "string",
							DataSchema: dataSchema.String{},
						},
						"var2": {
							Default:    "test",
							Type:       "string",
							DataSchema: dataSchema.String{},
						},
					},
				},
				Data: dataSchema.Data{
					// Const: (interface {}) <nil>,
					Default: false,
					Unit:    "",
					// +    Enum: ([]interface {}) <nil>,
					// Format:     "",
					Type:       "boolean",
					DataSchema: dataSchema.Boolean{},
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
					Forms: []*interaction.Form{
						{
							Href:        "http://127.0.0.1:8888/a",
							ContentType: "application/json",
							Op: []string{
								"invokeaction",
							},
							Supplement: map[string]interface{}{},
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
					// Format:     "",
					Type:       "string",
					DataSchema: dataSchema.String{},
				},
				//    Output: (*dataSchema.Data)(<nil>),
				Interaction: interaction.Interaction{
					Key: "",
					//    AtType: ([]string) <nil>,
					Title: "String Input, No Output",
					//    Titles: (map[string]string) <nil>,
					Description: (""),
					//    Descriptions: (map[string]string) <nil>,
					Forms: []*interaction.Form{
						{
							Href:        "http://127.0.0.1:8888/b",
							ContentType: "application/json",
							Op: []string{
								"invokeaction",
							},
							Supplement: map[string]interface{}{},
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
					// Format:     "",
					Type:       "string",
					DataSchema: dataSchema.String{},
				},
				Output: &dataSchema.Data{
					// Const: (interface {}) <nil>,
					Default: "",
					Unit:    "",
					// Enum: ([]interface {}) <nil>,
					// Format:     "",
					Type:       "string",
					DataSchema: dataSchema.String{},
				},
				Interaction: interaction.Interaction{
					Key: "",
					//    AtType: ([]string) <nil>,
					Title: "String Input, String Output",
					//    Titles: (map[string]string) <nil>,
					Description: (""),
					//    Descriptions: (map[string]string) <nil>,
					Forms: []*interaction.Form{
						{
							Href:        "http://127.0.0.1:8888/c",
							ContentType: "application/json",
							Op: []string{
								"invokeaction",
							},
							Supplement: map[string]interface{}{},
						},
					},
				},
			},
		},
	}

	tdGet := ts.consumedThing.GetThingDescription()
	ts.Equal(&td, tdGet, "they should be equal")
}
