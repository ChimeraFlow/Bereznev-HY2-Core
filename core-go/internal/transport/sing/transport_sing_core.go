// go:build android || ios || mobile_skel

package sing

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"sync/atomic"
	"time"
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

func newTransportSingHY2(cfg HY2Config) *transportSingHY2 {
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

	_ = startOnceSing(t, ctx)
	t.superWg.Add(1)
	go t.supervisor()

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

func toJSON(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}

// --- внутреннее ---

func (t *transportSingHY2) supervisor() {
	defer t.superWg.Done()
	bo := newBackoffState()

	for {
		if t.closed.Load() {
			return
		}
		if isAliveSing(t) { // в tests подменён, в prod — реальный критерий
			time.Sleep(50 * time.Millisecond)
			continue
		}

		next := bo.next()
		emit(EvtReconnecting, toJSON(evtReconnecting{
			Reason:  "lost",
			Attempt: bo.attempt,
			NextMs:  int(next.Milliseconds()),
		}))
		setLastBackoffMs(next.Milliseconds())

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
		if err := startOnceSing(t, t.ctx); err != nil {
			// фиксируем ошибку и ждём следующий backoff
			t.lastE.Store("reconnect: " + err.Error())
			setLastErrTs(time.Now().Unix())
			continue
		}
		bo.reset()
		reconnects.Add(1)
		emit(EvtReconnected, toJSON(evtReconnected{RttMs: t.rtt.Load()}))
	}
}

// startOnce пытается (пере)поднять клиент/сессию
func (t *transportSingHY2) startOnce(ctx context.Context) error {
	// TODO:
	// 1) собрать dialer с Protect(fd) через net.ListenConfig.Control / net.Dialer.Control
	// 2) инициализировать sing/hysteria2 outbound
	// 3) выполнить рукопожатие/проверку доступности (короткая операция с таймаутом)
	// 4) если ок — обновить t.rem, t.rtt.Store(…), запустить фоновый пульс RTT (goroutine)
	// 5) если не ок — вернуть err
	return errors.New("sing HY2 not wired yet")
}

func (t *transportSingHY2) isAlive() bool {
	// TODO: проверь состояние клиента/сокета/последний успешный пульс
	// на первом шаге: считаем живым, если rtt обновлялся < N секунд назад
	// (можно хранить atomic lastRTTts)
	return false
}

func (t *transportSingHY2) recordErr(stage string, err error) {
	if err == nil {
		return
	}
	t.lastE.Store(stage + ": " + err.Error())
	emitError(int(ErrEngineInitFailed), stage+": "+err.Error()) // или отдельный stage-код
}
