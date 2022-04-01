package interaction

import (
	"github.com/project-eria/go-wot/dataSchema"
)

// EventAffordance: An Interaction Affordance that describes
// an event source, which asynchronously pushes event data to
// Consumers (e.g., overheating alerts).
// https://w3c.github.io/wot-thing-description/#eventaffordance
type Event struct {
	Subscription *dataSchema.Data `json:"subscription,omitempty"` // (optional) Defines data that needs to be passed upon subscription, e.g., filters or message format for setting up Webhooks
	Data         *dataSchema.Data `json:"data,omitempty"`         // (Optional) Defines the data schema of the Event instance messages pushed by the Thing
	Cancellation *dataSchema.Data `json:"cancellation,omitempty"` // (Optional) Defines any data that needs to be passed to cancel a subscription, e.g., a specific message to remove a Webhook
	Interaction
}

func NewEvent(key string, title string, description string, data *dataSchema.Data) *Event {
	interaction := Interaction{
		Key:         key,
		Title:       title,
		Description: description,
		Forms:       []*Form{},
	}

	return &Event{
		Data:        data,
		Interaction: interaction,
	}
}
