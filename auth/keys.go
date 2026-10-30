// Package auth holds the cryptographic base shared by the ERIA auth service,
// the master, the relays and the Things (AUTH-DESIGN §1.1, §2.7): typed public
// keys (Ed25519 and ECDSA P-256), auth-signed JWTs and single-use challenges.
package auth

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"math/big"
	"slices"
)

// Alg is fixed when a public key is registered, never read from a signed
// message (AUTH-DESIGN §1.1).
type Alg string

const (
	// Ed25519 is used by every key kept in a file (deployment, auth, apps,
	// relays, ESP32) and by the Linux terminal identity key.
	Ed25519 Alg = "ed25519"
	// P256 is used by hardware-held keys (Secure Enclave, StrongBox, TPM,
	// YubiKey PIV): ECDSA over SHA-256(message), raw r ‖ s signature.
	P256 Alg = "p256"
)

// p256Size is the length of an uncompressed P-256 point: 0x04 ‖ X ‖ Y.
const p256Size = 65

// PublicKey is a registered public key. Each algorithm keeps its own Go type:
// a key is never a []byte interpreted according to its length.
type PublicKey struct {
	alg Alg
	raw []byte
	ed  ed25519.PublicKey
	ec  *ecdsa.PublicKey
}

// ParsePublicKey checks raw against alg: 32 bytes for Ed25519, 65 bytes
// 0x04 ‖ X ‖ Y on the P-256 curve for P256. Any other curve or length is
// refused at registration.
func ParsePublicKey(alg Alg, raw []byte) (PublicKey, error) {
	raw = slices.Clone(raw)
	switch alg {
	case Ed25519:
		if len(raw) != ed25519.PublicKeySize {
			return PublicKey{}, errors.New("ed25519 key must be 32 bytes")
		}
		return PublicKey{alg: alg, raw: raw, ed: ed25519.PublicKey(raw)}, nil
	case P256:
		if len(raw) != p256Size || raw[0] != 0x04 {
			return PublicKey{}, errors.New("p256 key must be 65 bytes, 0x04 ‖ X ‖ Y")
		}
		ec, err := ecdsa.ParseUncompressedPublicKey(elliptic.P256(), raw)
		if err != nil {
			return PublicKey{}, err
		}
		return PublicKey{alg: alg, raw: raw, ec: ec}, nil
	}
	return PublicKey{}, errors.New("unknown key algorithm")
}

func (k PublicKey) Alg() Alg      { return k.alg }
func (k PublicKey) Bytes() []byte { return slices.Clone(k.raw) }

// KID is sha256(alg ‖ key bytes), base64url: two algorithms never share a kid.
func (k PublicKey) KID() string {
	h := sha256.Sum256(append([]byte(k.alg), k.raw...))
	return base64.RawURLEncoding.EncodeToString(h[:])
}

// Verify checks sig over msg. Ed25519 signatures are 64 bytes; P-256
// signatures are r ‖ s, 32 bytes each, over SHA-256(msg). A high s is
// accepted: replay is detected on the challenge, never on the signature.
func (k PublicKey) Verify(msg, sig []byte) bool {
	switch k.alg {
	case Ed25519:
		return len(sig) == ed25519.SignatureSize && ed25519.Verify(k.ed, msg, sig)
	case P256:
		if len(sig) != 64 {
			return false
		}
		r := new(big.Int).SetBytes(sig[:32])
		s := new(big.Int).SetBytes(sig[32:])
		h := sha256.Sum256(msg)
		// ecdsa.Verify refuses r or s equal to 0 or ≥ n.
		return ecdsa.Verify(k.ec, h[:], r, s)
	}
	return false
}
