//go:build mobile_skel

package mobile

import (
	"strings"
	"sync"
	"testing"
)

// mock для логов
type testSink struct {
	mu   sync.Mutex
	logs []string
}

func (t *testSink) Log(level, msg string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.logs = append(t.logs, level+":"+msg)
}

func TestLogging_FilteringAndLevels(t *testing.T) {
	resetState()

	ts := &testSink{}
	SetLogger(ts)
	SetLogLevel("info")

	// debug не должен пройти при info
	logD("debug suppressed")
	logI("info visible")
	logW("warn visible")
	logE("error visible")

	ts.mu.Lock()
	got := append([]string(nil), ts.logs...)
	ts.mu.Unlock()

	if len(got) != 3 {
		t.Fatalf("expected 3 visible logs (info,warn,error), got %d: %#v", len(got), got)
	}
	if !strings.HasPrefix(got[0], "info:") || !strings.HasPrefix(got[1], "warn:") || !strings.HasPrefix(got[2], "error:") {
		t.Fatalf("unexpected log order/content: %#v", got)
	}
}

func TestLogging_DebugLevelShowsAll(t *testing.T) {
	resetState()

	ts := &testSink{}
	SetLogger(ts)
	SetLogLevel("debug")

	logD("debug 1")
	logI("info 2")
	logW("warn 3")
	logE("error 4")

	ts.mu.Lock()
	defer ts.mu.Unlock()

	if len(ts.logs) != 4 {
		t.Fatalf("expected 4 logs at debug level, got %d", len(ts.logs))
	}
	if !strings.HasPrefix(ts.logs[0], "debug:") {
		t.Fatalf("first log should be debug, got %#v", ts.logs)
	}
}

func TestLogging_NoSink_NoPanic(t *testing.T) {
	resetState()
	// ни одного SetLogger()
	logI("no sink")
	// просто проверяем, что не паника — если дошли сюда, тест пройден
}
