# Bereznev HY2 Core (gomobile AAR) — R8/ProGuard rules

# 1) Сгенерированные gomobile классы и мост Seq
-keep class go.** { *; }
-dontwarn go.**

# На старых gomobile встречается org.golang.**
-keep class org.golang.** { *; }
-dontwarn org.golang.**

# 2) Java-обёртка из AAR (пакет, который ты задаёшь флагом -javapkg)
-keep class tech.bereznev.hy2.** { *; }

# 3) Не ругаться на опциональные зависимости/неймспейсы из sing-стека
-dontwarn **.sing.**
-dontwarn **.quic.**
-dontwarn **.netlink.**
-dontwarn **.nftables.**
-dontwarn **.gvisor.**

# 4) Аннотации/инлайны (подстраховка для Android Gradle Plugin/R8)
-keepattributes *Annotation*, InnerClasses, EnclosingMethod, Signature

# 5) Если используешь Kotlin корутины/flow в host-приложении, R8 их часто трогает
# (правила ниже — мягкие, на случай, если приложение их требует)
-dontwarn kotlinx.coroutines.**
-keep class kotlinx.coroutines.** { *; }

# 6) Если приложение использует reflection на пакет tech.bereznev.hy2.*
# (обычно не нужно, но пусть будет)
-keepnames class tech.bereznev.hy2.** { *; }