// go:build android || ios || mobile_skel

package runtime

import (
	"context"
	"sync"
	"time"

	"github.com/ChimeraFlow/Bereznev-HY2-Core/core-go/internal/telemetry"
	"github.com/ChimeraFlow/Bereznev-HY2-Core/core-go/internal/transport"
	"github.com/ChimeraFlow/Bereznev-HY2-Core/core-go/pkg/config"
)

var (
	RtMu      sync.Mutex
	RtStarted bool
	RtCancel  context.CancelFunc
	RtTrans   transport.Transport
	RtUptime  time.Time
)

func RuntimeStart() error {
	if RtStarted {
		return nil
	}

	hc, err := config.ParseHY2Config()
	if err != nil {
		return err
	}

	// Выбор реализации переносим в selectTransport (см. ниже)
	RtTrans = transport.SelectTransport(hc)

	// контекст и запуск
	ctx, cancel := context.WithCancel(context.Background())
	if RtTrans != nil {
		if err := RtTrans.Start(ctx); err != nil {
			cancel()
			return err
		}
	}
	RtCancel = cancel
	RtStarted = true
	RtUptime = time.Now()
	telemetry.Emit(telemetry.EvtStarted, "{}")
	return nil
}

func RuntimeStop() {
	if !RtStarted {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if RtTrans != nil {
		_ = RtTrans.Stop(ctx)
		RtTrans = nil
	}
	RtCancel()
	RtStarted = false
	telemetry.Emit(telemetry.EvtStopped, "{}")
}

func RuntimeStatusInto(h *telemetry.Health) {
	if !RtStarted || RtTrans == nil {
		return
	}
	st := RtTrans.Status()
	if st.RTTms > 0 {
		h.QuicRttMs = st.RTTms
	}
}
