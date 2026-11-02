package auth

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceTokens(t *testing.T) {
	authPub, authPriv, _ := ed25519.GenerateKey(rand.Reader)
	authKey, _ := ParsePublicKey(Ed25519, authPub)
	_, appPriv, _ := ed25519.GenerateKey(rand.Reader)
	appKey, _ := ParsePublicKey(Ed25519, appPriv.Public().(ed25519.PublicKey))

	// A fake auth: a token of 3 s for a valid proof, refused until ready.
	var ready, issued atomic.Int32
	channel := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var in struct {
			Kid, Challenge string
			Sig            []byte
		}
		json.NewDecoder(r.Body).Decode(&in)
		switch {
		case ready.Load() == 0 || in.Kid != appKey.KID():
			w.WriteHeader(http.StatusForbidden)
		case r.URL.Path == "/v1/service/challenge":
			json.NewEncoder(w).Encode(map[string]string{"challenge": "ch"})
		case !appKey.Verify(ServiceProof(in.Challenge, authKey.KID()), in.Sig):
			w.WriteHeader(http.StatusForbidden)
		default:
			now := time.Now().Unix()
			tok, _ := Sign(authPriv, authKey.KID(), Claims{Iss: "auth", Aud: "a", Kind: KindService, Sub: "app:x",
				Iat: now, Exp: now + 3, Jti: "j"})
			issued.Add(1)
			json.NewEncoder(w).Encode(map[string]string{"token": tok})
		}
	}))
	t.Cleanup(channel.Close)
	saved := RetryAfterFailure
	RetryAfterFailure = 100 * time.Millisecond
	t.Cleanup(func() { RetryAfterFailure = saved })

	s, err := NewServiceTokens(channel.URL+"/", appPriv, Keyset{Keys: []KeysetKey{{Kid: authKey.KID(), Usage: UsageSign}}})
	require.NoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	first := make(chan struct{})
	go s.Run(ctx, first)
	<-first
	assert.Empty(t, s.Token(), "auth refused: no token yet")

	ready.Store(1)
	require.Eventually(t, func() bool { return s.Token() != "" }, 2*time.Second, 20*time.Millisecond, "retried after a failure")
	require.Eventually(t, func() bool { return issued.Load() >= 2 }, 4*time.Second, 50*time.Millisecond, "renewed at 2/3 of its life")
}
