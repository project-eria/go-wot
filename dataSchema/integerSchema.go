package dataSchema

import (
	"errors"
	"strconv"
)

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
