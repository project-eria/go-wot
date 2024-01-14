package dataSchema

import "errors"

type Boolean struct {
}

func NewBoolean(defaultValue bool) (Data, error) {
	d := Data{
		Default:    defaultValue,
		Type:       "boolean",
		DataSchema: Boolean{},
	}
	if err := d.Validate(d.Default); err != nil {
		return Data{}, errors.New("invalid default value: " + err.Error())
	}
	return d, nil
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
