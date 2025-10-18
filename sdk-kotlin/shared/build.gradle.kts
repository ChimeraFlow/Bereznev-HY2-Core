plugins {
    kotlin("multiplatform")
    id("com.android.library")
    id("maven-publish")
    id("signing")
    id("org.jetbrains.dokka")
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

    // Чтобы публиковался 1 релизный вариант и прикладывались sources
    publishing {
        singleVariant("release") {
            withSourcesJar()
        }
    }

    // если используешь Java 11
    compileOptions {
        sourceCompatibility = JavaVersion.VERSION_11
        targetCompatibility = JavaVersion.VERSION_11
    }
    kotlinOptions {
        jvmTarget = "11"
    }
}


val dokkaJavadoc by tasks.getting(org.jetbrains.dokka.gradle.DokkaTask::class)
val javadocJar by tasks.registering(Jar::class) {
    dependsOn(dokkaJavadoc)
    archiveClassifier.set("javadoc")
    from(dokkaJavadoc.outputDirectory)
}

publishing {
    publications {
        create<MavenPublication>("release") {
            from(components["release"])
            groupId = property("GROUP").toString()
            artifactId = property("POM_ARTIFACT_ID").toString()
            version = property("VERSION_NAME").toString()
            artifact(javadocJar.get())

            pom {
                name.set(property("POM_NAME").toString())
                description.set(property("POM_DESCRIPTION").toString())
                url.set(property("POM_URL").toString())
                licenses {
                    license {
                        name.set(property("POM_LICENSE_NAME").toString())
                        url.set(property("POM_LICENSE_URL").toString())
                        distribution.set(property("POM_LICENSE_DIST").toString())
                    }
                }
                scm {
                    url.set(property("POM_SCM_URL").toString())
                    connection.set(property("POM_SCM_CONNECTION").toString())
                    developerConnection.set(property("POM_SCM_DEV_CONNECTION").toString())
                }
                developers {
                    developer {
                        id.set(property("POM_DEVELOPER_ID").toString())
                        name.set(property("POM_DEVELOPER_NAME").toString())
                        url.set(property("POM_DEVELOPER_URL").toString())
                    }
                }
            }
        }
    }
    repositories {
        maven {
            name = "sonatype"
            val releases = uri("https://s01.oss.sonatype.org/service/local/staging/deploy/maven2/")
            val snapshots = uri("https://s01.oss.sonatype.org/content/repositories/snapshots/")
            url = if (version.toString().endsWith("SNAPSHOT")) snapshots else releases
            credentials {
                username = System.getenv("OSSRH_USERNAME") ?: findProperty("OSSRH_USERNAME")?.toString()
                password = System.getenv("OSSRH_PASSWORD") ?: findProperty("OSSRH_PASSWORD")?.toString()
            }
        }
    }
}

signing {
    val key = System.getenv("SIGNING_KEY") ?: findProperty("SIGNING_KEY")?.toString()
    val pass = System.getenv("SIGNING_PASSWORD") ?: findProperty("SIGNING_PASSWORD")?.toString()
    useInMemoryPgpKeys(key, pass)
    sign(publishing.publications)
}
