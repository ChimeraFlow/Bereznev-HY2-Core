//go:build (android || ios) && hc

package transport

func selectTransport(cfg HY2Config) Transport {
	switch cfg.Engine {
	case "hc", "hysteria_core":
		return newTransportHC(cfg) // определён в файлах под (android||ios)&&hc
	default:
		return newTransportSingHY2(cfg)
	}
}
