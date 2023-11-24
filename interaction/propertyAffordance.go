package interaction

import (
	"github.com/project-eria/go-wot/dataSchema"
)

// PropertyAffordance: An Interaction Affordance that exposes
// state of the Thing. This state can then be retrieved (read)
// and/or updated (write). Things can also choose to make Properties
// observable by pushing the new state after a change.
// https://w3c.github.io/wot-thing-description/#propertyaffordance
type Property struct {
	ReadOnly   bool `json:"readOnly"`   // (default = false) Boolean value that is a hint to indicate whether a property interaction / value is read only.
	WriteOnly  bool `json:"writeOnly"`  // (default = false) Boolean value that is a hint to indicate whether a property interaction / value is write only.
	Observable bool `json:"observable"` // A hint that indicates whether Servients hosting the Thing and Intermediaries should provide a Protocol Binding that supports the observeproperty and unobserveproperty operations for this Property.
	Interaction
	dataSchema.Data
}

func NewProperty(key string, title string, description string, readOnly bool, writeOnly bool, observable bool, uriVariables map[string]dataSchema.Data, data dataSchema.Data) *Property {
	return &Property{
		Interaction: Interaction{
			Key:          key,
			Title:        title,
			Description:  description,
			Forms:        []*Form{},
			UriVariables: uriVariables,
		},
		ReadOnly:   readOnly,
		WriteOnly:  writeOnly,
		Observable: observable,

		Data: data,
	}
}
