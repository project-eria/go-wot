package dataSchema

import "errors"

type Number struct {
	Minimum          uint16 `json:"minimum,omitempty"`          // (optional) Specifies a minimum numeric value, representing an inclusive lower limit. Only applicable for associated number or integer types.
	Maximum          uint16 `json:"maximum,omitempty"`          // (optional) Specifies a maximum numeric value, representing an inclusive upper limit. Only applicable for associated number or integer types.
	ExculsiveMinimum uint16 `json:"exclusiveMinimum,omitempty"` // (optional) Specifies a minimum numeric value, representing an exclusive lower limit. Only applicable for associated number or integer types.
	ExclusiveMaximum uint16 `json:"exclusiveMaximum,omitempty"` // (optional) Specifies a maximum numeric value, representing an exclusive upper limit. Only applicable for associated number or integer types.
	MultipleOf       uint16 `json:"multipleOf,omitempty"`       // (optional) Specifies the multipleOf value number. The value must strictly greater than 0. Only applicable for associated number or integer types.
}

func NewNumber(defaultValue float64, unit string, minimum uint16, maximum uint16) Data {
	return Data{
		Default: defaultValue,
		Type:    "number",
		Unit:    unit,
		DataSchema: Number{
			Minimum: minimum,
			Maximum: maximum,
		},
	}
}

func (n Number) Check(value interface{}) error {
	if _, ok := value.(float64); !ok {
		return errors.New("incorrect number value type")
	}
	return nil
}
