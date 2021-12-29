package thing

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/project-eria/go-wot/interaction"
	"github.com/project-eria/go-wot/securityScheme"
	"github.com/rs/zerolog/log"
)

// Thing resource provides a Thing Description for a device.
// A Thing Resource is considered the root resource of a Thing.
type Thing struct {
	ID           string   `json:"id,omitempty"`           // (optional) Identifier of the Thing in form of a URI [RFC3986] (e.g., stable URI, temporary and mutable URI, URI with local IP address, URN, etc.).
	AtContext    string   `json:"@context"`               // (mandatory) JSON-LD keyword to define short-hand names called terms that are used throughout a TD document. anyURI or Array
	AtTypes      []string `json:"@type,omitempty"`        // (optional) JSON-LD keyword to label the object with semantic tags (or types).
	Title        string   `json:"title"`                  // (mandatory) Provides a human-readable title (e.g., display a text for UI representation) based on a default language.
	Titles       []string `json:"titles,omitempty"`       // (optional) Provides multi-language human-readable titles (e.g., display a text for UI representation in different languages).
	Description  string   `json:"description,omitempty"`  // (optional) Provides additional (human-readable) information based on a default language.
	Descriptions []string `json:"descriptions,omitempty"` // Can be used to support (human-readable) information in different languages. Also see MultiLanguage.
	// version	Provides version information.	optional	VersionInfo
	// created	Provides information when the TD instance was created.	optional	dateTime
	// modified	Provides information when the TD instance was last modified.	optional	dateTime
	// support	Provides information about the TD maintainer as URI scheme (e.g., mailto [RFC6068], tel [RFC3966], https).	optional	anyURI
	// base	Define the base URI that is used for all relative URI references throughout a TD document. In TD instances, all relative URIs are resolved relative to the base URI using the algorithm defined in [RFC3986].

	// base does not affect the URIs used in @context and the IRIs used within Linked Data [LINKED-DATA] graphs that are relevant when semantic processing is applied to TD instances.	optional	anyURI
	Properties map[string]*interaction.Property `json:"properties,omitempty"` // (optional) All Property-based Interaction Affordances of the Thing.
	Actions    map[string]*interaction.Action   `json:"actions,omitempty"`    // (optional) All Action-based Interaction Affordances of the Thing.
	// events	All Event-based Interaction Affordances of the Thing.	optional	Map of EventAffordance
	// links	Provides Web links to arbitrary resources that relate to the specified Thing Description.	optional	Array of Link
	// forms	Set of form hypermedia controls that describe how an operation can be performed. Forms are serializations of Protocol Bindings. In this version of TD, all operations that can be described at the Thing level are concerning how to interact with the Thing's Properties collectively at once.	optional	Array of Form
	Security            []string                                 `json:"security"`            // (mandatory) Set of security definition names, chosen from those defined in securityDefinitions. These must all be satisfied for access to resources.
	SecurityDefinitions map[string]securityScheme.SecurityScheme `json:"securityDefinitions"` // (mandatory) Set of named security configurations (definitions only). Not actually applied unless names are used in a security name-value pair.
	// profile	Indicates the WoT Profile mechanisms followed by this Thing Description and the corresponding Thing implementation.

	MU sync.RWMutex `json:"-"`
	//	server      *Server
	//	*thingWSHandler
}

// New thing construct
func New(urn string, title string, description string, types []string) (*Thing, error) {
	if urn == "" {
		return nil, errors.New("Thing URN can't be empty")
	}

	thing := Thing{
		AtContext:           "http://www.w3.org/ns/td",
		AtTypes:             types,
		ID:                  "urn:" + urn,
		Title:               title,
		Description:         description,
		Security:            []string{},
		SecurityDefinitions: make(map[string]securityScheme.SecurityScheme),
		Properties:          make(map[string]*interaction.Property),
		Actions:             make(map[string]*interaction.Action),
	}

	if thing.AtTypes == nil {
		thing.AtTypes = make([]string, 0)
	}
	//	thing.thingWSHandler = &thingWSHandler{webSocketConnections: make(map[string]*wsConnection)}
	return &thing, nil
}

func (t *Thing) AddSecurity(key string, definition securityScheme.SecurityScheme) {
	t.Security = append(t.Security, key)
	t.SecurityDefinitions[key] = definition
}

// Ref: http://choly.ca/post/go-json-marshalling/
func (t *Thing) MarshalJSON() ([]byte, error) {
	type ThingOrigin Thing

	var modifiedSecurity interface{}

	// Security can be a string of array of string
	if len(t.Security) == 1 {
		modifiedSecurity = t.Security[0]
	} else {
		modifiedSecurity = t.Security
	}

	return json.Marshal(&struct {
		*ThingOrigin
		ModifiedSecurity interface{} `json:"security"`
	}{
		ThingOrigin:      (*ThingOrigin)(t),
		ModifiedSecurity: modifiedSecurity,
	})
}

func (t *Thing) UnmarshalJSON(data []byte) error {
	type ThingOrigin Thing
	mt := &struct {
		*ThingOrigin
		ModifiedSecurity interface{} `json:"security"`
	}{
		ThingOrigin: (*ThingOrigin)(t),
	}
	if err := json.Unmarshal(data, &mt); err != nil {
		return err
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

// // AddType add a new type
// func (t *Thing) AddType(label string) {
// 	if t == nil {
// 		log.Error().Msg("[thing:AddType] nil thing")
// 	}

// 	t.mu.RLock()
// 	defer t.mu.RUnlock()

// 	if !find(t.types, label) {
// 		t.types = append(t.types, label)
// 	}
// }

// //SetContext set the thing context for the capabilities types
// func (t *Thing) SetContext(context string) {
// 	if t == nil {
// 		log.Error().Msg("[thing:SetContext] nil thing")
// 	}

// 	t.mu.Lock()
// 	defer t.mu.Unlock()

// 	t.context = context
// }

// func find(source []string, value string) bool {
// 	for _, item := range source {
// 		if item == value {
// 			return true
// 		}
// 	}
// 	return false
// }

// func (t *Thing) ref() string {
// 	if t == nil {
// 		log.Error().Msg("[thing:ref] nil thing")
// 		return ""
// 	}

// 	t.mu.RLock()
// 	defer t.mu.RUnlock()

// 	return t.Ref
// }

// func (t *Thing) href() string {
// 	if t == nil {
// 		log.Error().Msg("[thing:href] nil thing")
// 		return ""
// 	}
// 	return t.ref() + "/"
// }

// // actionsList returns a full list of actions for external request
// func (t *Thing) actionsList() map[string]interface{} {
// 	result := make(map[string]interface{})
// 	if t == nil {
// 		log.Error().Msg("[thing:actionsList] nil thing")
// 		return result
// 	}

// 	t.mu.RLock()
// 	defer t.mu.RUnlock()

// 	for name, actions := range t.actions {
// 		result[name] = actions.description()
// 	}
// 	return result
// }

// // linksList returns a list of links
// // A link object represents a link relation
// func (t *Thing) linksList(secure bool, host string) []map[string]string {
// 	result := []map[string]string{}
// 	if t == nil {
// 		log.Error().Msg("[thing:linksList] nil thing")
// 		return result
// 	}

// 	schemeHTTP := "http"
// 	schemeWS := "ws"
// 	if secure {
// 		schemeHTTP = "https"
// 		schemeWS = "wss"
// 	}

// 	t.mu.RLock()
// 	defer t.mu.RUnlock()

// 	for _, name := range []string{"properties", "actions", "events"} {
// 		result = append(result, map[string]string{
// 			"rel":  name,
// 			"href": t.Ref + "/" + name,
// 		})
// 	}
// 	result = append(result, map[string]string{
// 		"rel":       "alternate",
// 		"mediaType": "text/html",
// 		"href":      schemeHTTP + "://" + host + "/" + t.Ref,
// 	})
// 	result = append(result, map[string]string{
// 		"rel":  "alternate",
// 		"href": schemeWS + "://" + host + "/" + t.Ref,
// 	})
// 	return result
// }

// // processGetProperties returns a list of thing properties values
// func (t *Thing) processGetProperties() (map[string]interface{}, error) {
// 	if t == nil {
// 		log.Error().Msg("[thing:processGetProperties] nil thing")
// 		return nil, errors.New("thing can't be nil")
// 	}
// 	content := make(map[string]interface{})

// 	t.mu.RLock()
// 	defer t.mu.RUnlock()

// 	for name, property := range t.properties {
// 		content[name] = property.processGetValue()
// 	}
// 	return content, nil
// }

// // processSetProperties batch update values for a list of thing properties, from an external request
// func (t *Thing) processSetProperties(data map[string]interface{}) (map[string]interface{}, error) {
// 	if t == nil {
// 		log.Error().Msg("[thing:processSetProperties] nil thing")
// 		return nil, errors.New("thing can't be nil")
// 	}

// 	t.mu.RLock()
// 	defer t.mu.RUnlock()
// 	content := make(map[string]interface{})

// 	eventMessage := wsMessage{MessageType: "propertyStatus", Data: make(map[string]interface{})}
// 	// process the list of names, from the request
// 	for name, value := range data {
// 		log.Trace().Str("property", name).Msg("[thing:processSetProperties] Processing value")
// 		// check is the property exists
// 		if property, ok := t.properties[name]; ok {
// 			newValue, err := property.processSetValue(value)
// 			if err != nil {
// 				log.Error().Err(err).Msg("[thing:processSetProperties]")
// 				content[name] = map[string]string{"error": err.Error()}
// 				continue
// 			}
// 			if newValue != nil {
// 				log.Trace().Str("property", name).Msg("[thing:processSetProperties] Value Changed")
// 				content[name] = map[string]string{"response": "ok"}
// 				eventMessage.Data[name] = value
// 			} else {
// 				log.Trace().Str("property", name).Msg("[thing:processSetProperties] Value Unchanged")
// 				content[name] = map[string]string{"response": "unchanged"}
// 			}
// 			continue
// 		}
// 		log.Error().Str("property", name).Msg("[thing:processSetProperties] Unknown property")
// 		content[name] = map[string]string{"error": "Unknown property"}
// 	}
// 	// If at least one value has changed, we broadcast the event
// 	if len(eventMessage.Data) > 0 {
// 		t.processTxMsg(&eventMessage)
// 	}
// 	return content, nil
// }

// // LocalSetValues updates a property value, for local changes, without checking ReadOnly flag
// // Brodcast the new value if changed
// func (t *Thing) LocalSetValues(data map[string]interface{}) (map[string]interface{}, error) {
// 	if t == nil {
// 		log.Error().Msg("[thing:LocalSetValues] nil thing")
// 		return nil, errors.New("thing can't be nil")
// 	}

// 	t.mu.RLock()
// 	defer t.mu.RUnlock()

// 	content := make(map[string]interface{})
// 	eventMessage := wsMessage{MessageType: "propertyStatus", Data: make(map[string]interface{})}
// 	for name, value := range data {
// 		// check is the property exists
// 		if property, ok := t.properties[name]; ok {
// 			newValue, err := property.setValue(value)
// 			if err != nil {
// 				log.Error().Err(err).Msg("[thing:LocalSetValue]")
// 				content[name] = map[string]string{"error": err.Error()}
// 				continue
// 			}
// 			if newValue != nil {
// 				log.Trace().Str("property", name).Msg("[thing:LocalSetValue] Value Changed")
// 				content[name] = map[string]string{"response": "ok"}
// 				eventMessage.Data[name] = value
// 			} else {
// 				log.Trace().Str("property", name).Msg("[thing:LocalSetValue] Value Unchanged")
// 				content[name] = map[string]string{"response": "unchanged"}
// 			}
// 			continue
// 		}
// 		log.Error().Str("property", name).Msg("[thing:LocalSetValue] Unknown property")
// 		content[name] = map[string]string{"error": "Unknown property"}
// 	}
// 	// If at least one value has changed, we broadcast the event
// 	if len(eventMessage.Data) > 0 {
// 		t.processTxMsg(&eventMessage)
// 	}

// 	return content, nil
// }

/*
 * Properties
 */

// AddProperty add property to a thing
func (t *Thing) AddProperty(property *interaction.Property) {
	if t == nil {
		log.Error().Msg("[thing:AddProperty] nil thing")
		return
	}

	t.MU.Lock()
	defer t.MU.Unlock()
	t.Properties[property.Key] = property
}

// WriteProperty set the property value
func (t *Thing) WriteProperty(key string, value interface{}) {
	if t == nil {
		log.Error().Msg("[thing:WriteProperty] nil thing")
		return
	}

	t.MU.Lock()
	defer t.MU.Unlock()
	property := t.Properties[key]
	property.SetValue(value)
}

// // GetProperty get an existing thing propertie
// func (t *Thing) GetProperty(name string) *property.Property {
// 	if t == nil {
// 		log.Error().Msg("[thing:GetProperty] nil thing")
// 		return nil
// 	}

// 	t.mu.RLock()
// 	defer t.mu.RUnlock()

// 	if property, ok := t.properties[name]; ok {
// 		return property
// 	}
// 	return nil
// }

func (t *Thing) AddAction(action *interaction.Action) {
	if t == nil {
		log.Error().Msg("[thing:AddAction] nil thing")
		return
	}

	t.MU.Lock()
	defer t.MU.Unlock()
	t.Actions[action.Key] = action
}
