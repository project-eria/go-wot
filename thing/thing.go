package thing

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/project-eria/go-wot/interaction"
	"github.com/project-eria/go-wot/securityScheme"
	zlog "github.com/rs/zerolog/log"
)

var namespace = "https://www.w3.org/2022/wot/td/v1.1"

// Thing resource provides a Thing Description for a device.
// A Thing Resource is considered the root resource of a Thing.
type Thing struct {
	ID           string            `json:"id,omitempty"`           // (optional) Identifier of the Thing in form of a URI [RFC3986] (e.g., stable URI, temporary and mutable URI, URI with local IP address, URN, etc.).
	AtContext    map[string]string `json:"@context"`               // (mandatory) JSON-LD keyword to define short-hand names called terms that are used throughout a TD document. anyURI or Array
	AtTypes      []string          `json:"@type,omitempty"`        // (optional) JSON-LD keyword to label the object with semantic tags (or types).
	Title        string            `json:"title"`                  // (mandatory) Provides a human-readable title (e.g., display a text for UI representation) based on a default language.
	Titles       []string          `json:"titles,omitempty"`       // (optional) Provides multi-language human-readable titles (e.g., display a text for UI representation in different languages).
	Description  string            `json:"description,omitempty"`  // (optional) Provides additional (human-readable) information based on a default language.
	Descriptions []string          `json:"descriptions,omitempty"` // Can be used to support (human-readable) information in different languages. Also see MultiLanguage.
	Version      map[string]string `json:"version,omitempty"`      // Provides version information.	optional	VersionInfo
	// created	Provides information when the TD instance was created.	optional	dateTime
	// modified	Provides information when the TD instance was last modified.	optional	dateTime
	// support	Provides information about the TD maintainer as URI scheme (e.g., mailto [RFC6068], tel [RFC3966], https).	optional	anyURI
	// base	Define the base URI that is used for all relative URI references throughout a TD document. In TD instances, all relative URIs are resolved relative to the base URI using the algorithm defined in [RFC3986].

	// base does not affect the URIs used in @context and the IRIs used within Linked Data [LINKED-DATA] graphs that are relevant when semantic processing is applied to TD instances.	optional	anyURI
	Properties map[string]*interaction.Property `json:"properties,omitempty"` // (optional) All Property-based Interaction Affordances of the Thing.
	Actions    map[string]*interaction.Action   `json:"actions,omitempty"`    // (optional) All Action-based Interaction Affordances of the Thing.
	Events     map[string]*interaction.Event    `json:"events,omitempty"`     // (optional) All Event-based Interaction Affordances of the Thing.
	// links	Provides Web links to arbitrary resources that relate to the specified Thing Description.	optional	Array of Link
	Forms               []interaction.Form                       `json:"forms,omitempty"`     // (optional) Set of form hypermedia controls that describe how an operation can be performed. Forms are serializations of Protocol Bindings. In this version of TD, all operations that can be described at the Thing level are concerning how to interact with the Thing's Properties collectively at once.
	Security            []string                                 `json:"security"`            // (mandatory) Set of security definition names, chosen from those defined in securityDefinitions. These must all be satisfied for access to resources.
	SecurityDefinitions map[string]securityScheme.SecurityScheme `json:"securityDefinitions"` // (mandatory) Set of named security configurations (definitions only). Not actually applied unless names are used in a security name-value pair.
	// profile	Indicates the WoT Profile mechanisms followed by this Thing Description and the corresponding Thing implementation.

	MU sync.RWMutex `json:"-"`
}

// New thing construct
func New(urn string, version string, title string, description string, types []string) (*Thing, error) {
	if urn == "" {
		return nil, errors.New("Thing URN can't be empty")
	}

	thing := Thing{
		AtContext:           map[string]string{"": namespace},
		AtTypes:             types,
		ID:                  "urn:" + urn,
		Version:             map[string]string{"instance": version},
		Title:               title,
		Description:         description,
		Security:            []string{},
		SecurityDefinitions: make(map[string]securityScheme.SecurityScheme),
		Properties:          make(map[string]*interaction.Property),
		Actions:             make(map[string]*interaction.Action),
		Events:              make(map[string]*interaction.Event),
	}

	if thing.AtTypes == nil {
		thing.AtTypes = make([]string, 0)
	}
	return &thing, nil
}

func (t *Thing) AddContext(key string, context string) {
	if key == "" {
		zlog.Error().Str("uri", context).Msg("[thing:AddContext] missing prefix for context")
		return
	}
	t.AtContext[key] = context
}

func (t *Thing) AddVersion(key string, version string) {
	t.Version[key] = version
}

func (t *Thing) AddSecurity(key string, definition securityScheme.SecurityScheme) {
	t.Security = append(t.Security, key)
	t.SecurityDefinitions[key] = definition
}

// Ref: http://choly.ca/post/go-json-marshalling/
func (t *Thing) MarshalJSON() ([]byte, error) {
	type ThingOrigin Thing

	var (
		modifiedSecurity  interface{}
		modifiedAtContext interface{}
	)

	// AtContext can be a string or an map of string
	if len(t.AtContext) == 1 {
		modifiedAtContext = t.AtContext[""]
	} else {
		modifiedAtContext = []interface{}{}
		atcontext := make(map[string]string)
		for k, v := range t.AtContext {
			if k == v || k == "" {
				modifiedAtContext = append(modifiedAtContext.([]interface{}), v)
			} else {
				atcontext[k] = v
			}
		}
		modifiedAtContext = append(modifiedAtContext.([]interface{}), atcontext)
	}

	// Security can be a string of array of string
	if len(t.Security) == 1 {
		modifiedSecurity = t.Security[0]
	} else {
		modifiedSecurity = t.Security
	}

	return json.Marshal(&struct {
		*ThingOrigin
		ModifiedAtContext interface{} `json:"@context"`
		ModifiedSecurity  interface{} `json:"security"`
	}{
		ThingOrigin:       (*ThingOrigin)(t),
		ModifiedAtContext: modifiedAtContext,
		ModifiedSecurity:  modifiedSecurity,
	})
}

func (t *Thing) UnmarshalJSON(data []byte) error {
	type ThingOrigin Thing
	mt := &struct {
		*ThingOrigin
		ModifiedAtContext interface{} `json:"@context"`
		ModifiedSecurity  interface{} `json:"security"`
	}{
		ThingOrigin: (*ThingOrigin)(t),
	}
	if err := json.Unmarshal(data, &mt); err != nil {
		return err
	}
	// AtContext can be a string or an map of string

	switch mt.ModifiedAtContext.(type) {
	case string:
		t.AtContext = map[string]string{"": mt.ModifiedAtContext.(string)}
	case []interface{}:
		t.AtContext = map[string]string{}
		for _, v := range mt.ModifiedAtContext.([]interface{}) {
			switch v.(type) {
			case string:
				t.AtContext[v.(string)] = v.(string)
			case map[string]interface{}:
				for k, v := range v.(map[string]interface{}) {
					t.AtContext[k] = v.(string)
				}
			}
		}
	}

	// Security can be a string of array of string
	switch mt.ModifiedSecurity.(type) {
	case string:
		t.Security = []string{mt.ModifiedSecurity.(string)}
	case []string:
		t.Security = mt.ModifiedSecurity.([]string)
	}
	return nil
}

// AddProperty add property to a thing
func (t *Thing) AddProperty(property *interaction.Property) {
	if t == nil {
		zlog.Error().Msg("[thing:AddProperty] nil thing")
		return
	}

	t.MU.Lock()
	defer t.MU.Unlock()
	t.Properties[property.Key] = property
}

// AddAction add action to a thing
func (t *Thing) AddAction(action *interaction.Action) {
	if t == nil {
		zlog.Error().Msg("[thing:AddAction] nil thing")
		return
	}

	t.MU.Lock()
	defer t.MU.Unlock()
	t.Actions[action.Key] = action
}

// AddEvent add event to a thing
func (t *Thing) AddEvent(event *interaction.Event) {
	if t == nil {
		zlog.Error().Msg("[thing:AddEvent] nil thing")
		return
	}

	t.MU.Lock()
	defer t.MU.Unlock()
	t.Events[event.Key] = event
}
