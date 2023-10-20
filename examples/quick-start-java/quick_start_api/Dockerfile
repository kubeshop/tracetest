FROM eclipse-temurin:17-jdk AS builder

WORKDIR /usr/src/app/

COPY ./gradlew* ./settings.gradle* .
COPY ./gradle ./gradle
COPY ./app ./app

RUN ./gradlew
RUN ./gradlew assemble

FROM eclipse-temurin:17-jre

WORKDIR /usr/src/app/

COPY --from=builder /usr/src/app/ ./

ENTRYPOINT [ "java", "-jar", "./app/build/libs/app.jar" ]
