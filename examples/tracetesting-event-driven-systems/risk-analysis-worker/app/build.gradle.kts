plugins {
    application
    id("java")
}

repositories {
    // Use Maven Central for resolving dependencies.
    mavenCentral()
}

dependencies {
    // This dependency is used by the application.
    implementation("com.google.guava:guava:32.1.1-jre")

    // OpenTelemetry core
    implementation(platform("io.opentelemetry:opentelemetry-bom:1.31.0"))
    implementation("io.opentelemetry:opentelemetry-api")
    implementation("io.opentelemetry:opentelemetry-sdk")
    implementation("io.opentelemetry:opentelemetry-exporter-otlp")

    // OpenTelemetry instrumentation
    implementation(platform("io.opentelemetry.instrumentation:opentelemetry-instrumentation-bom-alpha:1.31.0-alpha"))
    implementation("io.opentelemetry.instrumentation:opentelemetry-runtime-telemetry-java8")
    implementation("io.opentelemetry.instrumentation:opentelemetry-log4j-appender-2.17")
    implementation("io.opentelemetry.instrumentation:opentelemetry-kafka-clients-2.6")

    // Kafka
    implementation("org.apache.kafka:kafka-clients:3.0.0")

    // Jackson (JSON Parser)
    implementation("com.fasterxml.jackson.core:jackson-databind:2.8.9")
}

java {
    toolchain {
        languageVersion.set(JavaLanguageVersion.of(17))
    }
}

application {
    mainClass.set("risk.analysis.worker.App")
}
