package authz

import (
	"crypto/ed25519"
	"crypto/rand"
	"testing"
	"time"

	"github.com/project-eria/go-wot/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthenticate(t *testing.T) {
	pubRaw, priv, _ := ed25519.GenerateKey(rand.Reader)
	pub, err := auth.ParsePublicKey(auth.Ed25519, pubRaw)
	require.NoError(t, err)
	legacy := false
	keys := func(kid string) (auth.PublicKey, bool, bool) { return pub, legacy, kid == pub.KID() }
	now := time.Unix(1_800_000_000, 0)
	lvl := func(n int) *int { return &n }

	token := func(c auth.Claims) string {
		c.Iss, c.Aud, c.Iat, c.Jti = "auth", "eria-maison", now.Unix(), "j"
		if c.Exp == 0 {
			c.Exp = now.Add(time.Minute).Unix()
		}
		tok, err := auth.Sign(priv, pub.KID(), c)
		require.NoError(t, err)
		return tok
	}
	terminal := token(auth.Claims{Kind: auth.KindUser, Term: "salon", Lvl: lvl(9)})
	person := token(auth.Claims{Kind: auth.KindUser, Sub: "profil:lea", Term: "salon", Lvl: lvl(2)})
	admin := token(auth.Claims{Kind: auth.KindUser, Sub: "profil:admin", Term: "pc", Lvl: lvl(0)})
	service := token(auth.Claims{Kind: auth.KindService, Lvl: lvl(9)})
	confirm := token(auth.Claims{Kind: auth.KindConfirm, Sub: "profil:lea", Term: "salon", Act: "a"})

	run := func(v Verifier, m Mode, tok string) Decision {
		return Authenticate(Request{Mode: m, Verifier: v, Token: tok, Keys: keys, Aud: "eria-maison", Now: now})
	}
	allow := func(v Verifier, tok string, want auth.Profile) {
		d := run(v, Normal, tok)
		require.True(t, d.Allow, d.Reason)
		assert.Equal(t, want, d.Profile)
	}
	refuse := func(v Verifier, m Mode, tok string, step int) {
		d := run(v, m, tok)
		assert.False(t, d.Allow)
		assert.Equal(t, step, d.Step, d.Reason)
	}

	allow(Thing, terminal, auth.ProfileTerminal)
	allow(Thing, person, auth.ProfilePerson)
	allow(Thing, service, auth.ProfileService)
	allow(Relay, person, auth.ProfilePerson)
	allow(AuthService, admin, auth.ProfileAdmin)
	allow(Master, admin, auth.ProfileAdmin)

	refuse(Thing, Normal, admin, 3)   // an admin token never reaches a Thing
	refuse(Relay, Normal, admin, 3)   // nor a relay
	refuse(Master, Normal, person, 3) // the master takes the admin token only
	refuse(Thing, Normal, confirm, 3) // confirm is never a primary token
	refuse(Thing, Offline, terminal, 0)
	refuse(Thing, Normal, "x", 1)
	refuse(Thing, Minimal, person, 3)

	d := run(Thing, Minimal, terminal)
	assert.True(t, d.Allow, "a terminal token in minimal mode")

	legacy = true
	refuse(Thing, Normal, person, 3)
	assert.True(t, run(Thing, Normal, terminal).Allow, "a terminal token with a legacy key")
}
