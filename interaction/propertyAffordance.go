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
	Observable bool `json:"observable"` // A hint that indicates whether Servients hosting the Thing and Intermediaries should provide a Protocol Binding that supports the observeproperty and unobserveproperty operations for this Property.
	Interaction
	dataSchema.Data
}

func NewProperty(key string, title string, description string, readOnly bool, writeOnly bool, observable bool, data dataSchema.Data) *Property {
	interaction := Interaction{
		Key:         key,
		Title:       title,
		Description: description,
		Forms:       []*Form{},
	}
	data.ReadOnly = readOnly
	data.WriteOnly = writeOnly

	return &Property{
		Interaction: interaction,
		Observable:  observable,
		Data:        data,
	}
}
