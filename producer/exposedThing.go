package producer

import (
	"fmt"

	"github.com/project-eria/go-wot/thing"
	"github.com/rs/zerolog/log"
)

// https://w3c.github.io/wot-scripting-api/#the-exposedthing-interface
type ExposedThing struct {
	Td                        *thing.Thing
	propertiesReadHandlers    map[string]PropertyReadHandler
	propertiesWriteHandlers   map[string]PropertyWriteHandler
	propertiesObserveHandlers map[string]PropertyObserveHandler
	actionHandlers            map[string]ActionHandler
	eventListenerHandlers     map[string]EventListenerHandler
	eventSubscriptionHandlers map[string]EventSubscriptionHandler
}

func NewExposedThing(td *thing.Thing) *ExposedThing {
	return &ExposedThing{
		Td:                        td,
		propertiesReadHandlers:    map[string]PropertyReadHandler{},
		propertiesWriteHandlers:   map[string]PropertyWriteHandler{},
		actionHandlers:            map[string]ActionHandler{},
		eventListenerHandlers:     map[string]EventListenerHandler{},
		eventSubscriptionHandlers: map[string]EventSubscriptionHandler{},
	}
}

// https://www.w3.org/TR/wot-scripting-api/#the-getthingdescription-method-0
func (e *ExposedThing) GetThingDescription() *thing.Thing {
	return e.Td
}

// https://www.w3.org/TR/wot-scripting-api/#the-expose-method
func (e *ExposedThing) Expose() {
	// Todo
}

// https://www.w3.org/TR/wot-scripting-api/#the-destroy-method
func (e *ExposedThing) Destroy() {
	// Todo
}

/*
 * Properties
 */
// https://w3c.github.io/wot-scripting-api/#the-propertyreadhandler-callback
type PropertyReadHandler func() (interface{}, error)
type PropertyObserveHandler func() (interface{}, error)

//https://w3c.github.io/wot-scripting-api/#the-propertywritehandler-callback
type PropertyWriteHandler func(interface{}) error

// https://www.w3.org/TR/wot-scripting-api/#the-setpropertyreadhandler-method
func (e *ExposedThing) SetPropertyReadHandler(name string, handler PropertyReadHandler) error {
	if _, ok := e.Td.Properties[name]; ok {
		e.propertiesReadHandlers[name] = handler
		return nil
	}
	log.Debug().Str("property", name).Msg("[exposedThing:SetPropertyReadHandler] property not found")
	return fmt.Errorf("property %s not found", name)
}

// https://www.w3.org/TR/wot-scripting-api/#the-setpropertyobservehandler-method
func (e *ExposedThing) SetPropertyObserveHandler(name string, handler PropertyObserveHandler) error {
	if _, ok := e.Td.Properties[name]; ok {
		if e.Td.Properties[name].Observable {
			e.propertiesObserveHandlers[name] = handler
			return nil
		}
		log.Debug().Str("property", name).Msg("[exposedThing:SetPropertyObserveHandler] property not observable")
		return fmt.Errorf("property %s not observable", name)
	}
	log.Debug().Str("property", name).Msg("[exposedThing:SetPropertyObserveHandler] property not found")
	return fmt.Errorf("property %s not found", name)
}

// https://www.w3.org/TR/wot-scripting-api/#the-setpropertyunobservehandler-method
func (e *ExposedThing) SetPropertyUnobserveHandler(name string) error {
	if _, ok := e.Td.Properties[name]; ok {
		if e.Td.Properties[name].Observable {
			delete(e.propertiesObserveHandlers, name)
			return nil
		}
		log.Debug().Str("property", name).Msg("[exposedThing:SetPropertyUnobserveHandler] property not observable")
		return fmt.Errorf("property %s not observable", name)
	}
	log.Debug().Str("property", name).Msg("[exposedThing:SetPropertyUnobserveHandler] property not found")
	return fmt.Errorf("property %s not found", name)
}

// https://w3c.github.io/wot-scripting-api/#the-emitpropertychange-method
func (e *ExposedThing) EmitPropertyChange(name string) (interface{}, error) {
	if _, ok := e.Td.Properties[name]; ok {
		if handler, ok2 := e.propertiesObserveHandlers[name]; ok2 {
			value, err2 := handler()
			if err2 != nil {
				log.Debug().Str("property", name).Err(err2).Msg("[exposedThing:EmitPropertyChange] handler error for property")
				return nil, err2
			}
			return value, nil
		} else if handler, ok2 := e.propertiesReadHandlers[name]; ok2 {
			value, err2 := handler()
			if err2 != nil {
				log.Debug().Str("property", name).Err(err2).Msg("[exposedThing:EmitPropertyChange] handler error for property")
				return nil, err2
			}
			return value, nil
		} else {
			// No handler
			log.Debug().Str("property", name).Msg("[exposedThing:EmitPropertyChange] no handler available for property")
			return nil, nil
		}
	}
	log.Debug().Str("property", name).Msg("[exposedThing:EmitPropertyChange] property not found")
	return nil, fmt.Errorf("property %s not found", name)
}

// https://w3c.github.io/wot-scripting-api/#the-setpropertywritehandler-method
func (e *ExposedThing) SetPropertyWriteHandler(name string, handler PropertyWriteHandler) error {
	if _, ok := e.Td.Properties[name]; ok {
		e.propertiesWriteHandlers[name] = handler
		return nil
	}
	log.Debug().Str("property", name).Msg("[exposedThing:SetPropertyWriteHandler] property not found")
	return fmt.Errorf("property %s not found", name)
}

/*
 * Actions
 */

// https://w3c.github.io/wot-scripting-api/#the-actionhandler-callback
type ActionHandler func(interface{}) (interface{}, error)

// https://w3c.github.io/wot-scripting-api/#the-setactionhandler-method
func (e *ExposedThing) SetActionHandler(name string, handler ActionHandler) error {
	if _, ok := e.Td.Actions[name]; ok {
		e.actionHandlers[name] = handler
		return nil
	}
	log.Debug().Str("action", name).Msg("[exposedThing:SetActionHandler] action not found")
	return fmt.Errorf("action %s not found", name)
}

/*
 * Events
 */

// https://w3c.github.io/wot-scripting-api/#the-eventlistenerhandler-callback
type EventListenerHandler func()

// https://w3c.github.io/wot-scripting-api/#the-eventsubscriptionhandler-callback
type EventSubscriptionHandler func()

// https://w3c.github.io/wot-scripting-api/#the-seteventsubscribehandler-method
func (e *ExposedThing) SetEventSubscribeHandler() {
	// TODO
}

// https://w3c.github.io/wot-scripting-api/#the-seteventunsubscribehandler-method
func (e *ExposedThing) SetEventUnsubscribeHandler() {
	// TODO
}

// https://w3c.github.io/wot-scripting-api/#the-seteventhandler-method
func (e *ExposedThing) SetEventHandler() {
	// TODO
}

// https://w3c.github.io/wot-scripting-api/#the-emitevent-method
func (e *ExposedThing) EmitEvent() {
	// TODO
}
