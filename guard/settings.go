package guard

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/project-eria/go-wot/auth"
)

// The activation of a Go app (AUTH-DESIGN §4.7 step 5c): the master sends
// its node's agent the keyset (signed command setKeyset); the agent writes
// SettingsFile beside each app's key, and restarts the app, whose guard
// then requires tokens. Like setKeyset on an ESP32.

// SettingsFile is the file the agent writes in each app's data folder.
const SettingsFile = "auth-guard.json"

// Settings is what activates a Go app's guard.
type Settings struct {
	Keyset     string `json:"keyset"`     // auth-keyset.json, signed by the deployment key
	Channel    string `json:"channel"`    // auth's Things channel
	Deployment string `json:"deployment"` // aud of the tokens
	DeployKey  string `json:"deployKey"`  // the deployment key (hex): verifies the keyset
	Master     string `json:"master"`     // the master's TD URL: its inventory (F037)
}

// ReadSettings reads the settings beside the app key file; os.ErrNotExist
// when the agent wrote none (not activated).
func ReadSettings(appKeyFile string) (Settings, error) {
	var s Settings
	data, err := os.ReadFile(filepath.Join(filepath.Dir(appKeyFile), SettingsFile))
	if err != nil {
		return s, err
	}
	if err := json.Unmarshal(data, &s); err != nil {
		return s, err
	}
	return s, nil
}

// OpenKeyset checks the keyset against the deployment key.
func (s Settings) OpenKeyset() (auth.Keyset, auth.KeyLookup, error) {
	raw, err := hex.DecodeString(s.DeployKey)
	if err != nil {
		return auth.Keyset{}, nil, errors.New("guard: bad deployment key")
	}
	deploy, err := auth.ParsePublicKey(auth.Ed25519, raw)
	if err != nil {
		return auth.Keyset{}, nil, err
	}
	return auth.OpenKeyset([]byte(s.Keyset), []auth.PublicKey{deploy})
}
