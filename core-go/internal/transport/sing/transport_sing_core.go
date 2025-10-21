// go:build android || ios || mobile_skel

package sing

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ChimeraFlow/Bereznev-HY2-Core/core-go/internal/runtime"
	"github.com/ChimeraFlow/Bereznev-HY2-Core/core-go/internal/telemetry"
	"github.com/ChimeraFlow/Bereznev-HY2-Core/core-go/internal/transport"
	"github.com/ChimeraFlow/Bereznev-HY2-Core/core-go/internal/transport/hy2hc"
	ers "github.com/ChimeraFlow/Bereznev-HY2-Core/core-go/pkg/errors"
)

type transportSingHY2 struct {
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
}

func NewTransportSingHY2(cfg hy2hc.HY2Config) *transportSingHY2 {
	t := &transportSingHY2{sni: cfg.SNI}
	if len(cfg.ALPN) > 0 {
		t.alpn = cfg.ALPN[0]
	} else {
		t.alpn = "h3"
	}
	return t
}

func (t *transportSingHY2) Start(parent context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.ctx != nil {
		return nil
	}

	if parent == nil { // ⬅️ добавь это
		parent = context.Background()
	}

	ctx, cancel := context.WithCancel(parent)
	t.ctx, t.cancel = ctx, cancel
	t.closed.Store(false)

	_ = StartOnceSing(t, ctx)
	t.superWg.Add(1)
	go t.Supervisor()

	return nil
}

func (t *transportSingHY2) Stop(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background() // ☑️ защита от nil
	}

	// Снимем снапшоты под локом и НЕ зануляем t.ctx до завершения supervisor.
	t.mu.Lock()
	if t.ctx == nil {
		t.mu.Unlock()
		return nil
	}
	cancel := t.cancel
	t.closed.Store(true)
	t.mu.Unlock()

	if cancel != nil { // ☑️ защита от nil cancel
		cancel()
	}

	// Дадим supervisor корректно завершиться
	done := make(chan struct{})
	go func() { t.superWg.Wait(); close(done) }()

	select {
	case <-ctx.Done():
	case <-done:
	}

	// Теперь можно очистить поля
	t.mu.Lock()
	t.ctx = nil
	t.cancel = nil
	t.mu.Unlock()
	return nil
}

func (t *transportSingHY2) Status() transport.TransportStatus {
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

func ToJSON(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}

type ReconnectingPayload struct {
	Reason  string `json:"reason"`
	Attempt int    `json:"attempt"`
	NextMs  int    `json:"next_ms"`
}

type ReconnectedPayload struct {
	RttMs int64 `json:"rtt_ms"`
}

// --- внутреннее ---

func (t *transportSingHY2) Supervisor() {
	defer t.superWg.Done()
	bo := runtime.NewBackoffState()

	for {
		if t.closed.Load() {
			return
		}
		if IsAliveSing(t) {
			time.Sleep(50 * time.Millisecond)
			continue
		}

		// было: next := bo.next()
		next := bo.Next()

		telemetry.Emit(telemetry.EvtReconnecting, ToJSON(ReconnectingPayload{
			Reason: "lost",
			// было: bo.attempt
			Attempt: bo.Attempt(), // добавь этот геттер в runtime/backoff, если его ещё нет
			NextMs:  int(next.Milliseconds()),
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
		// было: StartOnceSing (функция есть, всё ок)
		if err := StartOnceSing(t, t.ctx); err != nil {
			t.lastE.Store("reconnect: " + err.Error())
			telemetry.SetLastErrTs(time.Now().Unix())
			continue
		}
		// было: bo.reset()
		bo.Reset()

		telemetry.Reconnects.Add(1)
		telemetry.Emit(telemetry.EvtReconnected, ToJSON(ReconnectedPayload{RttMs: t.rtt.Load()}))
	}
}

// startOnce пытается (пере)поднять клиент/сессию
func (t *transportSingHY2) StartOnce(ctx context.Context) error {
	// TODO:
	// 1) собрать dialer с Protect(fd) через net.ListenConfig.Control / net.Dialer.Control
	// 2) инициализировать sing/hysteria2 outbound
	// 3) выполнить рукопожатие/проверку доступности (короткая операция с таймаутом)
	// 4) если ок — обновить t.rem, t.rtt.Store(…), запустить фоновый пульс RTT (goroutine)
	// 5) если не ок — вернуть err
	return errors.New("sing HY2 not wired yet")
}

func (t *transportSingHY2) IsAlive() bool {
	// TODO: проверь состояние клиента/сокета/последний успешный пульс
	// на первом шаге: считаем живым, если rtt обновлялся < N секунд назад
	// (можно хранить atomic lastRTTts)
	return false
}

func (t *transportSingHY2) RecordErr(stage string, err error) {
	if err == nil {
		return
	}
	t.lastE.Store(stage + ": " + err.Error())
	telemetry.EmitError(int(ers.ErrEngineInitFailed), stage+": "+err.Error()) // или отдельный stage-код
}

func StartOnceSing(t *transportSingHY2, ctx context.Context) error {
	// TODO: собрать dialer с Protect(fd), поднять HY2, заполнить t.rem, t.rtt.Store(...)
	return nil
}

func IsAliveSing(t *transportSingHY2) bool {
	return t.rtt.Load() > 0 // временно
}
