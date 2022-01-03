package interaction

import (
	"github.com/project-eria/go-wot/dataSchema"
	"github.com/project-eria/go-wot/form"
)

type PropertyAffordance interface {
}

type Property struct {
	Observable bool `json:"observable"` // A hint that indicates whether Servients hosting the Thing and Intermediaries should provide a Protocol Binding that supports the observeproperty and unobserveproperty operations for this Property.
	Interaction
	dataSchema.Data
}

func NewProperty(key string, title string, description string, readOnly bool, writeOnly bool, observable bool, data dataSchema.Data) Property {
	interaction := Interaction{
		Key:         key,
		Title:       title,
		Description: description,
		Forms:       []form.Form{},
	}
	data.ReadOnly = readOnly
	data.WriteOnly = writeOnly

	return Property{
		Interaction: interaction,
		Observable:  observable,
		Data:        data,
	}
}
