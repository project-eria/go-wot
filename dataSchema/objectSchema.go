package dataSchema

import "errors"

type Object struct {
	/* TODO https://www.w3.org/TR/wot-thing-description11/#objectschema
	// properties 	Data schema nested definitions. 	optional 	Map of DataSchema
	// required 	Defines which members of the object type are mandatory, i.e. which members are mandatory in the payload that is to be sent (e.g. input of invokeaction, writeproperty) and what members will be definitely delivered in the payload that is being received (e.g. output of invokeaction, readproperty) 	optional
	*/
}

func NewObject(options ...ObjectOption) (Data, error) {
	opts := &ObjectOptions{
		Default: nil,
	}
	for _, option := range options {
		option(opts)
	}
	d := Data{
		Default:    opts.Default,
		Unit:       opts.Unit,
		Type:       "object",
		DataSchema: Object{},
	}
	if d.Default != nil {
		if err := d.Validate(d.Default); err != nil {
			return Data{}, errors.New("invalid default value: " + err.Error())
		}
	}
	return d, nil
}

type ObjectOption func(*ObjectOptions)

type ObjectOptions struct {
	Default interface{}
	Unit    string
}

func ObjectDefault(value map[string]interface{}) ObjectOption {
	return func(opts *ObjectOptions) {
		opts.Default = value
	}
}

func ObjectUnit(unit string) ObjectOption {
	return func(opts *ObjectOptions) {
		opts.Unit = unit
	}
}

func (b Object) Validate(value interface{}) error {
	if _, ok := value.(map[string]interface{}); !ok {
		return errors.New("incorrect object value type")
	}
	return nil
}

func (b Object) FromString(value string) (interface{}, error) {
	return nil, errors.New("not implemented")
}
