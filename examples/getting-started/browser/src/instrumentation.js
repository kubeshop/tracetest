import TracetestWebSDK from "@tracetest/opentelemetry-web";

const sdk = new TracetestWebSDK({
  serviceName: "browser-app",
  endpoint: "http://localhost:4318/v1/traces",
});

sdk.start();
