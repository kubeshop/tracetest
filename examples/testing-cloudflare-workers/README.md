# Testing Cloudflare Workers with Tracetest and OpenTelemetry

This example is from the article [**Crafting Observable Cloudflare Workers with OpenTelemetry**]() showing how to run integration and staging tests against Cloudflare Workers using [OpenTelemetry](https://opentelemetry.io/) and Tracetest.

This is a project bootstrapped with [`C3 (create-cloudflare-cli)`](https://developers.cloudflare.com/workers/get-started/guide/#1-create-a-new-worker-project).

It's using Cloudflare Workers with [OpenTelemetry configured with otel-cf-workers](https://github.com/evanderkoogh/otel-cf-workers).

## Prerequisites

- [Tracetest Account](https://app.tracetest.io/)
- [Tracetest Agent API Key](https://docs.tracetest.io/configuration/agent)
- [Tracetest Environment Token](https://docs.tracetest.io/concepts/environment-tokens)
- [Cloudflare Workers Account](https://workers.cloudflare.com/)
- [Cloudflare Database](https://developers.cloudflare.com/d1/get-started/)

Install npm modules:

```bash
npm i
```

If you do not have `npx`, install it:

```bash
npm install -g npx
```

Create a D1 database:

```bash
npx wrangler d1 create testing-cloudflare-workers

✅ Successfully created DB 'testing-cloudflare-workers' in region EEUR
Created your database using D1's new storage backend. The new storage backend is not yet recommended for production workloads, but backs up your data via point-in-time
restore.

[[d1_databases]]
binding = "DB" # i.e. available in your Worker on env.DB
database_name = "testing-cloudflare-workers"
database_id = "<your_database_id>"
```

Set the `database_id` credentials in your `wrangler.toml`.

Configure the D1 database:

```bash
npx wrangler d1 execute testing-cloudflare-workers --local --file=./schema.sql
npx wrangler d1 execute testing-cloudflare-workers --file=./schema.sql
```

## Getting Started with Docker

0. Set Tracetest Agent API Key in `docker-compose.yaml`, and set Tracetest Environment Token in `test/run.bash`, and `wrangler.toml`. Set the Cloudflare D1 credentials in the `wrangler.toml` as well.

1. Run Docker Compose

    ```bash
    docker compose up -d --build
    ```

2. Run Integration Tests

    ```bash
    docker compose run integration-tests
    ```

(Optional. Trigger Tracetest Tests via `app.tracetest.io` against `http://cloudflare-worker:8787`)

## Getting Started Locally

0. Set the Cloudflare D1 credentials in the `wrangler.toml` as well.

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
    tracetest run test -f ./api.pokemon.spec.development.yaml
    ```

(Optional. Trigger Tracetest Tests via `app.tracetest.io` against `http://localhost:8787`)

## Learn more

Feel free to check out the [docs](https://docs.tracetest.io/), and join our [Slack Community](https://dub.sh/tracetest-community) for more info!
