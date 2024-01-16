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

func NewNumber(options ...NumberOption) (Data, error) {
	opts := &NumberOptions{
		Default: nil,
	}
	for _, option := range options {
		option(opts)
	}
	d := Data{
		Default: opts.Default,
		Type:    "number",
		Unit:    opts.Unit,
		DataSchema: Number{
			Minimum: opts.Minimum,
			Maximum: opts.Maximum,
		},
	}
	if d.Default != nil {
		if err := d.Validate(d.Default); err != nil {
			return Data{}, errors.New("invalid default value: " + err.Error())
		}
	}
	return d, nil
}

type NumberOption func(*NumberOptions)

type NumberOptions struct {
	Default interface{}
	Unit    string
	Minimum *int
	Maximum *int
}

func NumberDefault(value float64) NumberOption {
	return func(opts *NumberOptions) {
		opts.Default = value
	}
}

func NumberUnit(unit string) NumberOption {
	return func(opts *NumberOptions) {
		opts.Unit = unit
	}
}

func NumberMin(min int) NumberOption {
	return func(opts *NumberOptions) {
		opts.Minimum = &min
	}
}

func NumberMax(max int) NumberOption {
	return func(opts *NumberOptions) {
		opts.Maximum = &max
	}
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
