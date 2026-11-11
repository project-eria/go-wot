package auth

import (
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"os"
	"strings"
)

// App keys and service proofs (AUTH-DESIGN §1.3, step 3a): shared by auth,
// the master that registers the keys and the apps that prove theirs.

// AppSubject prefixes the sub of a service token: app:<id>, the identity
// auth gave the app; the revocation subject is then sub:app:<id>.
const AppSubject = "app:"

// AppKeyPrefix is the domain of an app key registration, signed by a
// deployment key (the master).
const AppKeyPrefix = "eria-app-key-v1"

// AppKey binds an app to its key and its grants. Pubkey "" retires the
// app's key. Version is the time in ms of the master's message: auth keeps
// only a higher one, so an old message replayed never brings back a retired
// key, nor grants the admin removed.
type AppKey struct {
	App     string   `json:"app"`    // label: the app's name in the deployment
	Node    string   `json:"node"`   // where it runs
	Pubkey  string   `json:"pubkey"` // base64url, Ed25519
	Version int64    `json:"version"`
	Grants  []string `json:"grants,omitempty"` // accepted by the admin (§1.3), grammar of §2.6
}

// AdminKeysPrefix is the domain of the master's registry of admin keys,
// signed by a deployment key (revue de simplification, piste 7).
const AdminKeysPrefix = "eria-admin-keys-v1"

// AdminKeys is the master's registry of the keys that validate the admin's
// orders (admin-keys.json): the YubiKey (P-256) and the admin's phone (the
// kid of its biometric key). auth makes them its validators, never the
// other way round. Version is the registry's time in ms: an older one is
// refused.
type AdminKeys struct {
	Version int64           `json:"version"`
	Keys    []AdminKeyEntry `json:"keys"`
}

type AdminKeyEntry struct {
	Name string `json:"name"`
	// Kind is "yubikey" (a validator key of its own) or "phone" (the bio
	// key of an enrolled phone); both are P-256, the kind tells them apart.
	Kind string `json:"kind,omitempty"`
	Alg  string `json:"alg"`
	Key  string `json:"key"` // base64url
	Kid  string `json:"kid"`
}

// ServiceProofDomain is the domain of a service proof: the app's key signs
// challenge ‖ domain ‖ kid of auth (ServiceProof).
const ServiceProofDomain = "eria-service-proof-v1"

// ServiceProof is the message an app signs to get its service token: the
// challenge is bound to the app's key by auth, the kid of auth binds the
// proof to this auth.
func ServiceProof(challenge, authKid string) []byte {
	return []byte(challenge + ServiceProofDomain + authKid)
}

// AppKeyEnv names the file of the app's key, set by the agent (step 3b):
// the hex seed of an Ed25519 key, in the app's data folder.
const AppKeyEnv = "ERIA_APP_KEY"

// LoadAppKey reads the key named by ERIA_APP_KEY; nil when it is unset (an
// app the agent did not install, or a test).
func LoadAppKey() (ed25519.PrivateKey, error) {
	path := os.Getenv(AppKeyEnv)
	if path == "" {
		return nil, nil
	}
	return ReadAppKey(path)
}

// ReadAppKey reads an app key file.
func ReadAppKey(path string) (ed25519.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	seed, err := hex.DecodeString(strings.TrimSpace(string(data)))
	if err != nil || len(seed) != ed25519.SeedSize {
		return nil, errors.New("app key: expected a hex Ed25519 seed")
	}
	return ed25519.NewKeyFromSeed(seed), nil
}

// AppKID is the kid of an app key.
func AppKID(priv ed25519.PrivateKey) string {
	pub, _ := ParsePublicKey(Ed25519, priv.Public().(ed25519.PublicKey))
	return pub.KID()
}

// SignedBy is what an app sends auth on the Things channel (step 3d): a
// document signed by its key, with the kid auth finds the key by.
type SignedBy struct {
	Kid string `json:"kid"`
	Signed
}

// SignBy signs v with the app key priv.
func SignBy(priv ed25519.PrivateKey, prefix string, v any) (SignedBy, error) {
	d, err := SignDocument(priv, prefix, v)
	return SignedBy{Kid: AppKID(priv), Signed: d}, err
}

// TDPrefix is the domain of a TD a Thing sends auth: the Thing belongs to
// the app that signs it.
const TDPrefix = "eria-td-v1"

// StatesPrefix is the domain of the house states, published by
// manager-contexts.
const StatesPrefix = "eria-states-v1"

// States are the active states of the house. Version is the time in ms:
// auth keeps only a higher one, so an old list replayed changes nothing.
type States struct {
	States  []string `json:"states"`
	Version int64    `json:"version"`
}

// AckPrefix is the domain of a heartbeat acknowledgement (§4.4): a Thing
// tells auth it is alive, in which mode, and which versions it holds.
const AckPrefix = "eria-ack-v1"

// Ack acknowledges the heartbeat Beat (its version, ms): auth takes it only
// while that heartbeat is recent, so an old ack replayed proves nothing.
type Ack struct {
	Beat       int64            `json:"beat"`
	Mode       string           `json:"mode"`
	Tables     map[string]int64 `json:"tables"`            // held, by Thing
	Revocation int64            `json:"revocation"`        // held
	Refused    []string         `json:"refused,omitempty"` // documents refused: "table:<id>", "revocation"
}

// AppReleasePrefix is the domain of an app release's descriptor (AUTH-DESIGN
// §3.4): signed by a deployment key once the admin validated the
// publication; the app installs a package only when its descriptor allows
// it, and never a version at or below its own.
const AppReleasePrefix = "eria-app-release-v1"

// AppRelease is what a descriptor signs: a version of the app, for one
// updater target (darwin-aarch64, windows-x86_64…), and its package's sha256
// (hex), checked before the package reaches the updater.
type AppRelease struct {
	Version string `json:"version"`
	Target  string `json:"target"`
	Sha256  string `json:"sha256"`
}
