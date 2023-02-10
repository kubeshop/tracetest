using OpenTelemetry.Exporter;
using OpenTelemetry.Resources;
using OpenTelemetry.Trace;

var serviceName = "tracetest-example";
var serviceVersion = "1.0.0";

var builder = WebApplication.CreateBuilder(args);

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
          opt.Protocol = OtlpExportProtocol.Grpc;
          opt.Endpoint = new Uri("http://otel-collector:4317");
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
