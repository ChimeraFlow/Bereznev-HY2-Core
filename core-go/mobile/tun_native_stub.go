//go:build mobile_skel

package mobile

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var (
	tunMu      sync.Mutex
	tunStarted atomic.Bool
)

func StartWithTun(tunFd int, mtu int) string {
	tunMu.Lock()
	defer tunMu.Unlock()

	if tunStarted.Load() {
		logI("sing-tun already running")
		return ""
	}
	if mtu <= 0 {
		mtu = 1500
	}
	tunStarted.Store(true)
	logI(fmt.Sprintf("sing-tun created (mtu=%d)", mtu))
	emit(EvtStarted, `{"path":"tun","engine":"sing-tun"}`)
	return ""
}

func SetMTU(mtu int) {
	logI(fmt.Sprintf("SetMTU requested: %d (stub, no live update)", mtu))
}

func stopSingTun() {
	tunMu.Lock()
	defer tunMu.Unlock()
	if !tunStarted.Load() {
		return
	}
	tunStarted.Store(false)
	emit(EvtStopped, `{"path":"tun","engine":"sing-tun"}`)
	logI("sing-tun stopped")
}
