package auth

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	zlog "github.com/rs/zerolog/log"
)

// ServiceTokens gets the app's service tokens from auth's Things channel,
// proving its app key (AUTH-DESIGN §1.3, step 3c): one for the coordinator,
// renewed at 2/3 of its life by Run, and one per Thing (step 4e), fetched by
// For on first use and again past 2/3 of its life.
type ServiceTokens struct {
	channel string
	authKid string // the proof binds auth's key
	key     ed25519.PrivateKey
	client  *http.Client

	mu    sync.RWMutex
	token string
	exp   time.Time

	// ponytail: one lock across the Things' fetches (a few ms on the LAN, a
	// failure remembered RetryAfterFailure); per Thing if it ever waits.
	thingsMu sync.Mutex
	things   map[string]heldToken
}

type heldToken struct {
	token    string
	iat, exp time.Time
	failed   time.Time // last failed request
}

// RetryAfterFailure is the pause after a token request that failed.
var RetryAfterFailure = 30 * time.Second

// NewServiceTokens prepares the tokens of the app key, for the auth whose
// signing key is in ks.
func NewServiceTokens(channel string, key ed25519.PrivateKey, ks Keyset) (*ServiceTokens, error) {
	for _, k := range ks.Keys {
		if k.Usage == UsageSign {
			return &ServiceTokens{channel: strings.TrimSuffix(channel, "/"), authKid: k.Kid, key: key,
				client: &http.Client{Timeout: 10 * time.Second}, things: map[string]heldToken{}}, nil
		}
	}
	return nil, errors.New("no auth signing key in the keyset")
}

// For is the token for a Thing (its TD id, or its DNS-SD instance), ""
// while auth gives none: auth knows no such Thing, or can't be reached.
func (s *ServiceTokens) For(thing string) string {
	s.thingsMu.Lock()
	defer s.thingsMu.Unlock()
	h, now := s.things[thing], time.Now()
	fresh := now.Before(h.iat.Add(h.exp.Sub(h.iat) * 2 / 3))
	if !fresh && !now.Before(h.failed.Add(RetryAfterFailure)) {
		ctx, cancel := context.WithTimeout(context.Background(), s.client.Timeout)
		token, iat, exp, err := s.fetch(ctx, thing)
		cancel()
		if err != nil {
			zlog.Warn().Err(err).Str("thing", thing).Msg("[auth:ServiceTokens] No service token for this Thing")
			h.failed = now
		} else {
			h = heldToken{token: token, iat: iat, exp: exp}
		}
		s.things[thing] = h
	}
	if !now.Before(h.exp) {
		return ""
	}
	return h.token
}

// Token is the current token for the coordinator, "" before the first one or
// once expired.
func (s *ServiceTokens) Token() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if !time.Now().Before(s.exp) {
		return ""
	}
	return s.token
}

// Run renews the token until ctx ends; ready is closed after the first
// attempt, whether it worked or not.
func (s *ServiceTokens) Run(ctx context.Context, ready chan<- struct{}) {
	for first := true; ; first = false {
		wait := RetryAfterFailure
		token, iat, exp, err := s.fetch(ctx, "")
		if err != nil {
			zlog.Warn().Err(err).Msg("[auth:ServiceTokens] No service token, retrying")
		} else {
			s.mu.Lock()
			s.token, s.exp = token, exp
			s.mu.Unlock()
			wait = exp.Sub(iat) * 2 / 3
		}
		if first {
			close(ready)
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(wait):
		}
	}
}

// fetch asks a challenge for the app key and signs it (ServiceProof); the
// token is for thing, or for the coordinator when thing is "".
func (s *ServiceTokens) fetch(ctx context.Context, thing string) (token string, iat, exp time.Time, err error) {
	kid := AppKID(s.key)
	var ch struct{ Challenge string }
	if err = s.post(ctx, "/v1/service/challenge", map[string]string{"kid": kid}, &ch); err != nil {
		return
	}
	var out struct{ Token string }
	in := map[string]any{"kid": kid, "challenge": ch.Challenge, "sig": ed25519.Sign(s.key, ServiceProof(ch.Challenge, s.authKid))}
	if thing != "" {
		in["thing"] = thing
	}
	if err = s.post(ctx, "/v1/service/token", in, &out); err != nil {
		return
	}
	// The claims are read, not verified: a forged token only fails at the Things.
	parts := strings.Split(out.Token, ".")
	if len(parts) != 3 {
		return "", iat, exp, errors.New("not a token")
	}
	var c struct{ Iat, Exp int64 }
	payload, err := b64.DecodeString(parts[1])
	if err != nil || json.Unmarshal(payload, &c) != nil || c.Exp <= c.Iat {
		return "", iat, exp, errors.New("unreadable token")
	}
	return out.Token, time.Unix(c.Iat, 0), time.Unix(c.Exp, 0), nil
}

func (s *ServiceTokens) post(ctx context.Context, path string, in, out any) error {
	body, _ := json.Marshal(in)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, s.channel+path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s: auth answered %d", path, resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}
