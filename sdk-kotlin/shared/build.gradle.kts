plugins {
    kotlin("multiplatform")
    id("com.android.library")
    id("maven-publish")
    id("signing")
}

group = providers.gradleProperty("GROUP").get()
version = providers.gradleProperty("VERSION_NAME").get()

kotlin {
    androidTarget {
        publishLibraryVariants("release") // чтобы публиковался Android-артефакт
    }
    iosArm64()
    iosSimulatorArm64()
    iosX64()

    sourceSets {
        val commonMain by getting
        val commonTest by getting

        val androidMain by getting {
            // Локальная разработка: тянем AAR из dist/, если он есть.
            // На прод-публикациях — заведи maven-зависимость (см. ниже).
            dependencies {
                val localAar = rootProject.layout.projectDirectory
                    .dir("dist/android")
                    .file("hy2core.aar")
                    .asFile

                if (localAar.exists()) {
                    // локально
                    implementation(files(localAar))
                } else {
                    // из Maven (когда опубликуешь core-go AAR)
                    api("io.bereznev:hy2core-go-android:${project.version}")
                }
            }
        }

        val iosMain by getting {
            // когда будет готов XCFramework и .def — добавишь cinterop/бинарь.
            // Сейчас заглушка на уровне actual-реализации.
        }
    }
}

android {
    namespace = "tech.bereznev.hy2.sdk"
    compileSdk = providers.gradleProperty("android.compileSdk").get().toInt()

    defaultConfig {
        minSdk = providers.gradleProperty("android.minSdk").get().toInt()
        consumerProguardFiles("proguard-consumer.pro")
    }

    buildTypes {
        release {
            isMinifyEnabled = false
        }
    }

    // чтобы Gradle нашёл локальный AAR из dist/ при файловой зависимости
    // (если понадобятся еще локальные артефакты — раскомментируй и добавь flatDir)
    // repositories {
    //   flatDir { dirs(rootProject.layout.projectDirectory.dir("dist/android")) }
    // }
}

