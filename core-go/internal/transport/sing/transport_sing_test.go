//go:build mobile_skel

package sing

import "testing"

func TestTransportSing_StatusDefaults(t *testing.T) {
	cfg := HY2Config{SNI: "sni.test", ALPN: []string{"h3"}}
	tr := newTransportSingHY2(cfg)
	if tr == nil {
		t.Fatal("newTransportSingHY2 returned nil")
	}
	// Start/Stop — заглушки и должны возвращать nil
	if err := tr.Start(nil); err != nil {
		t.Fatalf("Start() unexpected error: %v", err)
	}
	st := tr.Status()
	if st.SNI != "sni.test" || st.ALPN != "h3" {
		t.Fatalf("unexpected status SNI/ALPN: %#v", st)
	}
	if err := tr.Stop(nil); err != nil {
		t.Fatalf("Stop() unexpected error: %v", err)
	}
}
