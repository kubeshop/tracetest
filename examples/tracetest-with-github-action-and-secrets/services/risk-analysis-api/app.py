from flask import Flask, request, jsonify

from opentelemetry import trace
from opentelemetry.trace import Status, StatusCode
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from opentelemetry.exporter.otlp.proto.grpc.trace_exporter import OTLPSpanExporter
from opentelemetry.sdk.resources import SERVICE_NAME, Resource

from opentelemetry.instrumentation.flask import FlaskInstrumentor

import config

current_config = config.get_current()

resource = Resource(attributes={
    SERVICE_NAME: current_config.otel_service_name
})

provider = TracerProvider(resource=resource)

processor = BatchSpanProcessor(OTLPSpanExporter(endpoint=current_config.otel_exporter_grpc_url, insecure=True))
provider.add_span_processor(processor)

trace.set_tracer_provider(provider)
tracer = trace.get_tracer(__name__)

instrumentor = FlaskInstrumentor()

app = Flask(__name__)
instrumentor.instrument_app(app)

@app.route('/computeRisk', methods=['POST'])
def compute_risk():
  risk_data = request.json

  amount = float(risk_data['amount'])
  age = int(risk_data['age'])
  score = risk_calculation(amount, age)

  return jsonify({ "score": score })

@tracer.start_as_current_span("risk_calculation")
def risk_calculation(amount, age):
  current_span = trace.get_current_span()

  current_span.set_attribute("amount", amount)
  current_span.set_attribute("age", age)

  try:
    score = amount / (2 * age)
  except Exception as ex:
    current_span.set_status(Status(StatusCode.ERROR))
    current_span.record_exception(ex)
    # poor error handling on purpose
    score = -1

  return score

if __name__ == '__main__':
  print('Running on port: ' + current_config.port)
  app.run(host='0.0.0.0', port=current_config.port)
