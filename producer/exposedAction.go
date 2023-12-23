package producer

import (
	"errors"
	"sync"

	"github.com/project-eria/go-wot/dataSchema"
	"github.com/project-eria/go-wot/interaction"
)

// https://w3c.github.io/wot-scripting-api/#the-exposedthing-interface
type ExposedAction interface {
	SetHandler(ActionHandler) error
	GetHandler() ActionHandler
	Input() *dataSchema.Data
	Output() *dataSchema.Data
	// Interaction
	CheckUriVariables(map[string]string) (map[string]interface{}, error)
}

type exposedAction struct {
	handler ActionHandler
	mu      sync.RWMutex
	*interaction.Action
}

func NewExposedAction(interaction *interaction.Action) ExposedAction {
	e := &exposedAction{
		Action: interaction,
	}
	return e
}

// https://w3c.github.io/wot-scripting-api/#the-actionhandler-callback
type ActionHandler func(interface{}, map[string]interface{}) (interface{}, error)

// https://w3c.github.io/wot-scripting-api/#the-setactionhandler-method
func (e *exposedAction) SetHandler(handler ActionHandler) error {
	if handler == nil {
		return errors.New("handler can't be nil")
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	e.handler = handler
	return nil
}

func (e *exposedAction) GetHandler() ActionHandler {
	return e.handler
}

func (e *exposedAction) Input() *dataSchema.Data {
	return e.Action.Input
}

func (e *exposedAction) Output() *dataSchema.Data {
	return e.Action.Output
}
