//go:build android || ios || mobile_skel

package protect

import (
	"sync"

	logpkg "github.com/ChimeraFlow/Bereznev-HY2-Core/core-go/pkg/logging"
)

var netHooks struct {
	mu      sync.RWMutex
	protect func(fd int) bool
}

// SetProtectHook — новая версия API (Kotlin → Go).
// На Android вызывается из VpnService.protect(fd).
func SetProtectHook(fn func(fd int) bool) {
	netHooks.mu.Lock()
	defer netHooks.mu.Unlock()
	netHooks.protect = fn
	logpkg.Info("Protect hook registered")
}

// 🔄 Backward-compat shim для старых тестов / API
// принимает объект с методом Protect(fd int) bool.
func SetNetHooks(h interface{}) {
	netHooks.mu.Lock()
	defer netHooks.mu.Unlock()

	switch v := h.(type) {
	case interface{ Protect(fd int) bool }:
		netHooks.protect = v.Protect
	default:
		netHooks.protect = nil
	}
	logpkg.Info("Legacy SetNetHooks() adapter called")
}

// protectFD — универсальная обёртка, вызываемая из ядра.
// Возвращает true, если fd защищён успешно.
func ProtectFD(fd int) bool {
	netHooks.mu.RLock()
	defer netHooks.mu.RUnlock()
	if netHooks.protect == nil {
		return false
	}
	ok := netHooks.protect(fd)
	if !ok {
		logpkg.Warn("protectFD failed")
	}
	return ok
}
