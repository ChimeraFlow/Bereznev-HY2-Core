//go:build mobile_skel

package mobile

import (
	"fmt"
	"sync/atomic"
)

var (
	t2sRunning atomic.Bool
)

// StartTun2Socks — тестовый stub-раннер: только логи/события/флаг.
func StartTun2Socks(tunFd int, socksHost string, socksPort int) string {
	logI("🧩 STUB version StartTun2Socks() called")
	if t2sRunning.Load() {
		logI("tun2socks already running")
		return ""
	}
	addr := fmt.Sprintf("%s:%d", socksHost, socksPort)
	logI("starting tun2socks → " + addr)

	t2sRunning.Store(true)
	logI("tun2socks started successfully")
	emit("tun2socks_started", fmt.Sprintf(`{"socks":"%s"}`, addr))
	return ""
}

func StopTun2Socks() {
	if !t2sRunning.Load() {
		return
	}
	t2sRunning.Store(false)
	emit("tun2socks_stopped", "{}")
	logI("tun2socks stopped")
}
