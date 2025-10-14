package tech.bereznev.hy2.sdk

import tech.bereznev.hy2.Mobile // класс с биндингами из gomobile (AAR)

/**
 * Android-слой actual для KMM-SDK.
 *
 * Назначение:
 * - Проксирует вызовы из Kotlin в Go (через gomobile bind).
 * - Регистрирует адаптеры логов/событий, переводя Kotlin-интерфейсы
 *   в gomobile-интерфейсы (которые являются абстрактными классами).
 *
 * Потоки:
 * - Все публичные методы безопасно вызывать из UI-потока.
 * - Внутри ядра (Go) операции выполняются в goroutine.
 * - Колбэки логов/событий могут приходить с фоновых потоков — учитывайте это в UI.
 *
 * Жизненный цикл:
 * - Вызовите [start] один раз с валидным JSON-конфигом.
 * - Для обновления конфигурации используйте [reload].
 * - Для выключения — [stop].
 *
 * Ошибки:
 * - Методы, возвращающие String?, возвращают `null`, если ошибок нет,
 *   и строку с текстом ошибки — если она произошла.
 */
actual object Hy2 {

    /**
     * Запуск HY2-ядра с JSON-конфигом.
     *
     * Делает:
     * - Регистрирует адаптеры логов и событий (если заданы ранее через [setLogSink] / [setEventSink]).
     * - Вызывает [Mobile.start] и возвращает `null` при успехе, либо текст ошибки.
     *
     * @param configJson JSON конфигурации (см. README: inbounds/outbounds/route).
     * @return `null` если всё ок, либо текст ошибки.
     */
    actual fun start(configJson: String): String? {
        Mobile.setLogger(AndroidSinks.goLogSink)
        Mobile.setEventSink(AndroidSinks.goEventSink)
        val err = Mobile.start(configJson)
        return if (err.isNullOrBlank()) null else err
    }

    /**
     * Горячая перезагрузка конфигурации без полного останова.
     *
     * @param configJson Новый JSON-конфиг.
     * @return `null` если всё ок, либо текст ошибки.
     */
    actual fun reload(configJson: String): String? {
        val err = Mobile.reload(configJson)
        return if (err.isNullOrBlank()) null else err
    }

    /** Останов ядра. Идempotентно. */
    actual fun stop() {
        Mobile.stop()
    }

    /** Краткий статус: `"running"` или `"stopped"`. */
    actual fun status(): String = Mobile.status()

    /** Версия SDK/движка, например: `"Bereznev-HY2-Core 0.1.0 (skeleton)"`. */
    actual fun version(): String = Mobile.version()

    /**
     * JSON со сводкой здоровья/метрик.
     * Минимально включает: running, engine, version.
     * В будущем: bytesIn/bytesOut, reconnects, quicRttMs.
     */
    actual fun healthJson(): String = Mobile.healthJSON()

    /**
     * Управление уровнем логирования ядра:
     * `"debug" | "info" | "warn" | "error"`.
     */
    actual fun setLogLevel(level: String) { Mobile.setLogLevel(level) }

    /**
     * Установить получателя событий ядра.
     *
     * События:
     * - "started", "stopped", "reloaded"
     * - "panic", "error", "warning"
     * - "metrics" (в будущем)
     *
     * Важно: колбэки приходят с фонового потока — синхронизируйте UI самостоятельно.
     *
     * @param sink Реализация интерфейса приёма событий; `null` — отписаться.
     */
    actual fun setEventSink(sink: Hy2EventSink?) {
        AndroidSinks.eventSink = sink
        if (sink == null) Mobile.setEventSink(null)
        else Mobile.setEventSink(AndroidSinks.goEventSink)
    }

    /**
     * Установить приёмник логов ядра.
     *
     * Важно: лог-колбэки приходят с фонового потока — если пишете в UI/файл,
     * учитывайте синхронизацию/буферизацию.
     *
     * @param sink Реализация интерфейса логов; `null` — отписаться.
     */
    actual fun setLogSink(sink: Hy2LogSink?) {
        AndroidSinks.logSink = sink
        if (sink == null) Mobile.setLogger(null)
        else Mobile.setLogger(AndroidSinks.goLogSink)
    }
}

/**
 * Внутренние адаптеры для сопряжения Kotlin-интерфейсов с gomobile-интерфейсами.
 *
 * Почему так:
 * - gomobile генерирует абстрактные классы (например, `Mobile.LogSink`) вместо интерфейсов;
 * - нам нужно обернуть наши `Hy2LogSink` / `Hy2EventSink` в эти классы.
 */
private object AndroidSinks {

    /** Текущий Kotlin-лог-sink, который назначает приложение. */
    var logSink: Hy2LogSink? = null

    /**
     * Адаптер gomobile → Kotlin для логов.
     *
     * Mobile вызывает этот объект, мы транслируем вызов в [logSink].
     * null-защита применяется: пустые строки вместо null.
     */
    val goLogSink = object : tech.bereznev.hy2.Mobile.LogSink() {
        override fun log(level: String?, msg: String?) {
            logSink?.log(level.orEmpty(), msg.orEmpty())
        }
    }

    /** Текущий Kotlin-event-sink, который назначает приложение. */
    var eventSink: Hy2EventSink? = null

    /**
     * Адаптер gomobile → Kotlin для событий.
     *
     * Mobile вызывает этот объект, мы транслируем вызов в [eventSink].
     * null-защита применяется: пустые строки вместо null.
     */
    val goEventSink = object : tech.bereznev.hy2.Mobile.EventSink() {
        override fun onEvent(name: String?, data: String?) {
            eventSink?.onEvent(name.orEmpty(), data.orEmpty())
        }
    }
}

/* ===== Интерфейсы SDK-слоя (expect объявлены в commonMain) =====
 *
 * interface Hy2LogSink {
 *   fun log(level: String, msg: String)
 * }
 *
 * interface Hy2EventSink {
 *   fun onEvent(name: String, data: String)
 * }
 *
 * Пример использования:
 *
 * Hy2.setLogSink(object : Hy2LogSink {
 *   override fun log(level: String, msg: String) {
 *     Log.d("HY2", "[$level] $msg")
 *   }
 * })
 *
 * Hy2.setEventSink(object : Hy2EventSink {
 *   override fun onEvent(name: String, data: String) {
 *     Log.d("HY2", "event=$name data=$data")
 *   }
 * })
 *
 * val err = Hy2.start(configJson)
 * if (err != null) { /* показать ошибку */ }
 * else { /* ok */ }
 */