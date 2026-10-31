package authz

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"slices"
	"strings"

	"github.com/project-eria/go-wot/auth"
)

// Op is the operation of a function, derived from the WoT forms (§2.4):
// readproperty, observeproperty, subscribeevent → read; writeproperty →
// write; invokeaction → invoke.
type Op string

const (
	Read   Op = "read"
	Write  Op = "write"
	Invoke Op = "invoke"
)

// Function is an (affordance, operation) pair. Kind is "property", "action",
// "event", or "td" for the Thing Description itself.
type Function struct {
	Kind string
	Name string
	Op   Op
}

// TD is the implicit row of every table (G03): 9 / 9 / 9, never overridden.
var TD = Function{Kind: "td", Op: Read}

// Key is the function's entry in a table: "property/on/write".
func (f Function) Key() string { return f.Kind + "/" + f.Name + "/" + string(f.Op) }

// Locality of a request (§2.2).
type Locality int

const (
	Local Locality = iota
	Remote
	Outside
)

func (l Locality) String() string { return [...]string{"local", "remote", "outside"}[l] }

// The table types live in package auth, which signs and parses the table
// documents.
type (
	Rule    = auth.Rule
	Table   = auth.Table
	Confirm = auth.Confirm
)

const (
	None        = auth.None
	Never       = auth.Never
	OutsideHome = auth.OutsideHome
	OutsideRoom = auth.OutsideRoom
	Always      = auth.Always
)

// confirmable tells whether the rule asks a confirmation in some locality.
func confirmable(r Rule) bool { return r.Confirm != "" && r.Confirm != Never }

// needsConfirm tells whether a write or invoke from locality l needs a
// confirmation (step 8).
func needsConfirm(r Rule, l Locality) bool {
	switch r.Confirm {
	case Always:
		return true
	case OutsideRoom:
		return l != Local
	case OutsideHome:
		return l == Outside
	}
	return false
}

// rule returns the requirement of f, from the minimal table derived from t
// in minimal mode (§4.4): 9 locally when f needs no confirmation and its
// local requirement is 9, None otherwise.
func rule(t Table, f Function, minimal bool) Rule {
	if f.Kind == "td" {
		return Rule{Levels: [3]int{9, 9, 9}}
	}
	r, ok := t[f.Key()]
	if !ok {
		return Rule{}
	}
	if minimal {
		if confirmable(r) || r.Levels[Local] != 9 {
			return Rule{}
		}
		return Rule{Levels: [3]int{9, None, None}}
	}
	return r
}

// covers tells whether a grant covers f on the Thing id with zones (§2.6):
// <scope>:<selector>:<op>, scope thing or room. A confirmable function is
// only covered by an exact thing:<id>:<op> entry. Thing ids may hold ":".
func covers(grant string, id string, zones []string, f Function, confirmable bool) bool {
	scope, rest, ok := strings.Cut(grant, ":")
	i := strings.LastIndexByte(rest, ':')
	if !ok || i < 0 {
		return false
	}
	selector, op := rest[:i], rest[i+1:]
	if op != string(f.Op) && (op != "*" || confirmable) {
		return false
	}
	switch scope {
	case "thing":
		return selector == id
	case "room":
		return !confirmable && slices.Contains(zones, selector)
	}
	return false
}

// Act is the hash a confirmation is bound to (§2.8):
// sha256(thing ‖ affordance ‖ canonical JSON of the input), the parts
// separated by a NUL byte, in base64url. The affordance is "kind/name"; an
// empty input is null. ok is false for an input that is not JSON.
//
// ponytail: the canonical form is encoding/json's (sorted keys, no HTML
// escaping, shortest numbers, integers beyond 2^53 rounded as in
// JavaScript); full RFC 8785 if an implementation ever disagrees on a vector.
func Act(thing string, f Function, input []byte) (string, bool) {
	var v any
	if len(strings.TrimSpace(string(input))) > 0 {
		if err := json.Unmarshal(input, &v); err != nil {
			return "", false
		}
	}
	var b strings.Builder
	enc := json.NewEncoder(&b)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return "", false
	}
	canonical := strings.TrimSuffix(b.String(), "\n")
	h := sha256.Sum256([]byte(thing + "\x00" + f.Kind + "/" + f.Name + "\x00" + canonical))
	return base64.RawURLEncoding.EncodeToString(h[:]), true
}
