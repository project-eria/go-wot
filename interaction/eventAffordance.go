package interaction

import (
	"github.com/project-eria/go-wot/dataSchema"
)

// EventAffordance: An Interaction Affordance that describes
// an event source, which asynchronously pushes event data to
// Consumers (e.g., overheating alerts).
// https://w3c.github.io/wot-thing-description/#eventaffordance
type Event struct {
	// TODO Subscription *dataSchema.Data `json:"subscription,omitempty"` // (optional) Defines data that needs to be passed upon subscription, e.g., filters or message format for setting up Webhooks
	Data *dataSchema.Data `json:"data,omitempty"` // (Optional) Defines the data schema of the Event instance messages pushed by the Thing
	// TODO Cancellation *dataSchema.Data `json:"cancellation,omitempty"` // (Optional) Defines any data that needs to be passed to cancel a subscription, e.g., a specific message to remove a Webhook
	Interaction
}

func NewEvent(key string, title string, description string, options ...EventOption) *Event {
	opts := &EventOptions{}
	for _, option := range options {
		if option != nil {
			option(opts)
		}
	}

	interaction := Interaction{
		Key:          key,
		Title:        title,
		Description:  description,
		Forms:        []*Form{},
		UriVariables: opts.UriVariables,
	}

	return &Event{
		Data:        opts.Data,
		Interaction: interaction,
	}
}

type EventOption func(*EventOptions)

type EventOptions struct {
	Data         *dataSchema.Data
	UriVariables map[string]dataSchema.Data
}

func EventData(data *dataSchema.Data) EventOption {
	return func(opts *EventOptions) {
		opts.Data = data
	}
}

func EventUriVariable(key string, data dataSchema.Data) EventOption {
	return func(opts *EventOptions) {
		if opts.UriVariables == nil {
			opts.UriVariables = make(map[string]dataSchema.Data)
		}
		opts.UriVariables[key] = data
	}
}
