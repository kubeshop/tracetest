Config = Struct.new(:port, :otel_service_name, :otel_exporter_http_url, :risk_analysis_url)

def current_config
  Config.new(
    ENV["PAYMENT_EXECUTOR_API_PORT"],
    ENV["OTEL_SERVICE_NAME"],
    ENV["OTEL_EXPORTER_HTTP_URL"],
    ENV["RISK_ANALYSIS_URL"],
  )
end
