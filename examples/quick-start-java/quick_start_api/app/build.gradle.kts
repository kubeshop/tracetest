plugins {
  id("java")
  id("org.springframework.boot") version "3.0.6"
  id("io.spring.dependency-management") version "1.1.0"
}

repositories {
  mavenCentral()
}

dependencies {
  implementation("org.springframework.boot:spring-boot-starter-web")

  // OpenTelemetry core
  implementation(platform("io.opentelemetry:opentelemetry-bom:1.31.0"))
  implementation("io.opentelemetry:opentelemetry-api")
  implementation("io.opentelemetry:opentelemetry-sdk")
  implementation("io.opentelemetry:opentelemetry-exporter-otlp")
  implementation("io.opentelemetry:opentelemetry-sdk-extension-autoconfigure")
  implementation("io.opentelemetry:opentelemetry-sdk-extension-autoconfigure-spi")

  // OpenTelemetry instrumentation
  implementation(platform("io.opentelemetry.instrumentation:opentelemetry-instrumentation-bom-alpha:1.31.0-alpha"))
  implementation("io.opentelemetry.instrumentation:opentelemetry-runtime-telemetry-java8")
  implementation("io.opentelemetry.instrumentation:opentelemetry-log4j-appender-2.17")
  implementation("io.opentelemetry.instrumentation:opentelemetry-spring-webmvc-6.0")

  // implementation(platform("io.opentelemetry:opentelemetry-bom:1.31.0"))
  // implementation("io.opentelemetry:opentelemetry-api:1.31.0")
  // implementation("io.opentelemetry:opentelemetry-sdk:1.31.0")
  // implementation("io.opentelemetry:opentelemetry-sdk-metrics:1.31.0")
  // implementation("io.opentelemetry.semconv:opentelemetry-semconv:1.21.0-alpha")
  // implementation("io.opentelemetry:opentelemetry-exporter-otlp:1.31.0")
  // implementation("io.opentelemetry:opentelemetry-sdk-extension-autoconfigure:1.31.0")
  // implementation("io.opentelemetry:opentelemetry-sdk-extension-autoconfigure-spi:1.31.0")

  // implementation("io.opentelemetry.instrumentation:opentelemetry-spring-webmvc-6.0")
}

java {
  toolchain {
    languageVersion.set(JavaLanguageVersion.of(17))
  }
}

sourceSets {
  main {
    java.setSrcDirs(setOf("."))
  }
}

ext {
    set("opentelemetry.version", "1.31.0")
}
