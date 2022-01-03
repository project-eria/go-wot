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
type EventListenerHandler func()

// https://w3c.github.io/wot-scripting-api/#the-eventsubscriptionhandler-callback
type EventSubscriptionHandler func()
