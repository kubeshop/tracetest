# Step-by-step

1. Create Maven Project

      ```bash
      mvn archetype:generate -DgroupId=com.example -DartifactId=hello-world-api -DarchetypeArtifactId=maven-archetype-quickstart -DinteractiveMode=false
      ```

2. Enter `hello-world-api` Directory

    ```bash
    cd ./hello-world-api
    ```

3. Add Spark Dependency in `pom.xml`

      ```xml
      <dependency>
        <groupId>com.sparkjava</groupId>
        <artifactId>spark-core</artifactId>
        <version>2.9.4</version>
      </dependency>
      ```

4. Install Maven Dependencies and Start App

      ```bash
      mvn clean install
      mvn exec:java -Dexec.mainClass="com.example.HelloWorldApi"
      ```

5. Run with Jar

      ```bash
      mvn clean package
      ```

      ```bash
      java -jar target/hello-world-api-1.0-SNAPSHOT.jar
      ```

6. Download OpenTelemetry Java Agent

    ```bash
    curl -L -O https://github.com/open-telemetry/opentelemetry-java-instrumentation/releases/latest/download/opentelemetry-javaagent.jar
    ```

7. Run with Jar and Include the OpenTelemetry Java Agent

    ```bash
    mvn clean package
    ```

    ```bash
    export OTEL_SERVICE_NAME=my-service-name
    export OTEL_EXPORTER_OTLP_PROTOCOL=http/protobuf
    export OTEL_EXPORTER_OTLP_ENDPOINT="http://localhost:4318"
    export OTEL_EXPORTER_OTLP_HEADERS="x-token=bla"

    java -javaagent:opentelemetry-javaagent.jar -jar target/hello-world-api-1.0-SNAPSHOT.jar
    ```
