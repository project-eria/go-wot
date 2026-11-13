// Package inventory is the coordinator's list of the Things of its network
// (AUTH-DESIGN §2.4, F037): the coordinator alone browses DNS-SD; the services and
// eria-ui-server read the list and send their token only to its URLs.
package inventory

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Entry is one Thing announced on the network.
type Entry struct {
	ID   string `json:"id"`             // DNS-SD instance: the TD id without "urn:"
	Host string `json:"host"`           // mDNS host label (an ESP32's own name)
	URL  string `json:"url"`            // its TD
	Node string `json:"node,omitempty"` // machine announcing it
}

// Read fetches the `inventory` property of the coordinator whose TD is at
// coordinatorURL, with the caller's service token.
func Read(ctx context.Context, client *http.Client, coordinatorURL, token string) ([]Entry, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, strings.TrimSuffix(coordinatorURL, "/")+"/inventory", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("coordinator answered %d", resp.StatusCode)
	}
	var list []Entry
	return list, json.NewDecoder(resp.Body).Decode(&list)
}

// Origin is scheme://host:port of a URL, "" when it has none: a token sent to
// a Thing's TD goes to its interactions too, on the same origin.
func Origin(u string) string {
	p, err := url.Parse(u)
	if err != nil || p.Host == "" {
		return ""
	}
	scheme := p.Scheme
	switch scheme { // a Thing's WebSocket shares its HTTP server
	case "ws":
		scheme = "http"
	case "wss":
		scheme = "https"
	}
	host := p.Host
	if p.Port() == "" {
		host += map[string]string{"http": ":80", "https": ":443"}[scheme]
	}
	return scheme + "://" + strings.ToLower(host)
}
