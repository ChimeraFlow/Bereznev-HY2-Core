package tech.bereznev.hy2.sdk

/**
 * 🧩 Kotlin-first expect-слой для Bereznev HY2 SDK.
 *
 * Это multiplatform-фасад, определяющий контракт SDK без привязки к платформе.
 * Реальные реализации (`actual`) находятся в:
 * - Android: `Hy2Android.kt` (через gomobile AAR)
 * - iOS: `Hy2Ios.kt` (через XCFramework)
 *
 * Идея:
 * - наружу отдаётся удобный Kotlin API;
 * - внутри — маршрутизация к нативному HY2-ядру (Go).
 *
 * Потоки:
 * - все публичные вызовы потокобезопасны;
 * - обратные вызовы логов и событий приходят из фоновых потоков.
 */
expect object Hy2 {

    /**
     * Запускает HY2-движок с заданной конфигурацией.
     *
     * - Инициализирует ядро и внутренние сервисы (tun, outbounds, маршрутизацию).
     * - Если ядро уже запущено, возвращает ошибку "already running".
     *
     * @param configJson JSON-строка с полной конфигурацией HY2
     * (пример в README: `"inbounds"`, `"outbounds"`, `"route"`).
     * @return `null`, если всё успешно, либо текст ошибки.
     */
    fun start(configJson: String): String?

    /**
     * Горячая перезагрузка конфигурации без полной остановки ядра.
     *
     * Позволяет обновить настройки маршрутизации, серверы или ключи,
     * не обрывая активные соединения.
     *
     * @param configJson Новый JSON-конфиг.
     * @return `null`, если успешно, иначе сообщение об ошибке.
     */
    fun reload(configJson: String): String?

    /**
     * Останавливает HY2-движок.
     *
     * Безопасно вызывать многократно — вызов идемпотентен.
     */
    fun stop()

    /**
     * Текущий статус ядра.
     *
     * Возможные значения:
     * - `"running"`
     * - `"stopped"`
     *
     * @return строка со статусом.
     */
    fun status(): String

    /**
     * Версия SDK/движка.
     *
     * Например: `"Bereznev-HY2-Core 0.1.0 (skeleton)"`.
     *
     * Используется для диагностики и отображения версии в UI.
     */
    fun version(): String

    /**
     * Возвращает JSON-строку со сводкой здоровья ядра.
     *
     * Минимальные поля:
     * ```
     * {
     *   "running": true,
     *   "engine": "hy2core",
     *   "version": "0.1.0"
     * }
     * ```
     *
     * В будущем может включать сетевые метрики (RTT, bytesIn/out, reconnects).
     */
    fun healthJson(): String

    /**
     * Устанавливает уровень логирования.
     *
     * Допустимые уровни:
     * - `"debug"`
     * - `"info"`
     * - `"warn"`
     * - `"error"`
     *
     * Применяется немедленно, без перезапуска ядра.
     *
     * @param level уровень логов.
     */
    fun setLogLevel(level: String)

    /**
     * Назначает обработчик событий ядра.
     *
     * События:
     * - `"started"`, `"stopped"`, `"reloaded"`
     * - `"panic"`, `"error"`, `"warning"`
     * - `"metrics"` (в будущем)
     *
     * При передаче `null` — отписывает обработчик.
     *
     * @param sink реализация [Hy2EventSink] или `null`.
     */
    fun setEventSink(sink: Hy2EventSink?)

    /**
     * Назначает обработчик логов.
     *
     * Колбэки приходят из фонового потока, если требуется — перенаправляйте в UI вручную.
     * При передаче `null` — отписывает обработчик.
     *
     * @param sink реализация [Hy2LogSink] или `null`.
     */
    fun setLogSink(sink: Hy2LogSink?)
}

/* === Справка ===
 *
 * Интерфейсы, ожидаемые expect-объявлением:
 *
 * interface Hy2LogSink {
 *     fun log(level: String, msg: String)
 * }
 *
 * interface Hy2EventSink {
 *     fun onEvent(name: String, data: String)
 * }
 *
 * Пример использования:
 *
 * Hy2.setLogSink(object : Hy2LogSink {
 *     override fun log(level: String, msg: String) {
 *         println("[$level] $msg")
 *     }
 * })
 *
 * Hy2.start(configJson)
 */