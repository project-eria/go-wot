package consumer

import (
	"context"
	"sync"

	"github.com/project-eria/go-wot/thing"
)

// https://w3c.github.io/wot-scripting-api/#the-consumedthing-interface
type ConsumedThing interface {
	GetThingDescription() *thing.Thing
	// Properties
	ReadProperty(string) (interface{}, error)
	ReadMultipleProperties()
	ReadAllProperties()
	WriteProperty(string, interface{}) (interface{}, error)
	WriteMultipleProperties()
	ObserveProperty(string, Listener) error
	// Actions
	InvokeAction(string, interface{}) (interface{}, error)
	// Events
	SubscribeEvent(string, Listener) error
}

type consumedThing struct {
	consumer  *Consumer
	td        *thing.Thing
	connected bool
	_ctx      context.Context // TO CHECK if needed
	_wait     *sync.WaitGroup // TO CHECK if needed
}

// https://w3c.github.io/wot-scripting-api/#the-getthingdescription-method
func (t *consumedThing) GetThingDescription() *thing.Thing {
	return t.td
}

type Subscription struct {
	Type string
	Name string
	//	Interaction *interaction.Interaction
	// Let subscription's [[thing]] be the value of thing.
	// Let subscription's [[form]] be the Form associated with formIndex in [[interaction]]'s forms array if option's formIndex is defined, otherwise let [[form]] be a Form in [[interaction]]'s forms array whose op is "observeproperty", selected by the implementation.
	// If subscription's [[form]] is failure, reject promise with a SyntaxError and abort these steps.
	// If subscription's [[interaction]] is undefined, reject promise with a NotFoundError and abort this steps.
}
