// Package protocolDnssd implements DNS-Based Service Discovery (DNS-SD)
// for the W3C Web of Things, scoped to the `_wot._tcp` service type, or to
// `_wot-<network>._tcp` for a named network (see ServiceType).
//
// Reference: https://www.w3.org/TR/wot-discovery/#introduction-dns-sd-sec
package protocolDnssd

import (
	"errors"
	"fmt"
	"strings"
)

const (
	ServiceWoTTCP       = "_wot._tcp"
	SubtypeDirectoryTCP = "_directory._sub._wot._tcp"

	EntityThing     = "Thing"
	EntityDirectory = "Directory"

	SchemeHTTP  = "http"
	SchemeHTTPS = "https"

	defaultEntityType = EntityThing
	defaultScheme     = SchemeHTTP

	maxTXTValueLen = 255
)

// TXT represents the parsed contents of a `_wot._tcp` TXT record.
//
// Per the spec, all three fields have defaults: `type` defaults to "Thing",
// `scheme` defaults to "http", and `td` is mandatory.
//
// Extra holds the non-standard keys (e.g. the `status` announced by ERIA
// ESP32 devices), so orchestration layers can react to re-announcements
// without polling the device over HTTP.
type TXT struct {
	TD     string
	Type   string
	Scheme string
	Extra  map[string]string
}

func validScheme(s string) bool {
	switch s {
	case SchemeHTTP, SchemeHTTPS:
		return true
	}
	return false
}

func validEntityType(s string) bool {
	switch s {
	case EntityThing, EntityDirectory:
		return true
	}
	return false
}

// maxNetworkLen keeps "wot-<network>" within the 15-character limit of a
// DNS-SD service name (RFC 6335 §5.1)
const maxNetworkLen = 11

// ServiceType returns the DNS-SD service type of a network: the standard
// `_wot._tcp` for the default network (""), `_wot-<network>._tcp` for a
// named one. Advertisers and browsers of different networks never see each
// other, so several ERIA houses can share a LAN.
//
// A network name is 1 to 11 characters among [a-z0-9-], with no hyphen at
// either end and no two in a row.
func ServiceType(network string) (string, error) {
	if network == "" {
		return ServiceWoTTCP, nil
	}
	if len(network) > maxNetworkLen {
		return "", fmt.Errorf("network %q: at most %d characters", network, maxNetworkLen)
	}
	if network[0] == '-' || network[len(network)-1] == '-' || strings.Contains(network, "--") {
		return "", fmt.Errorf("network %q: no hyphen at either end or twice in a row", network)
	}
	for _, c := range network {
		if (c < 'a' || c > 'z') && (c < '0' || c > '9') && c != '-' {
			return "", fmt.Errorf("network %q: only lowercase letters, digits and hyphens", network)
		}
	}
	return "_wot-" + network + "._tcp", nil
}

// BuildTXT serializes a TXT record for the `_wot._tcp` service type.
//
// `td` must be an absolute path beginning with "/". `entityType` and `scheme`
// may be empty, in which case the spec defaults are emitted explicitly so
// downstream tools see complete records.
func BuildTXT(td, entityType, scheme string) ([]string, error) {
	if td == "" {
		return nil, errors.New("td is required")
	}
	if !strings.HasPrefix(td, "/") {
		return nil, fmt.Errorf("td must be an absolute path, got %q", td)
	}
	if entityType == "" {
		entityType = defaultEntityType
	}
	if !validEntityType(entityType) {
		return nil, fmt.Errorf("type must be %q or %q, got %q", EntityThing, EntityDirectory, entityType)
	}
	if scheme == "" {
		scheme = defaultScheme
	}
	if !validScheme(scheme) {
		return nil, fmt.Errorf("scheme must be %q or %q, got %q", SchemeHTTP, SchemeHTTPS, scheme)
	}
	out := []string{
		"td=" + td,
		"type=" + entityType,
		"scheme=" + scheme,
	}
	for _, kv := range out {
		if len(kv) > maxTXTValueLen {
			return nil, fmt.Errorf("TXT entry exceeds 255 bytes: %q", kv)
		}
	}
	return out, nil
}

// ParseTXT parses a slice of TXT key=value strings, applies spec defaults,
// and validates the result. Unknown keys are ignored (forward compatibility).
func ParseTXT(records []string) (TXT, error) {
	t := TXT{}
	for _, r := range records {
		if len(r) > maxTXTValueLen {
			return TXT{}, fmt.Errorf("TXT entry exceeds 255 bytes")
		}
		eq := strings.IndexByte(r, '=')
		if eq < 0 {
			continue
		}
		key, val := r[:eq], r[eq+1:]
		switch key {
		case "td":
			t.TD = val
		case "type":
			t.Type = val
		case "scheme":
			t.Scheme = val
		default:
			// Keep unknown keys (forward compatibility, ERIA `status`...)
			if t.Extra == nil {
				t.Extra = map[string]string{}
			}
			t.Extra[key] = val
		}
	}
	if t.TD == "" {
		return TXT{}, errors.New("td field is required")
	}
	if !strings.HasPrefix(t.TD, "/") {
		return TXT{}, fmt.Errorf("td must be an absolute path, got %q", t.TD)
	}
	if t.Type == "" {
		t.Type = defaultEntityType
	}
	if !validEntityType(t.Type) {
		return TXT{}, fmt.Errorf("invalid type %q", t.Type)
	}
	if t.Scheme == "" {
		t.Scheme = defaultScheme
	}
	if !validScheme(t.Scheme) {
		return TXT{}, fmt.Errorf("invalid scheme %q", t.Scheme)
	}
	return t, nil
}
