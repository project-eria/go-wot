package auth

import (
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newP256(t *testing.T) (*ecdsa.PrivateKey, PublicKey) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	raw, err := priv.PublicKey.Bytes()
	require.NoError(t, err)
	pub, err := ParsePublicKey(P256, raw)
	require.NoError(t, err)
	return priv, pub
}

// signP256 signs msg until r matches want, and returns r ‖ s, each on
// exactly 32 bytes (the app's conversion from DER).
func signP256(t *testing.T, priv *ecdsa.PrivateKey, msg []byte, want func(r *big.Int) bool) []byte {
	h := sha256.Sum256(msg)
	for range 10000 {
		r, s, err := ecdsa.Sign(rand.Reader, priv, h[:])
		require.NoError(t, err)
		if want(r) {
			sig := make([]byte, 64)
			r.FillBytes(sig[:32])
			s.FillBytes(sig[32:])
			return sig
		}
	}
	t.Fatal("no signature with the wanted r")
	return nil
}

var anyR = func(*big.Int) bool { return true }

func TestP256Signatures(t *testing.T) {
	priv, pub := newP256(t)
	msg := []byte("challenge|deploy|{}")

	t.Run("r on 31 bytes", func(t *testing.T) {
		sig := signP256(t, priv, msg, func(r *big.Int) bool { return r.BitLen() <= 248 })
		assert.Zero(t, sig[0])
		assert.True(t, pub.Verify(msg, sig))
	})
	t.Run("r with the high bit set", func(t *testing.T) {
		sig := signP256(t, priv, msg, func(r *big.Int) bool { return r.BitLen() == 256 })
		assert.True(t, pub.Verify(msg, sig))
	})
	t.Run("high s accepted", func(t *testing.T) {
		sig := signP256(t, priv, msg, anyR)
		s := new(big.Int).SetBytes(sig[32:])
		s.Sub(elliptic.P256().Params().N, s)
		s.FillBytes(sig[32:])
		assert.True(t, pub.Verify(msg, sig))
	})
	t.Run("63 or 65 bytes refused", func(t *testing.T) {
		sig := signP256(t, priv, msg, anyR)
		assert.False(t, pub.Verify(msg, sig[:63]))
		assert.False(t, pub.Verify(msg, append(sig, 0)))
	})
	t.Run("other message refused", func(t *testing.T) {
		sig := signP256(t, priv, msg, anyR)
		assert.False(t, pub.Verify([]byte("challenge|deploy|{\"x\":1}"), sig))
	})
	t.Run("zero r refused", func(t *testing.T) {
		sig := signP256(t, priv, msg, anyR)
		clear(sig[:32])
		assert.False(t, pub.Verify(msg, sig))
	})
}

func TestParsePublicKey(t *testing.T) {
	edPub, _, _ := ed25519.GenerateKey(rand.Reader)
	ed, err := ParsePublicKey(Ed25519, edPub)
	require.NoError(t, err)

	_, err = ParsePublicKey(P256, edPub)
	assert.Error(t, err, "an Ed25519 key is not a P-256 key")
	_, err = ParsePublicKey(Ed25519, edPub[:31])
	assert.Error(t, err)

	p384, err := ecdh.P384().GenerateKey(rand.Reader)
	require.NoError(t, err)
	_, err = ParsePublicKey(P256, p384.PublicKey().Bytes())
	assert.Error(t, err, "P-384 refused at registration")

	notOnCurve := make([]byte, 65)
	notOnCurve[0] = 0x04
	notOnCurve[64] = 1
	_, err = ParsePublicKey(P256, notOnCurve)
	assert.Error(t, err)

	_, err = ParsePublicKey("rsa", edPub)
	assert.Error(t, err)

	_, p := newP256(t)
	assert.NotEqual(t, ed.KID(), p.KID())
}

func TestEd25519Signatures(t *testing.T) {
	pubRaw, priv, _ := ed25519.GenerateKey(rand.Reader)
	pub, err := ParsePublicKey(Ed25519, pubRaw)
	require.NoError(t, err)
	msg := []byte("hello")
	sig := ed25519.Sign(priv, msg)
	assert.True(t, pub.Verify(msg, sig))
	assert.False(t, pub.Verify(msg, sig[:63]))
	assert.False(t, pub.Verify([]byte("hellO"), sig))
}
