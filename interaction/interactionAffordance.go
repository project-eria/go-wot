package interaction

import (
	"fmt"

	"github.com/project-eria/go-wot/dataSchema"
)

// InteractionAffordance:
// Metadata of a Thing that shows the possible choices to Consumers,
// thereby suggesting how Consumers may interact with the Thing.
// There are many types of potential affordances, but W3C WoT defines
// three types of Interaction Affordances: Properties, Actions, and Events.
// https://w3c.github.io/wot-thing-description/#interactionaffordance
type Interaction struct {
	Key          string                     `json:"-"`
	AtType       []string                   `json:"@type,omitempty"`        // (optional) JSON-LD keyword to label the object with semantic tags (or types)
	Title        string                     `json:"title,omitempty"`        // (optional) Provides a human-readable title (e.g., display a text for UI representation) based on a default language.
	Titles       map[string]string          `json:"titles,omitempty"`       // (optional) Provides multi-language human-readable titles (e.g., display a text for UI representation in different languages).
	Description  string                     `json:"description,omitempty"`  // (optional) Provides additional (human-readable) information based on a default language.
	Descriptions map[string]string          `json:"descriptions,omitempty"` // (optional) Can be used to support (human-readable) information in different languages.
	Forms        []*Form                    `json:"forms"`                  // (mandatory) Set of form hypermedia controls that describe how an operation can be performed. Forms are serializations of Protocol Bindings.
	UriVariables map[string]dataSchema.Data `json:"uriVariables,omitempty"` // Define URI query template variables as collection based on DataSchema declarations. The individual variables DataSchema cannot be an ObjectSchema or an ArraySchema.	optional	Map of DataSchema
}

func (i Interaction) CheckUriVariables(variables map[string]string) (map[string]interface{}, error) {
	res := make(map[string]interface{})
	for varKey, dataDef := range i.UriVariables {
		if valueStr, ok := variables[varKey]; ok {
			value, err := dataDef.FromString(valueStr)
			if err != nil {
				return nil, fmt.Errorf("incorrect param `%s` value: %s", varKey, err.Error())
			}
			if err := dataDef.Validate(value); err != nil {
				return nil, fmt.Errorf("incorrect param `%s` value: %s", varKey, err.Error())
			}
			res[varKey] = value
		} else {
			return nil, fmt.Errorf("missing required parameter `%s`", varKey)
		}
	}
	return res, nil
}
