FROM python:3.10.1-slim

WORKDIR /opt/app
COPY . .

RUN pip install --no-cache-dir -r requirements.txt
RUN opentelemetry-bootstrap -a install

EXPOSE 8080

# CMD ["python3", "-m", "flask", "run", "--host=0.0.0.0", "--port=8080"]

# CMD ["opentelemetry-instrument --traces_exporter console", "flask", "run"]

# CMD ["opentelemetry-instrument", "--traces_exporter otlp", "--service_name app", "--exporter_otlp_endpoint otel-collector:4317", "python3", "-m", "flask", "run", "--host=0.0.0.0", "--port=8080"]

# CMD ["opentelemetry-instrument", "--traces_exporter otlp", "--service_name app", "--exporter_otlp_endpoint otel-collector:4317", "python", "app.py"]

# CMD ["start_app.sh"]

# CMD ["OTEL_SERVICE_NAME=app", "OTEL_TRACES_EXPORTER=otlp", "OTEL_METRICS_EXPORTER=none", "OTEL_EXPORTER_OTLP_TRACES_ENDPOINT=otel-collector:4317", "opentelemetry-instrument", "python", "app.py"]
