// go:build (android || ios || mobile_skel) && !hc

package transport

func selectTransport(cfg HY2Config) Transport {
	// Тег hc не активен — всегда используем sing-транспорт
	return newTransportSingHY2(cfg)
}
