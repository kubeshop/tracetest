# Integration Testing Vercel Functions with Tracetest and OpenTelemetry

This example is from the article [**Integration Testing Vercel Serverless Functions**](add link) showing how to run integration tests against Vercel Functions using [OpenTelemetry](https://opentelemetry.io/) and Tracetest.

This is a [Next.js](https://nextjs.org/) project bootstrapped with [`create-next-app`](https://github.com/vercel/next.js/tree/canary/packages/create-next-app).

It's using Vercel Functions via `/pages/api`, with [OpenTelemetry configured as explained in the Vercel docs](https://nextjs.org/docs/pages/building-your-application/optimizing/open-telemetry#manual-opentelemetry-configuration).

## Prerequisites

- [Tracetest Account](https://app.tracetest.io/)
- [Tracetest Agent API Key](https://docs.tracetest.io/configuration/agent)
- [Tracetest Environment Token](https://docs.tracetest.io/concepts/environment-tokens)
- [Vercel Account](https://vercel.com/)
- [Vercel Postgres Database](https://vercel.com/docs/storage/vercel-postgres)

## Getting Started with Docker

0. Set Tracetest Agent API Key in `docker-compose.yaml`, and set Tracetest Environment Token in `test/run.bash`. Set the Vercel Postgres credentials in the `.env*` files.

1. Run Docker Compose

    ```bash
    docker compose up -d --build
    ```

2. Run Integration Tests

    ```bash
    docker compose run integration-tests
    ```

(Optional. Trigger Tracetest Tests via `app.tracetest.io` against `http://next-app:3000`)

## Getting Started Locally

0. Set the Vercel Postgres credentials in the `.env*` files.

1. Install Node Packages

    ```bash
    npm i
    ```

2. Run Development Server

    ```bash
    npm run dev
    ```

3. Start Tracetest Agent

    ```bash
    tracetest start --api-key ttagent_<apikey>
    ```

4. Trigger Tracetest Tests via CLI

    ```bash
    tracetest run test -f ./test-api.development.yaml
    ```

(Optional. Trigger Tracetest Tests via `app.tracetest.io` against `http://localhost:3000`)

## Learn more

Feel free to check out the [docs](https://docs.tracetest.io/), and join our [Slack Community](https://dub.sh/tracetest-community) for more info!
