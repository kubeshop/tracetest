# xk6-distributed-tracing

This extension adds distributed tracing support to [k6](https://github.com/grafana/k6)! 

That means that if you're testing an instrumented system, you can use this extension to start the traces on k6. 

Currently, it supports HTTP requests and the following propagation formats: `w3c`, `b3`, and `jaeger`.

It is implemented using the [xk6](https://github.com/grafana/xk6) extension system.

## Build

To build a `k6` binary with this extension, first ensure you have the prerequisites:

- [Go toolchain](https://go101.org/article/go-toolchain.html)
- Git

Then:

1. Download `xk6`:

```bash
$ go install go.k6.io/xk6/cmd/xk6@latest
```

2. Build the binary:

```bash
$ xk6 build --with github.com/kubeshop/tracetest/extensions/k6@latest
```

## Example

```javascript
import { Http, Tracetest } from "k6/x/tracetest";
import { sleep } from "k6";

export const options = {
  vus: 1,
  duration: "5s",
};

const http = new Http();
const tracetest = Tracetest({
  serverUrl: "<your-tracetest-server-url>",
});
const testId = "<your-test-id>";

export default function () {
  /// successful test run
  http.get("http://localhost:8081/pokemon?take=5", {
    tracetest: {
      testId,
    },
  });

  sleep(1);
}

export function handleSummary() {
  return {
    stdout: tracetest.summary(),
  };
}
```

Result output:

```bash
$ ./k6 run examples/test-from-id.js -o xk6-tracetest

          /\      |‾‾| /‾‾/   /‾‾/   
     /\  /  \     |  |/  /   /  /    
    /  \/    \    |     (   /   ‾‾\  
   /          \   |  |\  \ |  (‾)  | 
  / __________ \  |__| \__\ \_____/ .io

  execution: local
     script: examples/test-from-id.js
     output: xk6-tracetest-output (TestRunID: 79680)

  scenarios: (100.00%) 1 scenario, 1 max VUs, 35s max duration (incl. graceful stop):
           * default: 1 looping VUs for 5s (gracefulStop: 30s)


running (05.1s), 0/1 VUs, 5 complete and 0 interrupted iterations
default ✓ [======================================] 1 VUs  5s
[TotalRuns=8, SuccessfulRus=4, FailedRuns=4] 
[FAILED] 
[Request=GET - http://localhost:8081/pokemon?take=10, TraceID=dc0718f6bd99b3dc30cc624b154beb23, RunState=FINISHED FailingSpecs=false, TracetestURL= http://localhost:11633/test/nDdBCnoVg/run/28] 
[Request=GET - http://localhost:8081/pokemon?take=10, TraceID=dc0718dacd99b3dc30dc0ed028659513, RunState=FINISHED FailingSpecs=false, TracetestURL= http://localhost:11633/test/nDdBCnoVg/run/25] 
[Request=GET - http://localhost:8081/pokemon?take=10, TraceID=dc071882b699b3dc30f993c5e0a8b330, RunState=FINISHED FailingSpecs=false, TracetestURL= http://localhost:11633/test/nDdBCnoVg/run/26] 
[Request=GET - http://localhost:8081/pokemon?take=10, TraceID=dc0718e3c599b3dc30c09b0fff1e83d7, RunState=FINISHED FailingSpecs=false, TracetestURL= http://localhost:11633/test/nDdBCnoVg/run/27] 
[SUCCESSFUL] 
[Request=GET - http://localhost:8081/pokemon?take=5, TraceID=dc0718cfcd99b3dc301f7fc40fa024a8, RunState=FINISHED FailingSpecs=false, TracetestURL= http://localhost:11633/test/J0d887oVR/run/145] 
[Request=GET - http://localhost:8081/pokemon?take=5, TraceID=dc0718f2bd99b3dc300da669a9c1d4b5, RunState=FINISHED FailingSpecs=false, TracetestURL= http://localhost:11633/test/J0d887oVR/run/142] 
[Request=GET - http://localhost:8081/pokemon?take=5, TraceID=dc0718e0c599b3dc30e774886605729d, RunState=FINISHED FailingSpecs=false, TracetestURL= http://localhost:11633/test/J0d887oVR/run/143] 
[Request=GET - http://localhost:8081/pokemon?take=5, TraceID=dc0718f7b599b3dc30f5fab29c1b01f4, RunState=FINISHED FailingSpecs=false, TracetestURL= http://localhost:11633/test/J0d887oVR/run/144] 
```
