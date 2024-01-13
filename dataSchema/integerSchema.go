package dataSchema

import (
	"errors"
	"strconv"
)

// "the field is empty and its tag specifies the "omitempty" option. The empty values are false, 0, any nil pointer or interface value, and any array, slice, map, or string of length zero. The object's default key string is the struct field name but can be specified in the struct field's tag value. The "json" key in the struct field's tag value is the key name, followed by an optional comma and options."
// TODO: Use a pointer for the fields, so that the zero value of the JSON type can be differentiated from the missing value.

type Integer struct {
	Minimum          int    `json:"minimum,omitempty"`          // (optional) Specifies a minimum numeric value, representing an inclusive lower limit. Only applicable for associated number or integer types.
	Maximum          int    `json:"maximum,omitempty"`          // (optional) Specifies a maximum numeric value, representing an inclusive upper limit. Only applicable for associated number or integer types.
	ExculsiveMinimum int    `json:"exclusiveMinimum,omitempty"` // (optional) Specifies a minimum numeric value, representing an exclusive lower limit. Only applicable for associated number or integer types.
	ExclusiveMaximum int    `json:"exclusiveMaximum,omitempty"` // (optional) Specifies a maximum numeric value, representing an exclusive upper limit. Only applicable for associated number or integer types.
	MultipleOf       uint16 `json:"multipleOf,omitempty"`       // (optional) Specifies the multipleOf value number. The value must strictly greater than 0. Only applicable for associated number or integer types.
}

func NewInteger(defaultValue int, unit string, minimum int, maximum int) Data {
	return Data{
		Default: defaultValue,
		Type:    "integer",
		Unit:    unit,
		DataSchema: Integer{
			Minimum: minimum,
			Maximum: maximum,
		},
	}
}

func (i Integer) Validate(value interface{}) error {
	if _, ok := value.(int); !ok {
		return errors.New("incorrect integer value type")
	}
	return nil
}

func (i Integer) FromString(value string) (interface{}, error) {
	return strconv.Atoi(value)
}
