// Package guard is the single guard of the Things (AUTH-DESIGN §2.7): a
// Fiber middleware, before the routes and the WebSocket upgrade, that runs
// the decision of package authz on every request.
//
// A Thing without an auth keyset is in generation 0: it keeps the nosec
// behaviour, with a warning. Once the keyset is there (the activation,
// §2.4), every request needs a token, except the reads of /eria (the deploy
// path) and the CORS preflights. The Thing sends its TD to auth's Things
// channel, then pulls its signed table and the revocation list.
//
// It decides with the signed time of auth's heartbeats (clock.go, §4.3), in
// the mode they allow (§4.4): normal; minimal (terminal tokens, local 9
// functions without confirmation) from the start until a heartbeat, after
// missedBeats pulls without one, or while a version a heartbeat announced
// is not received after versionWait; offline (every token refused) once
// isolated for more than isolatedMax.
//
// ponytail: the system time when there is neither floor nor heartbeat; the
// keyset is a file put by the deployment until the master delivers it;
// observations are checked at the handshake only, not closed at exp nor at
// a switch to the minimal mode.
package guard

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/project-eria/go-wot/auth"
	"github.com/project-eria/go-wot/auth/authz"
	"github.com/project-eria/go-wot/producer"
	"github.com/project-eria/go-wot/securityScheme"
	zlog "github.com/rs/zerolog/log"
)

// Config locates auth and the keyset of the deployment.
type Config struct {
	Channel    string // auth's Things channel, http://auth:8444
	Deployment string // aud of the tokens
	Keyset     string // auth keyset signed by the deployment key (eria-manager auth-keyset)
	Roots      string // deployment roots, PEM (eria-manager tls-root)
	// AppKey signs the TDs sent to auth (auth.LoadAppKey); nil sends them
	// unsigned, until the agent gives the app a key.
	AppKey ed25519.PrivateKey
	// Floor is the file keeping the time floor across restarts (§4.3); ""
	// keeps none.
	Floor string
}

// ControllerRef is the ref of the app controller, whose reads are the
// deploy path (/eria/*): open without token.
const ControllerRef = "eria"

// Pull is how often a Thing asks auth for a heartbeat and new versions.
var Pull = 30 * time.Second

// Thresholds of the modes (§4.4).
const (
	missedBeats = 3
	versionWait = time.Minute
	isolatedMax = 72 * time.Hour
	beatThings  = 16 // auth's bound on the tables a heartbeat announces
)

var errRefused = errors.New("document refused")

const maxSeen = 1024

type Guard struct {
	cfg    Config
	keys   auth.KeyLookup // nil: generation 0
	client *http.Client
	now    func() time.Time // before the first heartbeat
	clock  *clock

	mu      sync.RWMutex
	things  map[string]producer.ExposedThing // by lowercase ref
	tables  map[string]auth.Doc              // by TD id
	revoked auth.Doc
	seen    map[string]int64 // confirm jti → exp
	// Versions announced and not held, since when ("table:<id>",
	// "revocation"); a document refused is no isolation: kept out.
	pending  map[string]time.Time
	refused  map[string]bool
	mode     authz.Mode // the last one logged
	beatFrom int        // rotation over the Things, beatThings per heartbeat
}

// New loads the keyset. No keyset file: generation 0, nosec.
func New(cfg Config) (*Guard, error) {
	g := &Guard{cfg: cfg, client: &http.Client{Timeout: 10 * time.Second}, now: time.Now, clock: newClock(cfg.Floor),
		things: map[string]producer.ExposedThing{}, tables: map[string]auth.Doc{}, seen: map[string]int64{},
		pending: map[string]time.Time{}, refused: map[string]bool{}, mode: authz.Normal}
	data, err := os.ReadFile(cfg.Keyset)
	if errors.Is(err, os.ErrNotExist) || cfg.Keyset == "" {
		zlog.Warn().Msg("[guard] No auth keyset: generation 0, every request accepted without token")
		return g, nil
	}
	if err != nil {
		return nil, err
	}
	rootsPEM, err := os.ReadFile(cfg.Roots)
	if err != nil {
		return nil, err
	}
	deploy, err := auth.DeploymentKeys(rootsPEM)
	if err != nil {
		return nil, err
	}
	if _, g.keys, err = auth.OpenKeyset(data, deploy); err != nil {
		return nil, err
	}
	if cfg.Channel == "" || cfg.Deployment == "" {
		return nil, errors.New("guard: channel and deployment are required with a keyset")
	}
	zlog.Info().Msg("[guard] Activated: tokens required")
	return g, nil
}

// Active tells whether the Thing checks tokens.
func (g *Guard) Active() bool { return g.keys != nil }

// Expose registers a Thing before the server starts; an active guard
// declares the bearer scheme in its TD.
func (g *Guard) Expose(ref string, t producer.ExposedThing) {
	if g.Active() {
		td := t.GetThingDescription()
		td.Security = []string{"bearer_sc"}
		td.SecurityDefinitions = map[string]securityScheme.SecurityScheme{"bearer_sc": securityScheme.NewEriaBearer()}
	}
	g.mu.Lock()
	g.things[strings.ToLower(ref)] = t
	g.mu.Unlock()
}

// --- Decision ---

// Handler is the middleware.
func (g *Guard) Handler(c *fiber.Ctx) error {
	if !g.Active() || c.Method() == fiber.MethodOptions {
		return c.Next()
	}
	t, f, ok := g.resolve(c)
	if ok && t.Ref() == ControllerRef && f.Op == authz.Read {
		return c.Next() // deploy path: /eria reads
	}
	if !ok {
		return deny(c, 7, "unknown route", "")
	}
	id := t.GetThingDescription().ID
	g.mu.RLock()
	table, revoked := g.tables[id], g.revoked
	g.mu.RUnlock()
	var input []byte
	if f.Op != authz.Read {
		input = c.Body()
	}
	d := authz.Decide(authz.Check{
		Request: authz.Request{Mode: g.Mode(), Verifier: authz.Thing, Keys: g.keys, Aud: g.cfg.Deployment, Now: g.time(),
			Token: strings.TrimPrefix(c.Get(fiber.HeaderAuthorization), "Bearer ")},
		Thing:    id,
		Zones:    table.Zones,
		Table:    table.Table,
		Revoked:  &revoked,
		Function: f,
		Input:    input,
		Confirm:  c.Get("X-Eria-Confirm"),
		Seen:     g.seenJti,
	})
	if !d.Allow {
		return deny(c, d.Step, d.Reason, id)
	}
	return c.Next()
}

// deny answers 401 for a bad token (steps 0-3), 403 otherwise, without
// detail; the reason goes to the log.
func deny(c *fiber.Ctx, step int, reason, thing string) error {
	zlog.Info().Int("step", step).Str("reason", reason).Str("thing", thing).Str("path", c.Path()).Msg("[guard] Denied")
	status := fiber.StatusForbidden
	if step <= 3 {
		status = fiber.StatusUnauthorized
	}
	return c.Status(status).JSON(fiber.Map{"error": "denied"})
}

// resolve maps a request to its Thing and function the way Fiber routes it:
// case-insensitive, with or without a trailing slash; the main Thing has the
// ref "" and its affordances at the root.
func (g *Guard) resolve(c *fiber.Ctx) (producer.ExposedThing, authz.Function, bool) {
	segs := strings.Split(strings.Trim(strings.ToLower(c.Path()), "/"), "/")
	if segs[0] == "" {
		segs = nil
	}
	g.mu.RLock()
	t, ok := g.things[""]
	if len(segs) > 0 {
		if named, found := g.things[segs[0]]; found {
			t, ok, segs = named, true, segs[1:]
		}
	}
	g.mu.RUnlock()
	if !ok {
		return nil, authz.Function{}, false
	}
	method := c.Method()
	if len(segs) == 0 {
		return t, authz.TD, method == fiber.MethodGet
	}
	td := t.GetThingDescription()
	find := func(names []string) string {
		for _, n := range names {
			if strings.ToLower(n) == segs[0] {
				return n
			}
		}
		return ""
	}
	props, actions, events := keys(td.Properties), keys(td.Actions), keys(td.Events)
	switch method {
	case fiber.MethodGet:
		if n := find(props); n != "" {
			return t, authz.Function{Kind: "property", Name: n, Op: authz.Read}, true
		}
		if n := find(events); n != "" {
			return t, authz.Function{Kind: "event", Name: n, Op: authz.Read}, true
		}
	case fiber.MethodPut:
		if n := find(props); n != "" {
			return t, authz.Function{Kind: "property", Name: n, Op: authz.Write}, true
		}
	case fiber.MethodPost:
		if n := find(actions); n != "" {
			return t, authz.Function{Kind: "action", Name: n, Op: authz.Invoke}, true
		}
	}
	return t, authz.Function{}, false
}

func keys[V any](m map[string]V) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

// seenJti records a confirm token's jti until its exp; true when already
// used or the list is full (§2.8).
func (g *Guard) seenJti(jti string, exp int64) bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	now := g.time().Unix()
	for k, e := range g.seen {
		if e < now-int64(auth.ClockSkew.Seconds()) {
			delete(g.seen, k)
		}
	}
	if _, used := g.seen[jti]; used || len(g.seen) >= maxSeen {
		return true
	}
	g.seen[jti] = exp
	return false
}

// --- Documents from auth ---

// Run sends the TDs, then pulls the tables and the revocation list until ctx
// ends. Only an active guard talks to auth.
func (g *Guard) Run(ctx context.Context) {
	if !g.Active() {
		return
	}
	for {
		g.pullAll(ctx)
		select {
		case <-ctx.Done():
			return
		case <-time.After(Pull):
		}
	}
}

func (g *Guard) pullAll(ctx context.Context) {
	g.mu.RLock()
	things := make([]producer.ExposedThing, 0, len(g.things))
	ids := make([]string, 0, len(g.things))
	for _, t := range g.things {
		things = append(things, t)
		ids = append(ids, t.GetThingDescription().ID)
	}
	g.mu.RUnlock()
	v, beat := g.heartbeat(ctx, g.rotate(ids))
	for _, t := range things {
		td := t.GetThingDescription()
		g.mu.RLock()
		held, ok := g.tables[td.ID]
		g.mu.RUnlock()
		if !ok {
			g.sendTD(ctx, td.ID, td)
		}
		d, err := g.pull(ctx, "/v1/things/"+url.PathEscape(td.ID)+"/table", auth.DocTable, td.ID, held.Version)
		g.mu.Lock()
		g.refused["table:"+td.ID] = errors.Is(err, errRefused)
		if err == nil && d != nil {
			g.tables[td.ID] = *d
			g.seen = map[string]int64{} // F027: a new table, a fresh start
			zlog.Info().Str("thing", td.ID).Int64("version", d.Version).Msg("[guard] New table")
		}
		g.mu.Unlock()
	}
	g.mu.RLock()
	held := g.revoked.Version
	g.mu.RUnlock()
	d, err := g.pull(ctx, "/v1/revocation", auth.DocRevocation, "*", held)
	g.mu.Lock()
	g.refused["revocation"] = errors.Is(err, errRefused)
	if err == nil && d != nil {
		g.revoked = *d
	}
	if v != nil {
		g.announced("revocation", v.Revocation, g.revoked.Version)
		for id, version := range v.Tables {
			g.announced("table:"+id, version, g.tables[id].Version)
		}
	}
	ack := auth.Ack{Beat: beat, Tables: map[string]int64{}, Revocation: g.revoked.Version}
	for id, t := range g.tables {
		ack.Tables[id] = t.Version
	}
	for key, refused := range g.refused {
		if refused {
			ack.Refused = append(ack.Refused, key)
		}
	}
	g.mu.Unlock()
	g.switchMode()
	if v != nil {
		ack.Mode = g.Mode().String()
		g.acknowledge(ctx, ack)
	}
}

// acknowledge tells auth the heartbeat came (§4.4), signed by the app key:
// auth alerts when a Thing goes silent. Without a key, nothing is sent.
func (g *Guard) acknowledge(ctx context.Context, ack auth.Ack) {
	if g.cfg.AppKey == nil {
		return
	}
	signed, err := auth.SignBy(g.cfg.AppKey, auth.AckPrefix, ack)
	if err != nil {
		return
	}
	body, _ := json.Marshal(signed)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, g.cfg.Channel+"/v1/heartbeat/ack", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := g.client.Do(req)
	if err != nil {
		return // the next heartbeat tells
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		zlog.Warn().Int("status", resp.StatusCode).Msg("[guard] Acknowledgement refused by auth")
	}
}

// rotate picks the Things of this heartbeat: beatThings at most, in turn.
func (g *Guard) rotate(ids []string) []string {
	slices.Sort(ids)
	if len(ids) <= beatThings {
		return ids
	}
	g.beatFrom %= len(ids)
	chunk := append(ids[g.beatFrom:], ids[:g.beatFrom]...)[:beatThings]
	g.beatFrom += beatThings
	return chunk
}

// announced records a version a heartbeat announced against the held one;
// g.mu held.
func (g *Guard) announced(key string, version, held int64) {
	switch {
	case version <= held || g.refused[key]:
		delete(g.pending, key)
	case g.pending[key].IsZero():
		g.pending[key] = time.Now()
	}
}

// Mode is the mode of the decision now (§4.4).
func (g *Guard) Mode() authz.Mode {
	if g.clock.isolated() > isolatedMax {
		return authz.Offline
	}
	if g.clock.missed(missedBeats * Pull) {
		return authz.Minimal
	}
	g.mu.RLock()
	defer g.mu.RUnlock()
	for _, since := range g.pending {
		if time.Since(since) > versionWait {
			return authz.Minimal
		}
	}
	return authz.Normal
}

// switchMode logs a switch, and keeps the time of the last heartbeat at
// once; else the floor, hourly.
func (g *Guard) switchMode() {
	m := g.Mode()
	g.mu.Lock()
	switched := m != g.mode
	g.mode = m
	g.mu.Unlock()
	g.clock.save(switched)
	if switched {
		zlog.Warn().Str("mode", m.String()).Dur("isolated", g.clock.isolated()).Msg("[guard] Mode switched")
	}
}

// heartbeat asks auth the signed time for a fresh nonce (§4.3); it
// returns the versions announced and the heartbeat's version, nil when
// refused.
func (g *Guard) heartbeat(ctx context.Context, ids []string) (*auth.Versions, int64) {
	var b [16]byte
	rand.Read(b[:])
	nonce := base64.RawURLEncoding.EncodeToString(b[:])
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, g.cfg.Channel+"/v1/heartbeat?"+url.Values{"nonce": {nonce}, "thing": ids}.Encode(), nil)
	resp, err := g.client.Do(req)
	if err != nil {
		zlog.Warn().Err(err).Msg("[guard] auth unreachable")
		return nil, 0
	}
	defer resp.Body.Close()
	var body bytes.Buffer
	body.ReadFrom(http.MaxBytesReader(nil, resp.Body, auth.MaxDocSize+1))
	d, err := auth.ParseDoc(body.String(), g.keys, g.cfg.Deployment, auth.DocHeartbeat, "*")
	if err == nil && d.Nonce != nonce {
		err = errors.New("not our nonce")
	}
	if err == nil {
		err = g.clock.accept(d.Version)
	}
	if err != nil {
		zlog.Warn().Err(err).Int("status", resp.StatusCode).Msg("[guard] Heartbeat refused")
		return nil, 0
	}
	if d.Versions == nil {
		d.Versions = &auth.Versions{}
	}
	return d.Versions, d.Version
}

// time is the time to decide with: the signed one once a heartbeat came.
func (g *Guard) time() time.Time {
	if t, ok := g.clock.now(); ok {
		return t
	}
	return g.now()
}

func (g *Guard) sendTD(ctx context.Context, id string, td any) {
	var v any = td
	if g.cfg.AppKey != nil {
		signed, err := auth.SignBy(g.cfg.AppKey, auth.TDPrefix, td)
		if err != nil {
			return
		}
		v = signed
	}
	body, err := json.Marshal(v)
	if err != nil {
		return
	}
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, g.cfg.Channel+"/v1/things/td", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := g.client.Do(req)
	if err != nil {
		zlog.Warn().Err(err).Msg("[guard] auth unreachable")
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		why, _ := io.ReadAll(io.LimitReader(resp.Body, 256))
		zlog.Warn().Str("thing", id).Int("status", resp.StatusCode).Bytes("reason", why).Msg("[guard] TD refused by auth")
	}
}

// pull fetches a document newer than held; nil when there is none. A
// document refused (signature, size, recipient, older) keeps the held one.
func (g *Guard) pull(ctx context.Context, path, doc, thing string, held int64) (*auth.Doc, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, g.cfg.Channel+path, nil)
	if held > 0 {
		req.Header.Set("If-None-Match", `"`+strconv.FormatInt(held, 10)+`"`)
	}
	resp, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, nil // 304, or not adopted yet (404)
	}
	var b bytes.Buffer
	b.ReadFrom(http.MaxBytesReader(nil, resp.Body, auth.MaxDocSize+1))
	d, err := auth.ParseDoc(b.String(), g.keys, g.cfg.Deployment, doc, thing)
	if err != nil || d.Version < held {
		zlog.Warn().Err(err).Str("doc", doc).Str("thing", thing).Msg("[guard] Document refused, the held one kept")
		return nil, errRefused
	}
	if d.Version == held {
		return nil, nil // the held one again, without a 304
	}
	return &d, nil
}
