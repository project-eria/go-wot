package protocolDnssd

import (
	"context"
	"fmt"
	"net"
	"net/url"

	"github.com/grandcat/zeroconf"
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
func Browse(ctx context.Context, opts BrowseOptions) (<-chan Entry, error) {
	if opts.Domain == "" {
		opts.Domain = "local."
	}
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return nil, fmt.Errorf("zeroconf resolver: %w", err)
	}
	rawCh := make(chan *zeroconf.ServiceEntry, 16)
	out := make(chan Entry, 16)

	go func() {
		defer close(out)
		for raw := range rawCh {
			entry, err := convert(raw, opts)
			if err != nil {
				zlog.Debug().Err(err).Str("instance", raw.Instance).
					Msg("[protocolDnssd:Browse] dropping entry")
				continue
			}
			select {
			case <-ctx.Done():
				return
			case out <- entry:
			}
		}
	}()

	if err := resolver.Browse(ctx, ServiceWoTTCP, opts.Domain, rawCh); err != nil {
		return nil, fmt.Errorf("browse: %w", err)
	}
	return out, nil
}

func convert(raw *zeroconf.ServiceEntry, opts BrowseOptions) (Entry, error) {
	txt, err := ParseTXT(raw.Text)
	if err != nil {
		return Entry{}, err
	}
	addrs := append([]net.IP{}, raw.AddrIPv4...)
	addrs = append(addrs, raw.AddrIPv6...)
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
		Instance: raw.Instance,
		Host:     raw.HostName,
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
