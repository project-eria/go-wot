package guard

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/project-eria/go-wot/auth"
	"github.com/project-eria/go-wot/auth/authz"
	"github.com/project-eria/go-wot/dataSchema"
	"github.com/project-eria/go-wot/interaction"
	"github.com/project-eria/go-wot/producer"
	"github.com/project-eria/go-wot/thing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGuard(t *testing.T) {
	dir := t.TempDir()
	now := time.Now()

	// Deployment root and the auth keyset it signs.
	_, deployPriv, _ := ed25519.GenerateKey(rand.Reader)
	tmpl := &x509.Certificate{SerialNumber: big.NewInt(1), Subject: pkix.Name{CommonName: "root"}, IsCA: true,
		NotBefore: now.Add(-time.Hour), NotAfter: now.Add(time.Hour), BasicConstraintsValid: true}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, deployPriv.Public(), deployPriv)
	require.NoError(t, err)
	roots := filepath.Join(dir, "roots.pem")
	require.NoError(t, os.WriteFile(roots, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}), 0o600))
	authPub, authPriv, _ := ed25519.GenerateKey(rand.Reader)
	pub, _ := auth.ParsePublicKey(auth.Ed25519, authPub)
	signed, err := auth.SignDocument(deployPriv, auth.KeysetPrefix, auth.Keyset{Generation: 1, Issued: now.Unix(),
		Keys: []auth.KeysetKey{{Kid: pub.KID(), Pubkey: base64.RawURLEncoding.EncodeToString(authPub), Usage: auth.UsageSign}}})
	require.NoError(t, err)
	keyset := filepath.Join(dir, "keyset.json")
	data, _ := json.Marshal(signed)
	require.NoError(t, os.WriteFile(keyset, data, 0o600))

	// A fake auth channel: it takes the TD signed by the app and serves the table.
	_, appPriv, _ := ed25519.GenerateKey(rand.Reader)
	appPub, _ := auth.ParsePublicKey(auth.Ed25519, appPriv.Public().(ed25519.PublicKey))
	const id = "urn:eria:garage:1"
	var mu sync.Mutex
	gotTD := map[string]bool{}
	beat := now.UnixMilli() // the time of the heartbeats
	doc := func(d auth.Doc) string {
		d.Iss, d.Aud, d.Issued = "auth", "eria-maison", now.Unix()
		tok, err := auth.SignDoc(authPriv, pub.KID(), d)
		require.NoError(t, err)
		return tok
	}
	channel := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/v1/things/td":
			mu.Lock()
			var in auth.SignedBy
			json.NewDecoder(r.Body).Decode(&in)
			payload, err := in.Verify(auth.TDPrefix, []auth.PublicKey{appPub})
			var td struct{ ID string }
			json.Unmarshal(payload, &td)
			gotTD[td.ID] = err == nil && in.Kid == appPub.KID()
			mu.Unlock()
		case strings.HasSuffix(r.URL.Path, "/table"):
			if r.Header.Get("If-None-Match") == `"2"` {
				w.WriteHeader(304)
				return
			}
			w.Write([]byte(doc(auth.Doc{Doc: auth.DocTable, Thing: id, Version: 2, Zones: []string{"garage"},
				Table: auth.Table{"property/open/read": {Levels: [3]int{9, 1, 1}}, "action/open/invoke": {Levels: [3]int{1, 1, 1}, Confirm: auth.OutsideRoom}}})))
		case r.URL.Path == "/v1/heartbeat":
			mu.Lock()
			w.Write([]byte(doc(auth.Doc{Doc: auth.DocHeartbeat, Thing: "*", Version: beat, Nonce: r.URL.Query().Get("nonce")})))
			mu.Unlock()
		case r.URL.Path == "/v1/revocation":
			w.Write([]byte(doc(auth.Doc{Doc: auth.DocRevocation, Thing: "*", Version: 1})))
		}
	}))
	t.Cleanup(channel.Close)

	floor := filepath.Join(dir, "time-floor")
	g, err := New(Config{Channel: channel.URL, Deployment: "eria-maison", Keyset: keyset, Roots: roots, AppKey: appPriv, Floor: floor})
	require.NoError(t, err)
	require.True(t, g.Active())

	td, _ := thing.New(strings.TrimPrefix(id, "urn:"), "1", "Garage", "", nil)
	open, _ := dataSchema.NewBoolean()
	td.AddProperty(interaction.NewProperty("open", "", "", open))
	td.AddAction(interaction.NewAction("open", "", ""))
	g.Expose("", producer.NewExposedThing(td, "", &sync.WaitGroup{}))
	ctrl, _ := thing.New("eria:app:controler:1", "1", "Controller", "", nil)
	level, _ := dataSchema.NewBoolean()
	ctrl.AddProperty(interaction.NewProperty("logLevel", "", "", level))
	g.Expose(ControllerRef, producer.NewExposedThing(ctrl, ControllerRef, &sync.WaitGroup{}))
	assert.Equal(t, []string{"bearer_sc"}, td.Security, "the TD declares the bearer scheme")

	app := fiber.New()
	app.Use(g.Handler)
	app.All("/*", func(c *fiber.Ctx) error { return c.SendString("ok") })
	token := func(c auth.Claims) string {
		lvl := 9
		if c.Lvl == nil {
			c.Lvl = &lvl
		}
		c.Iss, c.Aud, c.Iat, c.Exp, c.Jti = "auth", "eria-maison", now.Unix(), now.Add(4*time.Minute).Unix(), "j"
		tok, err := auth.Sign(authPriv, pub.KID(), c)
		require.NoError(t, err)
		return tok
	}
	garage := token(auth.Claims{Kind: auth.KindUser, Term: "ecran-garage", Zone: "garage"})
	entrance := token(auth.Claims{Kind: auth.KindUser, Term: "ecran-entree", Zone: "entree"})
	status := func(method, path, tok string) int {
		req := httptest.NewRequest(method, path, nil)
		if tok != "" {
			req.Header.Set("Authorization", "Bearer "+tok)
		}
		resp, err := app.Test(req)
		require.NoError(t, err)
		return resp.StatusCode
	}

	// Before the table: the TD only.
	assert.Equal(t, 401, status("GET", "/", ""), "no token")
	assert.Equal(t, 200, status("GET", "/", entrance))
	assert.Equal(t, 403, status("GET", "/open", garage), "no table yet: None")
	assert.Equal(t, 200, status("GET", "/eria", ""), "deploy path")
	assert.Equal(t, 401, status("PUT", "/eria/logLevel", ""))

	g.pullAll(context.Background())
	mu.Lock()
	assert.True(t, gotTD[id], "the TD was sent to auth, signed by the app")
	mu.Unlock()

	assert.Equal(t, 200, status("GET", "/open", garage), "local 9")
	assert.Equal(t, 200, status("GET", "/OPEN/", garage), "routed like Fiber")
	assert.Equal(t, 403, status("GET", "/open", entrance), "remote 1")
	assert.Equal(t, 403, status("POST", "/open", garage), "level 9 under 1")
	assert.Equal(t, 403, status("DELETE", "/open", garage))
	assert.Equal(t, 200, status("OPTIONS", "/open", ""), "CORS preflight")

	// A parent at the garage screen; from the entrance, a confirmation.
	lvl1 := 1
	parent := token(auth.Claims{Kind: auth.KindUser, Sub: "profil:paul", Term: "ecran-garage", Zone: "garage", Lvl: &lvl1})
	assert.Equal(t, 200, status("POST", "/open", parent))
	away := token(auth.Claims{Kind: auth.KindUser, Sub: "profil:paul", Term: "tel-paul", Lvl: &lvl1})
	assert.Equal(t, 403, status("POST", "/open", away))
	act, _ := authz.Act(id, authz.Function{Kind: "action", Name: "open"}, nil)
	req := httptest.NewRequest("POST", "/open", nil)
	req.Header.Set("Authorization", "Bearer "+away)
	req.Header.Set("X-Eria-Confirm", confirmToken(t, authPriv, pub.KID(), now, act))
	resp, _ := app.Test(req)
	assert.Equal(t, 200, resp.StatusCode)
	resp, _ = app.Test(req)
	assert.Equal(t, 403, resp.StatusCode, "a confirmation is used once")

	// 304: the held table stays.
	g.pullAll(context.Background())
	assert.Equal(t, 200, status("GET", "/open", garage))

	// Signed time (§4.3): the floor was written; the guard decides with the
	// heartbeat's time, not the system's.
	written, _ := os.ReadFile(floor)
	assert.Equal(t, strconv.FormatInt(now.Unix(), 10), string(written))
	mu.Lock()
	beat = now.Add(50 * time.Minute).UnixMilli()
	mu.Unlock()
	g.pullAll(context.Background())
	assert.Equal(t, 401, status("GET", "/open", garage), "4 min tokens: expired at the signed time")
	mu.Lock()
	beat = now.Add(-time.Minute).UnixMilli()
	mu.Unlock()
	g.pullAll(context.Background())
	at, _ := g.clock.now()
	assert.WithinDuration(t, now.Add(50*time.Minute), at, 5*time.Second, "an older heartbeat is refused")
	mu.Lock()
	beat = now.Add(3 * time.Hour).UnixMilli()
	mu.Unlock()
	g.pullAll(context.Background())
	at, _ = g.clock.now()
	assert.WithinDuration(t, now.Add(50*time.Minute), at, 5*time.Second, "too far from the clock")
	// After a restart, a heartbeat before the floor is refused.
	c := newClock(floor)
	assert.Error(t, c.accept(now.Add(-time.Minute).UnixMilli()))
	assert.NoError(t, c.accept(now.Add(time.Minute).UnixMilli()))
}

// confirmToken is a confirm token (no lvl), 60 s.
func confirmToken(t *testing.T, priv ed25519.PrivateKey, kid string, now time.Time, act string) string {
	tok, err := auth.Sign(priv, kid, auth.Claims{Iss: "auth", Aud: "eria-maison", Kind: auth.KindConfirm, Sub: "profil:paul",
		Term: "tel-paul", Act: act, Iat: now.Unix(), Exp: now.Add(50 * time.Second).Unix(), Jti: "c1"})
	require.NoError(t, err)
	return tok
}

func TestGenerationZero(t *testing.T) {
	g, err := New(Config{Keyset: filepath.Join(t.TempDir(), "none.json")})
	require.NoError(t, err)
	assert.False(t, g.Active())
	app := fiber.New()
	app.Use(g.Handler)
	app.Get("/x", func(c *fiber.Ctx) error { return c.SendString("ok") })
	resp, _ := app.Test(httptest.NewRequest("GET", "/x", nil))
	assert.Equal(t, 200, resp.StatusCode, "nosec until the keyset")
}
