// Package authz is the decision function of AUTH-DESIGN §2.7, shared by the
// Things (go-wot), the relays and the auth service. Every early return is a
// refusal; Allow is only ever the last line.
//
// Authenticate runs the token steps (0 to 3), enough for the relays, the
// master and auth; Decide runs them all (0 to 9) for a Thing.
package authz

import (
	"fmt"
	"slices"
	"time"

	"github.com/project-eria/go-wot/auth"
)

// Mode of the verifier (§4.4).
type Mode int

const (
	Normal Mode = iota
	Minimal
	Offline
)

func (m Mode) String() string { return [...]string{"normal", "minimal", "offline"}[m] }

// Verifier is who checks the token: it decides which rows of the token table
// are accepted (the "accepted by" column of §2.7).
type Verifier int

const (
	Thing Verifier = iota
	Relay
	AuthService
	Master
)

var accepted = map[Verifier]map[auth.Profile]bool{
	Thing:       {auth.ProfileTerminal: true, auth.ProfilePerson: true, auth.ProfileService: true},
	Relay:       {auth.ProfileTerminal: true, auth.ProfilePerson: true, auth.ProfileService: true},
	AuthService: {auth.ProfileAdmin: true, auth.ProfileService: true},
	Master:      {auth.ProfileAdmin: true, auth.ProfileService: true}, // service: the inventory only (F037)
}

// Request holds what the token steps need.
type Request struct {
	Mode     Mode
	Verifier Verifier
	Token    string
	Keys     auth.KeyLookup
	Aud      string // deployment id
	Now      time.Time
}

// Decision is Allow, or the step that refused and why (logged; the client
// only gets 401 / 403).
type Decision struct {
	Allow   bool
	Step    int
	Reason  string
	Claims  auth.Claims
	Profile auth.Profile
	// Confirm: refused at step 8 for want of a confirmation that this user
	// could give (§2.8); the interfaces offer it instead of greying out.
	Confirm bool
	// Confirmable: an allowed command whose rule asks a confirmation
	// (journaled a year, §4.5).
	Confirmable bool
	// Of an allowed request: where from, and the grant that lifted the
	// level, "" without (the daily summary, §4.5).
	Locality Locality
	Grant    string
}

// ReasonReused is the reason of a confirm token used twice (§4.5 alert).
const ReasonReused = "confirm token already used"

func deny(step int, reason string) Decision { return Decision{Step: step, Reason: reason} }

// Authenticate runs steps 0 to 3: mode, structure, signature, claims.
func Authenticate(r Request) Decision {
	if r.Mode == Offline {
		return deny(0, "offline")
	}
	c, err := auth.Parse(r.Token, r.Keys)
	if err != nil {
		e := err.(*auth.TokenError)
		return deny(e.Step, e.Reason)
	}
	p, err := c.Check(r.Now, r.Aud)
	if err != nil {
		return deny(3, err.Error())
	}
	if !accepted[r.Verifier][p] {
		return deny(3, "token not accepted by this verifier")
	}
	if r.Mode == Minimal && p != auth.ProfileTerminal {
		return deny(3, "only a terminal token in minimal mode")
	}
	return Decision{Allow: true, Claims: c, Profile: p}
}

// Check holds what a Thing needs to decide on one request.
type Check struct {
	Request
	Thing    string    // id of the deciding Thing
	Zones    []string  // its room tags, without "room:"
	Table    Table     // effective table; the minimal one is derived from it
	Revoked  *auth.Doc // revocation list held, nil if none yet
	Function Function
	Input    []byte // body of a write or invoke, for act
	Confirm  string // X-Eria-Confirm header
	// Seen records the jti of a confirm token until exp; it returns true
	// when the jti was already used or the list is full.
	Seen func(jti string, exp int64) bool
}

// Decide runs steps 0 to 9 (§2.7) for a Thing.
func Decide(k Check) Decision {
	d := Authenticate(k.Request)
	if !d.Allow {
		return d
	}
	c := d.Claims
	// Refusals from here on keep the token's claims, for the journal.
	deny := func(step int, reason string) Decision {
		return Decision{Step: step, Reason: reason, Claims: c, Profile: d.Profile}
	}
	if c.Thing != "" && c.Thing != k.Thing {
		return deny(3, "token for another Thing")
	}
	// 4. Revocation.
	if k.Revoked != nil && k.Revoked.Revokes(c) {
		return deny(4, "revoked")
	}
	// 5. Locality.
	loc := Remote
	switch {
	case c.Out != nil && *c.Out:
		loc = Outside
	case c.Zone != "" && slices.Contains(k.Zones, c.Zone):
		loc = Local
	}
	if c.RoomOnly && loc != Local {
		return deny(5, "room only, from "+loc.String())
	}
	// 6. Requirement.
	r := rule(k.Table, k.Function, k.Mode == Minimal)
	e := r.Levels[loc]
	command := k.Function.Op == Write || k.Function.Op == Invoke
	confirmable := command && confirmable(r)
	// 7. Level, after grant.
	lvl := *c.Lvl
	grant := ""
	for _, g := range c.Grants {
		if covers(g, k.Thing, k.Zones, k.Function, confirmable) {
			grant = g
			lvl = 1
			break
		}
	}
	if e == None || lvl > e {
		return deny(7, fmt.Sprintf("level %d, %s requires %d", lvl, loc, e))
	}
	// 8. Proof. A service reaching here on a confirmable command holds an
	// exact grant (step 7: level 9 against a ceiling of 1) and a 5 min token
	// (step 3: its longest life).
	if command && needsConfirm(r, loc) {
		if d.Profile == auth.ProfileService {
			if r.Confirm == Always {
				return deny(8, "no service on an always confirmed command")
			}
		} else if reason := k.confirmed(c); reason != "" {
			d := deny(8, reason)
			d.Confirm = k.Confirm == ""
			return d
		}
	}
	// 9. Allow.
	d.Confirmable, d.Locality, d.Grant = confirmable, loc, grant
	return d
}

// confirmed checks the confirm token against the user token (§2.8), and
// returns why it fails, or "".
func (k Check) confirmed(user auth.Claims) string {
	if k.Confirm == "" {
		return "confirmation required"
	}
	c, err := auth.Parse(k.Confirm, k.Keys)
	if err != nil {
		return "confirm token unreadable"
	}
	if p, err := c.Check(k.Now, k.Aud); err != nil || p != auth.ProfileConfirm {
		return "not a valid confirm token"
	}
	act, ok := Act(k.Thing, k.Function, k.Input)
	if !ok || c.Act != act || c.Sub != user.Sub || c.Term != user.Term {
		return "confirm token for another command or person"
	}
	if k.Seen == nil || k.Seen(c.Jti, c.Exp) {
		return ReasonReused
	}
	return ""
}
