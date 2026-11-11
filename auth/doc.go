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
	Table   Table        `json:"table,omitempty"` // doc table: the effective table
	Zones   []string     `json:"zones,omitempty"` // doc table: the Thing's rooms
	// doc heartbeat (§4.3): the Thing's nonce, and the versions auth holds
	// for it (the Thing pulls those it lacks).
	Nonce    string    `json:"nonce,omitempty"`
	Versions *Versions `json:"versions,omitempty"`
}

// Versions are the versions a heartbeat announces: the table of each Thing
// asked for (0: none, not adopted) and the revocation list.
type Versions struct {
	Tables     map[string]int64 `json:"tables"`
	Revocation int64            `json:"revocation"`
}

// None forbids a function to everyone. Requirements go None < 1 < … < 9
// (§2.1): 0 is never a level presented to a Thing, so it stands for None.
const None = 0

// Confirm is confirm_when (§2.4), for write and invoke only.
type Confirm string

const (
	Never       Confirm = "never"
	OutsideHome Confirm = "outside"
	OutsideRoom Confirm = "remote"
	Always      Confirm = "always"
)

// Rule is a table entry: the level required per locality (local, remote,
// outside) and when a write or invoke needs a confirmation. An empty
// Confirm is never.
type Rule struct {
	Levels  [3]int  `json:"l"`
	Confirm Confirm `json:"c,omitempty"`
}

// Table is the effective table of a Thing, keyed by "<kind>/<name>/<op>"
// (authz.Function.Key). A function absent from it is None.
type Table map[string]Rule

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
	legacy, err := parseJWS(token, keys, MaxDocSize, 4, &d)
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
// falls under the list. "term:*" is every terminal (« revoke every
// terminal », §1.2).
func (d Doc) Revokes(c Claims) bool {
	for _, r := range d.Revoked {
		term := c.Term != "" && (r.Subject == "term:"+c.Term || r.Subject == "term:*")
		if c.Iat < r.Before && (term || r.Subject == "sub:"+c.Sub && c.Sub != "") {
			return true
		}
	}
	return false
}
