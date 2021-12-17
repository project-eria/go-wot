package interaction

import (
	"github.com/project-eria/go-wot/dataSchema"
	"github.com/project-eria/go-wot/form"
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

func (p *Property) AddHrefForm(host string, secure bool) {
	scheme := "http"
	if secure {
		scheme = "https"
	}
	op := []string{}
	if !p.ReadOnly {
		op = append(op, "writeproperty")
	}
	if !p.WriteOnly {
		op = append(op, "readproperty")
	}
	url := scheme + "://" + host
	p.Interaction.AddHrefForm(url,
		form.Form{
			ContentType: "application/json",
			Op:          op,
		},
	)
}

func (p *Property) GetValue() interface{} {
	return p.Value
}

func (p *Property) SetValue(value interface{}) error {
	if err := p.Data.Check(value); err != nil {
		return err
	}
	p.Value = value
	return nil
}
