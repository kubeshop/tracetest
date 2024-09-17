from opentelemetry import trace
from opentelemetry.exporter.otlp.proto.grpc.trace_exporter import OTLPSpanExporter
from opentelemetry.sdk.resources import SERVICE_NAME, Resource
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor

from traceloop.sdk import Traceloop

# import openlit
import os

otlp_endpoint = os.getenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
otlp_service_name = os.getenv("OTEL_SERVICE_NAME", "quick-start-llm")

resource = Resource(attributes={
  SERVICE_NAME: otlp_service_name
})

provider = TracerProvider(resource=resource)

processor = BatchSpanProcessor(
  OTLPSpanExporter(endpoint=otlp_endpoint, insecure=True)
)
provider.add_span_processor(processor)

trace.set_tracer_provider(provider)

def init():
  tracer = trace.get_tracer(otlp_service_name)

  Traceloop.init(
    exporter=OTLPSpanExporter(endpoint=otlp_endpoint, insecure=True)
  )
  # openlit.init(
  #   environment='development',
  #   application_name=otlp_service_name,
  #   tracer=tracer,
  #   disable_metrics=True,
  #   collect_gpu_stats=False
  # )

  return tracer

def heartbeat(tracer):
  with tracer.start_as_current_span("heartbeat"):
    current_span = trace.get_current_span()
    current_span.set_attribute("hello", "world")
