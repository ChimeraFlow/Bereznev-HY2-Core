//go:build mobile_skel

package mobile

import (
	"context"
	"sync/atomic"
)

// точная копия «скелета» из transport_sing.go, но под тегом mobile_skel,
// чтобы сборка и тесты проходили в окружении без android/ios.

type transportSingHY2 struct {
	rtt   atomic.Int64
	sni   string
	alpn  string
	rem   string
	lastE atomic.Value // string
}

func newTransportSingHY2(cfg HY2Config) *transportSingHY2 {
	t := &transportSingHY2{sni: cfg.SNI}
	if len(cfg.ALPN) > 0 {
		t.alpn = cfg.ALPN[0]
	}
	return t
}

func (t *transportSingHY2) Start(ctx context.Context) error {
	// заглушка
	return nil
}

func (t *transportSingHY2) Stop(ctx context.Context) error {
	// заглушка
	return nil
}

func (t *transportSingHY2) Status() TransportStatus {
	st := TransportStatus{
		RTTms:  t.rtt.Load(),
		Remote: t.rem,
		ALPN:   t.alpn,
		SNI:    t.sni,
	}
	if v := t.lastE.Load(); v != nil {
		st.LastErr = v.(string)
	}
	return st
}
