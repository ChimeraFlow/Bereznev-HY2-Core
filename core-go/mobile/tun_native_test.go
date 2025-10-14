//go:build mobile_skel

package mobile

import (
	"strings"
	"sync"
	"testing"
)

// mock лог-приёмник для проверки вывода
type tunTestLogSink struct {
	mu   sync.Mutex
	logs []string
}

func (l *tunTestLogSink) Log(level, msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logs = append(l.logs, level+":"+msg)
}

func TestStartWithTun_StubLogAndReturn(t *testing.T) {
	resetState()
	ls := &tunTestLogSink{}
	SetLogger(ls)

	err := StartWithTun(55, 1500)
	if err != "" {
		t.Fatalf("expected empty return (stub), got %q", err)
	}

	// проверяем, что записался info-лог
	ls.mu.Lock()
	defer ls.mu.Unlock()
	found := false
	for _, l := range ls.logs {
		if strings.Contains(l, "StartWithTun() not implemented yet") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected info log for StartWithTun(), got %#v", ls.logs)
	}
}

func TestSetMTU_StubLog(t *testing.T) {
	resetState()
	ls := &tunTestLogSink{}
	SetLogger(ls)

	SetMTU(1400)

	ls.mu.Lock()
	defer ls.mu.Unlock()
	found := false
	for _, l := range ls.logs {
		if strings.Contains(l, "SetMTU not implemented yet") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected log for SetMTU(), got %#v", ls.logs)
	}
}
