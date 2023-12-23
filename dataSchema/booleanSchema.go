package dataSchema

import "errors"

type Boolean struct {
}

func NewBoolean(defaultValue bool) Data {
	return Data{
		Default:    defaultValue,
		Type:       "boolean",
		DataSchema: Boolean{},
	}
}

func (b Boolean) Validate(value interface{}) error {
	if _, ok := value.(bool); !ok {
		return errors.New("incorrect boolean value type")
	}
	return nil
}

func (b Boolean) FromString(value string) (interface{}, error) {
	if value == "true" {
		return true, nil
	} else if value == "false" {
		return false, nil
	}
	return nil, errors.New("incorrect boolean value type")
}
