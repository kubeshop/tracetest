from dotenv import load_dotenv

# Load environment variables from .env file
load_dotenv()

# Initialize telemetry
from telemetry import init as telemetry_init
tracer = telemetry_init() # run telemetry.init() before loading any other modules to capture any module-level telemetry

from opentelemetry import trace
from opentelemetry.instrumentation.flask import FlaskInstrumentor

# from telemetry import heartbeat as telemetry_heartbeat
# telemetry_heartbeat(tracer)

from llm.providers import get_provider, get_providers
from flask import Flask, request, jsonify, make_response

instrumentor = FlaskInstrumentor()

app = Flask(__name__)
instrumentor.instrument_app(app)

api_port = '8800'

@app.route('/summarizeText', methods=['POST'])
def summarize_text():
  data = request.json

  provider_type = data['provider']

  providers = get_providers()
  has_provider = provider_type in providers

  if not has_provider:
    return make_response(jsonify({ "error": "Invalid provider" }), 400)

  source_text = data['text']

  provider = get_provider(provider_type)
  summarize_text =  provider.summarize(source_text)

  # Get trace ID from current span
  span = trace.get_current_span()
  trace_id = span.get_span_context().trace_id

  # Convert trace_id to a hex string
  trace_id_hex = format(trace_id, '032x')

  return jsonify({"summary": summarize_text, "trace_id": trace_id_hex})

if __name__ == '__main__':
  print('Running on port: ' + api_port)
  app.run(host='0.0.0.0', port=api_port)
