package producer

import (
	"sync"

	"github.com/project-eria/go-wot/interaction"
)

// https://w3c.github.io/wot-scripting-api/#the-exposedthing-interface
type ExposedEvent interface {
	SetSubscribeHandler(EventSubscriptionHandler)
	SetUnSubscribeHandler()
	SetEventHandler(EventListenerHandler)
	GetEventHandler() EventListenerHandler
	SetListenerSelectorHandler(ListenerSelectorHandler)
	GetListenerSelectorHandler() ListenerSelectorHandler
	// Interaction
	CheckUriVariables(map[string]string) error
}

type exposedEvent struct {
	eventListenerHandler     EventListenerHandler
	eventSubscriptionHandler EventSubscriptionHandler
	listenerSelectorHandler  ListenerSelectorHandler
	mu                       sync.RWMutex
	*interaction.Event
}

type Event struct {
	ThingRef string
	Name     string
	Value    interface{}
	Handler  ListenerSelectorHandler
	Options  map[string]string
}

func NewExposedEvent(interaction *interaction.Event) ExposedEvent {
	e := &exposedEvent{
		Event: interaction,
	}
	return e
}

// https://w3c.github.io/wot-scripting-api/#the-eventlistenerhandler-callback
type EventListenerHandler func() (interface{}, error)

// https://w3c.github.io/wot-scripting-api/#the-eventsubscriptionhandler-callback
type EventSubscriptionHandler func(ExposedThing, string, map[string]string) (interface{}, error)

type ListenerSelectorHandler func(map[string]string, map[string]string) bool

// https://w3c.github.io/wot-scripting-api/#the-seteventsubscribehandler-method
func (e *exposedEvent) SetSubscribeHandler(handler EventSubscriptionHandler) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.eventSubscriptionHandler = handler
}

// https://w3c.github.io/wot-scripting-api/#the-seteventunsubscribehandler-method
func (e *exposedEvent) SetUnSubscribeHandler() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.eventSubscriptionHandler = nil
}

// https://w3c.github.io/wot-scripting-api/#the-seteventhandler-method
func (e *exposedEvent) SetEventHandler(handler EventListenerHandler) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.eventListenerHandler = handler
}

func (e *exposedEvent) GetEventHandler() EventListenerHandler {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.eventListenerHandler
}

func (e *exposedEvent) SetListenerSelectorHandler(handler ListenerSelectorHandler) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.listenerSelectorHandler = handler
}

func (e *exposedEvent) GetListenerSelectorHandler() ListenerSelectorHandler {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.listenerSelectorHandler
}
