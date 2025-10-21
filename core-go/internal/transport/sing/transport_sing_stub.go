//go:build mobile_skel

package sing

import (
	"context"
	"errors"
	"sync/atomic"
)

var testStartOnceFailCount atomic.Int32

// тестовый хук — используем его в reconnect_test.go
func testSetStartFailCount(n int32) { testStartOnceFailCount.Store(n) }

func startOnceSing(t *transportSingHY2, ctx context.Context) error {
	if testStartOnceFailCount.Load() > 0 {
		testStartOnceFailCount.Add(-1)
		return errors.New("mock dial fail")
	}
	// успех: соединение «поднялось»
	t.rtt.Store(42)
	t.rem = "mock.remote:443"
	return nil
}

func isAliveSing(t *transportSingHY2) bool {
	return t.rtt.Load() > 0
}
