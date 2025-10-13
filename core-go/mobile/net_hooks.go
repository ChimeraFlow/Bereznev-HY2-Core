//go:build android || ios

// Package mobile — мобильный слой SDK (gomobile bind).
//
// Этот файл реализует систему сетевых хуков (NetHooks) —
// мост между Go-ядром и платформенной реализацией VPN-стека на Android/iOS.
//
// На Android это используется для вызова `VpnService.protect(fd)`,
// чтобы исключить внутренние сокеты SDK из VPN-туннеля и избежать
// рекурсивного зацикливания трафика.
//
// В общем виде NetHooks позволяют платформе брать под контроль низкоуровневые
// сетевые операции Go-ядра (например, защита файловых дескрипторов или
// применение собственных политик безопасности).
package mobile

// NetHooks — интерфейс, который реализуется на стороне платформы
// (обычно в Kotlin/Java или Swift) для перехвата сетевых операций.
//
// Применение на Android:
//
//	class Hy2Hooks : Hy2core.NetHooks {
//	    override fun Protect(fd: Int): Boolean {
//	        return vpnService.protect(fd)
//	    }
//	}
//
//	Hy2core.SetNetHooks(Hy2Hooks())
//
// После этого Go-ядро сможет вызывать protectFD(fd),
// и реальный вызов уйдёт в Java/Kotlin, где выполнится `VpnService.protect(fd)`.
//
// На iOS реализация обычно пустая (возвращает false),
// потому что системный VPN-фреймворк управляется по-другому.
type NetHooks interface {
	// Protect — вызывается из Go-кода для защиты сокета/FD от VPN-туннеля.
	//
	// Аргументы:
	//   • fd — файловый дескриптор сокета, который нужно исключить из VPN.
	//
	// Возвращает:
	//   • true — если операция защиты выполнена успешно;
	//   • false — если защита не поддерживается или произошла ошибка.
	Protect(fd int) bool
}

var netHooks NetHooks

// SetNetHooks регистрирует обработчик сетевых хуков, реализованный на платформе.
//
// Обычно вызывается сразу после инициализации SDK-движка,
// до запуска сети или VPN-сессии.
//
// Пример:
//
//	hooks := Hy2Hooks{vpnService: svc}
//	Hy2core.SetNetHooks(hooks)
func SetNetHooks(h NetHooks) { netHooks = h }

// protectFD — внутренняя вспомогательная функция для вызова NetHooks.Protect().
//
// Она безопасно проверяет, установлен ли обработчик, и вызывает Protect(fd)
// только при его наличии. Если хук не задан, функция возвращает false.
//
// Используется в нижележащем сетевом коде (tun/socks/conn).
func protectFD(fd int) bool {
	if netHooks != nil {
		return netHooks.Protect(fd)
	}
	return false
}
