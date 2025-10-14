//go:build mobile_skel

package mobile

import (
	"strings"
	"sync"
	"testing"
	"time"
)

type mockSink struct {
	mu     sync.Mutex
	events []string
	logs   []string
}

func (m *mockSink) OnEvent(name, data string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.events = append(m.events, name+"|"+data)
}

func (m *mockSink) Log(level, msg string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.logs = append(m.logs, level+":"+msg)
}

func TestSafeGo_PanicRecovered(t *testing.T) {
	resetState()

	m := &mockSink{}
	SetEventSink(m)
	SetLogger(m)

	ran := make(chan struct{})
	safeGo(func() {
		// этот defer выполнится при панике ещё до recover снаружи
		defer close(ran)
		panic("test crash")
	})

	// ждём подтверждения, что горутина реально запускалась
	select {
	case <-ran:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout: goroutine did not run")
	}

	// даём чуть времени на лог и событие из внешнего defer в safeGo
	time.Sleep(50 * time.Millisecond)

	m.mu.Lock()
	defer m.mu.Unlock()

	// проверяем, что есть error-лог с panic
	foundLog := false
	for _, l := range m.logs {
		if strings.HasPrefix(l, "error:panic:") {
			foundLog = true
			break
		}
	}
	if !foundLog {
		t.Fatalf("expected error log with panic, got %#v", m.logs)
	}

	// и событие "panic"
	foundEvent := false
	for _, e := range m.events {
		if strings.HasPrefix(e, "panic|") {
			foundEvent = true
			break
		}
	}
	if !foundEvent {
		t.Fatalf("expected panic event, got %#v", m.events)
	}
}

func TestSafeGo_NoPanic(t *testing.T) {
	resetState()

	m := &mockSink{}
	SetEventSink(m)
	SetLogger(m)

	done := make(chan struct{})
	safeGo(func() {
		logI("normal work")
		close(done)
	})

	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout: goroutine did not finish")
	}

	time.Sleep(50 * time.Millisecond)

	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.events) != 0 {
		t.Fatalf("expected no events, got %#v", m.events)
	}
	if len(m.logs) != 1 || !strings.HasPrefix(m.logs[0], "info:normal work") {
		t.Fatalf("expected info log, got %#v", m.logs)
	}
}
