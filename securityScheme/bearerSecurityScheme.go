package securityScheme

// BearerSecurity is a bearer token scheme (W3C WoT TD §5.3.3.5).
type BearerSecurity struct {
	Security
	In     string `json:"in,omitempty"`     // header, query, body, cookie
	Name   string `json:"name,omitempty"`   // name of the header
	Format string `json:"format,omitempty"` // jwt, cwt, jwe, jws
	Alg    string `json:"alg,omitempty"`    // ES256, EdDSA…
}

// NewEriaBearer is the scheme of the Things guarded by eria-auth tokens
// (AUTH-DESIGN §2.7): a JWT signed in EdDSA, in the Authorization header.
func NewEriaBearer() BearerSecurity {
	return BearerSecurity{Security: Security{Scheme: "bearer"}, In: "header", Name: "Authorization", Format: "jwt", Alg: "EdDSA"}
}
