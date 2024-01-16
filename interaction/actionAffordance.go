package interaction

import (
	"github.com/project-eria/go-wot/dataSchema"
)

// ActionAffordance: An Interaction Affordance that allows to invoke a function of the Thing,
// which manipulates state (e.g., toggling a lamp on or off) or triggers a
// process on the Thing (e.g., dim a lamp over time).
// https://w3c.github.io/wot-thing-description/#actionaffordance

// Note: Input and Output are pointer because they can be nil
type Action struct {
	Input  *dataSchema.Data `json:"input,omitempty"`  // (optional) Used to define the input data schema of the Action.
	Output *dataSchema.Data `json:"output,omitempty"` // (optional) Used to define the output data schema of the Action.
	// safe	Signals if the Action is safe (=true) or not. Used to signal if there is no internal state (cf. resource state) is changed when invoking an Action. In that case responses can be cached as example.	with default	boolean
	// idempotent	Indicates whether the Action is idempotent (=true) or not. Informs whether the Action can be called repeatedly with the same result, if present, based on the same input.	with default	boolean
	Interaction
}

type ActionHandler func(interface{}) (interface{}, error)

func NewAction(key string, title string, description string, options ...ActionOption) *Action {
	opts := &ActionOptions{}
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
	return &Action{
		Interaction: interaction,
		Input:       opts.Input,
		Output:      opts.Output,
	}
}

type ActionOption func(*ActionOptions)

type ActionOptions struct {
	Input        *dataSchema.Data
	Output       *dataSchema.Data
	UriVariables map[string]dataSchema.Data
}

func ActionInput(input *dataSchema.Data) ActionOption {
	return func(opts *ActionOptions) {
		opts.Input = input
	}
}

func ActionOutput(output *dataSchema.Data) ActionOption {
	return func(opts *ActionOptions) {
		opts.Output = output
	}
}

func ActionUriVariable(key string, data dataSchema.Data) ActionOption {
	return func(opts *ActionOptions) {
		if opts.UriVariables == nil {
			opts.UriVariables = make(map[string]dataSchema.Data)
		}
		opts.UriVariables[key] = data
	}
}
