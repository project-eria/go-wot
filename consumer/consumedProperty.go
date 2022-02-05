package consumer

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

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
		log.Error().Str("property", name).Msg("[consumer:ReadProperty] property not readable")
		return nil, fmt.Errorf("property %s not readable", name)
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
	if property, ok := c.td.Properties[name]; ok {
		for _, form := range property.Forms {
			for _, op := range form.Op {
				if op == "writeproperty" {
					_, err := putHTTPJSON(form.Href, value)
					log.Trace().Str("property", name).Str("url", form.Href).Interface("value", value).Msg("[consumer:WriteProperty] value sent")
					return err
				}
			}
		}
		log.Error().Str("property", name).Msg("[consumer:WriteProperty] property not writable")
		return fmt.Errorf("property %s not writable", name)
	}
	log.Error().Str("property", name).Msg("[consumer:WriteProperty] write property not found")
	return fmt.Errorf("property %s not found", name)
}

// https://w3c.github.io/wot-scripting-api/#the-writemultipleproperties-method
func (c *ConsumedThing) WriteMultipleProperties() {
	// TODO
}

type Listener func(value interface{}, err error)

// https://w3c.github.io/wot-scripting-api/#the-observeproperty-method
func (c *ConsumedThing) ObserveProperty(name string, listener Listener) error {
	if listener == nil {
		log.Error().Str("property", name).Msg("[consumer:ObserveProperty] missing listener")
		return fmt.Errorf("missing listener for property %s", name)
	}
	if property, ok := c.td.Properties[name]; ok {
		for _, form := range property.Forms {
			for _, op := range form.Op {
				if op == "observeproperty" {
					sub := &subscription{
						Type: "property",
						Name: name,
					}
					c._wait.Add(1)
					go func() {
						connectWebSocket(form.Href, sub, listener, c._ctx)
						c._wait.Done()
					}()
					log.Trace().Str("property", name).Str("url", form.Href).Msg("[consumer:ObserveProperty] listener added")
					return nil
				}
			}
		}
		log.Error().Str("property", name).Msg("[consumer:ObserveProperty] property not observable")
		return fmt.Errorf("property %s not observable", name)
	}
	log.Error().Str("property", name).Msg("[consumer:ObserveProperty] write property not found")
	return fmt.Errorf("property %s not found", name)
}
