package auth

import (
	"crypto/ed25519"
	"crypto/rand"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type signer struct {
	priv ed25519.PrivateKey
	pub  PublicKey
}

func newSigner(t *testing.T) signer {
	pubRaw, priv, _ := ed25519.GenerateKey(rand.Reader)
	pub, err := ParsePublicKey(Ed25519, pubRaw)
	require.NoError(t, err)
	return signer{priv, pub}
}

func (s signer) lookup() KeyLookup {
	return func(kid string) (PublicKey, bool) {
		return s.pub, kid == s.pub.KID()
	}
}

func lvl(n int) *int { return &n }

var now = time.Unix(1_800_000_000, 0)

func person() Claims {
	return Claims{Iss: "auth", Aud: "eria-maison", Kind: KindUser, Sub: "profil:lea",
		Term: "telephone-lea", Lvl: lvl(2), Grants: []string{"room:piscine:*"},
		Iat: now.Unix(), Exp: now.Add(5 * time.Minute).Unix(), Jti: "8f3c"}
}

// raw builds a token from hand-written header and claims, signed by s.
func (s signer) raw(header, claims string) string {
	signed := b64.EncodeToString([]byte(header)) + "." + b64.EncodeToString([]byte(claims))
	return signed + "." + b64.EncodeToString(ed25519.Sign(s.priv, []byte(signed)))
}

func step(t *testing.T, err error) int {
	var e *TokenError
	require.ErrorAs(t, err, &e)
	return e.Step
}

func TestParseRoundTrip(t *testing.T) {
	s := newSigner(t)
	tok, err := Sign(s.priv, s.pub.KID(), person())
	require.NoError(t, err)
	c, err := Parse(tok, s.lookup())
	require.NoError(t, err)
	assert.Equal(t, person(), c)
}

func TestParseStructure(t *testing.T) {
	s := newSigner(t)
	kid := s.pub.KID()
	hdr := `{"alg":"EdDSA","typ":"JWT","kid":"` + kid + `"}`
	claims := `{"iss":"auth","aud":"a","kind":"user","term":"t","lvl":9,"iat":1,"exp":2,"jti":"j"}`

	_, err := Parse(s.raw(hdr, claims), s.lookup())
	require.NoError(t, err)

	cases := map[string]string{
		"extra header field": s.raw(`{"alg":"EdDSA","typ":"JWT","kid":"`+kid+`","jku":"x"}`, claims),
		"other alg":          s.raw(`{"alg":"none","typ":"JWT","kid":"`+kid+`"}`, claims),
		"duplicate key":      s.raw(hdr, `{"iss":"auth","iss":"x","aud":"a","kind":"user","iat":1,"exp":2,"jti":"j"}`),
		"unknown claim":      s.raw(hdr, `{"iss":"auth","aud":"a","kind":"user","iat":1,"exp":2,"jti":"j","admin":true}`),
		"too deep":           s.raw(hdr, `{"iss":"auth","aud":"a","kind":"user","grants":[["x"]],"iat":1,"exp":2,"jti":"j"}`),
		"trailing data":      s.raw(hdr, claims+`{}`),
		"two parts":          "a.b",
		"empty":              "",
		"oversized":          strings.Repeat("a", MaxTokenSize+1),
		"padding":            s.raw(hdr, claims) + "=",
	}
	for name, tok := range cases {
		_, err := Parse(tok, s.lookup())
		assert.Equal(t, 1, step(t, err), name)
	}

	// Non-canonical base64url: a set unused trailing bit.
	tok := s.raw(hdr, claims)
	last := tok[len(tok)-1]
	alt := strings.IndexByte(b64Alphabet, last) ^ 1
	_, err = Parse(tok[:len(tok)-1]+string(b64Alphabet[alt]), s.lookup())
	assert.Error(t, err)
}

const b64Alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"

func TestParseSignature(t *testing.T) {
	s, other := newSigner(t), newSigner(t)
	tok, _ := Sign(s.priv, s.pub.KID(), person())

	_, err := Parse(tok, other.lookup())
	assert.Equal(t, 2, step(t, err), "unknown kid")

	forged, _ := Sign(other.priv, s.pub.KID(), person())
	_, err = Parse(forged, s.lookup())
	assert.Equal(t, 2, step(t, err), "bad signature")

	_, p256 := newP256(t)
	_, err = Parse(tok, func(string) (PublicKey, bool) { return p256, true })
	assert.Equal(t, 2, step(t, err), "tokens are never verified with P-256")
}

func TestClaimsCheck(t *testing.T) {
	ok := func(c Claims) Profile {
		p, err := c.Check(now, "eria-maison")
		require.NoError(t, err)
		return p
	}
	refused := func(c Claims, why string) {
		_, err := c.Check(now, "eria-maison")
		assert.Error(t, err, why)
	}

	assert.Equal(t, ProfilePerson, ok(person()))

	c := person()
	c.Sub, c.Lvl, c.Grants, c.Exp = "", lvl(9), nil, now.Add(72*time.Hour).Unix()
	assert.Equal(t, ProfileTerminal, ok(c))

	c = person()
	c.Lvl, c.Grants = lvl(0), nil
	assert.Equal(t, ProfileAdmin, ok(c))

	c = Claims{Iss: "auth", Aud: "eria-maison", Kind: KindService, Sub: "app:a1", Lvl: lvl(9),
		Iat: now.Unix(), Exp: now.Add(5 * time.Minute).Unix(), Jti: "j"}
	assert.Equal(t, ProfileService, ok(c))
	for _, sub := range []string{"", "app:", "a1"} {
		c.Sub = sub
		refused(c, "a service token names its app: "+sub)
	}

	c = Claims{Iss: "auth", Aud: "eria-maison", Kind: KindConfirm, Sub: "s", Term: "t", Act: "a",
		Iat: now.Unix(), Exp: now.Add(time.Minute).Unix(), Jti: "j"}
	assert.Equal(t, ProfileConfirm, ok(c))

	c = person()
	c.Lvl = nil
	refused(c, "absent lvl never reads as admin")
	c = person()
	c.Lvl = lvl(0)
	refused(c, "admin token with grants")
	c = person()
	c.Kind = "root"
	refused(c, "unknown kind")
	c = person()
	c.Act = "a"
	refused(c, "act on a user token")
	c = person()
	c.Exp = now.Add(6 * time.Minute).Unix()
	refused(c, "person token over 5 min")
	c = person()
	c.Exp = now.Add(-31 * time.Second).Unix()
	c.Iat = c.Exp - 60
	refused(c, "expired")
	c = person()
	c.Iat = now.Add(31 * time.Second).Unix()
	c.Exp = c.Iat + 60
	refused(c, "issued in the future")
	c = person()
	c.Aud = "other-house"
	refused(c, "other deployment")
	c = person()
	c.Grants = make([]string, MaxGrants+1)
	refused(c, "too many grants")
	c = person()
	c.Sub = strings.Repeat("x", MaxClaimSize+1)
	refused(c, "claim too long")
}

func TestValidGrants(t *testing.T) {
	assert.NoError(t, ValidGrants([]string{"room:piscine:*", "thing:urn:eria:garage:1:invoke", "thing:x:read", "room:salon:write"}))
	assert.NoError(t, ValidGrants(nil))
	for _, g := range []string{"", "room:piscine", "room::read", "house:x:read", "thing:x:delete", "thing:x:", ":x:read", "room:" + strings.Repeat("a", MaxClaimSize) + ":read"} {
		assert.Error(t, ValidGrants([]string{g}), g)
	}
	assert.Error(t, ValidGrants(make([]string, MaxGrants+1)))
}

// A keyset kept from before 2026-11-11 opens, without its legacy key.
func TestKeysetIgnoresLegacy(t *testing.T) {
	cur, old := newSigner(t), newSigner(t)
	entry := func(s signer, usage string) KeysetKey {
		return KeysetKey{Kid: s.pub.KID(), Pubkey: b64.EncodeToString(s.priv.Public().(ed25519.PublicKey)), Usage: usage}
	}
	lookup, err := Keyset{Keys: []KeysetKey{entry(cur, UsageSign), entry(old, "auth-sign-legacy")}}.Lookup()
	require.NoError(t, err)
	_, ok := lookup(cur.pub.KID())
	assert.True(t, ok)
	_, ok = lookup(old.pub.KID())
	assert.False(t, ok, "the legacy key signs nothing")
}
