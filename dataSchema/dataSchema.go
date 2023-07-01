package dataSchema

import "errors"

type SimpleType interface {
	bool | string | int8 | int16 | int32 | int64 | float32 | float64
}

type DataSchema interface {
	Check(interface{}) error
}

type Data struct {
	Title string `json:"title,omitempty"` // (optional) Provides a human-readable title (e.g., display a text for UI representation) based on a default language.
	// titles 	Provides multi-language human-readable titles (e.g., display a text for UI representation in different languages). Also see MultiLanguage. 	optional 	Map of MultiLanguage
	Description string `json:"description,omitempty"` // (optional)	Provides additional (human-readable) information based on a default language.
	// descriptions 	Can be used to support (human-readable) information in different languages. Also see MultiLanguage. 	optional 	Map of MultiLanguage

	Const   interface{} `json:"const,omitempty"`   // (optional) Provides a constant value.
	Default interface{} `json:"default,omitempty"` // (optional) Supply a default value. The value should validate against the data schema in which it resides.
	Unit    string      `json:"unit,omitempty"`    // (optional) Provides unit information that is used, e.g., in international science, engineering, and business.
	// oneOf	Used to ensure that the data is valid against one of the specified schemas in the array.	optional	Array of DataSchema
	Enum       []interface{} `json:"enum,omitempty"`   // (optional) Restricted set of values provided as an array.
	Format     string        `json:"format,omitempty"` // (optional) Allows validation based on a format pattern such as "date-time", "email", "uri", etc.
	Type       string        `json:"type,omitempty"`   // (optional) Assignment of JSON-based data types compatible with JSON Schema (one of boolean, integer, number, string, object, array, or null)
	DataSchema `json:"-"`
}

func (d *Data) Check(value interface{}) error {
	if value == nil {
		return errors.New("missing value")
	}

	// TODO check for enum, format,...
	var err error
	switch d.DataSchema.(type) {
	case Boolean:
		err = d.DataSchema.(Boolean).Check(value)
	case Integer:
		err = d.DataSchema.(Integer).Check(value)
	case Number:
		err = d.DataSchema.(Number).Check(value)
	case String:
		err = d.DataSchema.(String).Check(value)
	case Object:
		err = d.DataSchema.(Object).Check(value)
	case Array:
		err = d.DataSchema.(Array).Check(value)
	}
	if err != nil {
		return err
	}
	return nil
}

// Ref: https://stackoverflow.com/questions/47335352/converting-struct-with-embedded-interface-into-json
// func (p *Property) MarshalJSON() ([]byte, error) {
// 	type PropertyOrigin Property

// 	b1, err := json.Marshal((*PropertyOrigin)(p))
// 	if err != nil {
// 		return nil, err
// 	}

// 	b2, err := json.Marshal(p.DataSchema)
// 	if err != nil {
// 		return nil, err
// 	}

// 	s1 := string(b1[:len(b1)-1])
// 	s2 := string(b2[1:])

// 	return []byte(s1 + ", " + s2), nil
// }
