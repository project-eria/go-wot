package producer

import (
	"errors"
	"fmt"
	"sync"

	"github.com/project-eria/go-wot/thing"
	"github.com/rs/zerolog/log"
)

// https://w3c.github.io/wot-scripting-api/#the-exposedthing-interface
type ExposedThing interface {
	GetThingDescription() *thing.Thing
	TD() *thing.Thing
	Expose()
	Destroy()
	Ref() string
	ExposedProperty(string) (ExposedProperty, error)
	ExposedAction(string) (ExposedAction, error)
	ExposedEvent(string) (ExposedEvent, error)
	// Properties
	SetPropertyReadHandler(string, PropertyReadHandler) error
	SetPropertyObserveHandler(string, PropertyObserveHandler) error
	SetObserverSelectorHandler(string, ObserverSelectorHandler) error
	SetPropertyUnobserveHandler(string) error
	EmitPropertyChange(string, interface{}, map[string]string) error
	SetPropertyWriteHandler(string, PropertyWriteHandler) error
	// Actions
	SetActionHandler(string, ActionHandler) error
	// Events
	SetEventSubscribeHandler(string, EventSubscriptionHandler) error
	SetEventUnsubscribeHandler(string) error
	SetEventHandler(string, EventListenerHandler) error
	EmitEvent(string, map[string]string) error
	GetPropertyChangeChannel() <-chan PropertyChange
	GetEventChannel() <-chan Event
}

type exposedThing struct {
	td                     *thing.Thing
	ref                    string
	exposedProperties      map[string]ExposedProperty
	exposedActions         map[string]ExposedAction
	exposedEvents          map[string]ExposedEvent
	propertyChangeChannels []chan PropertyChange
	eventChannels          []chan Event
	_wait                  *sync.WaitGroup
}

func NewExposedThing(td *thing.Thing, ref string, wait *sync.WaitGroup) ExposedThing {
	t := &exposedThing{
		td:                     td,
		ref:                    ref,
		exposedProperties:      map[string]ExposedProperty{},
		exposedActions:         map[string]ExposedAction{},
		exposedEvents:          map[string]ExposedEvent{},
		propertyChangeChannels: []chan PropertyChange{},
		eventChannels:          []chan Event{},
		_wait:                  wait,
	}

	for key, property := range td.Properties {
		property := property // Copy https://go.dev/doc/faq#closures_and_goroutines
		t.exposedProperties[key] = NewExposedProperty(property)
	}
	for key, action := range td.Actions {
		action := action // Copy https://go.dev/doc/faq#closures_and_goroutines
		t.exposedActions[key] = NewExposedAction(action)
	}
	for key, event := range td.Events {
		event := event // Copy https://go.dev/doc/faq#closures_and_goroutines
		t.exposedEvents[key] = NewExposedEvent(event)
	}
	return t
}

// https://www.w3.org/TR/wot-scripting-api/#the-getthingdescription-method-0
func (t *exposedThing) GetThingDescription() *thing.Thing {
	return t.td
}

func (t *exposedThing) TD() *thing.Thing {
	return t.td
}

// https://www.w3.org/TR/wot-scripting-api/#the-expose-method
func (t *exposedThing) Expose() {
	// Todo
}

// https://www.w3.org/TR/wot-scripting-api/#the-destroy-method
func (t *exposedThing) Destroy() {
	// Todo
	// Close channels
}

func (t *exposedThing) Ref() string {
	return t.ref
}

func (t *exposedThing) ExposedProperty(key string) (ExposedProperty, error) {
	if property, ok := t.exposedProperties[key]; ok {
		return property, nil
	}
	return nil, errors.New("exposed property not found")
}

func (t *exposedThing) ExposedAction(key string) (ExposedAction, error) {
	if action, ok := t.exposedActions[key]; ok {
		return action, nil
	}
	return nil, errors.New("exposed action not found")
}

func (t *exposedThing) ExposedEvent(key string) (ExposedEvent, error) {
	if event, ok := t.exposedEvents[key]; ok {
		return event, nil
	}
	return nil, errors.New("exposed event not found")
}

/*
 * Properties
 */
// https://www.w3.org/TR/wot-scripting-api/#the-setpropertyreadhandler-method
func (t *exposedThing) SetPropertyReadHandler(name string, handler PropertyReadHandler) error {
	if _, ok := t.td.Properties[name]; ok {
		t.exposedProperties[name].SetReadHandler(handler)
		return nil
	}
	log.Trace().Str("property", name).Msg("[ExposedThing:SetPropertyReadHandler] property not found")
	return fmt.Errorf("property %s not found", name)
}

// https://www.w3.org/TR/wot-scripting-api/#the-setpropertyobservehandler-method
func (t *exposedThing) SetPropertyObserveHandler(name string, handler PropertyObserveHandler) error {
	if _, ok := t.td.Properties[name]; ok {
		if t.td.Properties[name].Observable {
			t.exposedProperties[name].SetObserveHandler(handler)
			return nil
		}
		log.Trace().Str("property", name).Msg("[ExposedThing:SetPropertyObserveHandler] property not observable")
		return fmt.Errorf("property %s not observable", name)
	}
	log.Trace().Str("property", name).Msg("[ExposedThing:SetPropertyObserveHandler] property not found")
	return fmt.Errorf("property %s not found", name)
}

func (t *exposedThing) SetObserverSelectorHandler(name string, handler ObserverSelectorHandler) error {
	if _, ok := t.td.Properties[name]; ok {
		if t.td.Properties[name].Observable {
			t.exposedProperties[name].SetObserverSelectorHandler(handler)
			return nil
		}
		log.Trace().Str("property", name).Msg("[ExposedThing:SetObserverSelectorHandler] property not observable")
		return fmt.Errorf("property %s not observable", name)
	}
	log.Trace().Str("property", name).Msg("[ExposedThing:SetObserverSelectorHandler] property not found")
	return fmt.Errorf("property %s not found", name)
}

// https://www.w3.org/TR/wot-scripting-api/#the-setpropertyunobservehandler-method
func (t *exposedThing) SetPropertyUnobserveHandler(name string) error {
	if _, ok := t.td.Properties[name]; ok {
		if t.td.Properties[name].Observable {
			t.exposedProperties[name].SetObserveHandler(nil)
			return nil
		}
		log.Trace().Str("property", name).Msg("[ExposedThing:SetPropertyUnobserveHandler] property not observable")
		return fmt.Errorf("property %s not observable", name)
	}
	log.Trace().Str("property", name).Msg("[ExposedThing:SetPropertyUnobserveHandler] property not found")
	return fmt.Errorf("property %s not found", name)
}

// https://w3c.github.io/wot-scripting-api/#the-emitpropertychange-method
func (t *exposedThing) EmitPropertyChange(name string, data interface{}, options map[string]string) error {
	if _, ok := t.td.Properties[name]; ok {
		p := t.exposedProperties[name]
		var value interface{}
		var err error
		if data != nil {
			value = data
		} else if handler := p.GetReadHandler(); handler != nil {
			if value, err = handler(t, name, options); err != nil {
				log.Error().Str("ThingRef", t.ref).Str("property", name).Err(err).Msg("[ExposedThing:EmitPropertyChange] read handler error for property")
				return err
			}
		} else {
			// No handler
			log.Trace().Str("ThingRef", t.ref).Str("property", name).Msg("[ExposedThing:EmitPropertyChange] no handler available for property")
			return fmt.Errorf("no handler available for property %s", name)
		}
		// Send the notification to all protocols, that requested a channel
		for _, c := range t.propertyChangeChannels {
			go func(c chan PropertyChange) {
				select {
				case c <- PropertyChange{ThingRef: t.ref, Name: name, Value: value, Handler: p.GetObserverSelectorHandler(), Options: options}:
					return
				default:
					log.Error().Msg("[ExposedThing:EmitPropertyChange] channel blocked (no reader?), can not write")
					return
				}
			}(c)
		}
		return nil
	}
	log.Trace().Str("property", name).Msg("[ExposedThing:EmitPropertyChange] property not found")
	return fmt.Errorf("property %s not found", name)
}

// https://w3c.github.io/wot-scripting-api/#the-setpropertywritehandler-method
func (t *exposedThing) SetPropertyWriteHandler(name string, handler PropertyWriteHandler) error {
	if _, ok := t.td.Properties[name]; ok {
		t.exposedProperties[name].SetWriteHandler(handler)
		return nil
	}
	log.Trace().Str("property", name).Msg("[ExposedThing:SetPropertyWriteHandler] property not found")
	return fmt.Errorf("property %s not found", name)
}

/*
 * Actions
 */
// https://w3c.github.io/wot-scripting-api/#the-setactionhandler-method
func (t *exposedThing) SetActionHandler(name string, handler ActionHandler) error {
	if _, ok := t.td.Actions[name]; ok {
		t.exposedActions[name].SetHandler(handler)
		return nil
	}
	log.Trace().Str("action", name).Msg("[ExposedThing:SetActionHandler] action not found")
	return fmt.Errorf("action %s not found", name)
}

/*
 * Events
 */
// https://w3c.github.io/wot-scripting-api/#the-seteventsubscribehandler-method
func (t *exposedThing) SetEventSubscribeHandler(name string, handler EventSubscriptionHandler) error {
	if _, ok := t.td.Events[name]; ok {
		t.exposedEvents[name].SetSubscribeHandler(handler)
		return nil
	}
	log.Trace().Str("event", name).Msg("[ExposedThing:SetEventSubscribeHandler] event not found")
	return fmt.Errorf("event %s not found", name)
}

// https://w3c.github.io/wot-scripting-api/#the-seteventunsubscribehandler-method
func (t *exposedThing) SetEventUnsubscribeHandler(name string) error {
	if _, ok := t.td.Events[name]; ok {
		t.exposedEvents[name].SetUnSubscribeHandler()
		return nil
	}
	log.Trace().Str("event", name).Msg("[ExposedThing:SetEventUnsubscribeHandler] event not found")
	return fmt.Errorf("event %s not found", name)
}

// https://w3c.github.io/wot-scripting-api/#the-seteventhandler-method
func (t *exposedThing) SetEventHandler(name string, handler EventListenerHandler) error {
	if _, ok := t.td.Events[name]; ok {
		t.exposedEvents[name].SetEventHandler(handler)
		return nil
	}
	log.Trace().Str("event", name).Msg("[ExposedThing:SetEventHandler] event not found")
	return fmt.Errorf("event %s not found", name)
}

// https://w3c.github.io/wot-scripting-api/#the-emitevent-method
func (t *exposedThing) EmitEvent(name string, options map[string]string) error {
	if _, ok := t.td.Events[name]; ok {
		e := t.exposedEvents[name]
		if handler := e.GetEventHandler(); handler != nil {
			var value interface{}
			var err error
			if value, err = handler(); err != nil {
				log.Trace().Str("ThingRef", t.ref).Str("event", name).Err(err).Msg("[ExposedThing:EmitEvent] handler error for event")
				return err
			}
			// Send the notification to all protocols, that requested a channel
			for _, c := range t.eventChannels {
				go func(c chan Event) {
					select {
					case c <- Event{ThingRef: t.ref, Name: name, Value: value, Handler: e.GetListenerSelectorHandler(), Options: options}:
						return
					default:
						log.Error().Msg("[ExposedThing:EmitEvente] channel blocked (no reader?), can not write")
						return
					}
				}(c)
			}
			return nil
		} else {
			// No handler
			log.Trace().Str("ThingRef", t.ref).Str("event", name).Msg("[ExposedThing:EmitEvent] no handler available for event")
			return fmt.Errorf("no handler available for event %s", name)
		}
	}
	log.Trace().Str("ThingRef", t.ref).Str("event", name).Msg("[ExposedThing:EmitEvent] event not found")
	return fmt.Errorf("event %s not found", name)
}

func (t *exposedThing) GetPropertyChangeChannel() <-chan PropertyChange {
	// Limit channel size to the number of property
	size := len(t.exposedProperties)
	c := make(chan PropertyChange, size)
	t.propertyChangeChannels = append(t.propertyChangeChannels, c)
	return c
}

func (t *exposedThing) GetEventChannel() <-chan Event {
	// Limit channel size to the number of property
	size := len(t.exposedEvents)
	c := make(chan Event, size)
	t.eventChannels = append(t.eventChannels, c)
	return c
}
