//go:build mobile_skel

package socks

import (
	"fmt"
	"sync/atomic"

	"github.com/ChimeraFlow/Bereznev-HY2-Core/core-go/internal/telemetry"
	logpkg "github.com/ChimeraFlow/Bereznev-HY2-Core/core-go/pkg/logging"
)

var (
	t2sRunning atomic.Bool
)

// StartTun2Socks — тестовый stub-раннер: только логи/события/флаг.
func StartTun2Socks(tunFd int, socksHost string, socksPort int) string {
	logpkg.LogI("🧩 STUB version StartTun2Socks() called")
	if t2sRunning.Load() {
		logpkg.LogI("tun2socks already running")
		return ""
	}
	addr := fmt.Sprintf("%s:%d", socksHost, socksPort)
	logpkg.LogI("starting tun2socks → " + addr)

	t2sRunning.Store(true)
	logpkg.LogI("tun2socks started successfully")
	telemetry.Emit("tun2socks_started", fmt.Sprintf(`{"socks":"%s"}`, addr))
	return ""
}

func StopTun2Socks() {
	if !t2sRunning.Load() {
		return
	}
	t2sRunning.Store(false)
	telemetry.Emit("tun2socks_stopped", "{}")
	logpkg.LogI("tun2socks stopped")
}
