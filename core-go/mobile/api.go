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
	sdkName = "Bereznev-HY2-Core"
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

	// идемпотентность: повторный Start не считается ошибкой
	if started {
		return ""
	}

	// 1) валидируем и сохраняем конфиг
	if err := cfgSet(configJSON); err != nil {
		return "invalid config: " + err.Error()
	}

	// 2) поднимаем рантайм (sing/hy2 транспорт и пр.)
	if err := runtimeStart(); err != nil {
		emitError(int(ErrEngineInitFailed), err.Error())
		return "engine init failed: " + err.Error()
	}

	started = true
	healthMarkStarted()
	logI("HY2 core started")
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
	if err := runtimeStart(); err != nil {
		emitError(int(ErrEngineInitFailed), err.Error())
		return ErrEngineInitFailed
	}
	started = true
	logI("HY2 core started; config accepted")
	return ErrOK
}

// Reload безопасно применяет новый конфиг во время работы.
// Возвращает пустую строку при успехе, иначе текст ошибки валидации.
// В skeleton-режиме просто заменяет сохранённый конфиг и шлёт "reloaded".
// Потокобезопасно.
func Reload(configJSON string) string {
	mu.Lock()
	defer mu.Unlock()

	// 1) валидируем и сохраняем конфиг (если невалиден — ничего не меняем)
	if err := cfgSet(configJSON); err != nil {
		return "invalid config: " + err.Error()
	}
	// 2) если не запущено — просто отдадим событие перезагрузки конфигурации
	if !started {
		emit(EvtReloaded, "{}")
		return ""
	}

	// 3) перезапускаем рантайм по новому конфигу
	// 3) если новый конфиг не задаёт HY2 — не трогаем транспорт
	if _, err := parseHY2Config(); err != nil {
		emit(EvtReloaded, "{}")
		logI("config reloaded (no HY2 changes)")
		return ""
	}
	// 4) иначе — перезапускаем рантайм
	runtimeStop()
	if err := runtimeStart(); err != nil {
		emitError(int(ErrEngineInitFailed), err.Error())
		started = false
		return "engine init failed: " + err.Error()
	}
	started = true
	emit(EvtReloaded, "{}")
	logI("HY2 core reloaded")
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
	runtimeStop()
	started = false
	healthMarkStopped()
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
