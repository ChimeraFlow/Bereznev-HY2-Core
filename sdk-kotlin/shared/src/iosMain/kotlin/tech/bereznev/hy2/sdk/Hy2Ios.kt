package tech.bereznev.hy2.sdk

/**
 * 🧩 iOS actual-реализация [Hy2] — пока заглушка (stub).
 *
 * Этот слой предназначен для дальнейшей интеграции с XCFramework,
 * собранным из Go-ядра (`core-go/mobile`).
 *
 * Сейчас методы возвращают статические значения или no-ops —
 * чтобы KMM-проект собирался и мог вызывать базовые API без крашей.
 *
 * 📘 Назначение:
 * - предоставить expect/actual совместимость для KMM;
 * - позволить собрать iOS-приложение до того, как нативная часть будет готова.
 *
 * ⚙️ Статус: `stub / placeholder`
 * После интеграции XCFramework сюда будет подключено `interop`-взаимодействие.
 */
actual object Hy2 {

    /**
     * Запуск HY2-ядра с JSON-конфигом.
     *
     * ⚠️ Сейчас возвращает строку `"ios: not wired yet"`.
     * После подключения XCFramework будет вызывать `Hy2Core.start(configJson)`.
     */
    actual fun start(configJson: String): String? = "ios: not wired yet"

    /**
     * Горячая перезагрузка конфигурации.
     *
     * ⚠️ Возвращает `"ios: not wired yet"`.
     */
    actual fun reload(configJson: String): String? = "ios: not wired yet"

    /**
     * Останавливает ядро.
     * ⚠️ Пока no-op.
     */
    actual fun stop() {}

    /**
     * Текущий статус ядра.
     * ⚠️ Возвращает `"stopped"`.
     */
    actual fun status(): String = "stopped"

    /**
     * Версия SDK для iOS-заглушки.
     * ⚠️ Возвращает `"Bereznev-HY2-Core (iOS stub)"`.
     */
    actual fun version(): String = "Bereznev-HY2-Core (iOS stub)"

    /**
     * JSON со статусом "здоровья".
     * ⚠️ Всегда возвращает фиктивное состояние (`running=false`).
     */
    actual fun healthJson(): String =
        """{"running":false,"engine":"stub","version":"0.0.0"}"""

    /**
     * Установка уровня логов.
     * ⚠️ Пока не делает ничего.
     */
    actual fun setLogLevel(level: String) {}

    /**
     * Назначение обработчика событий.
     * ⚠️ Пока не поддерживается.
     */
    actual fun setEventSink(sink: Hy2EventSink?) {}

    /**
     * Назначение обработчика логов.
     * ⚠️ Пока не поддерживается.
     */
    actual fun setLogSink(sink: Hy2LogSink?) {}
}
