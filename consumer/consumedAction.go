package consumer

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

// Makes a request for invoking an Action and return the result
// https://w3c.github.io/wot-scripting-api/#the-getthingdescription-method
func (c *ConsumedThing) InvokeAction(name string, params interface{}) (interface{}, error) {
	if action, ok := c.td.Actions[name]; ok {
		for _, form := range action.Forms {
			for _, op := range form.Op {
				if op == "invokeaction" {
					value, err := postHTTPJSON(form.Href, params)
					log.Trace().Str("action", name).Str("url", form.Href).Interface("value", value).Msg("[consumer:InvokeAction] Received value")
					return value, err
				}
			}
		}
		log.Error().Str("action", name).Msg("[consumer:InvokeAction] action not invokable")
		return nil, fmt.Errorf("action %s not invokable", name)
	}
	log.Error().Str("action", name).Msg("[consumer:InvokeAction] action not found")
	return nil, fmt.Errorf("action %s not found", name)
}
