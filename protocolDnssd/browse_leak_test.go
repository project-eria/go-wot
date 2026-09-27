package protocolDnssd

import (
	"context"
	"runtime"
	"testing"
	"time"
)

// TestBrowseNoLeak: brutella/dnssd stops its socket readers only on
// context.Canceled; a round ending on a deadline left two goroutines
// spinning forever on a closed socket (100% CPU per few hours of browsing).
// The parent here expires by deadline too, the worst case.
func TestBrowseNoLeak(t *testing.T) {
	defer func(d time.Duration) { browseRound = d }(browseRound)
	browseRound = 200 * time.Millisecond

	before := runtime.NumGoroutine()
	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()
	entries, err := Browse(ctx, BrowseOptions{AllowPrivate: true})
	if err != nil {
		t.Fatal(err)
	}
	for range entries {
	}
	deadline := time.Now().Add(3 * time.Second)
	for runtime.NumGoroutine() > before+1 && time.Now().Before(deadline) {
		time.Sleep(50 * time.Millisecond)
	}
	if n := runtime.NumGoroutine(); n > before+1 {
		buf := make([]byte, 1<<16)
		t.Fatalf("%d goroutines left after Browse (before: %d)\n%s", n, before, buf[:runtime.Stack(buf, true)])
	}
}
