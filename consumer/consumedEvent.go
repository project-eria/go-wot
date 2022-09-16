package consumer

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
)

// https://w3c.github.io/wot-scripting-api/#the-subscribeevent-method
func (t *ConsumedThing) SubscribeEvent(name string, listener Listener) error {
	if listener == nil {
		log.Error().Str("event", name).Msg("[consumer:SubscribeEvent] missing listener")
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
				err := client.SubscribeResource(form, sub, listener)
				log.Trace().Str("event", name).Str("url", form.Href).Msg("[consumer:SubscribeEvent] listener added")
				return err
			}
		}
		return errors.New("can't find client")
	}
	log.Error().Str("event", name).Msg("[consumer:SubscribeEvent] event not found")
	return fmt.Errorf("event %s not found", name)
}
