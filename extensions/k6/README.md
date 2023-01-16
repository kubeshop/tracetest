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
      // your tracetest test id
      testId: "<tracetest-test-id>",
    },
  });

  sleep(1);
}

```

Result output:

```bash
$ ./k6 run examples/async/test-from-id.js -o xk6-tracetest 

          /\      |‾‾| /‾‾/   /‾‾/   
     /\  /  \     |  |/  /   /  /    
    /  \/    \    |     (   /   ‾‾\  
   /          \   |  |\  \ |  (‾)  | 
  / __________ \  |__| \__\ \_____/ .io

  execution: local
     script: examples/async/test-from-id.js
     output: xk6-crocospans (TestRunID: 95314)

  scenarios: (100.00%) 1 scenario, 1 max VUs, 40s max duration (incl. graceful stop):
           * default: 1 looping VUs for 10s (gracefulStop: 30s)

INFO[0000] METADATA:  map[EndTimeUnixNano:1673907419319621000 Group: HTTPMethod:GET HTTPStatus:200 HTTPUrl:https://test-api.k6.io Scenario:default StartTimeUnixNano:1673907418899401000 TraceID:dc0718a2efa0e5db307bc9121208994a]  component=xk6-tracetest-output
INFO[0000] Queuing Job:  {dc0718a2efa0e5db307bc9121208994a EjnCE-2Vg  runTestFromId}  component=xk6-tracetest-tracing
INFO[0001] METADATA:  map[EndTimeUnixNano:1673907420455220000 Group: HTTPMethod:GET HTTPStatus:200 HTTPUrl:https://test-api.k6.io Scenario:default StartTimeUnixNano:1673907420321396000 TraceID:dc0718a1f9a0e5db30418e782d8f108c]  component=xk6-tracetest-output
INFO[0001] Queuing Job:  {dc0718a1f9a0e5db30418e782d8f108c EjnCE-2Vg  runTestFromId}  component=xk6-tracetest-tracing
INFO[0002] METADATA:  map[EndTimeUnixNano:1673907421593667000 Group: HTTPMethod:GET HTTPStatus:200 HTTPUrl:https://test-api.k6.io Scenario:default StartTimeUnixNano:1673907421456105000 TraceID:dc07189082a1e5db30236fd0f83dde51]  component=xk6-tracetest-output
INFO[0002] Queuing Job:  {dc07189082a1e5db30236fd0f83dde51 EjnCE-2Vg  runTestFromId}  component=xk6-tracetest-tracing
INFO[0003] METADATA:  map[EndTimeUnixNano:1673907422725890000 Group: HTTPMethod:GET HTTPStatus:200 HTTPUrl:https://test-api.k6.io Scenario:default StartTimeUnixNano:1673907422594271000 TraceID:dc0718828ba1e5db30c073f6d1d02933]  component=xk6-tracetest-output
INFO[0003] Queuing Job:  {dc0718828ba1e5db30c073f6d1d02933 EjnCE-2Vg  runTestFromId}  component=xk6-tracetest-tracing
INFO[0004] Test run path /test/EjnCE-2Vg/run/76          component=xk6-tracetest-tracing
INFO[0005] METADATA:  map[EndTimeUnixNano:1673907423858117000 Group: HTTPMethod:GET HTTPStatus:200 HTTPUrl:https://test-api.k6.io Scenario:default StartTimeUnixNano:1673907423727227000 TraceID:dc0718ef93a1e5db301c1b99b3d677d3]  component=xk6-tracetest-output
INFO[0005] Queuing Job:  {dc0718ef93a1e5db301c1b99b3d677d3 EjnCE-2Vg  runTestFromId}  component=xk6-tracetest-tracing
INFO[0005] Test run path /test/EjnCE-2Vg/run/77          component=xk6-tracetest-tracing
INFO[0006] METADATA:  map[EndTimeUnixNano:1673907424989942000 Group: HTTPMethod:GET HTTPStatus:200 HTTPUrl:https://test-api.k6.io Scenario:default StartTimeUnixNano:1673907424858904000 TraceID:dc0718da9ca1e5db30d806f8ec7a89ed]  component=xk6-tracetest-output
INFO[0006] Queuing Job:  {dc0718da9ca1e5db30d806f8ec7a89ed EjnCE-2Vg  runTestFromId}  component=xk6-tracetest-tracing
INFO[0006] Test run path /test/EjnCE-2Vg/run/78          component=xk6-tracetest-tracing
INFO[0007] METADATA:  map[EndTimeUnixNano:1673907426122543000 Group: HTTPMethod:GET HTTPStatus:200 HTTPUrl:https://test-api.k6.io Scenario:default StartTimeUnixNano:1673907425991332000 TraceID:dc0718c7a5a1e5db30972829a97e7370]  component=xk6-tracetest-output
INFO[0007] Queuing Job:  {dc0718c7a5a1e5db30972829a97e7370 EjnCE-2Vg  runTestFromId}  component=xk6-tracetest-tracing
INFO[0007] Test run path /test/EjnCE-2Vg/run/79          component=xk6-tracetest-tracing
INFO[0008] Test run path /test/EjnCE-2Vg/run/80          component=xk6-tracetest-tracing
INFO[0008] METADATA:  map[EndTimeUnixNano:1673907427255703000 Group: HTTPMethod:GET HTTPStatus:200 HTTPUrl:https://test-api.k6.io Scenario:default StartTimeUnixNano:1673907427123853000 TraceID:dc0718b3aea1e5db30a67c5296572b46]  component=xk6-tracetest-output
INFO[0008] Queuing Job:  {dc0718b3aea1e5db30a67c5296572b46 EjnCE-2Vg  runTestFromId}  component=xk6-tracetest-tracing
INFO[0009] Test run path /test/EjnCE-2Vg/run/81          component=xk6-tracetest-tracing
INFO[0009] METADATA:  map[EndTimeUnixNano:1673907428388287000 Group: HTTPMethod:GET HTTPStatus:200 HTTPUrl:https://test-api.k6.io Scenario:default StartTimeUnixNano:1673907428257135000 TraceID:dc0718a1b7a1e5db30e52e6a027d74d0]  component=xk6-tracetest-output
INFO[0009] Queuing Job:  {dc0718a1b7a1e5db30e52e6a027d74d0 EjnCE-2Vg  runTestFromId}  component=xk6-tracetest-tracing
INFO[0010] Test run path /test/EjnCE-2Vg/run/82          component=xk6-tracetest-tracing

running (10.3s), 0/1 VUs, 9 complete and 0 interrupted iterations
default ✓ [======================================] 1 VUs  10s

     data_received..................: 149 kB 14 kB/s
     data_sent......................: 2.3 kB 226 B/s
     http_req_blocked...............: avg=23.23ms  min=6µs     med=10µs     max=209.06ms p(90)=41.82ms  p(95)=125.44ms
     http_req_connecting............: avg=7.31ms   min=0s      med=0s       max=65.87ms  p(90)=13.17ms  p(95)=39.52ms 
     http_req_duration..............: avg=125.2ms  min=67.78ms med=131.2ms  max=137.55ms p(90)=134.56ms p(95)=136.06ms
       { expected_response:true }...: avg=125.2ms  min=67.78ms med=131.2ms  max=137.55ms p(90)=134.56ms p(95)=136.06ms
     http_req_failed................: 0.00%  ✓ 0        ✗ 9  
     http_req_receiving.............: avg=14.33ms  min=75µs    med=143µs    max=64.41ms  p(90)=63.65ms  p(95)=64.03ms 
     http_req_sending...............: avg=34.44µs  min=16µs    med=27µs     max=107µs    p(90)=52.6µs   p(95)=79.79µs 
     http_req_tls_handshaking.......: avg=8.61ms   min=0s      med=0s       max=77.5ms   p(90)=15.5ms   p(95)=46.5ms  
     http_req_waiting...............: avg=110.83ms min=67.32ms med=130.92ms max=133.66ms p(90)=132.07ms p(95)=132.87ms
     http_reqs......................: 9      0.869737/s
     iteration_duration.............: avg=1.14s    min=1.13s   med=1.13s    max=1.27s    p(90)=1.16s    p(95)=1.22s   
     iterations.....................: 9      0.869737/s
     vus............................: 1      min=1      max=1
     vus_max........................: 1      min=1      max=1
```
