package auth

import (
	"bytes"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

// MaxTokenSize bounds a token, as the ESP32 do (AUTH-DESIGN §2.7).
const MaxTokenSize = 2048

// Kind is one of the three token kinds (AUTH-DESIGN §2.7).
type Kind string

const (
	KindUser    Kind = "user"
	KindService Kind = "service"
	KindConfirm Kind = "confirm"
)

// Claims of an auth-signed token. Optional claims are pointers or omitted
// when empty; lvl is a pointer since an absent lvl must never read as 0
// (admin).
type Claims struct {
	Iss      string   `json:"iss"`
	Aud      string   `json:"aud"`
	Kind     Kind     `json:"kind"`
	Sub      string   `json:"sub,omitempty"`
	Term     string   `json:"term,omitempty"`
	Lvl      *int     `json:"lvl,omitempty"`
	Zone     string   `json:"zone,omitempty"`
	Out      *bool    `json:"out,omitempty"`
	RoomOnly bool     `json:"room_only,omitempty"`
	Grants   []string `json:"grants,omitempty"`
	Act      string   `json:"act,omitempty"`
	Iat      int64    `json:"iat"`
	Exp      int64    `json:"exp"`
	Jti      string   `json:"jti"`
}

// header is the only header accepted: exactly these three fields.
type header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
	Kid string `json:"kid"`
}

// TokenError tells which step of the decision refused the token
// (1: structure, 2: signature).
type TokenError struct {
	Step   int
	Reason string
}

func (e *TokenError) Error() string { return fmt.Sprintf("deny(%d): %s", e.Step, e.Reason) }

func structureErr(reason string) error { return &TokenError{Step: 1, Reason: reason} }
func signatureErr(reason string) error { return &TokenError{Step: 2, Reason: reason} }

// b64 decodes base64url without padding and refuses non-canonical input.
var b64 = base64.RawURLEncoding.Strict()

// Sign issues a token signed by the auth key priv, whose kid is kid.
func Sign(priv ed25519.PrivateKey, kid string, c Claims) (string, error) {
	return signJWS(priv, kid, c)
}

func signJWS(priv ed25519.PrivateKey, kid string, v any) (string, error) {
	h, err := json.Marshal(header{Alg: "EdDSA", Typ: "JWT", Kid: kid})
	if err != nil {
		return "", err
	}
	p, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	signed := b64.EncodeToString(h) + "." + b64.EncodeToString(p)
	sig := ed25519.Sign(priv, []byte(signed))
	return signed + "." + b64.EncodeToString(sig), nil
}

// KeyLookup returns the accepted auth key for kid (compared by strict
// equality, never used as a path) and whether it is a legacy key.
type KeyLookup func(kid string) (key PublicKey, legacy bool, ok bool)

// Parse checks the structure (step 1) and the signature (step 2) of a token,
// then returns its claims, still to be checked (step 3, see Claims.Check).
func Parse(token string, keys KeyLookup) (c Claims, legacy bool, err error) {
	legacy, err = parseJWS(token, keys, 2, &c)
	return c, legacy, err
}

// parseJWS checks a compact JWS signed by auth and decodes its payload
// strictly into v, nested at most maxDepth.
func parseJWS(token string, keys KeyLookup, maxDepth int, v any) (legacy bool, err error) {
	if token == "" || len(token) > MaxTokenSize {
		return false, structureErr("missing or oversized token")
	}
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return false, structureErr("not three parts")
	}
	rawHeader, err1 := b64.DecodeString(parts[0])
	rawClaims, err2 := b64.DecodeString(parts[1])
	sig, err3 := b64.DecodeString(parts[2])
	if err1 != nil || err2 != nil || err3 != nil {
		return false, structureErr("bad base64url")
	}
	var h header
	if err := strictDecode(rawHeader, 1, &h); err != nil {
		return false, structureErr("header: " + err.Error())
	}
	if h.Alg != "EdDSA" || h.Typ != "JWT" || h.Kid == "" {
		return false, structureErr("header not pinned")
	}
	if err := strictDecode(rawClaims, maxDepth, v); err != nil {
		return false, structureErr("claims: " + err.Error())
	}

	key, legacy, ok := keys(h.Kid)
	if !ok || key.Alg() != Ed25519 {
		return false, signatureErr("unknown kid")
	}
	if !key.Verify([]byte(parts[0]+"."+parts[1]), sig) {
		return false, signatureErr("bad signature")
	}
	return legacy, nil
}

// strictDecode decodes a JSON object into v, refusing unknown fields,
// duplicate keys, nesting deeper than maxDepth and trailing data.
func strictDecode(data []byte, maxDepth int, v any) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	if t, err := dec.Token(); err != nil || t != json.Delim('{') {
		return errors.New("not an object")
	}
	if err := walkObject(dec, 1, maxDepth); err != nil {
		return err
	}
	if _, err := dec.Token(); err != io.EOF {
		return errors.New("trailing data")
	}
	strict := json.NewDecoder(bytes.NewReader(data))
	strict.DisallowUnknownFields()
	return strict.Decode(v)
}

// walkObject reads the members of an object whose '{' was just read.
func walkObject(dec *json.Decoder, depth, maxDepth int) error {
	seen := map[string]bool{}
	for dec.More() {
		t, err := dec.Token()
		if err != nil {
			return err
		}
		key := t.(string)
		if seen[key] {
			return errors.New("duplicate key " + key)
		}
		seen[key] = true
		if err := walkValue(dec, depth, maxDepth); err != nil {
			return err
		}
	}
	_, err := dec.Token() // '}'
	return err
}

func walkValue(dec *json.Decoder, depth, maxDepth int) error {
	t, err := dec.Token()
	if err != nil {
		return err
	}
	d, ok := t.(json.Delim)
	if !ok {
		return nil
	}
	if depth+1 > maxDepth {
		return errors.New("too deep")
	}
	if d == '{' {
		return walkObject(dec, depth+1, maxDepth)
	}
	for dec.More() {
		if err := walkValue(dec, depth+1, maxDepth); err != nil {
			return err
		}
	}
	_, err = dec.Token() // ']'
	return err
}
