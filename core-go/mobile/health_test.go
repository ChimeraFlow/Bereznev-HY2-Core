//go:build mobile_skel

package mobile

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestHealthJSON_Stopped(t *testing.T) {
	resetState() // started=false, version/engine из пакета

	var h Health
	if err := json.Unmarshal([]byte(HealthJSON()), &h); err != nil {
		t.Fatalf("invalid JSON from HealthJSON(): %v", err)
	}
	if h.Running {
		t.Fatalf("expected Running=false, got true")
	}
	if h.Engine != EngineID {
		t.Fatalf("Engine=%q, want %q", h.Engine, EngineID)
	}
	if h.Version != SdkVersion {
		t.Fatalf("Version=%q, want %q", h.Version, SdkVersion)
	}
}

func TestHealthJSON_Running(t *testing.T) {
	resetState()
	if err := Start(`{}`); err != "" {
		t.Fatalf("Start() failed: %s", err)
	}

	var h Health
	if err := json.Unmarshal([]byte(HealthJSON()), &h); err != nil {
		t.Fatalf("invalid JSON from HealthJSON(): %v", err)
	}
	if !h.Running {
		t.Fatalf("expected Running=true, got false")
	}
	if h.Engine != EngineID || !strings.Contains(Version(), h.Version) {
		t.Fatalf("mismatch engine/version: %+v, Version()=%q", h, Version())
	}
}

func TestBytesStats_ResetAndRead(t *testing.T) {
	ResetBytesStats()
	bytesIn.Add(123)
	bytesOut.Add(456)
	in, out := BytesStats()
	if in != 123 || out != 456 {
		t.Fatalf("got in=%d out=%d, want 123/456", in, out)
	}
	ResetBytesStats()
	in, out = BytesStats()
	if in != 0 || out != 0 {
		t.Fatalf("after reset got in=%d out=%d, want 0/0", in, out)
	}
}
