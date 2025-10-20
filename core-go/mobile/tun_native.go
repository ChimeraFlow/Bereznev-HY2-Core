//go:build (android || ios) && !mobile_skel

package mobile

import (
	"context"
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

func StartWithTun(configJSON string) string {
	mu.Lock()
	defer mu.Unlock()

	if err := cfgSet(configJSON); err != nil {
		return "invalid config: " + err.Error()
	}
	hc, err := parseHY2Config()
	if err != nil {
		return "invalid config: " + err.Error()
	}

	// 1) Поднимаем транспорт по Engine
	tr := selectTransport(hc)
	ctx, cancel := context.WithCancel(context.Background())
	if tr != nil {
		if err := tr.Start(ctx); err != nil {
			cancel()
			emitError(int(ErrEngineInitFailed), "hc/sing start: "+err.Error())
			return "engine init failed: " + err.Error()
		}
	}

	// 2) Запускаем TUN-мост
	// 2a) Вариант через локальный SOCKS:
	//     port := 10808; StartLocalSocks(port); мост TUN→SOCKS на этот порт
	// 2б) Вариант без SOCKS: передать dial hooks на основе транспорта tr
	//     (понадобится расширить Transport интерфейс методами TCP/UDP dial).

	// здесь оставь твой текущий путь (у тебя уже есть TUN→SOCKS мост):
	// err = startTun2Socks(hc, port)
	// if err != nil { cancel(); tr.Stop(context.Background()); return "tun start failed: " + err.Error() }

	// 3) сохраним cancel в runtime, emitting done на верхнем уровне уже происходит
	rtCancel = cancel
	started = true
	healthMarkStarted()
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
