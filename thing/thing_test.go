package thing

import (
	"testing"

	"github.com/project-eria/go-wot/dataSchema"
	"github.com/project-eria/go-wot/interaction"
	"github.com/project-eria/go-wot/securityScheme"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	type args struct {
		urn         string
		title       string
		description string
		types       []string
	}
	tests := []struct {
		name  string
		args  args
		want  *Thing
		error string
	}{
		{
			name: "Thing, with no type",
			args: args{
				urn:         "dev:ops:my-actuator-1234",
				title:       "ActuatorExample",
				description: "An actuator example",
				types:       nil,
			},
			want: &Thing{
				AtContext:           "http://www.w3.org/ns/td",
				ID:                  "urn:dev:ops:my-actuator-1234",
				AtTypes:             make([]string, 0),
				Title:               "ActuatorExample",
				Description:         "An actuator example",
				Properties:          make(map[string]*interaction.Property),
				Actions:             make(map[string]*interaction.Action),
				Security:            []string{},
				SecurityDefinitions: make(map[string]securityScheme.SecurityScheme),
			},
		},
		{
			name: "Thing, with types",
			args: args{
				urn:         "dev:ops:my-actuator-1234",
				title:       "ActuatorExample",
				types:       []string{"OnOffSwitch", "Lamp"},
				description: "An actuator example",
			},
			want: &Thing{
				AtContext:           "http://www.w3.org/ns/td",
				ID:                  "urn:dev:ops:my-actuator-1234",
				AtTypes:             []string{"OnOffSwitch", "Lamp"},
				Title:               "ActuatorExample",
				Description:         "An actuator example",
				Properties:          make(map[string]*interaction.Property),
				Actions:             make(map[string]*interaction.Action),
				Security:            []string{},
				SecurityDefinitions: make(map[string]securityScheme.SecurityScheme),
			},
		},
		{
			name: "Thing, with empty URN",
			args: args{
				urn:         "",
				title:       "ActuatorExample",
				description: "An actuator example",
				types:       nil,
			},
			error: "Thing URN can't be empty",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.urn, tt.args.title, tt.args.description, tt.args.types)
			if tt.error != "" {
				assert.Error(t, err, "should return error")
				assert.Equal(t, tt.error, err.Error(), "they should be equal")
			} else {
				assert.NoError(t, err, "should not return error")
				assert.Equal(t, tt.want, got, "they should be equal")
			}
		})
	}
}

func TestAddSecurity(t *testing.T) {
	mything, err := New("dev:ops:my-actuator-1234", "ActuatorExample", "An actuator example", nil)
	noSecurityScheme := securityScheme.NewNoSecurity()
	mything.AddSecurity("no_sec", noSecurityScheme)
	want := &Thing{
		AtContext:   "http://www.w3.org/ns/td",
		ID:          "urn:dev:ops:my-actuator-1234",
		AtTypes:     []string{},
		Title:       "ActuatorExample",
		Description: "An actuator example",
		Properties:  make(map[string]*interaction.Property),
		Actions:     make(map[string]*interaction.Action),
		Security:    []string{"no_sec"},
		SecurityDefinitions: map[string]securityScheme.SecurityScheme{
			"no_sec": noSecurityScheme,
		},
	}
	assert.NoError(t, err, "should not return error")
	assert.Equal(t, want, mything, "they should be equal")

}

func TestAddProperty(t *testing.T) {
	mything, err := New("dev:ops:my-actuator-1234", "ActuatorExample", "An actuator example", nil)
	data := dataSchema.NewBoolean(false)
	property := interaction.NewProperty(
		"x",
		"y",
		"z",
		data,
	)
	mything.AddProperty(&property)
	want := &Thing{
		AtContext:           "http://www.w3.org/ns/td",
		ID:                  "urn:dev:ops:my-actuator-1234",
		AtTypes:             []string{},
		Title:               "ActuatorExample",
		Description:         "An actuator example",
		Properties:          map[string]*interaction.Property{"x": &property},
		Actions:             make(map[string]*interaction.Action),
		Security:            []string{},
		SecurityDefinitions: make(map[string]securityScheme.SecurityScheme),
	}
	assert.NoError(t, err, "should not return error")
	assert.Equal(t, want, mything, "they should be equal")
}

func TestAddAction(t *testing.T) {
	mything, err := New("dev:ops:my-actuator-1234", "ActuatorExample", "An actuator example", nil)
	aAction := interaction.NewAction(
		"a",
		"No Input, No Output",
		"",
		nil,
		nil,
	)
	mything.AddAction(aAction)
	want := &Thing{
		AtContext:           "http://www.w3.org/ns/td",
		ID:                  "urn:dev:ops:my-actuator-1234",
		AtTypes:             []string{},
		Title:               "ActuatorExample",
		Description:         "An actuator example",
		Properties:          make(map[string]*interaction.Property),
		Actions:             map[string]*interaction.Action{"a": aAction},
		Security:            []string{},
		SecurityDefinitions: make(map[string]securityScheme.SecurityScheme),
	}
	assert.NoError(t, err, "should not return error")
	assert.Equal(t, want, mything, "they should be equal")
}
