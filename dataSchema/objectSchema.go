package dataSchema

import "errors"

type Object struct {
	/* TODO https://www.w3.org/TR/wot-thing-description11/#objectschema
	properties map[string]DataSchema
	required []string
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
