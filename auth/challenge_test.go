package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChallenges(t *testing.T) {
	priv, pub := newP256(t)
	_, other := newP256(t)
	c := NewChallenges(time.Minute, 2)

	ch, err := c.Issue(pub.KID())
	require.NoError(t, err)
	msg := []byte(ch + "deploy{}")
	sig := signP256(t, priv, msg, anyR)

	assert.False(t, c.Redeem(ch, other, msg, sig), "bound to the key it was issued for")
	assert.False(t, c.Redeem(ch, pub, []byte("x"+ch), sig), "msg must start with the challenge")
	assert.True(t, c.Redeem(ch, pub, msg, sig))
	assert.False(t, c.Redeem(ch, pub, msg, sig), "valid signature on a consumed challenge refused")

	// A bad signature does not consume the challenge.
	ch, _ = c.Issue(pub.KID())
	msg = []byte(ch + "deploy{}")
	assert.False(t, c.Redeem(ch, pub, msg, make([]byte, 64)))
	assert.True(t, c.Redeem(ch, pub, msg, signP256(t, priv, msg, anyR)))
}

func TestChallengesExpireAndAreBounded(t *testing.T) {
	priv, pub := newP256(t)
	now := time.Unix(1_800_000_000, 0)
	c := NewChallenges(time.Minute, 2)
	c.now = func() time.Time { return now }

	ch, _ := c.Issue(pub.KID())
	_, err := c.Issue(pub.KID())
	require.NoError(t, err)
	_, err = c.Issue(pub.KID())
	assert.ErrorIs(t, err, ErrTooManyChallenges)

	now = now.Add(time.Minute)
	msg := []byte(ch + "x")
	assert.False(t, c.Redeem(ch, pub, msg, signP256(t, priv, msg, anyR)), "expired")
	_, err = c.Issue(pub.KID())
	assert.NoError(t, err, "expired challenges are purged")
}
