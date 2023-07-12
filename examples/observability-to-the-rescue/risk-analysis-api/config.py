import os

from collections import namedtuple

Config = namedtuple("Config", ["port", "otel_service_name", "otel_exporter_grpc_url"])

def get_current():
  return Config(
    port=os.getenv('RISK_ANALYSIS_API_PORT'),
    otel_service_name=os.getenv('OTEL_SERVICE_NAME'),
    otel_exporter_grpc_url=os.getenv('OTEL_EXPORTER_GRPC_URL')
  )
