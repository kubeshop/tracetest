import { check } from "k6";
import { textSummary } from "https://jslib.k6.io/k6-summary/0.0.2/index.js";
import { Http, Tracetest } from "k6/x/tracetest";
import { sleep } from "k6";

export const options = {
  vus: 1,
  duration: "6s",
};

const tracetest = Tracetest({
  serverUrl: "http://localhost:11633",
});
const testId = "kc_MgKoVR";
const pokemonId = 6; // charizad
const http = new Http();
const url = "http://localhost:8081/pokemon/import";

export default function () {
  const payload = JSON.stringify({
    id: pokemonId,
  });
  const params = {
    tracetest: {
      testId,
    },
    headers: {
      "Content-Type": "application/json",
    },
  };

  const response = http.post(url, payload, params);

  check(response, {
    "is status 200": (r) => r.status === 200,
    "body matches de id": (r) => JSON.parse(r.body).id === pokemonId,
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
