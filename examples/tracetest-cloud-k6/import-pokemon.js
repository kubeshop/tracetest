import { Http, Tracetest } from "k6/x/tracetest";
import { sleep } from "k6";

export const options = {
  vus: 1,
  duration: "5s",
};

const http = new Http();
const testId = "kc_MgKoVR";
const tracetest = Tracetest();

let pokemonId = 6;

export default function () {
  const url = "http://localhost:8081/pokemon/import";
  const payload = JSON.stringify({
    id: pokemonId++,
  });
  const params = {
    headers: {
      "Content-Type": "application/json",
    },
    tracetest: {
      testId,
    },
  };

  const response = http.post(url, payload, params);

  tracetest.runTest(
    response.trace_id,
    {
      test_id: testId,
      variable_name: "TRACE_ID",
      should_wait: true,
    },
    {
      id: "123",
      url,
      method: "GET",
    }
  );

  sleep(1);
}

export function handleSummary() {
  return {
    stdout: tracetest.summary(),
  };
}

export function teardown() {
  tracetest.validateResult();
}