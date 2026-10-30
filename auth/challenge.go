package auth

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"
	"time"
)

// Challenges issues single-use challenges, each bound to the key it was
// issued for (AUTH-DESIGN §1.1): it is verified with that key only, and
// consumed by the first valid signature.
type Challenges struct {
	mu   sync.Mutex
	ttl  time.Duration
	max  int
	now  func() time.Time
	open map[string]pending
}

type pending struct {
	kid string
	exp time.Time
}

// NewChallenges keeps at most max open challenges, each valid for ttl.
func NewChallenges(ttl time.Duration, max int) *Challenges {
	return &Challenges{ttl: ttl, max: max, now: time.Now, open: map[string]pending{}}
}

// ErrTooManyChallenges is returned when the bounded list is full.
var ErrTooManyChallenges = errors.New("too many open challenges")

// Issue returns a fresh challenge bound to the key kid.
func (c *Challenges) Issue(kid string) (string, error) {
	var b [32]byte
	rand.Read(b[:])
	ch := base64.RawURLEncoding.EncodeToString(b[:])
	c.mu.Lock()
	defer c.mu.Unlock()
	now := c.now()
	for k, p := range c.open {
		if !now.Before(p.exp) {
			delete(c.open, k)
		}
	}
	if len(c.open) >= c.max {
		return "", ErrTooManyChallenges
	}
	c.open[ch] = pending{kid: kid, exp: now.Add(c.ttl)}
	return ch, nil
}

// Redeem verifies sig over msg with key, where msg must start with the
// challenge ch issued for that very key. The challenge is consumed only by a
// valid signature, so a second signature on it (even a valid one) fails.
func (c *Challenges) Redeem(ch string, key PublicKey, msg, sig []byte) bool {
	if !bytes.HasPrefix(msg, []byte(ch)) {
		return false
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	p, ok := c.open[ch]
	if !ok || p.kid != key.KID() || !c.now().Before(p.exp) {
		return false
	}
	if !key.Verify(msg, sig) {
		return false
	}
	delete(c.open, ch)
	return true
}
