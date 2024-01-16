package dataSchema

import (
	"errors"
	"strconv"
)

// "the field is empty and its tag specifies the "omitempty" option. The empty values are false, 0, any nil pointer or interface value, and any array, slice, map, or string of length zero. The object's default key string is the struct field name but can be specified in the struct field's tag value. The "json" key in the struct field's tag value is the key name, followed by an optional comma and options."
// TODO: Use a pointer for the fields, so that the zero value of the JSON type can be differentiated from the missing value.

type Integer struct {
	Minimum *int `json:"minimum,omitempty"` // (optional) Specifies a minimum numeric value, representing an inclusive lower limit. Only applicable for associated number or integer types.
	Maximum *int `json:"maximum,omitempty"` // (optional) Specifies a maximum numeric value, representing an inclusive upper limit. Only applicable for associated number or integer types.
	// TODO ExculsiveMinimum *int    `json:"exclusiveMinimum,omitempty"` // (optional) Specifies a minimum numeric value, representing an exclusive lower limit. Only applicable for associated number or integer types.
	// TODO ExclusiveMaximum *int    `json:"exclusiveMaximum,omitempty"` // (optional) Specifies a maximum numeric value, representing an exclusive upper limit. Only applicable for associated number or integer types.
	// TODO MultipleOf *uint16 `json:"multipleOf,omitempty"` // (optional) Specifies the multipleOf value number. The value must strictly greater than 0. Only applicable for associated number or integer types.
}

func NewInteger(options ...IntegerOption) (Data, error) {
	opts := &IntegerOptions{
		Default: nil,
	}
	for _, option := range options {
		if option != nil {
			option(opts)
		}
	}
	d := Data{
		Default: opts.Default,
		Type:    "integer",
		Unit:    opts.Unit,
		DataSchema: Integer{
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

type IntegerOption func(*IntegerOptions)

type IntegerOptions struct {
	Default interface{}
	Unit    string
	Minimum *int
	Maximum *int
}

func IntegerDefault(value int) IntegerOption {
	return func(opts *IntegerOptions) {
		opts.Default = value
	}
}

func IntegerUnit(unit string) IntegerOption {
	return func(opts *IntegerOptions) {
		opts.Unit = unit
	}
}

func IntegerMin(min int) IntegerOption {
	return func(opts *IntegerOptions) {
		opts.Minimum = &min
	}
}

func IntegerMax(max int) IntegerOption {
	return func(opts *IntegerOptions) {
		opts.Maximum = &max
	}
}

func (i Integer) Validate(value interface{}) error {
	if _, ok := value.(int); !ok {
		return errors.New("incorrect integer value type")
	}
	if i.Minimum != nil && value.(int) < *i.Minimum {
		return errors.New("value is less than minimum")
	}
	if i.Maximum != nil && value.(int) > *i.Maximum {
		return errors.New("value is greater than maximum")
	}
	// TODO if i.MultipleOf != nil && value.(int)%int(*i.MultipleOf) != 0 {
	// 	return errors.New("value is not multiple of multipleOf")
	// }
	return nil
}

func (i Integer) FromString(value string) (interface{}, error) {
	return strconv.Atoi(value)
}
