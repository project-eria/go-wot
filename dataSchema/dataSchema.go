package dataSchema

import (
	"encoding/json"
	"errors"
)

type SimpleType interface {
	bool | string | int8 | int16 | int32 | int64 | float32 | float64
}

type DataSchema interface {
	FromString(string) (interface{}, error)
	Validate(interface{}) error
}

type Data struct {
	// (Redondency, with affordance title) Title string `json:"title,omitempty"` // (optional) Provides a human-readable title (e.g., display a text for UI representation) based on a default language.
	// // titles 	Provides multi-language human-readable titles (e.g., display a text for UI representation in different languages). Also see MultiLanguage. 	optional 	Map of MultiLanguage
	// (Redondency, with affordance title) Description string `json:"description,omitempty"` // (optional)	Provides additional (human-readable) information based on a default language.
	// // descriptions 	Can be used to support (human-readable) information in different languages. Also see MultiLanguage. 	optional 	Map of MultiLanguage

	Const   interface{} `json:"const,omitempty"`   // (optional) Provides a constant value.
	Default interface{} `json:"default,omitempty"` // (optional) Supply a default value. The value should validate against the data schema in which it resides.
	Unit    string      `json:"unit,omitempty"`    // (optional) Provides unit information that is used, e.g., in international science, engineering, and business.
	// oneOf	Used to ensure that the data is valid against one of the specified schemas in the array.	optional	Array of DataSchema
	Enum []interface{} `json:"enum,omitempty"` // (optional) Restricted set of values provided as an array.
	// TODO Format     string        `json:"format,omitempty"` // (optional) Allows validation based on a format pattern such as "date-time", "email", "uri", etc.
	Type       string `json:"type,omitempty"` // (optional) Assignment of JSON-based data types compatible with JSON Schema (one of boolean, integer, number, string, object, array, or null)
	DataSchema `json:"-"`
}

func (d *Data) Validate(value interface{}) error {
	if value == nil {
		return errors.New("missing value")
	}

	// TODO check for enum, format,...
	var err error
	switch d.DataSchema.(type) {
	case Boolean:
		err = d.DataSchema.(Boolean).Validate(value)
	case Integer:
		err = d.DataSchema.(Integer).Validate(value)
	case Number:
		err = d.DataSchema.(Number).Validate(value)
	case String:
		err = d.DataSchema.(String).Validate(value)
	case Object:
		err = d.DataSchema.(Object).Validate(value)
	case Array:
		err = d.DataSchema.(Array).Validate(value)
	}
	if err != nil {
		return err
	}
	return nil
}

// Ref: https://stackoverflow.com/questions/47335352/converting-struct-with-embedded-interface-into-json
// https://boldlygo.tech/posts/2020-06-26-go-json-tricks-embedded-marshaler/
func (d *Data) MarshalJSON() ([]byte, error) {
	type DataOrigin Data

	b1, err := json.Marshal((*DataOrigin)(d))
	if err != nil {
		return nil, err
	}

	b2, err := json.Marshal(d.DataSchema)
	if err != nil {
		return nil, err
	}
	if len(b2) > 2 { // '{}' is the empty object
		b3 := b1[:len(b1)-1] // remove last parenthesis
		b2[0] = ','          // replace first parenthesis, with a comma
		return append(b3, b2...), nil
	}
	return b1, nil // no DataSchema
}

func (d *Data) UnmarshalJSON(data []byte) error {
	type DataOrigin Data
	do := &struct {
		*DataOrigin
	}{
		DataOrigin: (*DataOrigin)(d),
	}
	if err := json.Unmarshal(data, do); err != nil {
		return err
	}

	switch d.Type {
	case "boolean":
		var ds = Boolean{}
		if err := json.Unmarshal(data, &ds); err != nil {
			return err
		}
		d.DataSchema = ds
	case "integer":
		var ds = Integer{}
		if err := json.Unmarshal(data, &ds); err != nil {
			return err
		}
		d.Default = int(do.Default.(float64))
		d.DataSchema = ds
	case "number":
		var ds = Number{}
		if err := json.Unmarshal(data, &ds); err != nil {
			return err
		}
		d.DataSchema = ds
	case "string":
		var ds = String{}
		if err := json.Unmarshal(data, &ds); err != nil {
			return err
		}
		d.DataSchema = ds
	case "object":
		var ds = Object{}
		if err := json.Unmarshal(data, &ds); err != nil {
			return err
		}
		d.DataSchema = ds
	case "array":
		var ds = Array{}
		if err := json.Unmarshal(data, &ds); err != nil {
			return err
		}
		d.DataSchema = ds
	}

	return nil
}
