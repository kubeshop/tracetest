module.exports = {
  port: process.env.WALLET_API_PORT,
  otelServiceName: process.env.OTEL_SERVICE_NAME || 'wallet-api',
  otelExporterGrpcUrl: process.env.OTEL_EXPORTER_GRPC_URL || 'otel-collector:4317',
}
