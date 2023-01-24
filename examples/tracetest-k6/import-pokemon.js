import { check } from "k6";
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
const pokemonId = 6;
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

export function handleSummary() {
  return {
    stdout: tracetest.summary(),
    "tracetest.json": tracetest.json(),
  };
}
