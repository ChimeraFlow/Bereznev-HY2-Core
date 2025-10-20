//go:build android || ios

package mobile

import (
	"context"
	"sync/atomic"
)

// transportSingHY2 — заглушка-реализация, чтобы компилилось и стыковалось.
// На втором шаге сюда приедет реальный sing/sing-box outbound:hysteria2.
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
	// TODO: инициализировать sing runtime, поднять HY2 клиент/обработчики,
	//        сохранить remote addr в t.rem, измерять RTT → t.rtt.Store()
	// emit(EvtStarted, `{"path":"transport","engine":"sing-hy2"}`)
	return nil
}

func (t *transportSingHY2) Stop(ctx context.Context) error {
	// TODO: остановить клиенты/коннекты, закрыть ресурсы
	// emit(EvtStopped, `{"path":"transport","engine":"sing-hy2"}`)
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
