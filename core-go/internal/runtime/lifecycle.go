// go:build android || ios || mobile_skel

// Package mobile — мобильный слой SDK (gomobile bind).
// Этот файл реализует утилиту safeGo — защищённый запуск горутин.
//
// safeGo используется для выполнения неблокирующих задач в фоне
// (например, пинга, обновления метрик, мониторинга соединений),
// при этом гарантируя, что паника внутри горутины не обрушит всё приложение.
//
// Концептуально safeGo = goroutine sandbox: любая panic() внутри функции
// будет перехвачена, залогирована и отправлена как событие "panic" через EventSink.
package runtime

import (
	"fmt"

	"github.com/ChimeraFlow/Bereznev-HY2-Core/core-go/internal/telemetry"
	logpkg "github.com/ChimeraFlow/Bereznev-HY2-Core/core-go/pkg/logging"
)

// safeGo запускает переданную функцию fn() в отдельной горутине
// с автоматическим перехватом panic и уведомлением через события SDK.
//
// Аргументы:
//   - fn — функция, которая выполняется в отдельной горутине.
//
// Поведение:
//   - Любая panic() внутри fn не прерывает работу ядра;
//   - Логируется через logE(...);
//   - Отправляется событие emit("panic", {"msg":"panic recovered"}).
//
// Пример использования:
//
//	safeGo(func() {
//	    doNetworkLoop()
//	})
//
// Потокобезопасность:
//   - safeGo не требует внешней синхронизации;
//   - Все вызовы logE() и emit() — потокобезопасны.
//
// Побочные эффекты:
//   - Создаёт новую горутину;
//   - В случае panic пишет лог уровня "error" и триггерит событие "panic".
func SafeGo(fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logpkg.LogE(fmt.Sprintf("panic: %v", r))
				telemetry.Emit("panic", `{"msg":"panic recovered"}`)
			}
		}()
		fn()
	}()
}
