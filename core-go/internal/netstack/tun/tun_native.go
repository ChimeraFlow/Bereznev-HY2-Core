//go:build (android || ios) && !mobile_skel

package tun

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	singtun "github.com/sagernet/sing-tun"

	// наши пакеты
	"github.com/ChimeraFlow/Bereznev-HY2-Core/core-go/internal/netstack/protect"
	"github.com/ChimeraFlow/Bereznev-HY2-Core/core-go/internal/runtime"
	"github.com/ChimeraFlow/Bereznev-HY2-Core/core-go/internal/telemetry"
	"github.com/ChimeraFlow/Bereznev-HY2-Core/core-go/internal/transport"
	"github.com/ChimeraFlow/Bereznev-HY2-Core/core-go/mobile"
	"github.com/ChimeraFlow/Bereznev-HY2-Core/core-go/pkg/config"
	"github.com/ChimeraFlow/Bereznev-HY2-Core/core-go/pkg/errors"
	logpkg "github.com/ChimeraFlow/Bereznev-HY2-Core/core-go/pkg/logging"
)

var (
	TunMu      sync.Mutex
	TunRunner  singtun.Tun // интерфейс из sing-tun
	TunStarted atomic.Bool
)

// Прокидываем Android VpnService.protect(fd) — если понадобится, можно впоследствии
// менять Options (см. примечание ниже).
func protectFn(fd int) bool { return protect.ProtectFD(fd) }

func StartWithTun(configJSON string) string {
	mobile.Mu.Lock()
	defer mobile.Mu.Unlock()

	if err := mobile.CfgSet(configJSON); err != nil {
		return "invalid config: " + err.Error()
	}
	hc, err := config.ParseHY2Config()
	if err != nil {
		return "invalid config: " + err.Error()
	}

	// 1) Поднимаем транспорт по Engine
	tr := transport.SelectTransport(hc)
	ctx, cancel := context.WithCancel(context.Background())
	if tr != nil {
		if err := tr.Start(ctx); err != nil {
			cancel()
			telemetry.EmitError(int(errors.ErrEngineInitFailed), "hc/sing start: "+err.Error())
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
	runtime.RtCancel = cancel
	mobile.Started = true
	telemetry.HealthMarkStarted()
	return ""
}

func SetMTU(mtu int) {
	// На v0.7.x live-MTU в публичном API отсутствует — делаем recreate-подход (в будущем).
	logpkg.Info(fmt.Sprintf("SetMTU requested: %d (not supported live; recreate required)", mtu))

}

func stopSingTun() {
	TunMu.Lock()
	defer TunMu.Unlock()
	if !TunStarted.Load() {
		return
	}
	TunStarted.Store(false)
	if TunRunner != nil {
		_ = TunRunner.Close()
		TunRunner = nil
	}
	telemetry.Emit(telemetry.EvtStopped, `{"path":"tun","engine":"sing-tun"}`)
	logpkg.Info("sing-tun stopped")
}
