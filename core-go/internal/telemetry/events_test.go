//go:build mobile_skel

package telemetry

import (
	"encoding/json"
	"testing"
)

/********* helpers *********/

// мок-структура для перехвата событий
type mockEventSink struct {
	received []struct {
		name string
		data string
	}
}

func (m *mockEventSink) OnEvent(name, data string) {
	m.received = append(m.received, struct {
		name string
		data string
	}{name: name, data: data})
}

/********* tests *********/

func TestSetEventSinkAndEmit(t *testing.T) {
	resetState() // сбрасывает evt через глобальный resetState из тестов

	m := &mockEventSink{}
	SetEventSink(m)

	emit("started", `{"ok":true}`)
	if len(m.received) != 1 {
		t.Fatalf("expected 1 event, got %d", len(m.received))
	}
	if m.received[0].name != "started" {
		t.Fatalf("expected name 'started', got %q", m.received[0].name)
	}
	if m.received[0].data != `{"ok":true}` {
		t.Fatalf("expected data '{\"ok\":true}', got %q", m.received[0].data)
	}
}

func TestEmit_NoSink_NoPanic(t *testing.T) {
	resetState()
	// emit при evt=nil не должен паниковать
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("emit() panicked: %v", r)
		}
	}()
	emit("nohandler", "{}")
}

func TestEmitStateConstants(t *testing.T) {
	resetState()
	m := &mockEventSink{}
	SetEventSink(m)

	emitState(EvtStarted)
	emitState(EvtStopped)
	emitState(EvtReloaded)

	if len(m.received) != 3 {
		t.Fatalf("expected 3 events, got %d", len(m.received))
	}
	want := []string{EvtStarted, EvtStopped, EvtReloaded}
	for i, ev := range m.received {
		if ev.name != want[i] {
			t.Errorf("event[%d] = %q, want %q", i, ev.name, want[i])
		}
		if ev.data != "{}" {
			t.Errorf("event[%d] data = %q, want '{}'", i, ev.data)
		}
	}
}

func TestEmitErrorPayload(t *testing.T) {
	resetState()
	m := &mockEventSink{}
	SetEventSink(m)

	emitError(42, "connection failed")
	if len(m.received) != 1 {
		t.Fatalf("expected 1 event, got %d", len(m.received))
	}
	ev := m.received[0]
	if ev.name != EvtError {
		t.Fatalf("expected event name %q, got %q", EvtError, ev.name)
	}

	var obj map[string]any
	if err := json.Unmarshal([]byte(ev.data), &obj); err != nil {
		t.Fatalf("invalid JSON in payload: %v", err)
	}
	if obj["code"] != float64(42) || obj["msg"] != "connection failed" {
		t.Fatalf("unexpected payload: %#v", obj)
	}
}
