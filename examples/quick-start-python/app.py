from flask import Flask, request
import json

from opentelemetry import trace
from opentelemetry.sdk.resources import Resource
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from opentelemetry.sdk.trace.export import ConsoleSpanExporter

provider = TracerProvider()
processor = BatchSpanProcessor(ConsoleSpanExporter())
provider.add_span_processor(processor)
trace.set_tracer_provider(provider)
tracer = trace.get_tracer(__name__)

app = Flask(__name__)

@app.route("/manual")
def manual():
    with tracer.start_as_current_span(
        "manual", 
        attributes={ "endpoint": "/manual", "foo": "bar" } 
    ):
        return "App works with a manual instrumentation."

@app.route('/automatic')
def automatic():
    return "App works with automatic instrumentation."

@app.route("/")
def home():
    return "App works."
