package main

import (
	"fmt"
)

// Demo is a scripted, offline simulation of a rapid takeover sweep.
//
// It exists so the operator console can be demonstrated (and screenshotted)
// without touching anyone's infrastructure: every host is a reserved
// example.com name and every verdict is a literal, not a measurement.
//
// Nothing here resolves DNS or opens a socket. That distinction matters —
// a demo that quietly performed real requests against real hostnames would
// be exactly the kind of surprise this project exists to avoid.
func Demo(color bool) {
	rows := []struct {
		host    string
		verdict Verdict
		detail  string
		note    string
	}{
		{"api.example.com", Alive, "nginx", ""},
		{"admin.example.com", Takeover, "Heroku", "live"},
		{"legacy.example.com", Takeover, "CloudFront", ""},
		{"staging.example.com", Likely, "generic", ""},
		{"cdn.example.com", NoDNS, "", ""},
		{"shop.example.com", Alive, "Cloudflare", ""},
	}

	fmt.Println()
	fmt.Println(Paint("  __  ______  ______________   ____________", Bold, color))
	fmt.Println(Paint("   / / / / __ \\/ ___/_  __/   | / ____/ ____/", Bold, color))
	fmt.Println(Paint("  / /_/ / / / /\\__ \\ / / / /| |/ / __/ __/", Bold, color))
	fmt.Println(Paint(" / __  / /_/ /___/ // / / ___ / /_/ / /___", Bold, color))
	fmt.Println(Paint("/_/ /_/\\____//____//_/ /_/  |_|\\____/_____/", Bold, color))
	fmt.Println()

	counts := map[Verdict]int{}
	for _, r := range rows {
		counts[r.verdict]++
	}
	alive := counts[Alive] + counts[Takeover] + counts[Likely]
	fmt.Printf("  %s\n", Paint("LEVIATHAN.AC // HOSTAGE LVX", Cyan, color))
	fmt.Printf("  %s\n\n", Paint("RAPID TAKEOVER CHECK (SIMULATED — no requests sent)", Dim, color))
	fmt.Printf("  scanned  %-3d  alive  %-3d  takeovers  %-3d  likely  %-3d\n",
		len(rows), alive, counts[Takeover], counts[Likely])
	fmt.Println()

	for _, r := range rows {
		m := marks[r.verdict]
		label := Pad(string(r.verdict), 9)
		detail := r.detail
		if r.note != "" {
			detail = r.detail + "  " + r.note
		}
		detail = Pad(detail, 22)
		fmt.Printf("  %s  %s  %s\n",
			Paint(Pad(r.host, 24), m.color, color),
			Paint(label, m.color, color),
			Paint(detail, Dim, color))
	}

	fmt.Println()
	last := rows[1]
	fmt.Printf("  %s\n", Paint(
		fmt.Sprintf("last event · %s · %s · %s", last.host, last.verdict, last.detail),
		Dim, color))
	fmt.Printf("  %s\n", Paint("offline simulation — this is not a scan result", Dim, color))
	fmt.Println()
}

// Pad left-aligns s in a fixed-width field, truncating if it is longer.
func Pad(s string, n int) string {
	r := []rune(s)
	if len(r) > n {
		r = r[:n]
	}
	return string(r) + Spaces(n-len(r))
}

func Spaces(n int) string {
	if n <= 0 {
		return ""
	}
	out := make([]byte, n)
	for i := range out {
		out[i] = ' '
	}
	return string(out)
}
