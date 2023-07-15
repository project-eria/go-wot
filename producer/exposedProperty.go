package producer

import (
	"errors"
	"sync"

	"github.com/project-eria/go-wot/interaction"
)

// https://w3c.github.io/wot-scripting-api/#the-exposedthing-interface
type ExposedProperty struct {
	mu                      sync.RWMutex
	propertyReadHandler     PropertyReadHandler
	propertyWriteHandler    PropertyWriteHandler
	propertyObserveHandler  PropertyObserveHandler
	observerSelectorHandler ObserverSelectorHandler
	*interaction.Property
}

type PropertyChange struct {
	ThingRef string
	Name     string
	Value    interface{}
	Handler  ObserverSelectorHandler
	Options  map[string]string
}

func NewExposedProperty(interaction *interaction.Property) *ExposedProperty {
	return &ExposedProperty{
		propertyReadHandler:     nil,
		propertyWriteHandler:    nil,
		propertyObserveHandler:  nil,
		observerSelectorHandler: nil,
		Property:                interaction,
	}
}

// https://w3c.github.io/wot-scripting-api/#the-propertyreadhandler-callback
type PropertyReadHandler func(*ExposedThing, string, map[string]string) (interface{}, error)
type PropertyObserveHandler func(*ExposedThing, string, map[string]string) (interface{}, error)

// https://w3c.github.io/wot-scripting-api/#the-propertywritehandler-callback
type PropertyWriteHandler func(*ExposedThing, string, interface{}, map[string]string) error

type ObserverSelectorHandler func(map[string]string, map[string]string) bool

func (p *ExposedProperty) SetReadHandler(handler PropertyReadHandler) error {
	if handler == nil {
		return errors.New("read handler can't be nil")
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.propertyReadHandler = handler
	return nil
}

func (p *ExposedProperty) GetReadHandler() PropertyReadHandler {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.propertyReadHandler
}

func (p *ExposedProperty) SetObserveHandler(handler PropertyObserveHandler) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.propertyObserveHandler = handler
}

func (p *ExposedProperty) GetObserveHandler() PropertyObserveHandler {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.propertyObserveHandler
}

func (p *ExposedProperty) SetObserverSelectorHandler(handler ObserverSelectorHandler) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.observerSelectorHandler = handler
}

func (p *ExposedProperty) GetObserverSelectorHandler() ObserverSelectorHandler {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.observerSelectorHandler
}

func (p *ExposedProperty) SetWriteHandler(handler PropertyWriteHandler) error {
	if handler == nil {
		return errors.New("write handler can't be nil")
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.propertyWriteHandler = handler
	return nil
}

func (p *ExposedProperty) GetWriteHandler() PropertyWriteHandler {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.propertyWriteHandler
}
