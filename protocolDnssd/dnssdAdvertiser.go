package protocolDnssd

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/grandcat/zeroconf"
	"github.com/project-eria/go-wot/producer"
	zlog "github.com/rs/zerolog/log"
)

// Advertiser is a producer.ProtocolServer that registers each exposed Thing
// as a `_wot._tcp` mDNS service. It does NOT serve the Thing Description —
// it only advertises the location of an HTTP server that does. The user is
// expected to add an HttpServer alongside this Advertiser on the same Producer.
type Advertiser struct {
	port      int
	scheme    string // "http" or "https"
	domain    string // typically "local."
	tdPathFor func(ref string) string
	skip      func(ref string) bool
	entries   []advertEntry
	servers   []*zeroconf.Server
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

// pickInterfaces returns the single interface used to reach the public
// internet (i.e. the default-route interface). Limiting to one interface
// avoids duplicate mDNS announcements on hosts connected to several LANs
// (Wi-Fi + Ethernet, bridge, etc.). Falls back to all multicast-capable
// non-loopback interfaces if the primary one cannot be detected.
func pickInterfaces() []net.Interface {
	if iface, _ := primary(); iface != nil {
		return []net.Interface{*iface}
	}
	all, err := net.Interfaces()
	if err != nil {
		return nil
	}
	var out []net.Interface
	for _, i := range all {
		if i.Flags&net.FlagUp == 0 ||
			i.Flags&net.FlagMulticast == 0 ||
			i.Flags&net.FlagLoopback != 0 ||
			i.Flags&net.FlagPointToPoint != 0 {
			continue
		}
		out = append(out, i)
	}
	return out
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
// label, e.g. 192.168.1.42 -> "192-168-1-42". The mDNS server then
// publishes an A record mapping <label>.local. -> <ip>, so the SRV
// Target shown by `dns-sd -L` and resolved by consumers maps directly
// to the IP without going through the host's .local name.
func hostnameForIP(ip net.IP) string {
	if v4 := ip.To4(); v4 != nil {
		return strings.ReplaceAll(v4.String(), ".", "-")
	}
	return strings.ReplaceAll(ip.String(), ":", "-")
}

// Start opens one mDNS registration per exposed Thing.
func (a *Advertiser) Start() {
	if len(a.entries) == 0 {
		zlog.Warn().Msg("[protocolDnssd:Start] no Things to advertise")
		return
	}
	iface, ip := primary()
	if ip == nil {
		zlog.Error().Msg("[protocolDnssd:Start] cannot determine primary IP")
		return
	}
	var ifaces []net.Interface
	if iface != nil {
		ifaces = []net.Interface{*iface}
	}
	host := hostnameForIP(ip)
	for _, e := range a.entries {
		txt, err := BuildTXT(e.tdPath, EntityThing, a.scheme)
		if err != nil {
			zlog.Error().Err(err).Str("ref", e.instance).Msg("[protocolDnssd:Start] BuildTXT")
			continue
		}
		srv, err := zeroconf.RegisterProxy(
			e.instance,
			ServiceWoTTCP,
			a.domain,
			a.port,
			host,
			[]string{ip.String()},
			txt,
			ifaces,
		)
		if err != nil {
			zlog.Error().Err(err).Str("ref", e.instance).Msg("[protocolDnssd:Start] register failed")
			continue
		}
		a.servers = append(a.servers, srv)
		zlog.Info().Str("instance", e.instance).Str("td", e.tdPath).Str("host", host).
			Str("ip", ip.String()).Int("port", a.port).Msg("[protocolDnssd:Start] advertising")
	}
}

// Stop tears down all active mDNS registrations. zeroconf v1.0.0's Shutdown
// can hang indefinitely on macOS because its recv goroutines do not always
// unblock when the UDP socket is closed (upstream issue). The Goodbye packet
// is multicasted at the start of Shutdown — well before the blocking Wait —
// so we run Shutdown in a goroutine and cap our wait at 1s: enough for the
// Goodbye to leave the host and for the kernel to flush, short enough to
// never freeze the process.
func (a *Advertiser) Stop() {
	servers := a.servers
	a.servers = nil
	zlog.Info().Int("count", len(servers)).Msg("[protocolDnssd:Stop] sending Goodbye")
	done := make(chan struct{})
	go func() {
		for _, s := range servers {
			s.Shutdown()
		}
		close(done)
	}()
	select {
	case <-done:
		zlog.Debug().Msg("[protocolDnssd:Stop] zeroconf Shutdown completed cleanly")
	case <-time.After(1 * time.Second):
		zlog.Debug().Msg("[protocolDnssd:Stop] zeroconf Shutdown still running after 1s, abandoning wait")
	}
}
