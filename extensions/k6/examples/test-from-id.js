import { Http, Tracetest } from "k6/x/tracetest";
import { sleep } from "k6";

export const options = {
  vus: 1,
  duration: "5s",
};

const http = new Http();
const tracetest = Tracetest({
  serverUrl: "http://localhost:3000",
});
const testId = "J0d887oVR";

export default function () {
  /// successful test run
  http.get("http://localhost:8081/pokemon?take=5", {
    tracetest: {
      testId,
    },
  });

  /// failed test specs test run
  http.get("http://localhost:8081/pokemon?take=10", {
    tracetest: {
      testId: "nDdBCnoVg",
    },
  });

  /// not existing test
  // http.get("http://localhost:8081/pokemon?take=10", {
  //   tracetest: {
  //     testId: "doesnt-exist",
  //   },
  // });

  /// wrong endpoint
  // http.get("http://localhost:8081/wrong", {
  //   tracetest: {
  //     testId: "nDdBCnoVg",
  //   },
  // });

  sleep(1);
}

export function handleSummary() {
  return {
    stdout: tracetest.summary(),
    'tracetest.json': tracetest.json(),
  };
}
