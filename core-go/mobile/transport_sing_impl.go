//go:build android || ios

package mobile

import "context"

// startOnceSing — реальный запуск sing/hysteria2 (пока TODO).
func startOnceSing(t *transportSingHY2, ctx context.Context) error {
	// TODO: собрать dialer с Protect(fd), поднять HY2, заполнить t.rem, t.rtt.Store(...)
	return nil
}

// isAliveSing — реальная проверка живости (например, по lastRTTAt или состоянию клиента).
func isAliveSing(t *transportSingHY2) bool {
	return t.rtt.Load() > 0 // временно
}
