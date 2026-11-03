package guard

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/project-eria/go-wot/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Step 5c: the settings the agent writes beside the app key activate the
// guard; the agent's own Thing is the deploy path.
func TestSettings(t *testing.T) {
	deployPub, deployPriv, _ := ed25519.GenerateKey(rand.Reader)
	authPub, _, _ := ed25519.GenerateKey(rand.Reader)
	pub, _ := auth.ParsePublicKey(auth.Ed25519, authPub)
	signed, err := auth.SignDocument(deployPriv, auth.KeysetPrefix, auth.Keyset{Generation: 3, Issued: 1,
		Keys: []auth.KeysetKey{{Kid: pub.KID(), Pubkey: base64.RawURLEncoding.EncodeToString(authPub), Usage: auth.UsageSign}}})
	require.NoError(t, err)
	keyset, _ := json.Marshal(signed)

	dir := t.TempDir()
	keyFile := filepath.Join(dir, "app.key")
	_, err = ReadSettings(keyFile)
	assert.ErrorIs(t, err, os.ErrNotExist, "not activated")
	s := Settings{Keyset: string(keyset), Channel: "http://auth:8444", Deployment: "maison", DeployKey: hex.EncodeToString(deployPub)}
	data, _ := json.Marshal(s)
	require.NoError(t, os.WriteFile(filepath.Join(dir, SettingsFile), data, 0o600))
	read, err := ReadSettings(keyFile)
	require.NoError(t, err)
	ks, _, err := read.OpenKeyset()
	require.NoError(t, err)
	assert.Equal(t, int64(3), ks.Generation)

	g, err := New(Config{Settings: &read})
	require.NoError(t, err)
	assert.True(t, g.Active())
	other, _, _ := ed25519.GenerateKey(rand.Reader)
	read.DeployKey = hex.EncodeToString(other)
	_, err = New(Config{Settings: &read})
	assert.Error(t, err, "a keyset signed by another key")

	g, err = New(Config{Settings: &s, DeployPath: true})
	require.NoError(t, err)
	assert.False(t, g.Active(), "the deploy path takes no token")
}
