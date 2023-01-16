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
$ xk6 build --with github.com/grafana/xk6-distributed-tracing@latest
```

## Example

```javascript
import tracing, { Http } from 'k6/x/tracing';
import { sleep } from 'k6';

export let options = {
  vus: 1,
  iterations: 10,
};

export function setup() {
  console.log(`Running xk6-distributed-tracing v${tracing.version}`, tracing);
}

export default function() {
  const http = new Http({
    propagator: "w3c",
  });
  const r = http.get('https://test-api.k6.io');
  console.log(`trace_id=${r.trace_id}`);
  sleep(1);
}
```

Result output:

```bash
$ ./k6 run script.js

          /\      |‾‾| /‾‾/   /‾‾/
     /\  /  \     |  |/  /   /  /
    /  \/    \    |     (   /   ‾‾\
   /          \   |  |\  \ |  (‾)  |
  / __________ \  |__| \__\ \_____/ .io

  execution: local
     script: script.js
     output: -

  scenarios: (100.00%) 1 scenario, 1 max VUs, 40s max duration (incl. graceful stop):
           * default: 1 looping VUs for 10s (gracefulStop: 30s)

INFO[0000] Running xk6-distributed-tracing v0.2.0  source=console
INFO[0000] trace_id=743fff0b96778539acb7139e72ea1e33
INFO[0001] trace_id=365f4637a52526db1de2d30a5568ca3a
INFO[0002] trace_id=c49e1df945049c5c3c8b59acc84d7d3b
INFO[0003] trace_id=53e1937d56aa172b46d2310e3380dfe9
INFO[0004] trace_id=d61e8757d35c9ca1780b88977ac56d72
INFO[0005] trace_id=358e794ed636d268a918dcd2f3f9db0a
INFO[0006] trace_id=992a959e09ee84f3905a215bec8b53a0
INFO[0007] trace_id=aee11c64de11744ab5b66d5dd8ed361b
INFO[0008] trace_id=c4dc45d857e99ede2bb902666457239d
INFO[0009] trace_id=7623d10293d9f03c15deb8055935664e

running (10.1s), 0/1 VUs, 10 complete and 0 interrupted iterations
default ✓ [======================================] 1 VUs  10s

     █ setup

     data_received..............: 1.6 kB 156 B/s
     data_sent..................: 1.7 kB 165 B/s
     http_req_blocked...........: avg=223.43µs min=146.53µs med=217.39µs max=314.54µs p(90)=276.68µs p(95)=295.61µs
     http_req_connecting........: avg=137.18µs min=87.22µs  med=130.17µs max=196.38µs p(90)=184.38µs p(95)=190.38µs
     http_req_duration..........: avg=6.58ms   min=5.07ms   med=6.45ms   max=7.91ms   p(90)=7.83ms   p(95)=7.87ms
     http_req_receiving.........: avg=187.27µs min=94.29µs  med=171.7µs  max=295.67µs p(90)=293.28µs p(95)=294.48µs
     http_req_sending...........: avg=128.07µs min=94.64µs  med=121.77µs max=175.65µs p(90)=160.41µs p(95)=168.03µs
     http_req_tls_handshaking...: avg=0s       min=0s       med=0s       max=0s       p(90)=0s       p(95)=0s
     http_req_waiting...........: avg=6.27ms   min=4.83ms   med=6.13ms   max=7.64ms   p(90)=7.56ms   p(95)=7.6ms
     http_reqs..................: 10     0.991797/s
     iteration_duration.........: avg=916.48ms min=65.67µs  med=1s       max=1s       p(90)=1s       p(95)=1s
     iterations.................: 10     0.991797/s
     vus........................: 1      min=1 max=1
     vus_max....................: 1      min=1 max=1

```
