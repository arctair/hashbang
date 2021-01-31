import org.jetbrains.kotlin.gradle.tasks.KotlinCompile

plugins {
  application
  kotlin("jvm") version "1.4.21"
  kotlin("plugin.serialization") version "1.4.21"
  id("com.github.johnrengelman.shadow") version "4.0.4"
  id("com.palantir.docker") version "0.25.0"
  id("com.painkillergis.stepper_client.stepperClient") version "1.0.18"
}

group = "com.arctair"

version =
  ProcessBuilder("sh", "-c", "git rev-list --count HEAD")
    .start()
    .apply { waitFor() }
    .let { it.inputStream.bufferedReader().readText().trim() }
    .let { "v1.0.$it" }

val sha =
  ProcessBuilder("sh", "-c", "git rev-parse HEAD")
    .start()
    .apply { waitFor() }
    .let { it.inputStream.bufferedReader().readText().trim() }

application {
  val safeName = rootProject.name.replace("-", "_")
  mainClassName = "$group.$safeName.ApplicationKt"
}

repositories {
  jcenter()
  mavenCentral()
  maven { url = uri("https://dl.bintray.com/kotlin/kotlinx") }
  maven { url = uri("https://dl.bintray.com/kotlin/ktor") }
}

dependencies {
  implementation("io.ktor:ktor-html-builder:+")
  implementation("io.ktor:ktor-jackson:+")
  implementation("io.ktor:ktor-server-netty:+")
  implementation("org.jetbrains.kotlin:kotlin-reflect:+")
  implementation("org.jetbrains.kotlin:kotlin-stdlib:+")
  implementation("org.jetbrains.kotlin:kotlin-stdlib:+")
  implementation("org.jetbrains.kotlinx:kotlinx-serialization-json:+")
  implementation("org.slf4j:slf4j-simple:+")
  testImplementation("io.ktor:ktor-client-apache:+")
  testImplementation("io.ktor:ktor-client-core:+")
  testImplementation("io.ktor:ktor-client-jackson:+")
  testImplementation("io.ktor:ktor-server-test-host:+")
  testImplementation("io.mockk:mockk:+")
  testImplementation("org.hamcrest:hamcrest-junit:+")
  testImplementation("org.jetbrains.kotlinx:kotlinx-coroutines-core:+")
  testImplementation("org.junit.jupiter:junit-jupiter-api:+")
  testImplementation(kotlin("test-junit5"))
  testRuntimeOnly("org.junit.jupiter:junit-jupiter-engine:+")
}

tasks.test {
  useJUnitPlatform()
  if (System.getenv("ktor_starter_baseUrl") != null) {
    include("**/bspec/")
  }
}

tasks.withType<KotlinCompile>() {
  kotlinOptions.jvmTarget = "11"
}

tasks.named("processResources") {
  doFirst {
    file("src/main/resources/version.properties")
      .writeText(listOf("sha=$sha", "version=$version").joinToString("\n"))
  }
}

configurations.all {
  resolutionStrategy {
    activateDependencyLocking()
    componentSelection
      .all(object : Action<ComponentSelection> {
        @Mutate
        override fun execute(selection: ComponentSelection) {
          val version = selection.candidate.version
          when {
            version.matches(
              Regex(
                ".*[-.]rc\\d*$",
                RegexOption.IGNORE_CASE
              )
            ) -> selection.reject("Release candidates are excluded")
            version.matches(Regex(".*-M\\d+$")) -> selection.reject("Milestones are excluded")
            version.matches(Regex(".*-alpha\\d+$")) -> selection.reject("Alphas are excluded")
            version.matches(Regex(".*-beta\\d+$")) -> selection.reject("Betas are excluded")
          }
        }
      })
  }
}

docker {
  name = "arctair/${rootProject.name}:$version"
  files("build/libs/${rootProject.name}-$version.jar")
}
