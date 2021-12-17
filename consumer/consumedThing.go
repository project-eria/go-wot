package consumer

import (
	"fmt"

	"github.com/project-eria/go-wot/thing"
	"github.com/rs/zerolog/log"
)

// https://w3c.github.io/wot-scripting-api/#the-consumedthing-interface
type ConsumedThing struct {
	td *thing.Thing
}

// https://w3c.github.io/wot-scripting-api/#the-getthingdescription-method
func (c *ConsumedThing) GetThingDescription() *thing.Thing {
	return c.td
}

/*
 * Properties
 */

// https://w3c.github.io/wot-scripting-api/#the-readproperty-method
func (c *ConsumedThing) ReadProperty(name string) (interface{}, error) {
	if property, ok := c.td.Properties[name]; ok {
		for _, form := range property.Forms {
			for _, op := range form.Op {
				if op == "readproperty" {
					value, err := getHTTPJSON(form.Href)
					log.Trace().Str("property", name).Str("url", form.Href).Interface("value", value).Msg("[consumer:ReadProperty] Received value")
					return value, err
				}
			}
		}
	}
	log.Error().Str("property", name).Msg("[consumer:ReadProperty] read property not found")
	return nil, fmt.Errorf("property %s not found", name)
}

// https://w3c.github.io/wot-scripting-api/#the-readmultipleproperties-method
func (c *ConsumedThing) ReadMultipleProperties() {
	// TODO
}

// https://w3c.github.io/wot-scripting-api/#the-readallproperties-method
func (c *ConsumedThing) ReadAllProperties() {
	// TODO
}

// https://w3c.github.io/wot-scripting-api/#the-writeproperty-method
func (c *ConsumedThing) WriteProperty(name string, value interface{}) error {
	// TODO
	return nil
}

// https://w3c.github.io/wot-scripting-api/#the-writemultipleproperties-method
func (c *ConsumedThing) WriteMultipleProperties() {
	// TODO
}
