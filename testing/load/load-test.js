
import { check } from "k6";
import { textSummary } from "https://jslib.k6.io/k6-summary/0.0.2/index.js";
import { Http, Tracetest } from "k6/x/tracetest";
import { sleep } from "k6";

export const options = {
  stages: [
    { duration: "5m", target: 30 },
  ],
  thresholds: {
    http_req_duration: ["p(95)<500"],
  },
};

const tracetest = Tracetest({
  serverUrl: "http://localhost:11633",
});
const testId = "kc_MgKoVR";
const http = new Http();
const url = "http://localhost:8081/pokemon?take=5";

export default function () {
  const params = {
    tracetest: {
      testId,
    },
    headers: {
      "Content-Type": "application/json",
    },
  };

  const response = http.get(url, params);

  check(response, {
    "is status 200": (r) => r.status === 200,
  });

  sleep(1);
}

// enable this to return a non-zero status code if a tracetest test fails
export function teardown() {
  tracetest.validateResult();
}

export function handleSummary(data) {
  // combine the default summary with the tracetest summary
  const tracetestSummary = tracetest.summary();
  const defaultSummary = textSummary(data);
  const summary = `
    ${defaultSummary}
    ${tracetestSummary}
  `;

  return {
    stdout: summary,
    "tracetest.json": tracetest.json(),
  };
}
