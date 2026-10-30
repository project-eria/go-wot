package auth

import (
	"errors"
	"time"
)

// Profile is a row of the token table (AUTH-DESIGN §2.7): the only claim
// combinations accepted.
type Profile string

const (
	ProfileTerminal Profile = "terminal" // signed-out terminal, level 9
	ProfilePerson   Profile = "person"   // signed-in person, level 1-8
	ProfileAdmin    Profile = "admin"    // admin session, level 0
	ProfileService  Profile = "service"
	ProfileConfirm  Profile = "confirm"
)

// MaxDuration is the longest exp − iat of each profile.
var MaxDuration = map[Profile]time.Duration{
	ProfileTerminal: 72 * time.Hour,
	ProfilePerson:   5 * time.Minute,
	ProfileAdmin:    5 * time.Minute,
	ProfileService:  4 * time.Hour,
	ProfileConfirm:  60 * time.Second,
}

// ClockSkew is the tolerance on exp and iat.
const ClockSkew = 30 * time.Second

// Bounds on claims (AUTH-DESIGN §2.7, ESP32 limits).
const (
	MaxGrants    = 32
	MaxClaimSize = 128
)

// Profile returns the table row the claims match, or an error for any other
// combination.
func (c Claims) Profile() (Profile, error) {
	hasLvl := c.Lvl != nil
	lvl := -1
	if hasLvl {
		lvl = *c.Lvl
	}
	switch c.Kind {
	case KindUser:
		switch {
		case c.Term == "" || c.Act != "":
		case c.Sub == "" && lvl == 9 && len(c.Grants) == 0:
			return ProfileTerminal, nil
		case c.Sub != "" && lvl >= 1 && lvl <= 8:
			return ProfilePerson, nil
		case c.Sub != "" && lvl == 0 && len(c.Grants) == 0:
			return ProfileAdmin, nil
		}
	case KindService:
		if c.Sub == "" && c.Term == "" && lvl == 9 && c.Act == "" {
			return ProfileService, nil
		}
	case KindConfirm:
		if c.Sub != "" && c.Term != "" && !hasLvl && len(c.Grants) == 0 && c.Act != "" {
			return ProfileConfirm, nil
		}
	}
	return "", errors.New("claim combination not allowed")
}

// Check is step 3 without the verifier-specific rules: required and bounded
// claims, issuer and audience, table row and times. It returns the row.
func (c Claims) Check(now time.Time, aud string) (Profile, error) {
	if c.Iss != "auth" || c.Aud != aud || c.Jti == "" {
		return "", errors.New("bad iss, aud or jti")
	}
	if len(c.Grants) > MaxGrants {
		return "", errors.New("too many grants")
	}
	for _, s := range append([]string{c.Sub, c.Term, c.Zone, c.Act, c.Jti}, c.Grants...) {
		if len(s) > MaxClaimSize {
			return "", errors.New("claim too long")
		}
	}
	p, err := c.Profile()
	if err != nil {
		return "", err
	}
	exp, iat := time.Unix(c.Exp, 0), time.Unix(c.Iat, 0)
	if !exp.After(now.Add(-ClockSkew)) || iat.After(now.Add(ClockSkew)) {
		return "", errors.New("expired or issued in the future")
	}
	if !exp.After(iat) || exp.Sub(iat) > MaxDuration[p] {
		return "", errors.New("lifetime too long")
	}
	return p, nil
}
