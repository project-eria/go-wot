package interaction

import (
	"github.com/project-eria/go-wot/dataSchema"
)

// ActionAffordance: An Interaction Affordance that allows to invoke a function of the Thing,
// which manipulates state (e.g., toggling a lamp on or off) or triggers a
// process on the Thing (e.g., dim a lamp over time).
// https://w3c.github.io/wot-thing-description/#actionaffordance
type Action struct {
	Input  *dataSchema.Data `json:"input,omitempty"`  // (optional) Used to define the input data schema of the Action.
	Output *dataSchema.Data `json:"output,omitempty"` // (optional) Used to define the output data schema of the Action.
	// safe	Signals if the Action is safe (=true) or not. Used to signal if there is no internal state (cf. resource state) is changed when invoking an Action. In that case responses can be cached as example.	with default	boolean
	// idempotent	Indicates whether the Action is idempotent (=true) or not. Informs whether the Action can be called repeatedly with the same result, if present, based on the same input.	with default	boolean
	Interaction
}

type ActionHandler func(interface{}) (interface{}, error)

func NewAction(key string, title string, description string, input *dataSchema.Data, output *dataSchema.Data) *Action {
	interaction := Interaction{
		Key:         key,
		Title:       title,
		Description: description,
		Forms:       []*Form{},
	}
	return &Action{
		Interaction: interaction,
		Input:       input,
		Output:      output,
	}
}
