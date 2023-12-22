# Tracetest + K6

This repository objective is to show how you can configure your Tracetest to run alongside your k6 load tests against an instrumented service.

For more detailed information about the K6 Tracetest Binary take a look a the [docs](https://docs.tracetest.io/tools-and-integrations/integrations/k6).

## Steps

1. [Install the Tracetest CLI](https://docs.tracetest.io/installing/)
2. Run `tracetest configure --endpoint http://localhost:11633` on a terminal
3. Run the project by using docker-compose: `docker-compose up -d` (Linux) or `docker compose up -d` (Mac)
4. Test if it works by running: `tracetest run test -f tests/test.yaml`. This will create and run a test with trace id as trigger
5. In as separate folder outside of the the Tracetest repo, build the k6 binary with the extension by using `xk6 build v0.42.0 --with github.com/kubeshop/xk6-tracetest`
6. Now you are ready to run your load test, you can achieve this by running the following command: `path/to/binary/k6 run import-pokemon.js -o xk6-tracetest`
7. After the load test finishes you should be able to see an output like the following:

  ```bash
./k6 run tracetest/examples/tracetest-k6/import-pokemon.js -o xk6-tracetest

          /\      |‾‾| /‾‾/   /‾‾/   
     /\  /  \     |  |/  /   /  /    
    /  \/    \    |     (   /   ‾‾\  
   /          \   |  |\  \ |  (‾)  | 
  / __________ \  |__| \__\ \_____/ .io

  execution: local
     script: tracetest/examples/tracetest-k6/import-pokemon.js
     output: xk6-tracetest-output (TestRunID: 93008)

  scenarios: (100.00%) 1 scenario, 1 max VUs, 36s max duration (incl. graceful stop):
           * default: 1 looping VUs for 6s (gracefulStop: 30s)

ERRO[0017] panic: Tracetest: 5 jobs failed

Goja stack:
native 

running (17.1s), 0/1 VUs, 6 complete and 0 interrupted iterations
default ✓ [======================================] 1 VUs  6s

         ✓ is status 200
     ✓ body matches de id

     █ teardown

     checks.........................: 100.00% ✓ 12       ✗ 0  
     data_received..................: 1.1 kB  67 B/s
     data_sent......................: 3.3 kB  190 B/s
     http_req_blocked...............: avg=89µs    min=3µs    med=12.5µs max=476µs  p(90)=249µs  p(95)=362.49µs
     http_req_connecting............: avg=37µs    min=0s     med=0s     max=222µs  p(90)=111µs  p(95)=166.49µs
     http_req_duration..............: avg=4.83ms  min=1.86ms med=5.35ms max=7.61ms p(90)=6.77ms p(95)=7.19ms  
       { expected_response:true }...: avg=4.83ms  min=1.86ms med=5.35ms max=7.61ms p(90)=6.77ms p(95)=7.19ms  
     http_req_failed................: 0.00%   ✓ 0        ✗ 6  
     http_req_receiving.............: avg=51µs    min=32µs   med=52.5µs max=74µs   p(90)=68µs   p(95)=71µs    
     http_req_sending...............: avg=47.83µs min=17µs   med=47µs   max=88µs   p(90)=71µs   p(95)=79.49µs 
     http_req_tls_handshaking.......: avg=0s      min=0s     med=0s     max=0s     p(90)=0s     p(95)=0s      
     http_req_waiting...............: avg=4.74ms  min=1.75ms med=5.23ms max=7.56ms p(90)=6.69ms p(95)=7.12ms  
     http_reqs......................: 6       0.350387/s
     iteration_duration.............: avg=2.44s   min=1s     med=1s     max=11.08s p(90)=5.03s  p(95)=8.06s   
     iterations.....................: 6       0.350387/s
     vus............................: 0       min=0      max=1
     vus_max........................: 1       min=1      max=1
    [TotalRuns=6, SuccessfulRus=1, FailedRuns=5] 
[FAILED] 
[Request=POST - http://localhost:8081/pokemon/import, TraceID=dc071893eaaca9de301f2147e2be372e, RunState=FINISHED FailingSpecs=true, TracetestURL= http://localhost:3000/test/kc_MgKoVR/run/272] 
[Request=POST - http://localhost:8081/pokemon/import, TraceID=dc0718fff1aca9de30b702c3a1bfad75, RunState=FINISHED FailingSpecs=true, TracetestURL= http://localhost:3000/test/kc_MgKoVR/run/275] 
[Request=POST - http://localhost:8081/pokemon/import, TraceID=dc0718b8daaca9de301e39889afca15b, RunState=FINISHED FailingSpecs=true, TracetestURL= http://localhost:3000/test/kc_MgKoVR/run/276] 
[Request=POST - http://localhost:8081/pokemon/import, TraceID=dc0718a7e2aca9de30955b5203b162a7, RunState=FINISHED FailingSpecs=true, TracetestURL= http://localhost:3000/test/kc_MgKoVR/run/273] 
[Request=POST - http://localhost:8081/pokemon/import, TraceID=dc0718edf9aca9de305916d7b1e7814c, RunState=FINISHED FailingSpecs=true, TracetestURL= http://localhost:3000/test/kc_MgKoVR/run/274] 
[SUCCESSFUL] 
[Request=POST - http://localhost:8081/pokemon/import, TraceID=dc0718c9d2aca9de3044a794f7248eab, RunState=FINISHED FailingSpecs=false, TracetestURL= http://localhost:3000/test/kc_MgKoVR/run/271] 

  ERRO[0017] a panic occurred during JS execution: Tracetest: 5 jobs failed
  ```
