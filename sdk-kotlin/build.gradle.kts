plugins {
    // версии можно держать в gradle.properties (см. ниже)
    id("com.android.library") version libs.versions.agp.get() apply false
    kotlin("multiplatform")   version libs.versions.kotlin.get() apply false
    id("org.jetbrains.dokka") version "1.9.20" apply false

}


allprojects {
    group = property("GROUP")!!
    version = property("VERSION_NAME")!!
}
