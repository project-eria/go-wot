package consumer

import (
	"errors"
	"fmt"

	zlog "github.com/rs/zerolog/log"
)

// https://w3c.github.io/wot-scripting-api/#the-subscribeevent-method
func (t *consumedThing) SubscribeEvent(name string, dataVariables map[string]interface{}, listener Listener) error {
	if listener == nil {
		zlog.Error().Str("event", name).Msg("[consumer:SubscribeEvent] missing listener")
		return fmt.Errorf("missing listener for event %s", name)
	}
	if event, ok := t.td.Events[name]; ok {
		for _, form := range event.Forms {
			form := form // Copy https://go.dev/doc/faq#closures_and_goroutines
			if client := t.consumer.GetClientFor(form); client != nil {
				sub := &Subscription{
					Type: "event",
					Name: name,
				}
				uriData, err := getUriVariables(event.UriVariables, dataVariables)
				if err != nil {
					return err
				}
				uri, err := client.SubscribeResource(form, uriData, sub, listener)
				zlog.Trace().Str("event", name).Str("uri", uri).Msg("[consumer:SubscribeEvent] listener added")
				return err
			}
		}
		return errors.New("can't find client")
	}
	zlog.Error().Str("event", name).Msg("[consumer:SubscribeEvent] event not found")
	return fmt.Errorf("event %s not found", name)
}
