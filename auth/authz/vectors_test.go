package authz

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"flag"
	"os"
	"testing"
	"time"

	"github.com/project-eria/go-wot/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The common vectors of §2.7: one file, run by go-wot here and later by
// rs-wot and the relays; a vector accepted by one implementation only blocks
// the merge. Regenerate with: go test ./auth/authz -run TestVectors -update
var update = flag.Bool("update", false, "regenerate testdata/vectors.json")

const vectorsPath = "testdata/vectors.json"

type vectorFunction struct {
	Kind string `json:"kind"`
	Name string `json:"name,omitempty"`
	Op   Op     `json:"op"`
}

type vector struct {
	Name     string          `json:"name"`
	Mode     string          `json:"mode"` // normal, minimal, offline
	Function vectorFunction  `json:"function"`
	Input    json.RawMessage `json:"input,omitempty"`
	Token    string          `json:"token"`
	Confirm  string          `json:"confirm,omitempty"`
	Seen     bool            `json:"seen,omitempty"` // the confirm jti was already used
	Allow    bool            `json:"allow"`
	Step     int             `json:"step"`
}

type vectorFile struct {
	Now        int64       `json:"now"`
	Aud        string      `json:"aud"`
	Thing      string      `json:"thing"`
	Zones      []string    `json:"zones"`
	Keyset     auth.Keyset `json:"keyset"`
	Table      Table       `json:"table"`
	Revocation string      `json:"revocation"` // signed document
	Vectors    []vector    `json:"vectors"`
}

func TestVectors(t *testing.T) {
	if *update {
		generate(t)
	}
	raw, err := os.ReadFile(vectorsPath)
	require.NoError(t, err)
	var f vectorFile
	require.NoError(t, json.Unmarshal(raw, &f))
	keys, err := f.Keyset.Lookup()
	require.NoError(t, err)
	now := time.Unix(f.Now, 0)
	revoked, err := auth.ParseDoc(f.Revocation, keys, f.Aud, auth.DocRevocation, f.Thing)
	require.NoError(t, err)
	modes := map[string]Mode{"normal": Normal, "minimal": Minimal, "offline": Offline}

	for _, v := range f.Vectors {
		t.Run(v.Name, func(t *testing.T) {
			d := Decide(Check{
				Request:  Request{Mode: modes[v.Mode], Verifier: Thing, Token: v.Token, Keys: keys, Aud: f.Aud, Now: now},
				Thing:    f.Thing,
				Zones:    f.Zones,
				Table:    f.Table,
				Revoked:  &revoked,
				Function: Function{Kind: v.Function.Kind, Name: v.Function.Name, Op: v.Function.Op},
				Input:    v.Input,
				Confirm:  v.Confirm,
				Seen:     func(string, int64) bool { return v.Seen },
			})
			assert.Equal(t, v.Allow, d.Allow, d.Reason)
			if !v.Allow {
				assert.Equal(t, v.Step, d.Step, d.Reason)
			}
		})
	}
}

// generate writes the vectors, signed by a fixed key so the file is stable.
func generate(t *testing.T) {
	priv := ed25519.NewKeyFromSeed(make([]byte, ed25519.SeedSize))
	pub, err := auth.ParsePublicKey(auth.Ed25519, priv.Public().(ed25519.PublicKey))
	require.NoError(t, err)
	now := time.Unix(1_800_000_000, 0)
	const thing = "urn:eria:garage:1"
	f := vectorFile{
		Now: now.Unix(), Aud: "eria-maison", Thing: thing, Zones: []string{"garage"},
		Keyset: auth.Keyset{Generation: 1, Issued: now.Unix(), Keys: []auth.KeysetKey{{
			Kid: pub.KID(), Pubkey: base64.RawURLEncoding.EncodeToString(pub.Bytes()), Usage: auth.UsageSign}}},
		Table: Table{
			"property/temperature/read":  {Levels: [3]int{9, 9, 8}},
			"property/light/write":       {Levels: [3]int{9, 8, None}, Confirm: Never},
			"property/state/read":        {Levels: [3]int{9, 1, 1}},
			"event/moved/read":           {Levels: [3]int{9, 1, 1}},
			"action/open/invoke":         {Levels: [3]int{1, 1, 1}, Confirm: OutsideRoom},
			"action/disarm/invoke":       {Levels: [3]int{1, 1, 1}, Confirm: Always},
			"action/unsetContext/invoke": {Levels: [3]int{1, 1, 1}, Confirm: OutsideHome},
			"action/setContext/invoke":   {Levels: [3]int{3, 3, 3}, Confirm: Never},
			"action/sealed/invoke":       {Levels: [3]int{None, None, None}},
		},
	}
	sign := func(c auth.Claims) string {
		c.Iss, c.Aud = "auth", f.Aud
		if c.Iat == 0 {
			c.Iat = now.Unix() - 10
		}
		if c.Exp == 0 {
			c.Exp = c.Iat + 240
		}
		if c.Jti == "" {
			c.Jti = "j-" + c.Term + c.Sub
		}
		tok, err := auth.Sign(priv, pub.KID(), c)
		require.NoError(t, err)
		return tok
	}
	lvl := func(n int) *int { return &n }
	yes := func() *bool { b := true; return &b }
	rev, err := auth.SignDoc(priv, pub.KID(), auth.Doc{Iss: "auth", Aud: f.Aud, Doc: auth.DocRevocation, Thing: "*",
		Version: now.UnixMilli(), Issued: now.Unix(), Revoked: []auth.Revocation{{Subject: "term:vole", Before: now.Unix()}}})
	require.NoError(t, err)
	f.Revocation = rev

	// Tokens: screens in the garage and the entrance, Léa (2) and a parent (1)
	// on a phone, from inside and from outside, and services.
	garageScreen := sign(auth.Claims{Kind: auth.KindUser, Term: "ecran-garage", Zone: "garage", Lvl: lvl(9)})
	roomOnly := sign(auth.Claims{Kind: auth.KindUser, Term: "ecran-salon", Zone: "salon", RoomOnly: true, Lvl: lvl(9)})
	stolen := sign(auth.Claims{Kind: auth.KindUser, Term: "vole", Zone: "garage", Lvl: lvl(9)})
	lea := sign(auth.Claims{Kind: auth.KindUser, Sub: "profil:lea", Term: "tel-lea", Lvl: lvl(2)})
	leaOut := sign(auth.Claims{Kind: auth.KindUser, Sub: "profil:lea", Term: "tel-lea", Lvl: lvl(2), Out: yes()})
	leaGrant := sign(auth.Claims{Kind: auth.KindUser, Sub: "profil:lea", Term: "tel-lea", Lvl: lvl(2), Grants: []string{"room:garage:*"}})
	leaExact := sign(auth.Claims{Kind: auth.KindUser, Sub: "profil:lea", Term: "tel-lea", Lvl: lvl(2), Grants: []string{"thing:" + thing + ":invoke"}})
	parentEntrance := sign(auth.Claims{Kind: auth.KindUser, Sub: "profil:paul", Term: "ecran-entree", Zone: "entree", Lvl: lvl(1)})
	parentGarage := sign(auth.Claims{Kind: auth.KindUser, Sub: "profil:paul", Term: "ecran-garage", Zone: "garage", Lvl: lvl(1)})
	parentOut := sign(auth.Claims{Kind: auth.KindUser, Sub: "profil:paul", Term: "tel-paul", Lvl: lvl(1), Out: yes()})
	parentPhone := sign(auth.Claims{Kind: auth.KindUser, Sub: "profil:paul", Term: "tel-paul", Lvl: lvl(1)})
	service := sign(auth.Claims{Kind: auth.KindService, Sub: "app:a1", Lvl: lvl(9), Grants: []string{"thing:" + thing + ":invoke"}, Iat: now.Unix() - 10, Exp: now.Unix() + 15*60 - 10})
	serviceShort := sign(auth.Claims{Kind: auth.KindService, Sub: "app:a1", Lvl: lvl(9), Grants: []string{"thing:" + thing + ":invoke"}, Iat: now.Unix() - 10, Exp: now.Unix() + 290})
	serviceRoom := sign(auth.Claims{Kind: auth.KindService, Sub: "app:a1", Lvl: lvl(9), Grants: []string{"room:garage:*"}})

	open := vectorFunction{Kind: "action", Name: "open", Op: Invoke}
	disarm := vectorFunction{Kind: "action", Name: "disarm", Op: Invoke}
	unset := vectorFunction{Kind: "action", Name: "unsetContext", Op: Invoke}
	set := vectorFunction{Kind: "action", Name: "setContext", Op: Invoke}
	confirm := func(sub, term string, fn vectorFunction, input string) string {
		act, ok := Act(thing, Function{Kind: fn.Kind, Name: fn.Name, Op: fn.Op}, []byte(input))
		require.True(t, ok)
		return sign(auth.Claims{Kind: auth.KindConfirm, Sub: sub, Term: term, Act: act, Exp: now.Unix() + 50, Jti: "c-" + term + fn.Name})
	}
	in := json.RawMessage(`{"b":1,"a":"x"}`)
	reordered := json.RawMessage(`{ "a": "x", "b": 1.0 }`)

	f.Vectors = []vector{
		{Name: "td read by a screen", Function: vectorFunction{Kind: "td", Op: Read}, Token: roomOnly, Allow: false, Step: 5},
		{Name: "td read by a screen in the room", Function: vectorFunction{Kind: "td", Op: Read}, Token: garageScreen, Allow: true},
		{Name: "td read in minimal mode from a phone", Mode: "minimal", Function: vectorFunction{Kind: "td", Op: Read}, Token: sign(auth.Claims{Kind: auth.KindUser, Term: "tel-x", Lvl: lvl(9)}), Allow: true},
		{Name: "offline refuses all", Mode: "offline", Function: vectorFunction{Kind: "td", Op: Read}, Token: garageScreen, Step: 0},
		{Name: "no token", Function: vectorFunction{Kind: "td", Op: Read}, Token: "", Step: 1},
		{Name: "revoked terminal", Function: vectorFunction{Kind: "property", Name: "temperature", Op: Read}, Token: stolen, Step: 4},
		{Name: "temperature, screen, local", Function: vectorFunction{Kind: "property", Name: "temperature", Op: Read}, Token: garageScreen, Allow: true},
		{Name: "temperature, outside, level 2", Function: vectorFunction{Kind: "property", Name: "temperature", Op: Read}, Token: leaOut, Allow: true},
		{Name: "light, remote level 2", Function: vectorFunction{Kind: "property", Name: "light", Op: Write}, Token: lea, Allow: true},
		{Name: "light, outside is None", Function: vectorFunction{Kind: "property", Name: "light", Op: Write}, Token: parentOut, Step: 7},
		{Name: "state, remote level 2", Function: vectorFunction{Kind: "property", Name: "state", Op: Read}, Token: lea, Step: 7},
		{Name: "state, room grant elevates", Function: vectorFunction{Kind: "property", Name: "state", Op: Read}, Token: leaGrant, Allow: true},
		{Name: "event, room grant elevates", Function: vectorFunction{Kind: "event", Name: "moved", Op: Read}, Token: leaGrant, Allow: true},
		{Name: "unknown function is None", Function: vectorFunction{Kind: "property", Name: "unknown", Op: Read}, Token: parentPhone, Step: 7},
		{Name: "None even with a grant", Function: vectorFunction{Kind: "action", Name: "sealed", Op: Invoke}, Token: leaExact, Step: 7},
		{Name: "open, parent local, no confirmation", Function: open, Token: parentGarage, Allow: true},
		{Name: "open, parent remote, no confirm", Function: open, Token: parentPhone, Step: 8},
		{Name: "open, parent remote, confirm", Function: open, Input: in, Token: parentPhone, Confirm: confirm("profil:paul", "tel-paul", open, `{"a":"x","b":1}`), Allow: true},
		{Name: "open, confirm on canonical input", Function: open, Input: reordered, Token: parentPhone, Confirm: confirm("profil:paul", "tel-paul", open, string(in)), Allow: true},
		{Name: "open, confirm for another input", Function: open, Input: json.RawMessage(`{"a":"y"}`), Token: parentPhone, Confirm: confirm("profil:paul", "tel-paul", open, string(in)), Step: 8},
		{Name: "open, confirm of another terminal", Function: open, Token: parentPhone, Confirm: confirm("profil:paul", "ecran-entree", open, ""), Step: 8},
		{Name: "open, confirm replayed", Function: open, Token: parentPhone, Confirm: confirm("profil:paul", "tel-paul", open, ""), Seen: true, Step: 8},
		{Name: "open, confirm as main token", Function: open, Token: confirm("profil:paul", "tel-paul", open, ""), Step: 3},
		{Name: "open, room grant does not cover a confirmable", Function: open, Token: leaGrant, Step: 7},
		{Name: "open, exact grant, remote, confirm", Function: open, Token: leaExact, Confirm: confirm("profil:lea", "tel-lea", open, ""), Allow: true},
		{Name: "disarm, parent local still needs confirm", Function: disarm, Token: parentGarage, Step: 8},
		{Name: "disarm, service never", Function: disarm, Token: serviceShort, Step: 8},
		{Name: "open, service 15 min token", Function: open, Token: service, Step: 8},
		{Name: "open, service room grant", Function: open, Token: serviceRoom, Step: 7},
		{Name: "open, service exact grant 5 min", Function: open, Token: serviceShort, Allow: true},
		// The confirm_when = outside vectors of §2.7.
		{Name: "unsetContext, service 15 min, remote", Function: unset, Token: service, Step: 8},
		{Name: "unsetContext, PIN level 1 at the entrance", Function: unset, Token: parentEntrance, Allow: true},
		{Name: "unsetContext, level 1 outside without confirm", Function: unset, Token: parentOut, Step: 8},
		{Name: "unsetContext, level 1 outside with confirm", Function: unset, Token: parentOut, Confirm: confirm("profil:paul", "tel-paul", unset, ""), Allow: true},
		{Name: "unsetContext, level 2 remote", Function: unset, Token: lea, Step: 7},
		{Name: "setContext, level 3 outside", Function: set, Token: sign(auth.Claims{Kind: auth.KindUser, Sub: "profil:ami", Term: "tel-ami", Lvl: lvl(3), Out: yes()}), Allow: true},
		// Minimal mode (§4.4): terminal tokens only, local 9 without confirmation.
		{Name: "minimal, temperature local", Mode: "minimal", Function: vectorFunction{Kind: "property", Name: "temperature", Op: Read}, Token: garageScreen, Allow: true},
		{Name: "minimal, light local", Mode: "minimal", Function: vectorFunction{Kind: "property", Name: "light", Op: Write}, Token: garageScreen, Allow: true},
		{Name: "minimal, nothing remote", Mode: "minimal", Function: vectorFunction{Kind: "property", Name: "temperature", Op: Read}, Token: sign(auth.Claims{Kind: auth.KindUser, Term: "ecran-entree", Zone: "entree", Lvl: lvl(9)}), Step: 7},
		{Name: "minimal, person refused", Mode: "minimal", Function: vectorFunction{Kind: "property", Name: "temperature", Op: Read}, Token: parentGarage, Step: 3},
	}
	for i := range f.Vectors {
		if f.Vectors[i].Mode == "" {
			f.Vectors[i].Mode = "normal"
		}
	}
	out, err := json.MarshalIndent(f, "", "  ")
	require.NoError(t, err)
	require.NoError(t, os.MkdirAll("testdata", 0o755))
	require.NoError(t, os.WriteFile(vectorsPath, append(out, '\n'), 0o644))
}
