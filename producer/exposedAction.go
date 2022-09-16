package producer

import (
	"errors"
	"sync"

	"github.com/project-eria/go-wot/interaction"
)

// https://w3c.github.io/wot-scripting-api/#the-exposedthing-interface
type ExposedAction struct {
	handler ActionHandler
	mu      sync.RWMutex
	*interaction.Action
}

func NewExposedAction(interaction *interaction.Action) *ExposedAction {
	e := &ExposedAction{
		Action: interaction,
	}
	return e
}

// https://w3c.github.io/wot-scripting-api/#the-actionhandler-callback
type ActionHandler func(interface{}) (interface{}, error)

// https://w3c.github.io/wot-scripting-api/#the-setactionhandler-method
func (e *ExposedAction) SetHandler(handler ActionHandler) error {
	if handler == nil {
		return errors.New("handler can't be nil")
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	e.handler = handler
	return nil
}

func (e *ExposedAction) GetHandler() ActionHandler {
	return e.handler
}
