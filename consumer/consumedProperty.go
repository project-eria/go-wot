package consumer

import (
	"errors"
	"fmt"

	zlog "github.com/rs/zerolog/log"
)

// https://w3c.github.io/wot-scripting-api/#the-readproperty-method
func (t *consumedThing) ReadProperty(name string, dataVariables map[string]interface{}) (interface{}, error) {
	if property, ok := t.td.Properties[name]; ok {
		for _, form := range property.Forms {
			form := form // Copy https://go.dev/doc/faq#closures_and_goroutines
			for _, op := range form.Op {
				op := op // Copy https://go.dev/doc/faq#closures_and_goroutines
				if op == "readproperty" {
					if client := t.consumer.GetClientFor(form); client != nil {
						uriData, err := getUriVariables(property.UriVariables, dataVariables)
						if err != nil {
							return nil, err
						}
						value, uri, err := client.ReadResource(form, uriData)
						zlog.Trace().Str("property", name).Str("uri", uri).Interface("value", value).Msg("[consumer:ReadProperty] Received value")
						return value, err
					}
					return nil, errors.New("can't find client")
				}
			}
		}
		zlog.Error().Str("property", name).Msg("[consumer:ReadProperty] property not readable")
		return nil, fmt.Errorf("property %s not readable", name)
	}
	zlog.Error().Str("property", name).Msg("[consumer:ReadProperty] read property not found")
	return nil, fmt.Errorf("property %s not found", name)
}

// https://w3c.github.io/wot-scripting-api/#the-readmultipleproperties-method
func (t *consumedThing) ReadMultipleProperties() {
	// TODO
}

// https://w3c.github.io/wot-scripting-api/#the-readallproperties-method
func (t *consumedThing) ReadAllProperties() {
	// TODO
}

// https://w3c.github.io/wot-scripting-api/#the-writeproperty-method
func (t *consumedThing) WriteProperty(name string, dataVariables map[string]interface{}, value interface{}) (interface{}, error) {
	if property, ok := t.td.Properties[name]; ok {
		for _, form := range property.Forms {
			form := form // Copy https://go.dev/doc/faq#closures_and_goroutines
			for _, op := range form.Op {
				op := op // Copy https://go.dev/doc/faq#closures_and_goroutines
				if op == "writeproperty" {
					if client := t.consumer.GetClientFor(form); client != nil {
						uriData, err := getUriVariables(property.UriVariables, dataVariables)
						if err != nil {
							return nil, err
						}
						value, uri, err := client.WriteResource(form, uriData, value)
						zlog.Trace().Str("property", name).Str("uri", uri).Interface("value", value).Msg("[consumer:WriteProperty] value sent")
						return value, err
					}
					return nil, errors.New("can't find client")
				}
			}
		}
		zlog.Error().Str("property", name).Msg("[consumer:WriteProperty] property not writable")
		return nil, fmt.Errorf("property %s not writable", name)
	}
	zlog.Error().Str("property", name).Msg("[consumer:WriteProperty] write property not found")
	return nil, fmt.Errorf("property %s not found", name)
}

// https://w3c.github.io/wot-scripting-api/#the-writemultipleproperties-method
func (t *consumedThing) WriteMultipleProperties() {
	// TODO
}

type Listener func(value interface{}, err error)

// https://w3c.github.io/wot-scripting-api/#the-observeproperty-method
func (t *consumedThing) ObserveProperty(name string, dataVariables map[string]interface{}, listener Listener) error {
	if listener == nil {
		zlog.Error().Str("property", name).Msg("[consumer:ObserveProperty] missing listener")
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
						uriData, err := getUriVariables(property.UriVariables, dataVariables)
						if err != nil {
							return err
						}
						uri, err := client.SubscribeResource(form, uriData, sub, listener)
						zlog.Trace().Str("property", name).Str("uri", uri).Msg("[consumer:ObserveProperty] listener added")
						return err
					}
					return errors.New("can't find client")
				}
			}
		}
		zlog.Error().Str("property", name).Msg("[consumer:ObserveProperty] property not observable")
		return fmt.Errorf("property %s not observable", name)
	}
	zlog.Error().Str("property", name).Msg("[consumer:ObserveProperty] property not found")
	return fmt.Errorf("property %s not found", name)
}
