//go:build mobile_skel

package mobile

import (
	"strings"
	"sync"
	"testing"
)

type socksTestLogSink struct {
	mu   sync.Mutex
	logs []string
}

func (l *socksTestLogSink) Log(level, msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logs = append(l.logs, level+":"+msg)
}

func TestStartTun2Socks_LogsAndOk(t *testing.T) {
	// общий ресет (из api_test.go), чтобы не тащить состояние между тестами
	resetState()

	ls := &socksTestLogSink{}
	SetLogger(ls)

	// любые значения — сейчас они не используются реальной реализацией
	err := StartTun2Socks(123, "127.0.0.1", 1080)
	if err != "" {
		t.Fatalf("expected empty error, got: %q", err)
	}

	// проверяем, что было info-сообщение из заглушки
	ls.mu.Lock()
	defer ls.mu.Unlock()

	found := false
	for _, l := range ls.logs {
		if strings.HasPrefix(l, "info:StartTun2Socks() requested") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected info log with StartTun2Socks(), got %#v", ls.logs)
	}
}
