// go:build android || ios || mobile_skel

package runtime

import (
	"context"
	"math/rand/v2"
	"time"
)

type backoffCfg struct {
	Base    time.Duration // 500ms
	Factor  float64       // 2.0
	Max     time.Duration // 30s
	Jitter  float64       // +-0.2
	FlapN   int           // 5 падений
	FlapWin time.Duration // 60s окно
	Cool    time.Duration // 60s пауза при флаппинге
}

type backoffState struct {
	cfg         backoffCfg
	attempt     int
	failTimes   []time.Time
	lastBackoff time.Duration
}

func newBackoffState() *backoffState {
	return &backoffState{
		cfg: backoffCfg{
			Base: 500 * time.Millisecond, Factor: 2.0, Max: 30 * time.Second,
			Jitter: 0.2, FlapN: 5, FlapWin: 60 * time.Second, Cool: 60 * time.Second,
		},
		failTimes: make([]time.Time, 0, 8),
	}
}

func (b *backoffState) next() time.Duration {
	// flapping guard
	now := time.Now()
	b.failTimes = append(b.failTimes, now)
	// drop old
	cut := now.Add(-b.cfg.FlapWin)
	n := 0
	for _, t := range b.failTimes {
		if t.After(cut) {
			b.failTimes[n] = t
			n++
		}
	}
	b.failTimes = b.failTimes[:n]

	if len(b.failTimes) >= b.cfg.FlapN {
		b.attempt = 0 // сбрасываем рост
		b.lastBackoff = b.cfg.Cool
		return b.lastBackoff
	}

	// expo
	var d time.Duration
	if b.attempt == 0 {
		d = b.cfg.Base
	} else {
		f := 1.0
		for i := 0; i < b.attempt; i++ {
			f *= b.cfg.Factor
		}
		d = time.Duration(float64(b.cfg.Base) * f)
	}
	if d > b.cfg.Max {
		d = b.cfg.Max
	}

	// jitter
	j := 1 + (rand.Float64()*2-1)*b.cfg.Jitter
	d = time.Duration(float64(d) * j)

	b.lastBackoff = d
	b.attempt++
	return d
}

func (b *backoffState) reset() {
	b.attempt = 0
	b.lastBackoff = 0
}

func (b *backoffState) last() time.Duration { return b.lastBackoff }

// waitNext blocks until either context done or duration elapsed
func (b *backoffState) waitNext(ctx context.Context) bool {
	d := b.next()
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}
