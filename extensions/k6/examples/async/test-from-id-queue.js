import { Http, Tracetest } from "k6/x/tracetest";
import { sleep } from "k6";

export const options = {
  vus: 1,
  duration: "10s",
};

const http = new Http({ propagator: ["w3c", "b3"] });
const testId = "EjnCE-2Vg";
const tracetest = new Tracetest();

export default function () {
  const response = http.get("https://test-api.k6.io");
  tracetest.runTest(testId, response.trace_id)

  sleep(1);
}
