module.exports = {
  port: process.env.YOUR_API_PORT,
  otelServiceName: process.env.OTEL_SERVICE_NAME,
  otelExporterGrpcUrl: process.env.OTEL_EXPORTER_GRPC_URL,
  walletAPIEndpoint: process.env.WALLET_API_ENDPOINT,
  paymentExecutorAPIEndpoint: process.env.PAYMENT_EXECUTOR_API_ENDPOINT,
}
