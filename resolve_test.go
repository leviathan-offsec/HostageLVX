package main

import (
	"net"
	"testing"
	"time"
)

// A custom resolver is the one input a user is most likely to paste in the form
// they see everywhere else, "1.1.1.1:53". Before this was handled, that produced
// net.JoinHostPort("1.1.1.1:53", "53") == "[1.1.1.1:53]:53", which cannot be
// dialed, so every target reported "probe failed" and the tool looked like it
// was broken rather than like it had been handed bad input.
func TestNormalizeResolverStripsPort(t *testing.T) {
	cases := []struct{ in, want string }{
		{"1.1.1.1:53", "1.1.1.1"},
		{"8.8.8.8:53", "8.8.8.8"},
		{"9.9.9.9:5353", "9.9.9.9"},
		{"1.1.1.1", "1.1.1.1"},
		{"", ""},
	}
	for _, c := range cases {
		if got := normalizeResolver(c.in); got != c.want {
			t.Errorf("normalizeResolver(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestNormalizeResolverKeepsBareIPv6(t *testing.T) {
	// A bare IPv6 literal has colons but no port. SplitHostPort reports "too
	// many colons" for it, which must not be mistaken for a port to strip.
	const bare = "2606:4700:4700::1111"
	if got := normalizeResolver(bare); got != bare {
		t.Errorf("normalizeResolver(%q) = %q, want it unchanged", bare, got)
	}
	if got := normalizeResolver("[" + bare + "]:53"); got != bare {
		t.Errorf("normalizeResolver(%q) = %q, want %q", "["+bare+"]:53", got, bare)
	}
}

// Whatever normalizeResolver returns must still be usable as the host half of
// JoinHostPort, which is the only thing it feeds. This is the property that
// actually broke: the old code produced an undialable address here.
func TestNormalizeResolverProducesDialableHost(t *testing.T) {
	for _, in := range []string{"1.1.1.1:53", "1.1.1.1", "8.8.8.8:53", "[2606:4700:4700::1111]:53"} {
		got := net.JoinHostPort(normalizeResolver(in), "53")
		host, port, err := net.SplitHostPort(got)
		if err != nil {
			t.Errorf("normalizeResolver(%q): JoinHostPort gave %q which does not split: %v", in, got, err)
			continue
		}
		if host == "" || port != "53" {
			t.Errorf("normalizeResolver(%q): got host %q port %q, want a non-empty host on 53", in, host, port)
		}
	}
}

// A custom server forces system mode regardless of -resolver, because asking
// for DoH and a specific UDP server at the same time is a contradiction and the
// flag the user typed explicitly should win.
func TestCustomServerForcesSystemMode(t *testing.T) {
	if r := NewResolver("1.1.1.1:53", "doh", defaultTestTimeout); r.method != "system" {
		t.Errorf("method = %q, want system when a custom server is given", r.method)
	}
	if r := NewResolver("", "doh", defaultTestTimeout); r.method != "doh" {
		t.Errorf("method = %q, want doh when no server is given", r.method)
	}
	if r := NewResolver("", "system", defaultTestTimeout); r.method != "system" {
		t.Errorf("method = %q, want system when explicitly asked for", r.method)
	}
}

// The resolver must not silently fall back to the default when a custom server
// is unusable. If the address cannot be dialed the user needs to be told, not
// quietly given a working resolver they did not ask for.
func TestCustomServerIsNotSwallowedOnDialFailure(t *testing.T) {
	// 203.0.113.0/24 is TEST-NET-3, reserved for documentation and not routed.
	r := NewResolver("203.0.113.1", "doh", defaultTestTimeout)
	if r.method != "system" || r.server != "203.0.113.1" {
		t.Fatalf("custom server not retained: method=%q server=%q", r.method, r.server)
	}
	if r.r == nil {
		t.Fatal("expected a resolver configured against the custom server")
	}
}

// defaultTestTimeout keeps the tests independent of whatever the CLI default
// happens to be. A network dial against a documentation-reserved address must
// fail fast rather than hold the suite open.
const defaultTestTimeout = 2 * time.Second
