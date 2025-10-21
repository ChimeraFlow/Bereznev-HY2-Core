// core-go/pkg/logging/logging.go
package protect

// НИКАКИХ build tags сверху — логгер чисто на Go.

// LogSink — внешний приёмник логов (Kotlin/Swift).
type LogSink interface{ Log(level, msg string) }

var (
	Sink     LogSink
	LogLevel = "info"
	Order    = map[string]int{"debug": 10, "info": 20, "warn": 30, "error": 40}
)

func SetLogger(s LogSink) { Sink = s }
func SetLogLevel(level string) {
	if _, ok := Order[level]; ok {
		LogLevel = level
	}
}

func Log(level, msg string) {
	if Sink == nil {
		return
	}
	if Order[level] < Order[LogLevel] {
		return
	}
	Sink.Log(level, msg)
}

// Экспортируемые функции — вызывай их из других пакетов:
func Debug(m string) { Log("debug", m) }
func Info(m string)  { Log("info", m) }
func Warn(m string)  { Log("warn", m) }
func Error(m string) { Log("error", m) }

// (опционально) оставь алиасы, если где-то внутри самого пакета ссылались на них:
func LogD(m string) { Debug(m) }
func LogI(m string) { Info(m) }
func LogW(m string) { Warn(m) }
func LogE(m string) { Error(m) }
