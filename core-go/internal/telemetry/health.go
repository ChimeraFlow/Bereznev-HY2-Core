// go:build android || ios || mobile_skel

// Package mobile — мобильный слой SDK (gomobile bind).
// Этот файл реализует системный модуль Health — точку доступа к состоянию ядра.
// Health используется для мониторинга, отладки и UI-диагностики,
// предоставляя информацию о статусе SDK, версии и сетевых метриках.
//
// Основное предназначение Health:
//   - вернуть агрегированное состояние ядра (запущено/остановлено);
//   - передать базовые метрики (версии, трафик, RTT и т.п.);
//   - использоваться в UI или логах SDK-уровня.
//
// HealthJSON — публичная функция, экспортируемая через gomobile.
// Она возвращает состояние в формате JSON-строки, готовой для парсинга
// на стороне Kotlin/Swift.
package telemetry

import (
	"encoding/json"
	"sync/atomic"
	"time"

	"github.com/ChimeraFlow/Bereznev-HY2-Core/core-go/internal/telemetry"
	"github.com/ChimeraFlow/Bereznev-HY2-Core/core-go/mobile"
	"github.com/ChimeraFlow/Bereznev-HY2-Core/core-go/pkg/version"
)

// Health — структура состояния SDK.
// Это универсальная форма для сериализации данных о работе HY2-ядра.
//
// JSON-теги заданы в формате snake_case для совместимости с Kotlin/Swift.
//
// Поля:
//   - Running — активно ли ядро (true / false);
//   - Engine — идентификатор движка ("skeleton", "sing-tun" и т.п.);
//   - Version — версия SDK (из sdkVersion);
//   - BytesIn / BytesOut — счётчики трафика (будущие поля);
//   - Reconnects — количество попыток переподключения (будущее);
//   - QuicRttMs — средний RTT в миллисекундах (будущее).
type Health struct {
	Running       bool   `json:"running"`
	Engine        string `json:"engine"`
	Version       string `json:"version"`
	BytesIn       uint64 `json:"bytes_in,omitempty"`
	BytesOut      uint64 `json:"bytes_out,omitempty"`
	Reconnects    uint32 `json:"reconnects,omitempty"`
	QuicRttMs     int64  `json:"quic_rtt_ms,omitempty"`
	UptimeS       int64  `json:"uptime_s,omitempty"`
	SNI           string `json:"sni,omitempty"`
	ALPN          string `json:"alpn,omitempty"`
	LastBackoffMs int64  `json:"last_backoff_ms"`
	LastErrorTs   int64  `json:"last_error_ts"`
}

// Глобальные счётчики (обновляются в tun2socks)
var (
	BytesIn    atomic.Uint64
	BytesOut   atomic.Uint64
	Reconnects atomic.Uint32
	QuicRttMs  atomic.Int64 // последняя измеренная оценка
	StartUnix  atomic.Int64
	SniValue   atomic.Value // string
	AlpnValue  atomic.Value // string

	// ⬇️ новое
	LastBackoffMs atomic.Int64
	LastErrTs     atomic.Int64
)

// BytesStats возвращает текущие счётчики трафика.
// Используется HealthJSON() для вывода метрик.
func BytesStats() (uint64, uint64) {
	return BytesIn.Load(), BytesOut.Load()
}

func SetLastBackoffMs(ms int64) { LastBackoffMs.Store(ms) }
func SetLastErrTs(ts int64)     { LastErrTs.Store(ts) }

// (опционально)
// ResetBytesStats сбрасывает счётчики — пригодится при Stop() или reload.

func ResetBytesStats() {
	BytesIn.Store(0)
	BytesOut.Store(0)
	Reconnects.Store(0)
	QuicRttMs.Store(0)
}

func HealthMarkStarted() {
	StartUnix.Store(time.Now().Unix())
}

func HealthMarkStopped() {
	StartUnix.Store(0)
}
func HealthSetIdentity(sni, alpn string) {
	if sni != "" {
		SniValue.Store(sni)
	}
	if alpn != "" {
		AlpnValue.Store(alpn)
	}
}

// HealthJSON возвращает агрегированное состояние ядра в виде JSON-строки.
//
// Возвращаемая строка готова к использованию на платформенном уровне —
// её можно напрямую декодировать в Kotlin/Swift или отобразить в UI.
//
// Пример JSON:
//
//	{
//	  "running": true,
//	  "engine": "skeleton",
//	  "version": "0.1.0"
//	}
//
// Потокобезопасность:
//   - Вызовы безопасны, так как HealthJSON обращается к IsRunning()
//     (который синхронизирован mutex'ом в api.go).
//
// Пример:
//
//	status := hy2core.HealthJSON()
//	println("SDK Status:", status)
func HealthJSON() string {
	in, out := BytesStats()
	h := Health{
		Running:    mobile.IsRunning(),
		Engine:     version.EngineID,
		Version:    version.SdkVersion,
		BytesIn:    in,
		BytesOut:   out,
		Reconnects: telemetry.Reconnects.Load(),
		QuicRttMs:  QuicRttMs.Load(),
	}
	h.LastBackoffMs = LastBackoffMs.Load()
	h.LastErrorTs = LastErrTs.Load()
	if su := StartUnix.Load(); su > 0 {
		now := time.Now().Unix()
		if now > su {
			h.UptimeS = now - su
		}
	}

	if v := SniValue.Load(); v != nil {
		if s, _ := v.(string); s != "" {
			h.SNI = s
		}
	}

	if v := AlpnValue.Load(); v != nil {
		if s, _ := v.(string); s != "" {
			h.ALPN = s
		}
	}

	b, _ := json.Marshal(h)
	return string(b)
}
