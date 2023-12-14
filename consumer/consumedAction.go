package consumer

import (
	"errors"
	"fmt"

	zlog "github.com/rs/zerolog/log"
)

// Makes a request for invoking an Action and return the result
// https://w3c.github.io/wot-scripting-api/#the-getthingdescription-method
func (t *consumedThing) InvokeAction(name string, dataVariables map[string]interface{}, params interface{}) (interface{}, error) {
	if action, ok := t.td.Actions[name]; ok {
		for _, form := range action.Forms {
			form := form // Copy https://go.dev/doc/faq#closures_and_goroutines
			for _, op := range form.Op {
				op := op // Copy https://go.dev/doc/faq#closures_and_goroutines
				if op == "invokeaction" {
					if client := t.consumer.GetClientFor(form); client != nil {
						uriData, err := getUriVariables(action.UriVariables, dataVariables)
						if err != nil {
							return nil, err
						}
						value, uri, err := client.InvokeResource(form, uriData, params)
						zlog.Trace().Str("action", name).Str("uri", uri).Interface("value", value).Msg("[consumer:InvokeAction] Received value")
						return value, err
					}
					return nil, errors.New("can't find client")
				}
			}
		}
		zlog.Error().Str("action", name).Msg("[consumer:InvokeAction] action not invokable")
		return nil, fmt.Errorf("action %s not invokable", name)
	}
	zlog.Error().Str("action", name).Msg("[consumer:InvokeAction] action not found")
	return nil, fmt.Errorf("action %s not found", name)
}
