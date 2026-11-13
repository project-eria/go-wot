package auth

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
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
	payload, err := d.Verify(prefix, keys)
	if err != nil {
		return err
	}
	return strictDecode(payload, 3, v)
}

// Verify checks d with any of keys for prefix and returns its payload, for
// a document too deep for Open (a TD).
func (d Signed) Verify(prefix string, keys []PublicKey) ([]byte, error) {
	payload, err1 := b64.DecodeString(d.Payload)
	sig, err2 := b64.DecodeString(d.Sig)
	if err1 != nil || err2 != nil {
		return nil, errors.New("bad base64url")
	}
	msg := append([]byte(prefix), payload...)
	for _, k := range keys {
		if k.Verify(msg, sig) {
			return payload, nil
		}
	}
	return nil, errors.New("bad signature")
}

// KeysetPrefix is the domain of the auth keyset, signed by a deployment key.
const KeysetPrefix = "eria-keyset-v1"

// CSRPrefix is the domain of auth's signature over its TLS certificate
// request: in a guided restore the coordinator signs only a CSR the verified
// auth key made (AUTH-DESIGN §4.2).
const CSRPrefix = "eria-auth-csr-v1"

// UsageSign is the usage of auth's key in a keyset. A key of another usage
// (auth-sign-legacy, removed on 2026-11-11) is ignored: a keyset kept from
// before still opens, without that key.
const UsageSign = "auth-sign"

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
	keys := map[string]PublicKey{}
	for _, k := range ks.Keys {
		raw, err := b64.DecodeString(k.Pubkey)
		if err != nil {
			return nil, err
		}
		pub, err := ParsePublicKey(Ed25519, raw)
		if err != nil {
			return nil, err
		}
		if pub.KID() != k.Kid {
			return nil, errors.New("bad keyset entry")
		}
		if k.Usage == UsageSign {
			keys[k.Kid] = pub
		}
	}
	return func(kid string) (PublicKey, bool) {
		pub, ok := keys[kid]
		return pub, ok
	}, nil
}

// DeploymentKeys returns the Ed25519 keys of the deployment root
// certificates (eria-coordinator tls-root): they sign the auth keyset.
func DeploymentKeys(rootsPEM []byte) ([]PublicKey, error) {
	var keys []PublicKey
	for rest := rootsPEM; ; {
		var block *pem.Block
		if block, rest = pem.Decode(rest); block == nil {
			break
		}
		c, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, err
		}
		if k, ok := c.PublicKey.(ed25519.PublicKey); ok {
			pub, err := ParsePublicKey(Ed25519, k)
			if err != nil {
				return nil, err
			}
			keys = append(keys, pub)
		}
	}
	if len(keys) == 0 {
		return nil, errors.New("no deployment root")
	}
	return keys, nil
}

// OpenKeyset checks a keyset signed by a deployment key (JSON of Signed) and
// returns it with its KeyLookup.
func OpenKeyset(signed []byte, deploy []PublicKey) (Keyset, KeyLookup, error) {
	var d Signed
	var ks Keyset
	if err := json.Unmarshal(signed, &d); err != nil {
		return ks, nil, err
	}
	if err := d.Open(KeysetPrefix, deploy, &ks); err != nil {
		return ks, nil, err
	}
	keys, err := ks.Lookup()
	return ks, keys, err
}
