import { Http } from "k6/x/tracetest";
import { sleep } from "k6";

export const options = {
  vus: 1,
  duration: "10s",
};

const http = new Http({ propagator: ["w3c", "b3"] });

export default function () {
  http.get("https://test-api.k6.io", {
    tracetest: {
      testId: "EjnCE-2Vg",
    },
  });

  sleep(1);
}
