//go:build android || ios || mobile_skel

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
package mobile

import (
	"encoding/json"
	"sync/atomic"
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
	Running    bool   `json:"running"`
	Engine     string `json:"engine"`
	Version    string `json:"version"`
	BytesIn    uint64 `json:"bytes_in,omitempty"`
	BytesOut   uint64 `json:"bytes_out,omitempty"`
	Reconnects uint32 `json:"reconnects,omitempty"`
	QuicRttMs  int64  `json:"quic_rtt_ms,omitempty"`
}

// Глобальные счётчики (обновляются в tun2socks)
var (
	bytesIn    atomic.Uint64
	bytesOut   atomic.Uint64
	reconnects atomic.Uint32
	quicRttMs  atomic.Int64 // последняя измеренная оценка
)

// BytesStats возвращает текущие счётчики трафика.
// Используется HealthJSON() для вывода метрик.
func BytesStats() (uint64, uint64) {
	return bytesIn.Load(), bytesOut.Load()
}

// (опционально)
// ResetBytesStats сбрасывает счётчики — пригодится при Stop() или reload.

func ResetBytesStats() {
	bytesIn.Store(0)
	bytesOut.Store(0)
	reconnects.Store(0)
	quicRttMs.Store(0)
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
		Running:    IsRunning(),
		Engine:     EngineID,
		Version:    SdkVersion,
		BytesIn:    in,
		BytesOut:   out,
		Reconnects: reconnects.Load(),
		QuicRttMs:  quicRttMs.Load(),
	}
	b, _ := json.Marshal(h)
	return string(b)
}
