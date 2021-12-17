package interaction

import "github.com/project-eria/go-wot/form"

type Interaction struct {
	Key          string            `json:"-"`
	AtType       []string          `json:"@type,omitempty"`        // (optional) JSON-LD keyword to label the object with semantic tags (or types)
	Title        string            `json:"title,omitempty"`        // (optional) Provides a human-readable title (e.g., display a text for UI representation) based on a default language.
	Titles       map[string]string `json:"titles,omitempty"`       // (optional) Provides multi-language human-readable titles (e.g., display a text for UI representation in different languages).
	Description  string            `json:"description,omitempty"`  // (optional) Provides additional (human-readable) information based on a default language.
	Descriptions map[string]string `json:"descriptions,omitempty"` // (optional) Can be used to support (human-readable) information in different languages.
	Forms        []form.Form       `json:"forms"`                  // (mandatory) Set of form hypermedia controls that describe how an operation can be performed. Forms are serializations of Protocol Bindings.
	//	uriVariables	Define URI query template variables as collection based on DataSchema declarations. The individual variables DataSchema cannot be an ObjectSchema or an ArraySchema.	optional	Map of DataSchema
}

func (i *Interaction) AddHrefForm(url string, interactionForm form.Form) {
	interactionForm.Href = url + "/" + i.Key
	i.Forms = append(i.Forms, interactionForm)
}
