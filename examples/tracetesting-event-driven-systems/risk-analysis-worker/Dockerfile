FROM eclipse-temurin:17-jdk AS builder

WORKDIR /usr/src/app/

COPY ./gradlew* ./settings.gradle* .
COPY ./gradle ./gradle
COPY ./app ./app

RUN ./gradlew

ENTRYPOINT [ "./gradlew", "run" ]
