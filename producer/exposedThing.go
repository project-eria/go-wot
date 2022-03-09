package producer

import (
	"fmt"
	"sync"

	"github.com/project-eria/go-wot/thing"
	"github.com/rs/zerolog/log"
)

// https://w3c.github.io/wot-scripting-api/#the-exposedthing-interface
type ExposedThing struct {
	Td                 *thing.Thing
	ExposedProperties  map[string]*ExposedProperty
	ExposedActions     map[string]*ExposedAction
	ExposedEvents      map[string]*ExposedEvent
	PropertyChangeChan chan PropertyChange
	EventChan          chan Event
	_wait              *sync.WaitGroup
}

func NewExposedThing(td *thing.Thing, wait *sync.WaitGroup) *ExposedThing {
	t := &ExposedThing{
		Td:                 td,
		ExposedProperties:  map[string]*ExposedProperty{},
		ExposedActions:     map[string]*ExposedAction{},
		ExposedEvents:      map[string]*ExposedEvent{},
		PropertyChangeChan: make(chan PropertyChange),
		EventChan:          make(chan Event),
		_wait:              wait,
	}

	for key, property := range td.Properties {
		t.ExposedProperties[key] = NewExposedProperty(property)
	}
	for key, action := range td.Actions {
		t.ExposedActions[key] = NewExposedAction(action)
	}
	for key, event := range td.Events {
		t.ExposedEvents[key] = NewExposedEvent(event)
	}
	return t
}

// https://www.w3.org/TR/wot-scripting-api/#the-getthingdescription-method-0
func (t *ExposedThing) GetThingDescription() *thing.Thing {
	return t.Td
}

// https://www.w3.org/TR/wot-scripting-api/#the-expose-method
func (t *ExposedThing) Expose() {
	// Todo
}

// https://www.w3.org/TR/wot-scripting-api/#the-destroy-method
func (t *ExposedThing) Destroy() {
	// Todo
	// Close channels
}

/*
 * Properties
 */
// https://www.w3.org/TR/wot-scripting-api/#the-setpropertyreadhandler-method
func (t *ExposedThing) SetPropertyReadHandler(name string, handler PropertyReadHandler) error {
	if _, ok := t.Td.Properties[name]; ok {
		t.ExposedProperties[name].SetReadHandler(handler)
		return nil
	}
	log.Debug().Str("property", name).Msg("[ExposedThing:SetPropertyReadHandler] property not found")
	return fmt.Errorf("property %s not found", name)
}

// https://www.w3.org/TR/wot-scripting-api/#the-setpropertyobservehandler-method
func (t *ExposedThing) SetPropertyObserveHandler(name string, handler PropertyObserveHandler) error {
	if _, ok := t.Td.Properties[name]; ok {
		if t.Td.Properties[name].Observable {
			t.ExposedProperties[name].SetObserveHandler(handler)
			return nil
		}
		log.Debug().Str("property", name).Msg("[ExposedThing:SetPropertyObserveHandler] property not observable")
		return fmt.Errorf("property %s not observable", name)
	}
	log.Debug().Str("property", name).Msg("[ExposedThing:SetPropertyObserveHandler] property not found")
	return fmt.Errorf("property %s not found", name)
}

// https://www.w3.org/TR/wot-scripting-api/#the-setpropertyunobservehandler-method
func (t *ExposedThing) SetPropertyUnobserveHandler(name string) error {
	if _, ok := t.Td.Properties[name]; ok {
		if t.Td.Properties[name].Observable {
			t.ExposedProperties[name].SetObserveHandler(nil)
			return nil
		}
		log.Debug().Str("property", name).Msg("[ExposedThing:SetPropertyUnobserveHandler] property not observable")
		return fmt.Errorf("property %s not observable", name)
	}
	log.Debug().Str("property", name).Msg("[ExposedThing:SetPropertyUnobserveHandler] property not found")
	return fmt.Errorf("property %s not found", name)
}

// https://w3c.github.io/wot-scripting-api/#the-emitpropertychange-method
func (t *ExposedThing) EmitPropertyChange(name string) error {
	if _, ok := t.Td.Properties[name]; ok {
		p := t.ExposedProperties[name]
		var value interface{}
		var err error
		if handler := p.GetObserveHandler(); handler != nil {
			if value, err = handler(t, name); err != nil {
				log.Debug().Str("property", name).Err(err).Msg("[ExposedThing:EmitPropertyChange] handler error for property")
				return err
			}
		} else if handler := p.GetReadHandler(); handler != nil {
			if value, err = handler(t, name); err != nil {
				log.Debug().Str("property", name).Err(err).Msg("[ExposedThing:EmitPropertyChange] handler error for property")
				return err
			}
		} else {
			// No handler
			log.Debug().Str("property", name).Msg("[ExposedThing:EmitPropertyChange] no handler available for property")
			return fmt.Errorf("no handler available for property %s", name)
		}
		t.PropertyChangeChan <- PropertyChange{name, value}
		return nil
	}
	log.Debug().Str("property", name).Msg("[ExposedThing:EmitPropertyChange] property not found")
	return fmt.Errorf("property %s not found", name)
}

// https://w3c.github.io/wot-scripting-api/#the-setpropertywritehandler-method
func (t *ExposedThing) SetPropertyWriteHandler(name string, handler PropertyWriteHandler) error {
	if _, ok := t.Td.Properties[name]; ok {
		t.ExposedProperties[name].SetWriteHandler(handler)
		return nil
	}
	log.Debug().Str("property", name).Msg("[ExposedThing:SetPropertyWriteHandler] property not found")
	return fmt.Errorf("property %s not found", name)
}

/*
 * Actions
 */
// https://w3c.github.io/wot-scripting-api/#the-setactionhandler-method
func (t *ExposedThing) SetActionHandler(name string, handler ActionHandler) error {
	if _, ok := t.Td.Actions[name]; ok {
		t.ExposedActions[name].SetHandler(handler)
		return nil
	}
	log.Debug().Str("action", name).Msg("[ExposedThing:SetActionHandler] action not found")
	return fmt.Errorf("action %s not found", name)
}

/*
 * Events
 */
// https://w3c.github.io/wot-scripting-api/#the-seteventsubscribehandler-method
func (t *ExposedThing) SetEventSubscribeHandler(name string, handler EventSubscriptionHandler) error {
	if _, ok := t.Td.Events[name]; ok {
		t.ExposedEvents[name].SetSubscribeHandler(handler)
		return nil
	}
	log.Debug().Str("event", name).Msg("[ExposedThing:SetEventSubscribeHandler] event not found")
	return fmt.Errorf("event %s not found", name)
}

// https://w3c.github.io/wot-scripting-api/#the-seteventunsubscribehandler-method
func (t *ExposedThing) SetEventUnsubscribeHandler(name string) error {
	if _, ok := t.Td.Events[name]; ok {
		t.ExposedEvents[name].SetUnSubscribeHandler()
		return nil
	}
	log.Debug().Str("event", name).Msg("[ExposedThing:SetEventUnsubscribeHandler] event not found")
	return fmt.Errorf("event %s not found", name)
}

// https://w3c.github.io/wot-scripting-api/#the-seteventhandler-method
func (t *ExposedThing) SetEventHandler(name string, handler EventListenerHandler) error {
	if _, ok := t.Td.Events[name]; ok {
		t.ExposedEvents[name].SetEventHandler(handler)
		return nil
	}
	log.Debug().Str("event", name).Msg("[ExposedThing:SetEventHandler] event not found")
	return fmt.Errorf("event %s not found", name)
}

// https://w3c.github.io/wot-scripting-api/#the-emitevent-method
func (t *ExposedThing) EmitEvent(name string) error {
	if _, ok := t.Td.Events[name]; ok {
		if handler := t.ExposedEvents[name].GetEventHandler(); handler != nil {
			var value interface{}
			var err error
			if value, err = handler(); err != nil {
				log.Debug().Str("event", name).Err(err).Msg("[ExposedThing:EmitEvent] handler error for event")
				return err
			}
			t.EventChan <- Event{name, value}
			return nil
		} else {
			// No handler
			log.Debug().Str("event", name).Msg("[ExposedThing:EmitEvent] no handler available for event")
			return fmt.Errorf("no handler available for event %s", name)
		}
	}
	log.Debug().Str("event", name).Msg("[ExposedThing:EmitEvent] event not found")
	return fmt.Errorf("event %s not found", name)
}
