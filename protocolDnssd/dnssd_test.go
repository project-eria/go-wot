package protocolDnssd

import (
	"strings"
	"testing"
)

func TestBuildTXT_Defaults(t *testing.T) {
	got, err := BuildTXT("/thingA", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"td=/thingA", "type=Thing", "scheme=http"}
	if len(got) != len(want) {
		t.Fatalf("len mismatch: got %v want %v", got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("at %d: got %q want %q", i, got[i], want[i])
		}
	}
}

func TestBuildTXT_Errors(t *testing.T) {
	tests := []struct {
		name           string
		td, etype, sch string
	}{
		{"empty td", "", "", ""},
		{"relative td", "thingA", "", ""},
		{"bad type", "/x", "Robot", ""},
		{"bad scheme", "/x", "", "ftp"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := BuildTXT(tc.td, tc.etype, tc.sch); err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

func TestBuildTXT_OversizeRejected(t *testing.T) {
	long := "/" + strings.Repeat("a", 300)
	if _, err := BuildTXT(long, "", ""); err == nil {
		t.Fatal("expected oversize error")
	}
}

func TestParseTXT_Defaults(t *testing.T) {
	got, err := ParseTXT([]string{"td=/x"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.TD != "/x" || got.Type != "Thing" || got.Scheme != "http" {
		t.Errorf("defaults not applied: %+v", got)
	}
}

func TestParseTXT_FullRecord(t *testing.T) {
	got, err := ParseTXT([]string{"td=/dir", "type=Directory", "scheme=https", "extra=ignored"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.TD != "/dir" || got.Type != "Directory" || got.Scheme != "https" {
		t.Errorf("got %+v", got)
	}
}

func TestParseTXT_MissingTD(t *testing.T) {
	if _, err := ParseTXT([]string{"type=Thing"}); err == nil {
		t.Fatal("expected error for missing td")
	}
}

func TestParseTXT_RelativeTD(t *testing.T) {
	if _, err := ParseTXT([]string{"td=relative"}); err == nil {
		t.Fatal("expected error for relative td")
	}
}

func TestParseTXT_InvalidScheme(t *testing.T) {
	if _, err := ParseTXT([]string{"td=/x", "scheme=coap"}); err == nil {
		t.Fatal("expected error for non-TCP scheme")
	}
}

func TestRoundTrip(t *testing.T) {
	rec, err := BuildTXT("/things/sensor", EntityThing, SchemeHTTPS)
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := ParseTXT(rec)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.TD != "/things/sensor" || parsed.Type != "Thing" || parsed.Scheme != "https" {
		t.Errorf("round-trip mismatch: %+v", parsed)
	}
}
