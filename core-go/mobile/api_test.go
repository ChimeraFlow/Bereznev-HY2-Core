//go:build mobile_skel

package mobile

import (
	"strings"
	"sync"
	"testing"
)

/********* helpers *********/

type testEventSink struct {
	mu     sync.Mutex
	events []string
	data   []string
}

func (s *testEventSink) OnEvent(name, data string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events = append(s.events, name)
	s.data = append(s.data, data)
}

type testLogSink struct {
	mu   sync.Mutex
	logs []string
}

func (l *testLogSink) Log(level, msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logs = append(l.logs, level+":"+msg)
}

// сброс глобального состояния между тестами
func resetState() {
	mu.Lock()
	defer mu.Unlock()

	// best-effort stop рантайма, если вдруг поднят
	runtimeStop()

	started = false
	cfgRaw = nil
	logLevel = "info"
	evt = nil
	sink = nil

	netHooks = struct {
		mu      sync.RWMutex
		protect func(fd int) bool
	}{}
}

/********* tests *********/

const validCfg = `{
  "server":   "example.com:443",
  "password": "testpass",
  "sni":      "example.com",
  "alpn":     ["h3"],
  "mode":     "tun2socks"
}`

func TestStartStopStatus(t *testing.T) {
	resetState()

	es := &testEventSink{}
	SetEventSink(es)

	if Status() != "stopped" {
		t.Fatalf("expected initial status 'stopped', got %q", Status())
	}

	err := Start(validCfg)
	if err != "" {
		t.Fatalf("Start() unexpected error: %s", err)
	}
	if Status() != "running" {
		t.Fatalf("expected status 'running', got %q", Status())
	}

	hasStarted := false
	for _, ev := range es.events {
		if ev == EvtStarted || ev == "started" {
			hasStarted = true
			break
		}
	}
	if !hasStarted {
		t.Fatalf("expected at least one 'started' event, got %#v", es.events)
	}

	err2 := Start(validCfg)
	if err2 != "" {
		t.Fatalf("Start() second call unexpected error: %s", err2)
	}
	if len(es.events) != 1 {
		t.Fatalf("expected still one 'started' event, got %#v", es.events)
	}

	Stop()
	if Status() != "stopped" {
		t.Fatalf("expected status 'stopped' after Stop, got %q", Status())
	}
	if len(es.events) < 2 || (es.events[len(es.events)-1] != EvtStopped && es.events[len(es.events)-1] != "stopped") {
		t.Fatalf("expected 'stopped' event, got %#v", es.events)
	}
}

func TestStartWithCode(t *testing.T) {
	resetState()
	code := StartWithCode(validCfg)
	if code != ErrOK {
		t.Fatalf("StartWithCode(valid) = %d, want %d", code, ErrOK)
	}
	code2 := StartWithCode(validCfg)
	if code2 != ErrOK {
		t.Fatalf("StartWithCode(second) = %d, want %d", code2, ErrOK)
	}
}

func TestInvalidConfig(t *testing.T) {
	resetState()
	bad := `{ invalid json`
	err := Start(bad)
	if err == "" || !strings.HasPrefix(err, "invalid config:") {
		t.Fatalf("expected invalid config error, got %q", err)
	}
	if IsRunning() {
		t.Fatalf("expected not running after invalid config")
	}
}

func TestReload(t *testing.T) {
	resetState()
	es := &testEventSink{}
	SetEventSink(es)

	if err := Start(validCfg); err != "" {
		t.Fatalf("Start() unexpected error: %s", err)
	}
	if err := Reload(`{"log":{"level":"info"}}`); err != "" {
		t.Fatalf("Reload(valid) unexpected error: %s", err)
	}
	found := false
	for _, ev := range es.events {
		if ev == EvtReloaded || ev == "reloaded" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected 'reloaded' event, got %#v", es.events)
	}

	if err := Reload(`{ bad json`); err == "" || !strings.HasPrefix(err, "invalid config:") {
		t.Fatalf("expected invalid config error on Reload, got %q", err)
	}
}

func TestVersionFormat(t *testing.T) {
	resetState()
	v := Version()
	if !strings.Contains(v, sdkName) || !strings.Contains(v, SdkVersion) || !strings.Contains(v, EngineID) {
		t.Fatalf("Version() %q must contain sdkName, sdkVersion and engineID", v)
	}
}

func TestLoggingLevel(t *testing.T) {
	resetState()
	ls := &testLogSink{}
	SetLogger(ls)

	logD("debug message") // по умолчанию скрыт (logLevel=info)
	logI("info message")

	ls.mu.Lock()
	if len(ls.logs) != 1 || !strings.HasPrefix(ls.logs[0], "info:") {
		ls.mu.Unlock()
		t.Fatalf("expected only info log, got %#v", ls.logs)
	}
	ls.mu.Unlock()

	SetLogLevel("debug")
	logD("now visible")

	ls.mu.Lock()
	defer ls.mu.Unlock()
	if len(ls.logs) != 2 || !strings.HasPrefix(ls.logs[1], "debug:") {
		t.Fatalf("expected debug log after SetLogLevel, got %#v", ls.logs)
	}
}
