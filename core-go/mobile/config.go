//go:build android || ios || mobile_skel

// Package mobile — мобильный слой SDK (gomobile bind).
// Этот файл отвечает за работу с конфигурацией HY2 Core:
// хранение, валидацию и безопасное чтение JSON-конфига.
//
// Здесь используется расширенный парсер sing/common/json, который:
//   - поддерживает комментарии в JSON (//, /* */),
//   - допускает нестрогие ключи,
//   - совместим с синтаксисом sing-box/sing-tun.
//
// Потокобезопасность обеспечивается на уровне вызывающих функций (см. api.go).
package mobile

import sjson "github.com/sagernet/sing/common/json"

var (
	// cfgRaw — последний успешно применённый конфиг в виде байт.
	// Он хранится в памяти, чтобы можно было быстро получить текущее состояние.
	cfgRaw []byte
)

// cfgSet выполняет валидацию и сохраняет конфигурацию SDK.
// Использует sjson.UnmarshalExtended для поддержки комментариев и расширенного синтаксиса.
//
// Аргументы:
//   - jsonStr — строка JSON-конфига.
//
// Возвращает:
//   - nil — если JSON корректен;
//   - error — если формат некорректный или невалидный.
//
// Побочные эффекты:
//   - перезаписывает cfgRaw (глобальное состояние).
//
// Пример:
//
//	err := cfgSet(`{ "inbounds": [], "outbounds": [] }`)
//	if err != nil {
//	    log.Println("invalid config:", err)
//	}
func cfgSet(jsonStr string) error {
	// валидация расширенным парсером (комментарии, расширенные поля и т.п.)
	if _, err := sjson.UnmarshalExtended[map[string]any]([]byte(jsonStr)); err != nil {
		return err
	}
	cfgRaw = []byte(jsonStr)
	return nil
}

// cfgGet возвращает текущий сохранённый конфиг в виде байтового среза.
// Возвращаемое значение — копия (append к nil), чтобы исключить гонки данных.
//
// Используется при отладке, health-запросах или экспорте состояния
// в сторону Kotlin/Swift (через HealthJSON или debug-интерфейсы).
//
// Пример:
//
//	current := string(cfgGet())
//	fmt.Println("Active config:", current)
func cfgGet() []byte {
	return append([]byte(nil), cfgRaw...)
}
