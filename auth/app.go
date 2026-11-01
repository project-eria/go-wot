package auth

// App keys and service proofs (AUTH-DESIGN §1.3, step 3a): shared by auth,
// the master that registers the keys and the apps that prove theirs.

// AppSubject prefixes the sub of a service token: app:<id>, the identity
// auth gave the app; the revocation subject is then sub:app:<id>.
const AppSubject = "app:"

// AppKeyPrefix is the domain of an app key registration, signed by a
// deployment key (the master).
const AppKeyPrefix = "eria-app-key-v1"

// AppKey binds an app to its key. Pubkey "" retires the app's key. Version
// is the time in ms of the master's message: auth keeps only a higher one,
// so an old message replayed never brings back a retired key.
type AppKey struct {
	App     string `json:"app"`    // label: the app's name in the deployment
	Node    string `json:"node"`   // where it runs
	Pubkey  string `json:"pubkey"` // base64url, Ed25519
	Version int64  `json:"version"`
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
