//go:build mobile_skel

package mobile

import "testing"

type captureSink struct{ names []string }

func (c *captureSink) OnEvent(name, data string) { c.names = append(c.names, name) }

func TestRuntime_StartStop_StatusHook(t *testing.T) {
	resetState()
	_ = cfgSet(`{"server":"s:443","password":"p"}`)
	sink := &captureSink{}
	SetEventSink(sink)

	if err := runtimeStart(); err != nil {
		t.Fatalf("runtimeStart: %v", err)
	}
	defer runtimeStop()

	if len(sink.names) == 0 || sink.names[0] != "started" {
		t.Fatalf("expected first event 'started', got %#v", sink.names)
	}

	// Проставим RTT и проверим, что runtimeStatusInto прокидывает его в Health
	if rtTrans == nil {
		t.Fatal("rtTrans is nil after runtimeStart")
	}
	tr, ok := rtTrans.(*transportSingHY2)
	if !ok {
		t.Fatalf("unexpected transport type: %T", rtTrans)
	}
	tr.rtt.Store(42)

	var h Health
	runtimeStatusInto(&h)
	if h.QuicRttMs != 42 {
		t.Fatalf("expected QuicRttMs=42, got %d", h.QuicRttMs)
	}

	runtimeStop()
	if sink.names[len(sink.names)-1] != "stopped" {
		t.Fatalf("expected last event 'stopped', got %#v", sink.names)
	}
}
