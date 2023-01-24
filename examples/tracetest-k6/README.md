# Tracetest + K6

This repository objective is to show how you can configure your tracetest to run alongside your k6 load tests against an instrumented service

## Steps

1. [Install the tracetest CLI](https://docs.tracetest.io/installing/)
2. Run `tracetest configure --endpoint http://localhost:11633` on a terminal
3. Run the project by using docker-compose: `docker-compose up` (Linux) or `docker compose up` (Mac)
4. Test if it works by running: `tracetest test run -d tests/test.yaml`. This will create and run a test with trace id as trigger
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
      output: xk6-tracetest-output (TestRunID: 46145)

    scenarios: (100.00%) 1 scenario, 1 max VUs, 36s max duration (incl. graceful stop):
            * default: 1 looping VUs for 6s (gracefulStop: 30s)


  running (06.0s), 0/1 VUs, 6 complete and 0 interrupted iterations
  default ✓ [======================================] 1 VUs  6s
  [TotalRuns=6, SuccessfulRus=6, FailedRuns=0] 
  [SUCCESSFUL] 
  [Request=POST - http://localhost:8081/pokemon/import, TraceID=dc0718dae7fd8dde305bb00c9a3ff4d3, RunState=FINISHED FailingSpecs=false, TracetestURL= http://localhost:11633/test/kc_MgKoVR/run/6] 
  [Request=POST - http://localhost:8081/pokemon/import, TraceID=dc0718eddffd8dde3070aa912cf76fff, RunState=FINISHED FailingSpecs=false, TracetestURL= http://localhost:11633/test/kc_MgKoVR/run/7] 
  [Request=POST - http://localhost:8081/pokemon/import, TraceID=dc0718b2f7fd8dde3023788d4498a1a9, RunState=FINISHED FailingSpecs=false, TracetestURL= http://localhost:11633/test/kc_MgKoVR/run/3] 
  [Request=POST - http://localhost:8081/pokemon/import, TraceID=dc071892d0fd8dde30da7fa44694f64c, RunState=FINISHED FailingSpecs=false, TracetestURL= http://localhost:11633/test/kc_MgKoVR/run/2] 
  [Request=POST - http://localhost:8081/pokemon/import, TraceID=dc0718c6effd8dde30f35102cd0afcd0, RunState=FINISHED FailingSpecs=false, TracetestURL= http://localhost:11633/test/kc_MgKoVR/run/4] 
  [Request=POST - http://localhost:8081/pokemon/import, TraceID=dc071880d8fd8dde3049cbe7f9bb5367, RunState=FINISHED FailingSpecs=false, TracetestURL= http://localhost:11633/test/kc_MgKoVR/run/5]   
  ```
