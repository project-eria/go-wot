package interaction

import (
	"encoding/json"
	"strings"
)

type Form struct {
	Href        string `json:"href,omitempty"`        // (optional) Target IRI of a link or submission target of a form.
	ContentType string `json:"contentType,omitempty"` // (default: "application/json") Assign a content type based on a media type (e.g., text/plain) and potential parameters (e.g., charset=utf-8) for the media type [RFC2046].
	// contentCoding	Content coding values indicate an encoding transformation that has been or can be applied to a representation. Content codings are primarily used to allow a representation to be compressed or otherwise usefully transformed without losing the identity of its underlying media type and without loss of information. Examples of content coding include "gzip", "deflate", etc. .	optional	string
	// security	Set of security definition names, chosen from those defined in securityDefinitions. These must all be satisfied for access to resources.	optional	string or Array of string
	// scopes	Set of authorization scope identifiers provided as an array. These are provided in tokens returned by an authorization server and associated with forms in order to identify what resources a client may access and how. The values associated with a form should be chosen from those defined in an OAuth2SecurityScheme active on that form.	optional	string or Array of string
	// response	This optional term can be used if, e.g., the output communication metadata differ from input metadata (e.g., output contentType differ from the input contentType). The response name contains metadata that is only valid for the primary response messages.	optional	ExpectedResponse
	// additionalResponses	This optional term can be used if additional expected responses are possible, e.g. for error reporting. Each additional response needs to be distinguished from others in some way (for example, by specifying a protocol-specific error code), and may also have its own data schema.	optional	AdditionalExpectedResponse or Array of AdditionalExpectedResponse
	Subprotocol string   `json:"subprotocol,omitempty"` // (Optional) (e.g., longpoll, websub, or sse) Indicates the exact mechanism by which an interaction will be accomplished for a given protocol when there are multiple options. For example, for HTTP and Events, it indicates which of several available mechanisms should be used for asynchronous notifications such as long polling (longpoll), WebSub [websub] (websub), Server-Sent Events (sse) [html] (also known as EventSource). Please note that there is no restriction on the subprotocol selection and other mechanisms can also be announced by this subprotocol term.
	Op          []string `json:"op,omitempty"`          // (optional) Indicates the semantic intention of performing the operation(s) described by the form. For example, the Property interaction allows get and set operations. The protocol binding may contain a form for the get operation and a different form for the set operation. The op attribute indicates which form is for which and allows the client to select the correct form for the operation required. op can be assigned one or more interaction verb(s) each representing a semantic intention of an operation. Array of string (one of readproperty, writeproperty, observeproperty, unobserveproperty, invokeaction, subscribeevent, unsubscribeevent, readallproperties, writeallproperties, readmultipleproperties, writemultipleproperties, observeallproperties, or unobserveallproperties)

	Supplement map[string]interface{} `json:"-"`

	UrlBuilder func(string, bool) string `json:"-"` // host, secure
}

// Ref: http://choly.ca/post/go-json-marshalling/
func (f *Form) MarshalJSON() ([]byte, error) {
	type FormOrigin Form

	b1, err := json.Marshal((*FormOrigin)(f))
	if err != nil {
		return nil, err
	}

	// Supplement
	if len(f.Supplement) > 0 {
		b2, err := json.Marshal(f.Supplement)
		if err != nil {
			return nil, err
		}

		s1 := string(b1[:len(b1)-1])
		s2 := string(b2[1:])
		return []byte(s1 + ", " + s2), nil
	} else {
		return b1, nil
	}
}

func (f *Form) UnmarshalJSON(data []byte) error {
	type FormOrigin Form
	fo := &struct {
		*FormOrigin
	}{
		FormOrigin: (*FormOrigin)(f),
	}
	if err := json.Unmarshal(data, &fo); err != nil {
		return err
	}

	// Supplement
	fs := map[string]interface{}{}
	f.Supplement = map[string]interface{}{}
	if err := json.Unmarshal(data, &fs); err != nil {
		return err
	}
	for key, value := range fs {
		if strings.Contains(key, ":") {
			f.Supplement[key] = value
		}
	}
	return nil
}
