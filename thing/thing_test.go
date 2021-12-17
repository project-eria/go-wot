package thing

import (
	"testing"

	"github.com/project-eria/go-wot/action"
	"github.com/project-eria/go-wot/property"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	type args struct {
		urn         string
		ref         string
		title       string
		types       []string
		description string
	}
	tests := []struct {
		name string
		args args
		want *Thing
	}{
		{
			name: "Thing, with no type",
			args: args{
				urn:         "dev:ops:my-actuator-1234",
				ref:         "example",
				title:       "ActuatorExample",
				types:       nil,
				description: "An actuator example",
			},
			want: &Thing{
				AtContext:   "http://www.w3.org/ns/td",
				ID:          "urn:dev:ops:my-actuator-1234",
				AtTypes:     make([]string, 0),
				Ref:         "example",
				Title:       "ActuatorExample",
				Description: "An actuator example",
				Properties:  make(map[string]*property.Property),
				Actions:     make(map[string]*action.Action),
			},
		},
		{
			name: "Thing, with types",
			args: args{
				urn:         "dev:ops:my-actuator-1234",
				ref:         "example",
				title:       "ActuatorExample",
				types:       []string{"OnOffSwitch", "Lamp"},
				description: "An actuator example",
			},
			want: &Thing{
				AtContext:   "http://www.w3.org/ns/td",
				ID:          "urn:dev:ops:my-actuator-1234",
				AtTypes:     []string{"OnOffSwitch", "Lamp"},
				Title:       "ActuatorExample",
				Description: "An actuator example",
				Ref:         "example",
				Properties:  make(map[string]*property.Property),
				Actions:     make(map[string]*action.Action),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.urn, tt.args.ref, tt.args.title, tt.args.description, tt.args.types)
			assert.NoError(t, err, "should not return error")
			assert.Equal(t, tt.want, got, "they should be equal")
		})
	}
}

// var thingA = &Thing{urn: "dev:ops:my-actuator-1234",
// 	title:       "ActuatorExample",
// 	Ref:         "",
// 	properties:  make(map[string]*Property),
// 	types:       []string{"OnOffSwitch", "Lamp"},
// 	description: "An actuator example",
// 	context:     "http://www.w3.org/ns/td",
// }
