package producer

import (
	"errors"
	"reflect"
	"sync"

	"github.com/project-eria/go-wot/dataSchema"
	"github.com/project-eria/go-wot/interaction"
	zlog "github.com/rs/zerolog/log"
)

// https://w3c.github.io/wot-scripting-api/#the-exposedthing-interface
type ExposedAction interface {
	SetHandler(ActionHandler) error
	Input() *dataSchema.Data
	Output() *dataSchema.Data
	// Interaction
	CheckUriVariables(map[string]string) (map[string]interface{}, error)
	Run(ExposedThing, string, interface{}, map[string]string) (interface{}, error)
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

func (e *exposedAction) Input() *dataSchema.Data {
	return e.Action.Input
}

func (e *exposedAction) Output() *dataSchema.Data {
	return e.Action.Output
}

func (e *exposedAction) Run(t ExposedThing, key string, data interface{}, parameters map[string]string) (interface{}, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.handler != nil {
		// Check the params (uriVariables) data
		options, err := e.CheckUriVariables(parameters)
		if err != nil {
			return nil, &DataError{
				Message: err.Error(),
			}
		}

		var input interface{}
		// Check the input data
		if e.Input() != nil {
			if data == nil { // TODO: check if it's required
				return nil, &DataError{
					Message: "input value is required",
				}
			}
			if reflect.TypeOf(data).Kind() == reflect.String {
				if input, err = e.Input().FromString(data.(string)); err != nil {
					return nil, &DataError{
						Message: err.Error(),
					}
				}
			} else {
				input = data
			}

			if err := e.Input().Validate(input); err != nil {
				return nil, &DataError{
					Message: "incorrect input value: " + err.Error(),
				}
			}
		} else {
			// Check if provided but not needed
			if data != nil {
				zlog.Warn().Str("action", key).Msg("input value provided, but not declared in the schema")
			}
		}
		// Execute the action requests
		output, err := e.handler(input, options)
		if err != nil {
			return nil, &UnknownError{
				Message: err.Error(),
			}
		}

		// Check the output data
		if e.Output() != nil {
			if err := e.Output().Validate(output); err != nil {
				return nil, &UnknownError{
					Message: "Incorrect handler returned value: " + err.Error(),
				}
			}
			return output, nil // Return the output data
		}
		return nil, nil // No output
	} else {
		return nil, &NotImplementedError{
			Message: "No handler function for the action",
		}
	}
}
