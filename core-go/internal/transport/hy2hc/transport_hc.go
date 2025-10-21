//go:build (android || ios) && hc

package hy2hc

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ChimeraFlow/Bereznev-HY2-Core/core-go/internal/netstack/protect"
	"github.com/ChimeraFlow/Bereznev-HY2-Core/core-go/internal/runtime"
	"github.com/ChimeraFlow/Bereznev-HY2-Core/core-go/internal/telemetry"
	"github.com/ChimeraFlow/Bereznev-HY2-Core/core-go/internal/transport"
	"github.com/ChimeraFlow/Bereznev-HY2-Core/core-go/internal/transport/sing"
	"github.com/ChimeraFlow/Bereznev-HY2-Core/core-go/pkg/config"
	hcclient "github.com/apernet/hysteria/core/client"
)

type transportHC struct {
	mu     sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc

	rtt   atomic.Int64
	sni   string
	alpn  string
	rem   string
	lastE atomic.Value // string

	superWg sync.WaitGroup
	closed  atomic.Bool

	cli   *hcclient.Client
	pconn net.PacketConn
	cfg   config.HY2Config
}

func NewTransportHC(cfg config.HY2Config) transport.Transport {
	t := &transportHC{cfg: cfg, sni: cfg.SNI}
	if len(cfg.ALPN) > 0 {
		t.alpn = cfg.ALPN[0]
	} else {
		t.alpn = "h3"
	}
	return t
}

func (t *transportHC) Start(parent context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.ctx != nil {
		return nil
	}
	if parent == nil {
		parent = context.Background()
	}

	ctx, cancel := context.WithCancel(parent)
	t.ctx, t.cancel = ctx, cancel
	t.closed.Store(false)

	_ = t.StartOnce(ctx)

	t.superWg.Add(1)
	go t.Supervisor()
	return nil
}

func (t *transportHC) Stop(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	t.mu.Lock()
	if t.ctx == nil {
		t.mu.Unlock()
		return nil
	}
	cancel := t.cancel
	t.closed.Store(true)
	t.mu.Unlock()

	if cancel != nil {
		cancel()
	}
	if t.cli != nil {
		_ = t.cli.Close()
	}
	if t.pconn != nil {
		_ = t.pconn.Close()
	}

	done := make(chan struct{})
	go func() { t.superWg.Wait(); close(done) }()
	select {
	case <-ctx.Done():
	case <-done:
	}

	t.mu.Lock()
	t.ctx = nil
	t.cancel = nil
	t.mu.Unlock()
	return nil
}

func (t *transportHC) Status() transport.TransportStatus {
	st := transport.TransportStatus{
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

func (t *transportHC) Supervisor() {
	defer t.superWg.Done()
	bo := runtime.NewBackoffState()

	for {
		if t.closed.Load() {
			return
		}
		if sing.IsAliveSing() {
			time.Sleep(200 * time.Millisecond)
			continue
		}

		next := bo.Next()
		telemetry.Emit(
			telemetry.EvtReconnecting,
			sing.ToJSON(sing.ReconnectingPayload{
				Reason: "lost", Attempt: bo.Attempt, NextMs: int(next.Milliseconds()),
			}))
		telemetry.SetLastBackoffMs(next.Milliseconds())

		timer := time.NewTimer(next)
		select {
		case <-t.ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
		}

		if t.closed.Load() {
			return
		}
		if err := t.StartOnce(t.ctx); err != nil {
			t.lastE.Store("reconnect: " + err.Error())
			telemetry.SetLastErrTs(time.Now().Unix())
			continue
		}
		bo.Reset()
		telemetry.Reconnects.Add(1)
		telemetry.Emit(telemetry.EvtReconnected, sing.ToJSON(sing.ReconnectedPayload{RttMs: t.rtt.Load()}))
	}
}

func (t *transportHC) IsAlive() bool {
	// достаточно простого критерия; позже можно добавить lastRTTAt с таймаутом
	return t.rtt.Load() > 0
}

// --- ключевая точка: запуск Hysteria2 Core + Protect(fd) ---

func (t *transportHC) StartOnce(ctx context.Context) error {
	// 1) UDP PacketConn с Protect(fd) — это твой protect_dial.go
	if t.pconn != nil {
		_ = t.pconn.Close()
		t.pconn = nil
	}
	pc, err := protect.ProtectedPacketConn(ctx)
	if err != nil {
		return fmt.Errorf("udp listen: %w", err)
	}

	// 2) TLS (SNI/ALPN)
	tlsConf := &tls.Config{
		ServerName: t.cfg.SNI,
		NextProtos: func() []string {
			if len(t.cfg.ALPN) > 0 {
				return t.cfg.ALPN
			}
			return []string{"h3"}
		}(),
		MinVersion: tls.VersionTLS13,
	}

	// 3) Конфиг клиента HC.
	// ВАЖНО: конкретные поля hcclient.Config могут отличаться по версии.
	// Скелет ниже — подгони названия полей под твою версию либы.
	cconf := &hcclient.Config{
		Server:    t.cfg.Server,   // "host:port"
		Auth:      t.cfg.Password, // базовая auth (PSK)
		TLSConfig: tlsConf,
		// Если библиотека позволяет — передай готовый PacketConn внутрь QUIC:
		// QUIC: hcclient.QUICConfig{ PacketConn: pc }, // ← проверь имя поля в твоей версии
	}

	// Альтернатива, если PacketConn не поддерживается: установить Control-hook через Dialer.
	// dialer := protectedTCPDialer()  // твой helper
	// и передать его в поля клиента/QUIC (если API даёт такой хук).

	cli, err := hcclient.New(ctx, cconf)
	if err != nil {
		_ = pc.Close()
		return fmt.Errorf("hc new: %w", err)
	}

	// 4) Быстрая проверка канала и первичный RTT (замени на реальный Ping, если есть)
	t.rtt.Store(25) // TODO: заменить на значение из клиента (Ping/RTT)
	t.rem = t.cfg.Server
	t.cli = cli
	t.pconn = pc
	telemetry.HealthSetIdentity(t.cfg.SNI, t.alpn)

	// (опционально) пульс RTT раз в несколько секунд:
	// go t.rttProbeLoop()

	return nil
}
