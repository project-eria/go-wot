package interaction

import (
	"github.com/project-eria/go-wot/dataSchema"
	"github.com/project-eria/go-wot/form"
	"github.com/rs/zerolog/log"
)

type PropertyAffordance interface {
}

type Property struct {
	Value      interface{} `json:"-"`
	Observable bool        `json:"observable"` // A hint that indicates whether Servients hosting the Thing and Intermediaries should provide a Protocol Binding that supports the observeproperty and unobserveproperty operations for this Property.
	Interaction
	dataSchema.Data
}

func NewProperty(key string, title string, description string, data dataSchema.Data) Property {
	interaction := Interaction{
		Key:         key,
		Title:       title,
		Description: description,
		Forms:       []form.Form{},
	}
	return Property{
		Interaction: interaction,
		Data:        data,
	}
}

func (p *Property) GetValue() interface{} {
	log.Trace().Str("property", p.Interaction.Key).Interface("value", p.Value).Msg("[property:GetValue] Value get")
	return p.Value
}

func (p *Property) SetValue(value interface{}) error {
	log.Trace().Str("property", p.Interaction.Key).Interface("value", value).Msg("[property:SetValue] Value set")

	if err := p.Data.Check(value); err != nil {
		log.Error().Str("property", p.Interaction.Key).Interface("value", value).Err(err).Msg("[property:SetValue]")
		return err
	}
	p.Value = value
	return nil
}
