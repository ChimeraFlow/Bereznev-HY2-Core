//go:build (android || ios) && !mobile_skel

package mobile

import (
	"fmt"
	"sync"
	"sync/atomic"

	tun "github.com/sagernet/sing-tun"
)

var (
	tunMu      sync.Mutex
	tunRunner  tun.Tun // интерфейс из sing-tun
	tunStarted atomic.Bool
)

// Прокидываем Android VpnService.protect(fd) — если понадобится, можно впоследствии
// менять Options (см. примечание ниже).
func protectFn(fd int) bool { return protectFD(fd) }

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

	// Важно: в sing-tun v0.7.x публичный Options минималистичный.
	// В нём нет полей TunFD/Protect. Биндим минимально-жизнеспособный инстанс:
	opts := tun.Options{
		Name: "bereznev-tun", // имя интерфейса (sing-tun создаст/поднимет)
		MTU:  uint32(mtu),
	}

	r, err := tun.New(opts)
	if err != nil {
		logE("sing-tun init failed: " + err.Error())
		return "sing-tun init failed: " + err.Error()
	}

	// Запуск
	if err := r.Start(); err != nil {
		_ = r.Close()
		logE("sing-tun start failed: " + err.Error())
		return "sing-tun start failed: " + err.Error()
	}

	tunRunner = r
	tunStarted.Store(true)
	logI(fmt.Sprintf("sing-tun started (name=%s, mtu=%d)", opts.Name, mtu))
	emit(EvtStarted, `{"path":"tun","engine":"sing-tun"}`)

	return ""
}

func SetMTU(mtu int) {
	// На v0.7.x live-MTU в публичном API отсутствует — делаем recreate-подход (в будущем).
	logI(fmt.Sprintf("SetMTU requested: %d (not supported live; consider recreate)", mtu))
}

func stopSingTun() {
	tunMu.Lock()
	defer tunMu.Unlock()
	if !tunStarted.Load() {
		return
	}
	tunStarted.Store(false)
	if tunRunner != nil {
		_ = tunRunner.Close()
		tunRunner = nil
	}
	emit(EvtStopped, `{"path":"tun","engine":"sing-tun"}`)
	logI("sing-tun stopped")
}
