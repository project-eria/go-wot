package dataSchema

import "errors"

type Boolean struct {
}

func NewBoolean(readOnly bool) Data {
	return Data{
		Default:    false,
		Type:       "boolean",
		ReadOnly:   readOnly,
		DataSchema: Boolean{},
	}
}

func (b Boolean) Check(value interface{}) error {
	if _, ok := value.(bool); !ok {
		return errors.New("incorrect boolean value type")
	}
	return nil
}
