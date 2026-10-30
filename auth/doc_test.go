package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRevocationDoc(t *testing.T) {
	s := newSigner(t)
	d := Doc{Iss: "auth", Aud: "eria-maison", Doc: DocRevocation, Thing: "*", Version: 1_800_000_000_000,
		Revoked: []Revocation{{Subject: "term:salon", Before: 100}}}
	tok, err := SignDoc(s.priv, s.pub.KID(), d)
	require.NoError(t, err)

	got, err := ParseDoc(tok, s.lookup(false), "eria-maison", DocRevocation, "urn:eria:garage")
	require.NoError(t, err, "* is for every recipient")
	assert.True(t, got.Revokes(Claims{Term: "salon", Iat: 99}))
	assert.False(t, got.Revokes(Claims{Term: "salon", Iat: 100}))
	assert.False(t, got.Revokes(Claims{Term: "cuisine", Iat: 99}))

	_, err = ParseDoc(tok, s.lookup(true), "eria-maison", DocRevocation, "x")
	assert.Error(t, err, "never signed by a legacy key")
	_, err = ParseDoc(tok, s.lookup(false), "eria-maison", DocHeartbeat, "x")
	assert.Error(t, err, "another kind of document")

	d.Thing = "urn:eria:garage"
	tok, _ = SignDoc(s.priv, s.pub.KID(), d)
	_, err = ParseDoc(tok, s.lookup(false), "eria-maison", DocRevocation, "urn:eria:portail")
	assert.Error(t, err, "document for another recipient")

	// A token is not a document.
	tok, _ = Sign(s.priv, s.pub.KID(), person())
	_, err = ParseDoc(tok, s.lookup(false), "eria-maison", DocRevocation, "x")
	assert.Error(t, err)
}
