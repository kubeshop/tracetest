# Step-by-step

1. Add Dependencies

    ```bash
    dotnet add package OpenTelemetry
    dotnet add package OpenTelemetry.Extensions.Hosting
    dotnet add package OpenTelemetry.Exporter.OpenTelemetryProtocol
    dotnet add package OpenTelemetry.Instrumentation.AspNetCore
    dotnet add package OpenTelemetry.Instrumentation.Http
    ```

2. Configure Tracing

    ```cs
    // Import OpenTelemetry SDK
    using OpenTelemetry.Trace;
    var builder = WebApplication.CreateBuilder(args);
    builder.Services.AddControllers();
    // Configure OpenTelemetry Tracing
    builder.Services.AddOpenTelemetry().WithTracing(builder =>
    {
      builder
          // Configure ASP.NET Core Instrumentation
          .AddAspNetCoreInstrumentation()
          // Configure HTTP Client Instrumentation
          .AddHttpClientInstrumentation()
          // Configure OpenTelemetry Protocol (OTLP) Exporter
          .AddOtlpExporter();
    });
    ```

3. Start the App

    ```bash
    export OTEL_SERVICE_NAME=my-service-name
    export OTEL_EXPORTER_OTLP_PROTOCOL=http/protobuf
    export OTEL_EXPORTER_OTLP_ENDPOINT="http://localhost:4318"
    export OTEL_EXPORTER_OTLP_HEADERS="x-tracetest-token=bla"

    dotnet run
    ```
