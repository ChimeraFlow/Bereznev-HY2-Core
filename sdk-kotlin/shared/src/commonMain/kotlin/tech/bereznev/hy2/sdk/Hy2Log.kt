package tech.bereznev.hy2.sdk

/**
 * 🔍 Hy2Log — модель и утилиты для логирования в SDK.
 *
 * Этот слой служит для унификации логов между Kotlin-частью и Go-ядром.
 * На Android логи приходят из AAR (через gomobile.LogSink),
 * на iOS — позже через interop с XCFramework.
 *
 * Используется двумя путями:
 *  - SDK внутри вызывает [Hy2LogSink.log];
 *  - пользователь может реализовать свой sink и подписаться через [Hy2.setLogSink].
 *
 * Формат уровней полностью совпадает с Go-реализацией:
 *  "debug" | "info" | "warn" | "error".
 */
object Hy2Log {

    /** Возможные уровни логирования SDK. */
    enum class Level(val tag: String) {
        DEBUG("debug"),
        INFO("info"),
        WARN("warn"),
        ERROR("error");

        companion object {
            /** Поиск уровня по строке (безопасный). */
            fun from(tag: String?): Level = when (tag?.lowercase()) {
                "debug" -> DEBUG
                "info"  -> INFO
                "warn"  -> WARN
                "error" -> ERROR
                else    -> INFO
            }
        }
    }

    /**
     * Преобразует уровень в emoji + тэг — для наглядности в логах.
     * Пример: "⚙️ [INFO] HY2 core started"
     */
    fun format(level: Level, message: String): String = when (level) {
        Level.DEBUG -> "🧩 [DEBUG] $message"
        Level.INFO  -> "⚙️ [INFO] $message"
        Level.WARN  -> "⚠️ [WARN] $message"
        Level.ERROR -> "❌ [ERROR] $message"
    }

    /**
     * Быстрая печать в консоль (для отладки).
     * Работает независимо от зарегистрированного sink.
     */
    fun print(level: Level, message: String) {
        println(format(level, message))
    }

    /** Утилиты для коротких вызовов. */
    fun d(msg: String) = print(Level.DEBUG, msg)
    fun i(msg: String) = print(Level.INFO, msg)
    fun w(msg: String) = print(Level.WARN, msg)
    fun e(msg: String) = print(Level.ERROR, msg)
}
