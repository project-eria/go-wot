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
	Master:      {auth.ProfileAdmin: true},
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
}

func deny(step int, reason string) Decision { return Decision{Step: step, Reason: reason} }

// Authenticate runs steps 0 to 3: mode, structure, signature, claims.
func Authenticate(r Request) Decision {
	if r.Mode == Offline {
		return deny(0, "offline")
	}
	c, legacy, err := auth.Parse(r.Token, r.Keys)
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
	if (r.Mode == Minimal || legacy) && p != auth.ProfileTerminal {
		return deny(3, "only a terminal token in minimal mode or with a legacy key")
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
	granted := false
	for _, g := range c.Grants {
		if covers(g, k.Thing, k.Zones, k.Function, confirmable) {
			granted = true
			lvl = 1
			break
		}
	}
	if e == None || lvl > e {
		return deny(7, fmt.Sprintf("level %d, %s requires %d", lvl, loc, e))
	}
	// 8. Proof.
	if confirmable && d.Profile == auth.ProfileService && (!granted || c.Exp-c.Iat > 300) {
		return deny(8, "service on a confirmable command: exact grant and 5 min token")
	}
	if command && needsConfirm(r, loc) {
		if d.Profile == auth.ProfileService {
			if r.Confirm == Always {
				return deny(8, "no service on an always confirmed command")
			}
		} else if reason := k.confirmed(c); reason != "" {
			return deny(8, reason)
		}
	}
	// 9. Allow.
	return d
}

// confirmed checks the confirm token against the user token (§2.8), and
// returns why it fails, or "".
func (k Check) confirmed(user auth.Claims) string {
	if k.Confirm == "" {
		return "confirmation required"
	}
	c, legacy, err := auth.Parse(k.Confirm, k.Keys)
	if err != nil || legacy {
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
		return "confirm token already used"
	}
	return ""
}
