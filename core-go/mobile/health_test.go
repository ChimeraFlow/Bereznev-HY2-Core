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
	if h.Engine != engineID {
		t.Fatalf("Engine=%q, want %q", h.Engine, engineID)
	}
	if h.Version != sdkVersion {
		t.Fatalf("Version=%q, want %q", h.Version, sdkVersion)
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
	if h.Engine != engineID || !strings.Contains(Version(), h.Version) {
		t.Fatalf("mismatch engine/version: %+v, Version()=%q", h, Version())
	}
}
