//go:build mobile_skel

package tun

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/apernet/quic-go/logging"
)

var (
	tunMu      sync.Mutex
	tunStarted atomic.Bool
)

func StartWithTun(tunFd int, mtu int) string {
	tunMu.Lock()
	defer tunMu.Unlock()

	if tunStarted.Load() {
		logging.logI("sing-tun already running")
		return ""
	}
	if mtu <= 0 {
		mtu = 1500
	}
	tunStarted.Store(true)
	logging.logI(fmt.Sprintf("sing-tun created (mtu=%d)", mtu))
	logging.emit(EvtStarted, `{"path":"tun","engine":"sing-tun"}`)
	return ""
}

func SetMTU(mtu int) {
	logging.logI(fmt.Sprintf("SetMTU requested: %d (stub, no live update)", mtu))
}

func stopSingTun() {
	tunMu.Lock()
	defer tunMu.Unlock()
	if !tunStarted.Load() {
		return
	}
	tunStarted.Store(false)
	logging.Emit(EvtStopped, `{"path":"tun","engine":"sing-tun"}`)
	logging.LogI("sing-tun stopped")
}
