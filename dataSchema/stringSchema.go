package dataSchema

import "errors"

type String struct {
	MinLength uint16 `json:"minLength,omitempty"` // (optional) Specifies the minimum length of a string.
	MaxLength uint16 `json:"maxLength,omitempty"` // (optional) Specifies the maximum length of a string.package dataSchema
	Pattern   string `json:"pattern,omitempty"`   // (optional) Provides a regular expressions to express constraints of the string value. The regular expression must follow the [ECMA-262] dialect.
}

func NewString(defaultValue string, minLength uint16, maxLength uint16, pattern string) Data {
	return Data{
		Default: defaultValue,
		Type:    "string",
		DataSchema: String{
			MinLength: minLength,
			MaxLength: maxLength,
			Pattern:   pattern,
		},
	}
}

func (s String) Check(value interface{}) error {
	if _, ok := value.(string); !ok {
		return errors.New("incorrect string value type")
	}
	return nil
}
