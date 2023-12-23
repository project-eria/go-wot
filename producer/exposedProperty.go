package producer

import (
	"errors"
	"sync"

	"github.com/project-eria/go-wot/dataSchema"
	"github.com/project-eria/go-wot/interaction"
)

// https://w3c.github.io/wot-scripting-api/#the-exposedthing-interface
type ExposedProperty interface {
	SetReadHandler(PropertyReadHandler) error
	GetReadHandler() PropertyReadHandler
	SetObserveHandler(PropertyObserveHandler)
	GetObserveHandler() PropertyObserveHandler
	SetObserverSelectorHandler(ObserverSelectorHandler)
	GetObserverSelectorHandler() ObserverSelectorHandler
	SetWriteHandler(PropertyWriteHandler) error
	GetWriteHandler() PropertyWriteHandler
	Data() dataSchema.DataSchema
	IsObservable() bool
	// Interaction
	CheckUriVariables(map[string]string) (map[string]interface{}, error)
}

type exposedProperty struct {
	mu                      sync.RWMutex
	propertyReadHandler     PropertyReadHandler
	propertyWriteHandler    PropertyWriteHandler
	propertyObserveHandler  PropertyObserveHandler
	observerSelectorHandler ObserverSelectorHandler
	*interaction.Property
}

type PropertyChange struct {
	ThingRef       string
	Name           string
	Value          interface{}
	Handler        ObserverSelectorHandler
	EmitParameters map[string]interface{} // Parameters sent via the emit method
}

func NewExposedProperty(interaction *interaction.Property) ExposedProperty {
	return &exposedProperty{
		propertyReadHandler:     nil,
		propertyWriteHandler:    nil,
		propertyObserveHandler:  nil,
		observerSelectorHandler: nil,
		Property:                interaction,
	}
}

// https://w3c.github.io/wot-scripting-api/#the-propertyreadhandler-callback
type PropertyReadHandler func(ExposedThing, string, map[string]interface{}) (interface{}, error)
type PropertyObserveHandler func(ExposedThing, string, map[string]interface{}) (interface{}, error)

// https://w3c.github.io/wot-scripting-api/#the-propertywritehandler-callback
type PropertyWriteHandler func(ExposedThing, string, interface{}, map[string]interface{}) error

type ObserverSelectorHandler func(map[string]interface{}, map[string]interface{}) bool

// emitOptions / listenerOptions

func (p *exposedProperty) SetReadHandler(handler PropertyReadHandler) error {
	if handler == nil {
		return errors.New("read handler can't be nil")
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.propertyReadHandler = handler
	return nil
}

func (p *exposedProperty) GetReadHandler() PropertyReadHandler {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.propertyReadHandler
}

func (p *exposedProperty) SetObserveHandler(handler PropertyObserveHandler) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.propertyObserveHandler = handler
}

func (p *exposedProperty) GetObserveHandler() PropertyObserveHandler {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.propertyObserveHandler
}

func (p *exposedProperty) SetObserverSelectorHandler(handler ObserverSelectorHandler) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.observerSelectorHandler = handler
}

func (p *exposedProperty) GetObserverSelectorHandler() ObserverSelectorHandler {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.observerSelectorHandler
}

func (p *exposedProperty) SetWriteHandler(handler PropertyWriteHandler) error {
	if handler == nil {
		return errors.New("write handler can't be nil")
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.propertyWriteHandler = handler
	return nil
}

func (p *exposedProperty) GetWriteHandler() PropertyWriteHandler {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.propertyWriteHandler
}

func (p *exposedProperty) Data() dataSchema.DataSchema {
	return p.DataSchema
}

func (p *exposedProperty) IsObservable() bool {
	return p.Property.Observable
}
