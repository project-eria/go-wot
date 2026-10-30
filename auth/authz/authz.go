// Package authz is the decision function of AUTH-DESIGN §2.7, shared by the
// Things (go-wot), the relays and the auth service. Every early return is a
// refusal; Allow is only ever the last line.
//
// This first part covers the token steps (0 to 3). The rights steps (4 to 9:
// revocation, locality, requirement, level, proof) arrive with the per-Thing
// tables (§4.7, step 2) and build on Authenticate.
package authz

import (
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
