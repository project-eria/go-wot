package dataSchema

import "errors"

type Boolean struct {
}

func NewBoolean(options ...BooleanOption) (Data, error) {
	opts := &BooleanOptions{
		Default: nil,
	}
	for _, option := range options {
		option(opts)
	}
	d := Data{
		Default:    opts.Default,
		Unit:       opts.Unit,
		Type:       "boolean",
		DataSchema: Boolean{},
	}
	if d.Default != nil {
		if err := d.Validate(d.Default); err != nil {
			return Data{}, errors.New("invalid default value: " + err.Error())
		}
	}
	return d, nil
}

type BooleanOption func(*BooleanOptions)

type BooleanOptions struct {
	Default interface{}
	Unit    string
}

func BooleanDefault(value bool) BooleanOption {
	return func(opts *BooleanOptions) {
		opts.Default = value
	}
}

func BooleanUnit(unit string) BooleanOption {
	return func(opts *BooleanOptions) {
		opts.Unit = unit
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
