# Tracetest + k6

This example's objective is to show how you can run load tests enhanced with trace-based testing using Tracetest Cloud and k6 against an instrumented service (Pokeshop API).

For more detailed information about the K6 Tracetest Binary take a look a the [docs](https://docs.tracetest.io/tools-and-integrations/integrations/k6).

## Prerequisites

1. Signing up to [app.tracetest.io](https://app.tracetest.io).
2. Creating an [environment](https://docs.tracetest.io/concepts/environments).
3. Having access to the environment's [agent token](https://docs.tracetest.io/configuration/agent).

## Steps

1. [Install the Tracetest CLI](https://docs.tracetest.io/installing/).
2. Copy the `.env.template` file into `.env` and add the `TRACETEST_API_KEY`. This is the Agent API token from your environment.
3. Create a [token from your environment](https://docs.tracetest.io/concepts/environment-tokens).
4. Run the project by using docker-compose: `docker-compose run k6-tracetest` (Linux) or `docker compose run k6-tracetest` (Mac).
5. After the load test finishes you should be able to see an output like the following:

```bash
docker compose run k6-tracetest
WARN[0000] The "TRACETEST_SERVER_URL" variable is not set. Defaulting to a blank string.
[+] Running 1/1
 ✔ demo-api Pulled                                                                                                                                                1.4s
[+] Creating 5/0
 ✔ Container tracetest-k6-tracetest-agent-1  Running                                                                                                              0.0s
 ✔ Container tracetest-k6-cache-1            Running                                                                                                              0.0s
 ✔ Container tracetest-k6-queue-1            Running                                                                                                              0.0s
 ✔ Container tracetest-k6-postgres-1         Running                                                                                                              0.0s
 ✔ Container tracetest-k6-demo-api-1         Running                                                                                                              0.0s
[+] Running 3/3
 ✔ Container tracetest-k6-postgres-1  Healthy                                                                                                                     0.5s
 ✔ Container tracetest-k6-cache-1     Healthy                                                                                                                     0.5s
 ✔ Container tracetest-k6-queue-1     Healthy                                                                                                                     0.5s
[+] Running 1/1
 ✔ demo-api Pulled                                                                                                                                                1.4s

          /\      |‾‾| /‾‾/   /‾‾/
     /\  /  \     |  |/  /   /  /
    /  \/    \    |     (   /   ‾‾\
   /          \   |  |\  \ |  (‾)  |
  / __________ \  |__| \__\ \_____/ .io

     execution: local
        script: /import-pokemon.js
        output: xk6-tracetest-output (TestRunID: 39663)

     scenarios: (100.00%) 1 scenario, 1 max VUs, 35s max duration (incl. graceful stop):
              * default: 1 looping VUs for 5s (gracefulStop: 30s)

Goja stack:
[RunGroup=#BlfA2YxIg, Status=failed] - https://app.tracetest.io/organizations/ttorg_ced62e34638d965e/environments/ttenv_b42fa137465c6e04/run/BlfA2YxIg
[TotalRuns=1, SuccessfulRus=0, FailedRuns=1]
[FAILED Request=GET - http://demo-api:8081/pokemon/import, TraceID=dc071892e690abea3133062b5d9aa298, RunState=FINISHED FailingSpecs=true, TracetestURL= https://app.tracetest.io/organizations/ttorg_ced62e34638d965e/environments/ttenv_b42fa137465c6e04/test/import-pokemon-k6/run/11]
[FAILED Request=GET - http://demo-api:8081/pokemon/import, TraceID=dc0718ffed90abea3162c624004519f2, RunState=FINISHED FailingSpecs=true, TracetestURL= https://app.tracetest.io/organizations/ttorg_ced62e34638d965e/environments/ttenv_b42fa137465c6e04/test/import-pokemon-k6/run/12]
[FAILED Request=GET - http://demo-api:8081/pokemon/import, TraceID=dc0718e9f590abea3134c5682550bbd5, RunState=FINISHED FailingSpecs=true, TracetestURL= https://app.tracetest.io/organizations/ttorg_ced62e34638d965e/environments/ttenv_b42fa137465c6e04/test/import-pokemon-k6/run/13]
[FAILED Request=GET - http://demo-api:8081/pokemon/import, TraceID=dc0718d3fd90abea31a06e8b40545477, RunState=FINISHED FailingSpecs=true, TracetestURL= https://app.tracetest.io/organizations/ttorg_ced62e34638d965e/environments/ttenv_b42fa137465c6e04/test/import-pokemon-k6/run/14]
[FAILED Request=GET - http://demo-api:8081/pokemon/import, TraceID=dc0718bd8591abea31c86274fb4b7e0d, RunState=FINISHED FailingSpecs=true, TracetestURL= https://app.tracetest.io/organizations/ttorg_ced62e34638d965e/environments/ttenv_b42fa137465c6e04/test/import-pokemon-k6/run/15]

running (16.3s), 0/1 VUs, 5 complete and 0 interrupted iterations
default ✓ [======================================] 1 VUs  5s
ERRO[0020] a panic occurred during JS execution: Tracetest: 1 jobs failed
```

## What's Next?

After running the initial set of tests, you can click the run link for any of them, update the assertions and run the scripts once more. This flow enables complete a trace-based TDD flow.

![assertions](assets/assertions.gif)
