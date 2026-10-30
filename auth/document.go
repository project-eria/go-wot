package auth

import (
	"crypto/ed25519"
	"encoding/json"
	"errors"
)

// Signed is an object signed by an app, terminal, relay or deployment key
// (AUTH-DESIGN §2.7): the signature covers "eria-<type>-v1" ‖ payload, the
// exact JSON bytes, so no canonical form is ever needed.
type Signed struct {
	Payload string `json:"payload"` // base64url of the JSON bytes
	Sig     string `json:"sig"`     // base64url
}

// SignDocument marshals v and signs it for the given domain prefix.
func SignDocument(priv ed25519.PrivateKey, prefix string, v any) (Signed, error) {
	payload, err := json.Marshal(v)
	if err != nil {
		return Signed{}, err
	}
	sig := ed25519.Sign(priv, append([]byte(prefix), payload...))
	return Signed{Payload: b64.EncodeToString(payload), Sig: b64.EncodeToString(sig)}, nil
}

// Open verifies d with any of keys for prefix, then decodes it strictly into v.
func (d Signed) Open(prefix string, keys []PublicKey, v any) error {
	payload, err1 := b64.DecodeString(d.Payload)
	sig, err2 := b64.DecodeString(d.Sig)
	if err1 != nil || err2 != nil {
		return errors.New("bad base64url")
	}
	msg := append([]byte(prefix), payload...)
	for _, k := range keys {
		if k.Verify(msg, sig) {
			return strictDecode(payload, 3, v)
		}
	}
	return errors.New("bad signature")
}

// KeysetPrefix is the domain of the auth keyset, signed by a deployment key.
const KeysetPrefix = "eria-keyset-v1"

// Key usages in a keyset: a legacy key is accepted for terminal tokens only.
const (
	UsageSign       = "auth-sign"
	UsageSignLegacy = "auth-sign-legacy"
)

// Keyset lists the auth public keys (AUTH-DESIGN §1.1). Generation is the
// time in ms; a verifier keeps only a higher generation.
type Keyset struct {
	Generation int64       `json:"generation"`
	Keys       []KeysetKey `json:"keys"`
	Issued     int64       `json:"issued"`
}

type KeysetKey struct {
	Kid    string `json:"kid"`
	Pubkey string `json:"pubkey"` // base64url, Ed25519
	Usage  string `json:"usage"`
}

// Lookup turns a verified keyset into the KeyLookup of Parse.
func (ks Keyset) Lookup() (KeyLookup, error) {
	type entry struct {
		key    PublicKey
		legacy bool
	}
	keys := map[string]entry{}
	for _, k := range ks.Keys {
		raw, err := b64.DecodeString(k.Pubkey)
		if err != nil {
			return nil, err
		}
		pub, err := ParsePublicKey(Ed25519, raw)
		if err != nil {
			return nil, err
		}
		if pub.KID() != k.Kid || (k.Usage != UsageSign && k.Usage != UsageSignLegacy) {
			return nil, errors.New("bad keyset entry")
		}
		keys[k.Kid] = entry{pub, k.Usage == UsageSignLegacy}
	}
	return func(kid string) (PublicKey, bool, bool) {
		e, ok := keys[kid]
		return e.key, e.legacy, ok
	}, nil
}
