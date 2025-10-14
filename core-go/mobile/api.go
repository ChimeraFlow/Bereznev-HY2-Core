//go:build android || ios || mobile_skel

// Package mobile — публичная прослойка для мобильных платформ (Android/iOS),
// которая будет биндиться через gomobile в AAR/XCFramework.
// Здесь живёт стабильный API жизненного цикла HY2-движка для вызова из
// Kotlin/Java/Swift. Текущая реализация — skeleton (без сетевого стека):
//   - валидация и хранение конфигурации,
//   - управление состоянием (start/stop/reload),
//   - события (started/stopped/reloaded),
//   - базовое логирование и health (в других файлах пакета).
//
// Потокобезопасность: все экспортируемые функции потокобезопасны
// и защищены внутренним sync.Mutex.
package mobile

import "sync"

var (
	// mu защищает глобальное состояние (started, cfgRaw и пр.)
	mu sync.Mutex

	// started — признак запущенного ядра (skeleton).
	started bool

	// sdkName/sdVersion/engineID — метаданные SDK, видимые в Version().
	sdkName    = "Bereznev-HY2-Core"
	sdkVersion = "0.1.0"
	// engineID: в skeleton — "skeleton"; после интеграции sing/sing-tun поменяем.
	engineID = "skeleton"
)

// Коды ошибок для мобильных биндингов (удобны для Kotlin/Swift):

// Start запускает HY2-ядро с конфигурацией configJSON.
// Возвращает пустую строку при успехе, иначе — текст ошибки.
// Используется там, где удобнее строковая ошибка (например, простая интеграция в Kotlin).
//
// Побочные эффекты:
//   - сохраняет валидный конфиг,
//   - включает состояние started,
//   - эмитит событие "started".
//
// Потокобезопасно.
func Start(configJSON string) string {
	mu.Lock()
	defer mu.Unlock()
	if started {
		return ""
	}
	if err := cfgSet(configJSON); err != nil {
		return "invalid config: " + err.Error()
	}
	started = true
	logI("HY2 core started; config accepted")
	emit("started", "{}")
	return ""
}

// StartWithCode — то же, что Start(), но возвращает числовой код ошибки.
// Удобно для строго типизированных мобильных вызовов (Kotlin/Swift),
// чтобы не парсить строки.
//
// Коды возврата: ErrOK, ErrAlreadyRunning, ErrInvalidConfig, ErrEngineInitFailed (на будущее).
// Потокобезопасно.
func StartWithCode(configJSON string) ErrCode {
	mu.Lock()
	defer mu.Unlock()

	if started {
		return ErrOK // считаем idempotent запуск «не ошибкой»
	}
	if err := cfgSet(configJSON); err != nil {
		return ErrInvalidConfig
	}
	started = true
	logI("HY2 core started; config accepted")
	emit(EvtStarted, "{}")
	return ErrOK
}

// Reload безопасно применяет новый конфиг во время работы.
// Возвращает пустую строку при успехе, иначе текст ошибки валидации.
// В skeleton-режиме просто заменяет сохранённый конфиг и шлёт "reloaded".
// Потокобезопасно.
func Reload(configJSON string) string {
	mu.Lock()
	defer mu.Unlock()

	if err := cfgSet(configJSON); err != nil {
		return "invalid config: " + err.Error()
	}
	emit("reloaded", "{}")
	return ""
}

// Stop останавливает ядро (в skeleton — только снимает флаг started).
// Эмитит событие "stopped". Потокобезопасно.
func Stop() {
	mu.Lock()
	defer mu.Unlock()

	if !started {
		return
	}
	started = false
	emit("stopped", "{}")
	logI("HY2 core stopped")
}

// Status возвращает строковый статус: "running" или "stopped".
// Подходит для простого опроса состояния из мобильного кода.
// Потокобезопасно.
func Status() string {
	if IsRunning() {
		return "running"
	}
	return "stopped"
}

// IsRunning — булев вариант Status(); пригодится в JNI/Swift мостах.
// Потокобезопасно.
func IsRunning() bool {
	mu.Lock()
	defer mu.Unlock()
	return started
}

// Version возвращает человекочитаемую строку вида:
// "Bereznev-HY2-Core 0.1.0 (skeleton)".
// Полезно для логов, health-эндпоинтов и экранов «О приложении».
func Version() string { return sdkName + " " + sdkVersion + " (" + engineID + ")" }
