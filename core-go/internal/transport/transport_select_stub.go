// go:build (android || ios || mobile_skel) && !hc

package transport

import (
	"github.com/ChimeraFlow/Bereznev-HY2-Core/core-go/internal/transport/hy2hc"
	"github.com/ChimeraFlow/Bereznev-HY2-Core/core-go/internal/transport/sing"
)

func SelectTransport(cfg hy2hc.HY2Config) Transport {
	// Тег hc не активен — всегда используем sing-транспорт
	return sing.NewTransportSingHY2(cfg)
}
