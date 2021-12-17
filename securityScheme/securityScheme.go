package securityScheme

type SecurityScheme interface{}

type Security struct {
	AtType       []string `json:"@type,omitempty"`        // (optional) JSON-LD keyword to label the object with semantic tags (or types).
	Description  string   `json:"description,omitempty"`  // (optional) Provides additional (human-readable) information based on a default language.
	Descriptions []string `json:"descriptions,omitempty"` // (optional) Can be used to support (human-readable) information in different languages.
	Proxy        string   `json:"proxy,omitempty"`        // (optional) URI of the proxy server this security configuration provides access to. If not given, the corresponding security configuration is for the endpoint.
	Scheme       string   `json:"scheme,omitempty"`       // (optional)	Identification of the security mechanism being configured. Any type (one of nosec, combo, basic, digest, bearer, psk, oauth2, or apikey)
}
