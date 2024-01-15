package interaction

import (
	"encoding/json"

	"github.com/project-eria/go-wot/dataSchema"
)

// PropertyAffordance: An Interaction Affordance that exposes
// state of the Thing. This state can then be retrieved (read)
// and/or updated (write). Things can also choose to make Properties
// observable by pushing the new state after a change.
// https://w3c.github.io/wot-thing-description/#propertyaffordance

// Note: Data is not a pointer because it's not optional

type Property struct {
	ReadOnly    bool            `json:"readOnly"`   // (default = false) Boolean value that is a hint to indicate whether a property interaction / value is read only.
	WriteOnly   bool            `json:"writeOnly"`  // (default = false) Boolean value that is a hint to indicate whether a property interaction / value is write only.
	Observable  bool            `json:"observable"` // A hint that indicates whether Servients hosting the Thing and Intermediaries should provide a Protocol Binding that supports the observeproperty and unobserveproperty operations for this Property.
	Data        dataSchema.Data `json:"-"`          // not embedded, because of MarshalJSON
	Interaction `json:"-"`
}

func NewProperty(key string, title string, description string, readOnly bool, writeOnly bool, observable bool, uriVariables map[string]dataSchema.Data, data dataSchema.Data) *Property {
	// TODO readOnly, writeOnly can't be true at the same time
	return &Property{
		Interaction: Interaction{
			Key:          key,
			Title:        title,
			Description:  description,
			Forms:        []*Form{},
			UriVariables: uriVariables,
		},
		Data:       data,
		ReadOnly:   readOnly,
		WriteOnly:  writeOnly,
		Observable: observable,
	}
}

func (p *Property) MarshalJSON() ([]byte, error) {
	type PropertyOrigin Property

	b1, err := json.Marshal(p.Interaction)
	if err != nil {
		return nil, err
	}
	final := b1[:len(b1)-1] // remove last parenthesis

	b2, err := p.Data.MarshalJSON()
	if err != nil {
		return nil, err
	}

	b2 = b2[:len(b2)-1] // remove last parenthesis
	b2[0] = ','         // replace first parenthesis, with a comma
	final = append(final, b2...)

	b3, err := json.Marshal((*PropertyOrigin)(p))
	if err != nil {
		return nil, err
	}

	b3[0] = ',' // replace first parenthesis, with a comma
	final = append(final, b3...)

	return final, nil
}

func (p *Property) UnmarshalJSON(data []byte) error {
	type PropertyOrigin Property
	po := &struct {
		*PropertyOrigin
	}{
		PropertyOrigin: (*PropertyOrigin)(p),
	}
	if err := json.Unmarshal(data, po); err != nil {
		return err
	}

	i := Interaction{}
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	p.Interaction = i

	d := dataSchema.Data{}
	if err := d.UnmarshalJSON(data); err != nil {
		return err
	}
	p.Data = d

	return nil
}
