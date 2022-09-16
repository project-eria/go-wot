package consumer

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
)

// Makes a request for invoking an Action and return the result
// https://w3c.github.io/wot-scripting-api/#the-getthingdescription-method
func (t *ConsumedThing) InvokeAction(name string, params interface{}) (interface{}, error) {
	if action, ok := t.td.Actions[name]; ok {
		for _, form := range action.Forms {
			form := form // Copy https://go.dev/doc/faq#closures_and_goroutines
			for _, op := range form.Op {
				op := op // Copy https://go.dev/doc/faq#closures_and_goroutines
				if op == "invokeaction" {
					if client := t.consumer.GetClientFor(form); client != nil {
						value, err := client.InvokeResource(form, params)
						log.Trace().Str("action", name).Str("url", form.Href).Interface("value", value).Msg("[consumer:InvokeAction] Received value")
						return value, err
					}
					return nil, errors.New("can't find client")
				}
			}
		}
		log.Error().Str("action", name).Msg("[consumer:InvokeAction] action not invokable")
		return nil, fmt.Errorf("action %s not invokable", name)
	}
	log.Error().Str("action", name).Msg("[consumer:InvokeAction] action not found")
	return nil, fmt.Errorf("action %s not found", name)
}
