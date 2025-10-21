//go:build mobile_skel

package protect

import (
	"sync/atomic"
	"testing"
)

type fakeHooks struct {
	calls int32
	last  int
	ok    bool
}

func (f *fakeHooks) Protect(fd int) bool {
	atomic.AddInt32(&f.calls, 1)
	f.last = fd
	return f.ok
}
func TestNetHooks_SetAndProtect(t *testing.T) {
	resetState()

	h := &fakeHooks{ok: true}
	SetNetHooks(h)

	if !protectFD(42) {
		t.Fatal("expected protectFD() to return true with valid hook")
	}
	if atomic.LoadInt32(&h.calls) != 1 {
		t.Fatalf("expected Protect() to be called once, got %d", h.calls)
	}
	if h.last != 42 {
		t.Fatalf("expected fd=42, got %d", h.last)
	}
}

func TestNetHooks_NotSet(t *testing.T) {
	resetState()
	// netHooks должен быть nil
	if protectFD(10) {
		t.Fatal("expected false when no hook set")
	}
}

func TestNetHooks_FalseReturn(t *testing.T) {
	resetState()
	h := &fakeHooks{ok: false}
	SetNetHooks(h)
	if protectFD(7) {
		t.Fatal("expected false when hook returns false")
	}
	if atomic.LoadInt32(&h.calls) != 1 {
		t.Fatalf("expected Protect() to be called once, got %d", h.calls)
	}
}
