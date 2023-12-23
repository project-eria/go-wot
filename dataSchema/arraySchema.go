package dataSchema

import (
	"errors"
	"reflect"
)

type Array struct {
	/* TODO https://www.w3.org/TR/wot-thing-description11/#arrayschema*/

	// Items    DataSchema `json:"items,omitempty"`    // Array of DataSchema // (optional) Used to define the characteristics of an array.
	MinItems uint `json:"minItems,omitempty"` // (optional) Defines the minimum number of items that have to be in the array.
	MaxItems uint `json:"maxItems,omitempty"` // (optional) Defines the maximum number of items that have to be in the array.
}

func NewArray[T SimpleType](defaultValue []T, minItems uint, maxItems uint) Data {
	return Data{
		Default: defaultValue,
		Type:    "array",
		DataSchema: Array{
			MinItems: minItems,
			MaxItems: maxItems,
		},
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
