//go:build mobile_skel

package mobile

import "testing"

func TestParseHY2Config_Defaults(t *testing.T) {
	resetState()
	// минимальный валидный конфиг с корневыми полями HY2
	jsonStr := `{"server":"example.com:443","password":"secret"}`
	if err := cfgSet(jsonStr); err != nil {
		t.Fatalf("cfgSet(valid) unexpected error: %v", err)
	}
	cfg, err := parseHY2Config()
	if err != nil {
		t.Fatalf("parseHY2Config: %v", err)
	}
	if cfg.Server != "example.com:443" || cfg.Password != "secret" {
		t.Fatalf("unexpected server/password: %#v", cfg)
	}
	if len(cfg.ALPN) == 0 || cfg.ALPN[0] != "h3" {
		t.Fatalf("default ALPN not applied: %#v", cfg.ALPN)
	}
	if cfg.Mode != "tun2socks" {
		t.Fatalf("default Mode must be tun2socks, got %q", cfg.Mode)
	}
}

func TestParseHY2Config_Invalid(t *testing.T) {
	resetState()
	// cfgSet проверяет синтаксис JSON, поэтому он может не упасть.
	_ = cfgSet(`{"server":"bad","password":""}`)
	_, err := parseHY2Config()
	if err == nil {
		t.Fatal("expected error for invalid host:port and empty password")
	}
}
