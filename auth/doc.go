package auth

import (
	"crypto/ed25519"
	"errors"
)

// Documents signed by auth (AUTH-DESIGN §2.7, F011): JWS with the same
// pinned header and strict parser as the tokens, a doc claim and the id of
// the Thing they are for ("*" for revocations and heartbeats).

const (
	DocRevocation = "revocation"
	DocHeartbeat  = "heartbeat"
	DocTable      = "table"
)

// Doc is a signed document. Version is the time in ms: a verifier keeps
// only a higher version per (doc, thing).
type Doc struct {
	Iss     string       `json:"iss"`
	Aud     string       `json:"aud"`
	Doc     string       `json:"doc"`
	Thing   string       `json:"thing"`
	Version int64        `json:"version"`
	Issued  int64        `json:"issued"`
	Revoked []Revocation `json:"revoked,omitempty"`
}

// Revocation makes the tokens of subject (term:<id>, sub:<profile>,
// service:<kid>) issued before Before (unix seconds) fall.
type Revocation struct {
	Subject string `json:"subject"`
	Before  int64  `json:"before"`
}

// SignDoc signs d with the auth key.
func SignDoc(priv ed25519.PrivateKey, kid string, d Doc) (string, error) {
	return signJWS(priv, kid, d)
}

// ParseDoc checks a document of the kind doc for the expected recipient
// thing (a Thing passes its own id, a relay the id of the Thing it observes).
// A document for another recipient, or signed by a legacy key, is refused.
func ParseDoc(token string, keys KeyLookup, aud, doc, thing string) (Doc, error) {
	var d Doc
	legacy, err := parseJWS(token, keys, 3, &d)
	if err != nil {
		return d, err
	}
	if legacy || d.Iss != "auth" || d.Aud != aud || d.Doc != doc || d.Version <= 0 {
		return d, errors.New("not the expected document")
	}
	if d.Thing != thing && d.Thing != "*" {
		return d, errors.New("document for another recipient")
	}
	return d, nil
}

// Revokes tells whether a token issued at iat for term, sub or service kid
// falls under the list.
func (d Doc) Revokes(c Claims) bool {
	for _, r := range d.Revoked {
		if c.Iat < r.Before && (r.Subject == "term:"+c.Term && c.Term != "" || r.Subject == "sub:"+c.Sub && c.Sub != "") {
			return true
		}
	}
	return false
}
