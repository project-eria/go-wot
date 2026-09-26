package protocolDnssd

import (
	"context"
	"fmt"
	"net"
	"os"
	"runtime"
	"testing"
	"time"
)

// TestAdvertiseBrowseResolve runs on the real network (multicast needed,
// skipped with -short). Each OS checks what it can:
//   - macOS: the SRV host name must resolve through the SYSTEM resolver
//     (mDNSResponder) — the case grandcat/zeroconf failed, it never answered
//     A queries for the host. In-process Browse is not checked there:
//     mDNSResponder keeps the multicast replies (zeroconf had the same limit).
//   - Linux (production): the advertised Thing must be found by Browse. The
//     Go resolver does no mDNS there, so name resolution is not checked.
func TestAdvertiseBrowseResolve(t *testing.T) {
	if testing.Short() {
		t.Skip("needs multicast on the real network")
	}
	_, ip := primary()
	if ip == nil {
		t.Skip("no primary interface")
	}
	ref := fmt.Sprintf("dnssd-test-%d", os.Getpid())

	a, err := NewAdvertiser(Options{Port: 18080})
	if err != nil {
		t.Fatal(err)
	}
	a.Expose(ref, nil) // nil Thing: the instance name falls back to ref
	a.Start()
	defer a.Stop()

	if runtime.GOOS == "darwin" {
		checkSystemResolution(t, ip)
		return
	}

	// Browse starts while the advertiser is still probing: the Thing must
	// show up at a later round at the latest (see browseRound)
	ctx, cancel := context.WithTimeout(context.Background(), 3*browseRound)
	defer cancel()
	entries, err := Browse(ctx, BrowseOptions{PreferIPv4: true, AllowPrivate: true})
	if err != nil {
		t.Fatal(err)
	}
	var found *Entry
	for e := range entries {
		if e.Instance == ref {
			found = &e
			break
		}
	}
	if found == nil {
		t.Fatal("advertised Thing not found by Browse")
	}
	if found.Port != 18080 || found.TXT.TD != "/"+ref || found.TDURL.Hostname() != ip.String() {
		t.Fatalf("unexpected entry: port %d td %q url %s", found.Port, found.TXT.TD, found.TDURL)
	}
}

// checkSystemResolution: the SRV target resolves by name, as a browser or
// reqwest does (waits for the probe + announce first)
func checkSystemResolution(t *testing.T, ip net.IP) {
	t.Helper()
	time.Sleep(3 * time.Second)
	host := hostnameForIP(ip) + ".local"
	rctx, rcancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer rcancel()
	addrs, err := net.DefaultResolver.LookupHost(rctx, host)
	if err != nil {
		t.Fatalf("host %s does not resolve: %v", host, err)
	}
	for _, a := range addrs {
		if a == ip.String() {
			return
		}
	}
	t.Fatalf("host %s resolved to %v, want %s", host, addrs, ip)
}
