//go:build mobile_skel && hc

package hy2hc

import (
	"context"
	"testing"
)

func TestSelectTransport_HCChosen(t *testing.T) {
	cfg := HY2Config{
		Engine:   "hc",
		Server:   "example.com:443",
		Password: "secret",
		SNI:      "example.com",
		ALPN:     []string{"h3"},
	}
	tr := selectTransport(cfg)
	if tr == nil {
		t.Fatalf("selectTransport returned nil for Engine=hc")
	}
	// транспорт стартует и корректно останавливается
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := tr.Start(ctx); err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	if err := tr.Stop(ctx); err != nil {
		t.Fatalf("Stop failed: %v", err)
	}
}
