package producer

import (
	"errors"
	"reflect"
	"sync"

	"github.com/project-eria/go-wot/dataSchema"
	"github.com/project-eria/go-wot/interaction"
)

// https://w3c.github.io/wot-scripting-api/#the-exposedthing-interface
type ExposedProperty interface {
	SetReadHandler(PropertyReadHandler) error
	SetObserveHandler(PropertyObserveHandler)
	GetObserveHandler() PropertyObserveHandler
	SetObserverSelectorHandler(ObserverSelectorHandler)
	GetObserverSelectorHandler() ObserverSelectorHandler
	SetWriteHandler(PropertyWriteHandler) error
	Data() dataSchema.DataSchema
	IsObservable() bool
	Read(ExposedThing, string, map[string]string) (interface{}, error)
	Write(ExposedThing, string, interface{}, map[string]string) error
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

func (p *exposedProperty) Read(t ExposedThing, key string, parameters map[string]string) (interface{}, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.propertyReadHandler != nil {
		// Check the options (uriVariables) data
		options, err := p.CheckUriVariables(parameters)
		if err != nil {
			return nil, &DataError{
				Message: err.Error(),
			}
		}
		// Call the function that handle the property read
		content, err := p.propertyReadHandler(t, key, options)
		if err != nil {
			return nil, &UnknownError{
				Message: err.Error(),
			}
		}
		// TODO: check output??
		return content, nil
	}
	return nil, &NotImplementedError{
		Message: "No handler function for reading the property",
	}
}

func (p *exposedProperty) Write(t ExposedThing, key string, data interface{}, parameters map[string]string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.propertyWriteHandler != nil {
		// Check the options (uriVariables) data
		options, err := p.CheckUriVariables(parameters)
		if err != nil {
			return &DataError{
				Message: err.Error(),
			}
		}

		// Check if data has been provided
		if data == nil {
			return &DataError{
				Message: "No data provided",
			}
		}

		var input interface{}
		if reflect.TypeOf(data).Kind() == reflect.String {
			if input, err = p.Data().FromString(data.(string)); err != nil {
				return &DataError{
					Message: err.Error(),
				}
			}
		} else {
			input = data
		}

		// Check the input data
		if p.Data() != nil { // TODO should not be nil
			if err := p.Data().Validate(input); err != nil {
				return &DataError{
					Message: "incorrect input value: " + err.Error(),
				}
			}
		}

		// Call the function that handle the property write
		err = p.propertyWriteHandler(t, key, input, options)
		if err != nil {
			return &UnknownError{
				Message: err.Error(),
			}
		}

		// Notify all listeners that the property changed
		if err := t.EmitPropertyChange(key, input, options); err != nil {
			return &UnknownError{
				Message: err.Error(),
			}
		}

		return nil
	}
	return &NotImplementedError{
		Message: "No handler function for writing the property",
	}
}
