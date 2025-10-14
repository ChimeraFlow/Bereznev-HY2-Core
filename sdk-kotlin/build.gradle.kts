plugins {
    // версии можно держать в gradle.properties (см. ниже)
    id("com.android.library") version libs.versions.agp.get() apply false
    kotlin("multiplatform")   version libs.versions.kotlin.get() apply false
}
