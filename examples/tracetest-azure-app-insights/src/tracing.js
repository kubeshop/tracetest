const {
  ApplicationInsightsClient,
  ApplicationInsightsConfig,
} = require("applicationinsights");
const {
  ExpressInstrumentation,
} = require("@opentelemetry/instrumentation-express");
const { HttpInstrumentation } = require("@opentelemetry/instrumentation-http");

const config = new ApplicationInsightsConfig();
config.azureMonitorExporterConfig.connectionString = process.env.CONNECTION_STRING;

const appInsights = new ApplicationInsightsClient(config);

const traceHandler = appInsights.getTraceHandler();
traceHandler.addInstrumentation(new ExpressInstrumentation());
traceHandler.addInstrumentation(new HttpInstrumentation());
