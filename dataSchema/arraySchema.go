package dataSchema

import (
	"errors"
	"reflect"
)

type Array struct {
	/* TODO https://www.w3.org/TR/wot-thing-description11/#arrayschema*/

	// Items    DataSchema `json:"items,omitempty"`    // Array of DataSchema // (optional) Used to define the characteristics of an array.
	MinItems *uint `json:"minItems,omitempty"` // (optional) Defines the minimum number of items that have to be in the array.
	MaxItems *uint `json:"maxItems,omitempty"` // (optional) Defines the maximum number of items that have to be in the array.
}

func NewArray[T SimpleType](options ...ArrayOption) (Data, error) {
	opts := &ArrayOptions{
		Default: nil,
	}
	for _, option := range options {
		option(opts)
	}

	d := Data{
		Default: opts.Default,
		Unit:    opts.Unit,
		Type:    "array",
		DataSchema: Array{
			MinItems: opts.MinItems,
			MaxItems: opts.MaxItems,
		},
	}
	if d.Default != nil {
		if err := d.Validate(d.Default); err != nil {
			return Data{}, errors.New("invalid default value: " + err.Error())
		}
	}
	return d, nil
}

type ArrayOption func(*ArrayOptions)

type ArrayOptions struct {
	Default  interface{}
	Unit     string
	MinItems *uint
	MaxItems *uint
}

func ArrayDefault[T SimpleType](value []T) ArrayOption {
	return func(opts *ArrayOptions) {
		opts.Default = value
	}
}

func ArrayUnit(unit string) ArrayOption {
	return func(opts *ArrayOptions) {
		opts.Unit = unit
	}
}

func (b Array) Validate(value interface{}) error {
	typ := reflect.TypeOf(value)
	if typ.Kind() != reflect.Slice {
		return errors.New("incorrect array value type")
	}
	return nil
}

func (b Array) FromString(value string) (interface{}, error) {
	return nil, errors.New("not implemented")
}
