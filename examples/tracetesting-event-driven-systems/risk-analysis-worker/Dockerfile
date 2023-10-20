FROM eclipse-temurin:17-jdk AS builder

WORKDIR /usr/src/app/

COPY ./gradlew* ./settings.gradle* .
COPY ./gradle ./gradle
COPY ./app ./app

RUN ./gradlew
RUN ./gradlew assemble

FROM eclipse-temurin:17-jre

WORKDIR /usr/src/app/

COPY --from=builder /usr/src/app/app/build/distributions/app.tar ./app.tar
RUN tar -xvf ./app.tar
RUN rm ./app.tar
RUN mv ./app ./risk-analysis-worker

ADD https://github.com/open-telemetry/opentelemetry-java-instrumentation/releases/download/v1.31.0/opentelemetry-javaagent.jar /usr/src/app/opentelemetry-javaagent.jar
RUN chmod 644 /usr/src/app/opentelemetry-javaagent.jar
ENV JAVA_TOOL_OPTIONS=-javaagent:/usr/src/app/opentelemetry-javaagent.jar

ENTRYPOINT [ "./risk-analysis-worker/bin/app" ]
