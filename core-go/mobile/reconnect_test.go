//go:build mobile_skel

package mobile

import (
	"context"
	"sync"
	"testing"
	"time"
)

type sinkFunc func(name, payload string)

func (f sinkFunc) OnEvent(name, payload string) { f(name, payload) }

func TestReconnect_Backoff_Events(t *testing.T) {
	ResetBytesStats()
	setLastBackoffMs(0)
	setLastErrTs(0)
	reconnects.Store(0)

	var mu sync.Mutex
	evs := make([]string, 0, 8)
	SetEventSink(sinkFunc(func(name, payload string) {
		mu.Lock()
		evs = append(evs, name)
		mu.Unlock()
	}))
	defer SetEventSink(nil)

	// упадём 2 раза, на 3-й — ОК (см. stub)
	testStartOnceFailCount.Store(2)

	tr := newTransportSingHY2(HY2Config{Server: "s", Password: "p"})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := tr.Start(ctx); err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	defer tr.Stop(context.Background())

	dead := time.Now().Add(4 * time.Second)
	var gotR bool
	var rc int
	for time.Now().Before(dead) {
		time.Sleep(25 * time.Millisecond)
		mu.Lock()
		rc = 0
		gotR = false
		for _, e := range evs {
			if e == EvtReconnecting {
				rc++
			}
			if e == EvtReconnected {
				gotR = true
			}
		}
		mu.Unlock()
		if gotR && rc >= 2 {
			break
		}
	}

	if rc < 2 || !gotR {
		t.Fatalf("expected >=2 '%s' and 1 '%s', got events=%v", EvtReconnecting, EvtReconnected, evs)
	}
	if v := reconnects.Load(); v < 1 {
		t.Fatalf("expected reconnects >= 1, got %d", v)
	}
	if ms := lastBackoffMs.Load(); ms <= 0 {
		t.Fatalf("expected lastBackoffMs > 0, got %d", ms)
	}
}
