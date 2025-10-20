//go:build android || ios || mobile_skel

// Package mobile — мобильный слой SDK (gomobile bind).
// Этот файл реализует систему событий SDK — EventSink и emit(),
// которая используется для связи между Go-ядром и платформенным кодом (Kotlin/Swift).
//
// EventSink — это универсальный канал коммуникации из Go в сторону UI/SDK.
// Через него SDK сообщает приложению о ключевых состояниях и ошибках:
// запуск, остановка, перезагрузка, паника, reconnect и т.д.
//
// События используются внутри api.go, lifecycle.go и будущих модулей (HY2 runtime).
package mobile

import "fmt"

// EventSink — интерфейс для передачи событий из Go в Kotlin/Swift.
// Реализуется на стороне платформенного кода (например, в Kotlin SDK).
//
// Каждый вызов OnEvent(name, data) получает:
//   - name — строку с типом события ("started", "stopped", "error", ...)
//   - data — JSON-полезную нагрузку, описывающую детали.
//
// Пример (Kotlin):
// ```kotlin
//
//	class MyEventHandler : EventSink {
//	    override fun OnEvent(name: String, data: String) {
//	        println("Event: $name, Data: $data")
//	    }
//	}
//
// hy2core.SetEventSink(MyEventHandler())
// ```
type EventSink interface{ OnEvent(name, data string) }

func (f EventSinkFunc) OnEvent(name, payload string) { f(name, payload) }

var (
	// evt — текущий зарегистрированный EventSink.
	// Если не установлен, события будут проигнорированы.
	evt EventSink
)

// SetEventSink регистрирует внешний обработчик событий SDK.
// Обычно вызывается из Kotlin/Swift после инициализации ядра.
//
// Аргументы:
//   - s — объект, реализующий интерфейс EventSink.
//
// Побочные эффекты:
//   - сохраняет ссылку в глобальную переменную evt.
func SetEventSink(s EventSink) { evt = s }

// emit — внутренняя функция для отправки события (если EventSink установлен).
// Используется другими модулями SDK (api.go, lifecycle.go, logging.go).
//
// Аргументы:
//   - name — идентификатор события (строка);
//   - data — JSON-полезная нагрузка (может быть "{}").
//
// Безопасна для вызова из любых горутин.
func emit(name, data string) {
	if evt != nil {
		evt.OnEvent(name, data)
	}
}

// ===============================
// 🧭 Типизированные события
// ===============================

// Константы имён событий.
// Определены централизованно, чтобы избежать расхождений между слоями.
const (
	EvtStarted      = "started"   // ядро запущено
	EvtStopped      = "stopped"   // ядро остановлено
	EvtReloaded     = "reloaded"  // конфигурация перезагружена
	EvtPanic        = "panic"     // panic() перехвачена
	EvtError        = "error"     // ошибка выполнения
	EvtReconnect    = "reconnect" // переподключение / попытка восстановления
	EvtReconnecting = "reconnecting"
	EvtReconnected  = "reconnected"
)

type evtReconnecting struct {
	Reason  string `json:"reason"`
	Attempt int    `json:"attempt"`
	NextMs  int    `json:"next_ms"`
}

type evtReconnected struct {
	RttMs int64 `json:"rtt_ms"`
}

type EventSinkFunc func(name, payload string)

// emitError — отправляет стандартное событие об ошибке.
// Используется при ошибках конфигурации, сетевых сбоях и т.п.
//
// Аргументы:
//   - code — числовой код ошибки;
//   - msg — текст ошибки (короткое описание).
//
// Пример:
//
//	emitError(1001, "connection timeout")
func emitError(code int, msg string) {
	emit(EvtError, fmt.Sprintf(`{"code":%d,"msg":%q}`, code, msg))
}

// emitState — отправляет простые служебные события
// (например, started, stopped, reloaded).
//
// Аргументы:
//   - event — одно из значений EvtStarted, EvtStopped, EvtReloaded.
//
// Пример:
//
//	emitState(EvtStarted)
func emitState(event string) {
	emit(event, "{}")
}
