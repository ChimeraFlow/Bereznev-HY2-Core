//go:build android || ios || mobile_skel

package mobile

import (
	"context"
	"sync"
	"time"
)

var (
	rtMu      sync.Mutex
	rtStarted bool
	rtCancel  context.CancelFunc
	rtTrans   Transport
	rtUptime  time.Time
)

func runtimeStart() error {
	if rtStarted {
		return nil
	}

	hc, err := parseHY2Config()
	if err != nil {
		return err
	}

	// Выбор реализации переносим в selectTransport (см. ниже)
	rtTrans = selectTransport(hc)

	// контекст и запуск
	ctx, cancel := context.WithCancel(context.Background())
	if rtTrans != nil {
		if err := rtTrans.Start(ctx); err != nil {
			cancel()
			return err
		}
	}
	rtCancel = cancel
	rtStarted = true
	rtUptime = time.Now()
	emit(EvtStarted, "{}")
	return nil
}

func runtimeStop() {
	if !rtStarted {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if rtTrans != nil {
		_ = rtTrans.Stop(ctx)
		rtTrans = nil
	}
	rtCancel()
	rtStarted = false
	emit(EvtStopped, "{}")
}

func runtimeStatusInto(h *Health) {
	if !rtStarted || rtTrans == nil {
		return
	}
	st := rtTrans.Status()
	if st.RTTms > 0 {
		h.QuicRttMs = st.RTTms
	}
}
