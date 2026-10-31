package auth

import (
	"fmt"
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

	// A document may exceed a token's size, up to MaxDocSize: a revocation
	// per terminal after « revoke every terminal ».
	d.Thing = "*"
	d.Revoked = nil
	for i := 0; i < 30; i++ {
		d.Revoked = append(d.Revoked, Revocation{Subject: fmt.Sprintf("term:%022d", i), Before: 100})
	}
	tok, _ = SignDoc(s.priv, s.pub.KID(), d)
	require.Greater(t, len(tok), MaxTokenSize)
	require.LessOrEqual(t, len(tok), MaxDocSize)
	got, err = ParseDoc(tok, s.lookup(false), "eria-maison", DocRevocation, "x")
	require.NoError(t, err)
	assert.Len(t, got.Revoked, 30)
	_, _, err = Parse(tok, s.lookup(false))
	assert.Error(t, err, "a token stays within MaxTokenSize")

	// A token is not a document.
	tok, _ = Sign(s.priv, s.pub.KID(), person())
	_, err = ParseDoc(tok, s.lookup(false), "eria-maison", DocRevocation, "x")
	assert.Error(t, err)
}
