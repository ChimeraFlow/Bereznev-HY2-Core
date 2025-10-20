//go:build mobile_skel && hc

package mobile

// test shim: на юнитах под mobile_skel возвращаем sing-транспорт вместо реального HC.
// Это позволяет проверить выбор транспорта и StartWithTun без Android/iOS окружения.
func newTransportHC(cfg HY2Config) Transport {
	return newTransportSingHY2(cfg)
}
