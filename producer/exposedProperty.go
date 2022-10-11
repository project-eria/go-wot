package producer

import (
	"errors"
	"fmt"
	"sync"

	"github.com/project-eria/go-wot/interaction"
	"github.com/rs/zerolog/log"
)

// https://w3c.github.io/wot-scripting-api/#the-exposedthing-interface
type ExposedProperty struct {
	Value                  interface{}
	propertyReadHandler    PropertyReadHandler
	propertyWriteHandler   PropertyWriteHandler
	propertyObserveHandler PropertyObserveHandler
	mu                     sync.RWMutex
	*interaction.Property
}

type PropertyChange struct {
	ThingRef string
	Name     string
	Value    interface{}
}

func NewExposedProperty(interaction *interaction.Property) *ExposedProperty {
	return &ExposedProperty{
		Value:                  interaction.Default,
		propertyReadHandler:    defaultPropertyReadHandler,
		propertyWriteHandler:   defaultPropertyWriteHandler,
		propertyObserveHandler: nil,
		Property:               interaction,
	}
}

// https://w3c.github.io/wot-scripting-api/#the-propertyreadhandler-callback
type PropertyReadHandler func(*ExposedThing, string) (interface{}, error)
type PropertyObserveHandler func(*ExposedThing, string) (interface{}, error)

// https://w3c.github.io/wot-scripting-api/#the-propertywritehandler-callback
type PropertyWriteHandler func(*ExposedThing, string, interface{}) error

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

func defaultPropertyReadHandler(t *ExposedThing, name string) (interface{}, error) {
	if property, ok := t.ExposedProperties[name]; ok {
		log.Trace().Str("property", name).Interface("value", property.Value).Msg("[exposedProperty:defaultPropertyReadHandler] Value get")
		property.mu.Lock()
		defer property.mu.Unlock()
		return property.Value, nil
	}
	return nil, fmt.Errorf("property %s not found", name)
}

func defaultPropertyWriteHandler(t *ExposedThing, name string, value interface{}) error {
	if property, ok := t.ExposedProperties[name]; ok {
		if err := property.Data.Check(value); err != nil {
			log.Error().Str("property", name).Interface("value", value).Err(err).Msg("[exposedProperty:defaultPropertyWriteHandler]")
			return err
		}
		property.mu.Lock()
		property.Value = value
		property.mu.Unlock()
		log.Trace().Str("property", name).Interface("value", value).Msg("[exposedProperty:defaultPropertyWriteHandler] Value set")
		if err := t.EmitPropertyChange(name); err != nil {
			log.Error().Str("property", name).Interface("value", value).Err(err).Msg("[exposedProperty:defaultPropertyWriteHandler]")
			return err
		}
		return nil
	}
	return fmt.Errorf("property %s not found", name)
}
