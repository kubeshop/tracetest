# Tracetest Cloud + K6

This example objective is to show how you can run load tests enhanced with trace-based testing using Tracetest Cloud and k6 against an instrumented service (Pokeshop API).

For more detailed information about the K6 Tracetest Binary take a look a the [docs](https://docs.tracetest.io/tools-and-integrations/integrations/k6).

## Prerequisites

1. Signing up to [app.tracetest.io](https://app.tracetest.io)
2. Creating an [environment](https://docs.tracetest.io/concepts/environments)
3. Having access to the environment's [agent token](https://docs.tracetest.io/configuration/agent)

## Steps

1. [Install the Tracetest CLI](https://docs.tracetest.io/installing/)
2. Copy the `.env.template` file into `.env` and add the `TRACETEST_API_KEY`. This is the Agent API token from your environment.
3. Create a [token from your environment](https://docs.tracetest.io/concepts/environment-tokens).
4. Run `tracetest configure` on a terminal and select the environment in use
5. Run the project by using docker-compose: `docker-compose up -d` (Linux) or `docker compose up -d` (Mac)
6. Test if it works by running: `tracetest run test -f tests/test.yaml`. This will create and run a test with trace id as trigger
7. Build the k6 binary with the extension by using `xk6 build v0.42.0 --with github.com/kubeshop/xk6-tracetest`
8. Now you are ready to run your load test, you can achieve this by running the following command: `XK6_TRACETEST_API_TOKEN=<your-environment-token> ./k6 run ./import-pokemon.js -o xk6-tracetest`
9. After the load test finishes you should be able to see an output like the following:

<iframe src="https://app.warp.dev/block/embed/aPBVIne2tE8fp4Y27xVrmd" title="xk6 trigger" style="width: 1892px; height: 796px; border:0; overflow:hidden;" allow="clipboard-read; clipboard-write"></iframe>

## What's Next?

After running the initial set of tests, you can click the run link for any of them, update the assertions and run the scripts once more. This flow enables complete a trace-based TDD flow.

![assertions](assets/assertions.gif)
