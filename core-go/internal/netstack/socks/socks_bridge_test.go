//go:build mobile_skel

package socks

import (
	"strings"
	"sync"
	"testing"
)

// socksTestLogSink — тестовый логгер для перехвата логов из tun2socks.
type socksTestLogSink struct {
	mu   sync.Mutex
	logs []string
}

func (l *socksTestLogSink) Log(level, msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logs = append(l.logs, level+":"+msg)
}

// TestStartTun2Socks_LogsAndOk — проверяет базовый запуск и логи.
func TestStartTun2Socks_LogsAndOk(t *testing.T) {
	t.Log("🚀 TestStartTun2Socks_LogsAndOk started")
	resetState()
	ls := &socksTestLogSink{}
	SetLogger(ls)

	// ⚙️ тестовый вызов (tunFd фейковый, но функция не должна упасть)
	err := StartTun2Socks(123, "127.0.0.1", 1080)
	if err != "" {
		t.Fatalf("unexpected error from StartTun2Socks: %q", err)
	}

	// 🧩 Проверяем, что есть лог начала и успешного запуска
	var logsCopy []string
	ls.mu.Lock()
	logsCopy = append([]string(nil), ls.logs...)
	ls.mu.Unlock()

	var started, success bool
	for _, line := range logsCopy {
		if strings.Contains(line, "starting tun2socks") {
			started = true
		}
		if strings.Contains(line, "tun2socks started successfully") {
			success = true
		}
	}

	if !started {
		t.Fatalf("expected log 'starting tun2socks', got logs: %#v", ls.logs)
	}
	if !success {
		t.Fatalf("expected log 'tun2socks started successfully', got logs: %#v", ls.logs)
	}

	// 🧩 Проверяем, что статус t2sRunning = true
	if !t2sRunning.Load() {
		t.Fatalf("expected t2sRunning=true after start")
	}

	// ⚙️ Останавливаем и проверяем, что статус сбросился
	StopTun2Socks()
	if t2sRunning.Load() {
		t.Fatalf("expected t2sRunning=false after stop")
	}
}
