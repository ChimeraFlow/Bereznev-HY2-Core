//go:build mobile_skel

package mobile

import (
	"strings"
	"testing"
)

func TestCfgSetAndGet_ValidConfig(t *testing.T) {
	resetState()

	jsonStr := `{
  // minimal HY2 config (sing-box style allowed: comments ok)
  "server":   "example.com:443",
  "password": "testpass",
  "alpn":     ["h3"],
  "mode":     "tun2socks"
}`

	if err := cfgSet(jsonStr); err != nil {
		t.Fatalf("cfgSet(valid) unexpected error: %v", err)
	}

	got := string(cfgGet())
	if !strings.Contains(got, `"server"`) {
		t.Fatalf("cfgGet() = %q, want it to contain 'server'", got)
	}

	// Проверим, что возвращается копия, а не ссылка
	buf := cfgGet()
	buf[0] = 'X'
	if string(cfgGet())[0] == 'X' {
		t.Fatalf("cfgGet() returned shared slice, should be copy")
	}
}

func TestCfgSet_InvalidConfig(t *testing.T) {
	resetState()

	badJSON := `{ invalid }`
	err := cfgSet(badJSON)
	if err == nil {
		t.Fatal("cfgSet(invalid) expected error, got nil")
	}

	if cfgRaw != nil && len(cfgRaw) > 0 {
		t.Fatalf("cfgRaw should not be modified on invalid config, got %s", string(cfgRaw))
	}
}

func TestCfgGet_Empty(t *testing.T) {
	resetState()

	if got := cfgGet(); len(got) != 0 {
		t.Fatalf("expected empty cfgGet() initially, got %q", string(got))
	}
}
