package producer

import (
	"sync"
)

// https://w3c.github.io/wot-scripting-api/#the-exposedthing-interface
type ExposedEvent struct {
	eventListenerHandler     EventListenerHandler
	eventSubscriptionHandler EventSubscriptionHandler

	mu sync.RWMutex
}

// https://w3c.github.io/wot-scripting-api/#the-eventlistenerhandler-callback
type EventListenerHandler func(*ExposedThing, string) (interface{}, error)

// https://w3c.github.io/wot-scripting-api/#the-eventsubscriptionhandler-callback
type EventSubscriptionHandler func(*ExposedThing, string) (interface{}, error)

// https://w3c.github.io/wot-scripting-api/#the-seteventsubscribehandler-method
func (e *ExposedEvent) SetSubscribeHandler(handler EventSubscriptionHandler) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.eventSubscriptionHandler = handler
}

// https://w3c.github.io/wot-scripting-api/#the-seteventunsubscribehandler-method
func (e *ExposedEvent) SetUnSubscribeHandler() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.eventSubscriptionHandler = nil
}

// https://w3c.github.io/wot-scripting-api/#the-seteventhandler-method
func (e *ExposedEvent) SetEventHandler(handler EventListenerHandler) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.eventListenerHandler = handler
}
