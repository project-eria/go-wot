package dataSchema

import (
	"errors"
	"strconv"
)

type Number struct {
	Minimum *int `json:"minimum,omitempty"` // (optional) Specifies a minimum numeric value, representing an inclusive lower limit. Only applicable for associated number or integer types.
	Maximum *int `json:"maximum,omitempty"` // (optional) Specifies a maximum numeric value, representing an inclusive upper limit. Only applicable for associated number or integer types.
	// TODO ExculsiveMinimum *int16  `json:"exclusiveMinimum,omitempty"` // (optional) Specifies a minimum numeric value, representing an exclusive lower limit. Only applicable for associated number or integer types.
	// TODO ExclusiveMaximum *int16  `json:"exclusiveMaximum,omitempty"` // (optional) Specifies a maximum numeric value, representing an exclusive upper limit. Only applicable for associated number or integer types.
	// TODO MultipleOf       *uint16 `json:"multipleOf,omitempty"`       // (optional) Specifies the multipleOf value number. The value must strictly greater than 0. Only applicable for associated number or integer types.
}

func NewNumber(defaultValue float64, unit string, minimum *int, maximum *int) (Data, error) {
	d := Data{
		Default: defaultValue,
		Type:    "number",
		Unit:    unit,
		DataSchema: Number{
			Minimum: minimum,
			Maximum: maximum,
		},
	}
	if err := d.Validate(d.Default); err != nil {
		return Data{}, errors.New("invalid default value: " + err.Error())
	}
	return d, nil
}

func (n Number) Validate(value interface{}) error {
	if _, ok := value.(float64); !ok {
		return errors.New("incorrect number value type")
	}
	if n.Minimum != nil && value.(float64) < float64(*n.Minimum) {
		return errors.New("value is less than minimum")
	}
	if n.Maximum != nil && value.(float64) > float64(*n.Maximum) {
		return errors.New("value is greater than maximum")
	}
	return nil
}

func (n Number) FromString(value string) (interface{}, error) {
	return strconv.ParseFloat(value, 64)
}
