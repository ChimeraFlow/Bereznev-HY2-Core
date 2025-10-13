//go:build android || ios

// Package mobile — мобильный слой SDK (gomobile bind).
//
// Этот файл реализует систему логирования SDK-уровня, служащую мостом
// между Go-ядром и платформенным слоем (Kotlin/Swift).
//
// Модуль logging отвечает за:
//   - регистрацию внешнего лог-приёмника (LogSink),
//   - фильтрацию по уровню важности,
//   - отправку логов из Go в мобильное приложение,
//   - удобные шорткаты logD / logI / logW / logE.
//
// В отличие от стандартных логгеров Go, данный модуль не пишет в stdout —
// он полностью делегирует вывод в платформенный слой через интерфейс LogSink.
// Это необходимо, чтобы интегрировать SDK с системным логом Android (Logcat)
// или Swift Console на iOS.
package mobile

// LogSink — интерфейс, который должен реализовать принимающая сторона (Kotlin/Swift).
// Он получает сообщения из Go-ядра через gomobile.
//
// Пример Kotlin-реализации:
//
//	class Hy2Logger : Hy2core.LogSink {
//	    override fun Log(level: String, msg: String) {
//	        Log.i("HY2", "[$level] $msg")
//	    }
//	}
//
//	Hy2core.SetLogger(Hy2Logger())
//
// После этого вызовы logI(), logE() и т.д. из Go будут появляться в Logcat.
//
// Аргументы метода:
//   - level — текстовый уровень ("debug", "info", "warn", "error")
//   - msg — сообщение
type LogSink interface{ Log(level, msg string) }

var (
	sink     LogSink
	logLevel = "info" // Текущий уровень фильтрации логов
	order    = map[string]int{
		"debug": 10,
		"info":  20,
		"warn":  30,
		"error": 40,
	}
)

// SetLogger регистрирует внешний лог-приёмник (Kotlin/Swift).
//
// После вызова этой функции все внутренние вызовы log*()
// будут перенаправлены в sink.Log(level, msg).
func SetLogger(s LogSink) { sink = s }

// SetLogLevel устанавливает текущий уровень фильтрации логов.
//
// Допустимые уровни: "debug", "info", "warn", "error".
// Все сообщения ниже текущего уровня будут игнорироваться.
//
// Пример:
//
//	SetLogLevel("debug") // включает все логи
func SetLogLevel(level string) {
	if _, ok := order[level]; ok {
		logLevel = level
	}
}

// log — базовая внутренняя функция отправки логов.
//
// Она проверяет наличие зарегистрированного sink и уровень фильтрации.
// Если фильтр пропускает сообщение — оно передаётся в LogSink.Log().
func log(level, msg string) {
	if sink == nil {
		return
	}
	if order[level] < order[logLevel] {
		return
	}
	sink.Log(level, msg)
}

// logD отправляет лог уровня DEBUG (диагностика, отладка).
func logD(m string) { log("debug", m) }

// logI отправляет лог уровня INFO (стандартные сообщения).
func logI(m string) { log("info", m) }

// logW отправляет лог уровня WARN (предупреждения).
func logW(m string) { log("warn", m) }

// logE отправляет лог уровня ERROR (ошибки, сбои, panic-guard).
func logE(m string) { log("error", m) }
