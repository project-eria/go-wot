package protocolDnssd

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/brutella/dnssd"
	"github.com/project-eria/go-wot/producer"
	zlog "github.com/rs/zerolog/log"
)

// Advertiser is a producer.ProtocolServer that registers each exposed Thing
// as a `_wot._tcp` mDNS service. It does NOT serve the Thing Description —
// it only advertises the location of an HTTP server that does. The user is
// expected to add an HttpServer alongside this Advertiser on the same Producer.
//
// The responder (brutella/dnssd) answers queries for the service AND for
// the host name, and probes before announcing (RFC 6762 §8-9): an instance
// name already taken on the network is renamed ("name (2)"), so two
// producers advertising the same TD id no longer clash silently.
type Advertiser struct {
	port      int
	scheme    string // "http" or "https"
	domain    string // typically "local."
	tdPathFor func(ref string) string
	skip      func(ref string) bool
	mu        sync.Mutex // guards entries, cancel and done
	entries   []advertEntry
	cancel    context.CancelFunc // stops the responder (sends Goodbye)
	done      chan struct{}      // closed when the responder has returned
}

type advertEntry struct {
	instance string
	tdPath   string
}

// Options configures the Advertiser.
type Options struct {
	// Port the HTTP server listens on. Required.
	Port int
	// Scheme advertised in the TXT record. Defaults to "http".
	Scheme string
	// Domain for mDNS, defaults to "local."
	Domain string
	// TDPathFor maps a Thing ref to an absolute TD path. Defaults to "/<ref>",
	// matching the route registered by HttpServer.Expose.
	TDPathFor func(ref string) string
	// Skip, when set, returns true for Things that must not be advertised.
	Skip func(ref string) bool
}

// NewAdvertiser builds an Advertiser. The Port is required.
func NewAdvertiser(opts Options) (*Advertiser, error) {
	if opts.Port <= 0 {
		return nil, fmt.Errorf("Port is required")
	}
	if opts.Scheme == "" {
		opts.Scheme = SchemeHTTP
	}
	if !validScheme(opts.Scheme) {
		return nil, fmt.Errorf("invalid scheme %q", opts.Scheme)
	}
	if opts.Domain == "" {
		opts.Domain = "local."
	}
	if opts.TDPathFor == nil {
		opts.TDPathFor = func(ref string) string { return "/" + ref }
	}
	return &Advertiser{
		port:      opts.Port,
		scheme:    opts.Scheme,
		domain:    opts.Domain,
		tdPathFor: opts.TDPathFor,
		skip:      opts.Skip,
	}, nil
}

// Expose records the Thing for later advertisement. The actual mDNS
// registration happens in Start(). The mDNS instance name is the Thing's TD
// id (URN) when available, falling back to the ref, then to "wot".
//
// Expose is idempotent per (instance, tdPath): repeated calls for the same
// Thing produce a single mDNS registration.
func (a *Advertiser) Expose(ref string, t producer.ExposedThing) {
	if a.skip != nil && a.skip(ref) {
		return
	}
	instance := ""
	if t != nil && t.TD() != nil {
		instance = strings.TrimPrefix(t.TD().ID, "urn:")
	}
	if instance == "" {
		instance = ref
	}
	if instance == "" {
		instance = "wot"
	}
	tdPath := a.tdPathFor(ref)
	a.mu.Lock()
	defer a.mu.Unlock()
	for _, e := range a.entries {
		if e.instance == instance && e.tdPath == tdPath {
			return
		}
	}
	a.entries = append(a.entries, advertEntry{
		instance: instance,
		tdPath:   tdPath,
	})
}

// primary returns the interface that owns the source IP the kernel would
// pick to reach the public internet, plus that IP. No packet is actually
// sent — UDP Dial only resolves the route.
func primary() (*net.Interface, net.IP) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, nil
	}
	defer conn.Close()
	localIP := conn.LocalAddr().(*net.UDPAddr).IP
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, localIP
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}
		for _, a := range addrs {
			ipnet, ok := a.(*net.IPNet)
			if ok && ipnet.IP.Equal(localIP) {
				return &i, localIP
			}
		}
	}
	return nil, localIP
}

// hostnameForIP turns a primary IPv4 address into a DNS-safe hostname
// label, e.g. 192.168.1.42 -> "192-168-1-42". The responder answers A
// queries for <label>.local. with the IP, so the SRV Target resolves to the
// primary IP only (never a Docker bridge or a second NIC). Every process on
// a host publishes the same label with the same address: identical records,
// hence no probing conflict between them.
func hostnameForIP(ip net.IP) string {
	if v4 := ip.To4(); v4 != nil {
		return strings.ReplaceAll(v4.String(), ".", "-")
	}
	return strings.ReplaceAll(ip.String(), ":", "-")
}

// Start registers every exposed Thing on one mDNS responder and runs it in
// the background until Stop.
func (a *Advertiser) Start() {
	a.mu.Lock()
	defer a.mu.Unlock()
	if len(a.entries) == 0 {
		zlog.Warn().Msg("[protocolDnssd:Start] no Things to advertise")
		return
	}
	iface, ip := primary()
	if ip == nil {
		zlog.Error().Msg("[protocolDnssd:Start] cannot determine primary IP")
		return
	}
	var ifaces []string
	if iface != nil {
		ifaces = []string{iface.Name}
	}
	rp, err := dnssd.NewResponder()
	if err != nil {
		zlog.Error().Err(err).Msg("[protocolDnssd:Start] cannot create the mDNS responder")
		return
	}
	host := hostnameForIP(ip)
	for _, e := range a.entries {
		records, err := BuildTXT(e.tdPath, EntityThing, a.scheme)
		if err != nil {
			zlog.Error().Err(err).Str("ref", e.instance).Msg("[protocolDnssd:Start] BuildTXT")
			continue
		}
		txt := map[string]string{}
		for _, kv := range records {
			k, v, _ := strings.Cut(kv, "=")
			txt[k] = v
		}
		srv, err := dnssd.NewService(dnssd.Config{
			Name:   e.instance,
			Type:   ServiceWoTTCP,
			Domain: strings.TrimSuffix(a.domain, "."),
			Host:   host,
			Text:   txt,
			IPs:    []net.IP{ip}, // only the primary IP in the A record
			Port:   a.port,
			Ifaces: ifaces,
		})
		if err != nil {
			zlog.Error().Err(err).Str("ref", e.instance).Msg("[protocolDnssd:Start] service")
			continue
		}
		if _, err := rp.Add(srv); err != nil {
			zlog.Error().Err(err).Str("ref", e.instance).Msg("[protocolDnssd:Start] register failed")
			continue
		}
		zlog.Info().Str("instance", e.instance).Str("td", e.tdPath).Str("host", host).
			Str("ip", ip.String()).Int("port", a.port).Msg("[protocolDnssd:Start] advertising")
	}

	ctx, cancel := context.WithCancel(context.Background())
	a.cancel = cancel
	a.done = make(chan struct{})
	go func(done chan struct{}) {
		defer close(done)
		// Probes, announces, answers queries; on cancel sends Goodbye
		if err := rp.Respond(ctx); err != nil && ctx.Err() == nil {
			zlog.Error().Err(err).Msg("[protocolDnssd] responder stopped")
		}
	}(a.done)
}

// Stop cancels the responder, which multicasts Goodbye packets, and waits
// for it (bounded: a process exit must never hang on mDNS).
func (a *Advertiser) Stop() {
	a.mu.Lock()
	cancel, done := a.cancel, a.done
	a.cancel, a.done = nil, nil
	a.mu.Unlock()
	if cancel == nil {
		return
	}
	zlog.Info().Msg("[protocolDnssd:Stop] sending Goodbye")
	cancel()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		zlog.Debug().Msg("[protocolDnssd:Stop] responder still running after 2s, abandoning wait")
	}
}
