package consumer

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
)

// https://w3c.github.io/wot-scripting-api/#the-readproperty-method
func (t *ConsumedThing) ReadProperty(name string) (interface{}, error) {
	if property, ok := t.td.Properties[name]; ok {
		for _, form := range property.Forms {
			form := form // Copy https://go.dev/doc/faq#closures_and_goroutines
			for _, op := range form.Op {
				op := op // Copy https://go.dev/doc/faq#closures_and_goroutines
				if op == "readproperty" {
					if client := t.consumer.GetClientFor(form); client != nil {
						value, err := client.ReadResource(form)
						log.Trace().Str("property", name).Str("url", form.Href).Interface("value", value).Msg("[consumer:ReadProperty] Received value")
						return value, err
					}
					return nil, errors.New("can't find client")
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
func (t *ConsumedThing) ReadMultipleProperties() {
	// TODO
}

// https://w3c.github.io/wot-scripting-api/#the-readallproperties-method
func (t *ConsumedThing) ReadAllProperties() {
	// TODO
}

// https://w3c.github.io/wot-scripting-api/#the-writeproperty-method
func (t *ConsumedThing) WriteProperty(name string, value interface{}) (interface{}, error) {
	if property, ok := t.td.Properties[name]; ok {
		for _, form := range property.Forms {
			form := form // Copy https://go.dev/doc/faq#closures_and_goroutines
			for _, op := range form.Op {
				op := op // Copy https://go.dev/doc/faq#closures_and_goroutines
				if op == "writeproperty" {
					if client := t.consumer.GetClientFor(form); client != nil {
						value, err := client.WriteResource(form, value)
						log.Trace().Str("property", name).Str("url", form.Href).Interface("value", value).Msg("[consumer:WriteProperty] value sent")
						return value, err
					}
					return nil, errors.New("can't find client")
				}
			}
		}
		log.Error().Str("property", name).Msg("[consumer:WriteProperty] property not writable")
		return nil, fmt.Errorf("property %s not writable", name)
	}
	log.Error().Str("property", name).Msg("[consumer:WriteProperty] write property not found")
	return nil, fmt.Errorf("property %s not found", name)
}

// https://w3c.github.io/wot-scripting-api/#the-writemultipleproperties-method
func (t *ConsumedThing) WriteMultipleProperties() {
	// TODO
}

type Listener func(value interface{}, err error)

// https://w3c.github.io/wot-scripting-api/#the-observeproperty-method
func (t *ConsumedThing) ObserveProperty(name string, listener Listener) error {
	if listener == nil {
		log.Error().Str("property", name).Msg("[consumer:ObserveProperty] missing listener")
		return fmt.Errorf("missing listener for property %s", name)
	}
	if property, ok := t.td.Properties[name]; ok {
		for _, form := range property.Forms {
			form := form // Copy https://go.dev/doc/faq#closures_and_goroutines
			for _, op := range form.Op {
				op := op // Copy https://go.dev/doc/faq#closures_and_goroutines
				if op == "observeproperty" {
					if client := t.consumer.GetClientFor(form); client != nil {
						sub := &Subscription{
							Type: "property",
							Name: name,
						}
						err := client.SubscribeResource(form, sub, listener)
						log.Trace().Str("property", name).Str("url", form.Href).Msg("[consumer:ObserveProperty] listener added")
						return err
					}
					return errors.New("can't find client")
				}
			}
		}
		log.Error().Str("property", name).Msg("[consumer:ObserveProperty] property not observable")
		return fmt.Errorf("property %s not observable", name)
	}
	log.Error().Str("property", name).Msg("[consumer:ObserveProperty] property not found")
	return fmt.Errorf("property %s not found", name)
}
