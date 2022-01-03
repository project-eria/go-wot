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

func (b Boolean) Check(value interface{}) error {
	if _, ok := value.(bool); !ok {
		return errors.New("incorrect boolean value type")
	}
	return nil
}
