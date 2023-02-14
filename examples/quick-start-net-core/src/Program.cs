using OpenTelemetry.Exporter;
using OpenTelemetry.Resources;
using OpenTelemetry.Trace;

var serviceName = "tracetest-example";
var serviceVersion = "1.0.0";

var builder = WebApplication.CreateBuilder(args);

AppContext.SetSwitch("System.Net.Http.SocketsHttpHandler.Http2UnencryptedSupport", true);

// Add services to the container.

builder.Services.AddControllers();

var appResourceBuilder = ResourceBuilder.CreateDefault()
    .AddService(serviceName: serviceName, serviceVersion: serviceVersion);

builder.Services.AddOpenTelemetryTracing(tracerProviderBuilder =>
{
  tracerProviderBuilder
      .AddConsoleExporter()
      .AddOtlpExporter(opt =>
        {
          opt.BatchExportProcessorOptions = new OpenTelemetry.BatchExportProcessorOptions<System.Diagnostics.Activity> { MaxQueueSize = 1000, ExporterTimeoutMilliseconds = 100, MaxExportBatchSize = 1 };

          // Using Grpc Protocol
          opt.Protocol = OtlpExportProtocol.Grpc;
          opt.Endpoint = new Uri("http://otel-collector:4317");

          // Using Http Protocol
          // opt.Protocol = OtlpExportProtocol.HttpProtobuf;
          // opt.Endpoint = new Uri("http://otel-collector:4318/v1/traces");
        })
      .AddSource(serviceName)
      .SetResourceBuilder(appResourceBuilder)
      .AddHttpClientInstrumentation()
      .AddAspNetCoreInstrumentation();
});

var app = builder.Build();

app.UseHttpsRedirection();
app.UseAuthorization();
app.MapControllers();

app.Run();
