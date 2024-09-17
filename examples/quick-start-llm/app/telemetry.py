from opentelemetry import trace
from opentelemetry.exporter.otlp.proto.grpc.trace_exporter import OTLPSpanExporter

from traceloop.sdk import Traceloop

# import openlit
import os

otlp_endpoint = os.getenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT", "localhost:4317")
otlp_service_name = os.getenv("OTEL_SERVICE_NAME", "quick-start-llm")

def init():
  tracer = trace.get_tracer(otlp_service_name)

  Traceloop.init(
    exporter=OTLPSpanExporter(endpoint=otlp_endpoint, insecure=True),
    disable_batch=True,
    should_enrich_metrics=True
  )

  return tracer

# Test method to guarantee that the telemetry is working
def heartbeat(tracer):
  with tracer.start_as_current_span("heartbeat"):
    current_span = trace.get_current_span()
    current_span.set_attribute("hello", "world")
