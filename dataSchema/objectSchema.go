package dataSchema

import "errors"

type Object struct {
	/* TODO https://www.w3.org/TR/wot-thing-description11/#objectschema
	// properties 	Data schema nested definitions. 	optional 	Map of DataSchema
	// required 	Defines which members of the object type are mandatory, i.e. which members are mandatory in the payload that is to be sent (e.g. input of invokeaction, writeproperty) and what members will be definitely delivered in the payload that is being received (e.g. output of invokeaction, readproperty) 	optional
	*/
}

func NewObject(defaultValue map[string]interface{}) Data {
	return Data{
		Default:    defaultValue,
		Type:       "object",
		DataSchema: Object{},
	}
}

func (b Object) Check(value interface{}) error {
	if _, ok := value.(map[string]interface{}); !ok {
		return errors.New("incorrect object value type")
	}
	return nil
}
