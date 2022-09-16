package dataSchema

import "errors"

type DataSchema interface {
	Check(interface{}) error
}

type Data struct {
	Const   interface{} `json:"const,omitempty"`   // (optional) Provides a constant value.
	Default interface{} `json:"default,omitempty"` // (optional) Supply a default value. The value should validate against the data schema in which it resides.
	Unit    string      `json:"unit,omitempty"`    // (optional) Provides unit information that is used, e.g., in international science, engineering, and business.
	// oneOf	Used to ensure that the data is valid against one of the specified schemas in the array.	optional	Array of DataSchema
	Enum             []interface{} `json:"enum,omitempty"`             // (optional) Restricted set of values provided as an array.
	ReadOnly         bool          `json:"readOnly"`                   // (default = false) Boolean value that is a hint to indicate whether a property interaction / value is read only.
	WriteOnly        bool          `json:"writeOnly"`                  // (default = false) Boolean value that is a hint to indicate whether a property interaction / value is write only.
	Format           string        `json:"format,omitempty"`           // (optional) Allows validation based on a format pattern such as "date-time", "email", "uri", etc.
	ContentEncoding  string        `json:"contentEncoding,omitempty"`  // (optional) Specifies the encoding used to store the contents, as specified in RFC 2054. The values that are accepted: "7bit", "8bit", "binary", "quoted-printable" and "base64".
	ContentMediaType string        `json:"contentMediaType,omitempty"` // (optional) Specifies the MIME type (e.g., image/png, audio/mpeg) of the contents of a string value, as described in RFC 2046.
	Type             string        `json:"type,omitempty"`             // (optional) Assignment of JSON-based data types compatible with JSON Schema (one of boolean, integer, number, string, object, array, or null)
	DataSchema       `json:"-"`
}

func (d *Data) Check(value interface{}) error {
	if value == nil {
		return errors.New("missing value")
	}
	switch d.DataSchema.(type) {
	case Boolean:
		err := d.DataSchema.(Boolean).Check(value)
		if err != nil {
			return err
		}
	case Integer:
		err := d.DataSchema.(Integer).Check(value)
		if err != nil {
			return err
		}
	case Number:
		err := d.DataSchema.(Number).Check(value)
		if err != nil {
			return err
		}
	case String:
		err := d.DataSchema.(String).Check(value)
		if err != nil {
			return err
		}
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
