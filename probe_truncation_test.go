package main

import (
    "context"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
    "time"
)

// A body that cuts off mid-stream must be reported, not silently passed on as
// a complete read.
func TestProbeRecordsTruncatedBody(t *testing.T) {
    srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Length", "4096")
        w.WriteHeader(http.StatusOK)
        if _, err := w.Write([]byte("short")); err != nil {
            return
        }
        // Connection dies before Content-Length bytes are sent.
        if hj, ok := w.(http.Hijacker); ok {
            conn, _, err := hj.Hijack()
            if err != nil {
                return
            }
            conn.Close()
        }
    }))
    defer srv.Close()

    host := strings.TrimPrefix(srv.URL, "http://")
    res := Probe(context.Background(), host, 5*time.Second)
    if !res.OK {
        t.Fatalf("expected OK, got Err=%q", res.Err)
    }
    if res.BodyErr == "" {
        t.Errorf("BodyErr empty: truncated body silently accepted (Body=%q)", res.Body)
    }
}

// The normal path must stay clean: no false truncation warnings.
func TestProbeNoFalseTruncation(t *testing.T) {
    srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if _, err := w.Write([]byte("complete body")); err != nil {
            t.Errorf("write: %v", err)
        }
    }))
    defer srv.Close()

    host := strings.TrimPrefix(srv.URL, "http://")
    res := Probe(context.Background(), host, 5*time.Second)
    if !res.OK {
        t.Fatalf("expected OK, got Err=%q", res.Err)
    }
    if res.BodyErr != "" {
        t.Errorf("BodyErr=%q on a complete read", res.BodyErr)
    }
    if res.Body != "complete body" {
        t.Errorf("Body=%q", res.Body)
    }
}

// An Alive finding must disclose a truncated body so the reader does not read
// it as a clean full-body confirmation.
func TestClassifyDisclosesTruncatedBody(t *testing.T) {
    f := Classify("x.example.com", &DNSResult{IPs: []string{"1.2.3.4"}}, &ProbeResult{
        OK: true, Status: 200, BodyErr: "unexpected EOF",
    })
    if !strings.Contains(f.Evidence, "body truncated") {
        t.Errorf("evidence does not disclose truncation: %q", f.Evidence)
    }
}