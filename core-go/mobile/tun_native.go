//go:build android || ios

// Package mobile — мобильный слой SDK (gomobile bind).
//
// Этот файл описывает прямую интеграцию с TUN-интерфейсом через sing-tun (вариант B архитектуры).
//
// В отличие от socks_bridge.go (вариант A, где трафик идёт через локальный SOCKS),
// здесь мы используем "чистый" стек — прямое подключение TUN ↔ HY2 без промежуточных прокси.
//
// Это даёт более низкие задержки, меньше накладных расходов и лучшее управление MTU,
// но требует полноценной инициализации sing-tun.Options, а также контроля Offload/Protect(fd).
//
// На Android это выполняется в рамках VPNService, на iOS — в NetworkExtension.
package mobile

// StartWithTun запускает HY2-ядро с уже установленным TUN-интерфейсом.
//
// Этот метод используется, если приложение поднимает собственный VPNService и
// передаёт файловый дескриптор (tunFd) напрямую в SDK, минуя локальный SOCKS.
//
// Аргументы:
//   - tunFd — файловый дескриптор системного TUN-интерфейса;
//   - mtu   — желаемое значение MTU (обычно 1500 или 1400 для мобильных сетей).
//
// Возвращает:
//   - "" (пустая строка) — успешный запуск или stub на этапе skeleton;
//   - иначе — текст ошибки (для отображения в UI или логах).
//
// Пример (Kotlin):
//
//	val code = Hy2core.StartWithTun(tunFd, 1500)
//	if (code.isNotEmpty()) Log.e("HY2", "Error: $code")
//
// В прод-реализации функция будет:
//  1. Инициализировать sing-tun.Options с нужным MTU, ProtectFn, и offload-настройками.
//  2. Подключать обработку пакетов, создавать адаптеры потоков.
//  3. Управлять shutdown через Stop() и emit("stopped").
//  4. Отправлять статистику (bytes_in/out, reconnects) в HealthJSON().
//
// TODO:
//   - Реализовать singtun.New(options) из github.com/sagernet/sing-tun
//   - Передавать Protect(fd) через NetHooks (Android VPNService)
//   - Добавить обработку ошибок с emit(EvtError, ...)
//   - Интегрировать метрики и graceful shutdown
func StartWithTun(tunFd int, mtu int) string {
	logI("StartWithTun() not implemented yet")
	return ""
}

// SetMTU изменяет MTU активного TUN-интерфейса.
//
// Используется для динамической подстройки размера пакета в runtime,
// если обнаруживаются фрагментации или сетевые деградации.
//
// Аргументы:
//   - mtu — новое значение MTU.
//
// TODO:
//   - Реализовать реальное обновление параметра через sing-tun (SetMTU())
//   - Обновлять HealthJSON() при изменении
func SetMTU(mtu int) { logI("SetMTU not implemented yet") }
