package protocolDnssd

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/brutella/dnssd"
	zlog "github.com/rs/zerolog/log"
)

// Entry is a single discovered `_wot._tcp` service after TXT parsing.
type Entry struct {
	Instance string
	Host     string
	Port     int
	Addrs    []net.IP
	TXT      TXT
	// TDURL is the fully-qualified URL of the Thing Description, built from
	// the TXT scheme + resolved address + port + td path.
	TDURL *url.URL
}

// BrowseOptions configures a discovery session.
type BrowseOptions struct {
	// Domain to browse, defaults to "local."
	Domain string
	// PreferIPv4 selects the first IPv4 address (if any) when building TDURL.
	// When false, IPv6 is used if present.
	PreferIPv4 bool
	// AllowPrivate allows private/loopback/link-local addresses in TDURL.
	// Defaults to false to mitigate SSRF on untrusted networks.
	AllowPrivate bool
}

// Browse runs an mDNS browse for `_wot._tcp` until ctx is cancelled,
// emitting validated Entry values on the returned channel. The channel
// is closed when ctx is done.
//
// A process that both advertises (Advertiser) and browses will discover
// its own announcements: mDNS has no notion of "self". Filter by instance
// name on the caller side if self-discovery is unwanted.
// browseRound bounds one lookup: each round starts from a fresh cache and
// sends a new PTR query. brutella/dnssd reports an instance only once per
// lookup, and the first sighting of a newly starting service is often its
// probe (SRV without TXT) — unusable and never reported again. Restarting
// the lookup makes such a service visible at the next round at the latest,
// with the complete records a direct query returns.
var browseRound = 10 * time.Second // var: shortened by tests

// Browse emits every `_wot._tcp` instance, once per round (consumers must
// accept repeats), until ctx is done; the channel is then closed.
func Browse(ctx context.Context, opts BrowseOptions) (<-chan Entry, error) {
	if opts.Domain == "" {
		opts.Domain = "local."
	}
	service := ServiceWoTTCP + "." + strings.TrimSuffix(opts.Domain, ".") + "."
	out := make(chan Entry, 16)

	go func() {
		defer close(out)
		add := func(raw dnssd.BrowseEntry) {
			entry, err := convert(raw, opts)
			if err != nil {
				zlog.Debug().Err(err).Str("instance", raw.Name).
					Msg("[protocolDnssd:Browse] dropping entry")
				return
			}
			select {
			case <-ctx.Done():
			case out <- entry:
			}
		}
		rmv := func(dnssd.BrowseEntry) {} // consumers track liveness themselves
		for ctx.Err() == nil {
			// brutella/dnssd stops its socket readers only on context.Canceled:
			// a round ending on a deadline (its own, or the caller's) leaves
			// them spinning forever on a closed socket. So each round runs on
			// a context detached from the caller's error, always ended by cancel.
			round, cancel := context.WithCancel(context.Background())
			stopParent := context.AfterFunc(ctx, cancel)
			timer := time.AfterFunc(browseRound, cancel)
			err := dnssd.LookupType(round, service, add, rmv)
			early := round.Err() == nil
			timer.Stop()
			stopParent()
			cancel()
			if err != nil && early {
				// Failed before its time (no network...): do not spin
				zlog.Warn().Err(err).Msg("[protocolDnssd:Browse] lookup failed, retrying")
				select {
				case <-ctx.Done():
				case <-time.After(5 * time.Second):
				}
			}
		}
	}()
	return out, nil
}

func convert(raw dnssd.BrowseEntry, opts BrowseOptions) (Entry, error) {
	records := make([]string, 0, len(raw.Text))
	for k, v := range raw.Text {
		records = append(records, k+"="+v)
	}
	txt, err := ParseTXT(records)
	if err != nil {
		return Entry{}, err
	}
	addrs := append([]net.IP{}, raw.IPs...)
	if len(addrs) == 0 {
		return Entry{}, fmt.Errorf("no addresses")
	}
	pick := pickAddr(addrs, opts.PreferIPv4)
	if !opts.AllowPrivate && isPrivate(pick) {
		return Entry{}, fmt.Errorf("private/loopback address %s rejected", pick)
	}
	host := pick.String()
	if pick.To4() == nil {
		host = "[" + host + "]"
	}
	u := &url.URL{
		Scheme: txt.Scheme,
		Host:   fmt.Sprintf("%s:%d", host, raw.Port),
		Path:   txt.TD,
	}
	return Entry{
		Instance: raw.Name,
		Host:     raw.Host,
		Port:     raw.Port,
		Addrs:    addrs,
		TXT:      txt,
		TDURL:    u,
	}, nil
}

func pickAddr(addrs []net.IP, preferIPv4 bool) net.IP {
	if preferIPv4 {
		for _, a := range addrs {
			if a.To4() != nil {
				return a
			}
		}
	} else {
		for _, a := range addrs {
			if a.To4() == nil {
				return a
			}
		}
	}
	return addrs[0]
}

// isPrivate reports whether the address is loopback, link-local, private, or
// otherwise unsuitable for cross-host TD fetching on an untrusted network.
func isPrivate(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() ||
		ip.IsInterfaceLocalMulticast() || ip.IsUnspecified() {
		return true
	}
	if ip.IsPrivate() {
		return true
	}
	return false
}
