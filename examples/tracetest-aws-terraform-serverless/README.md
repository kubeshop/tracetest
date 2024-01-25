# Quick Start - Tracetest with Terraform, AWS Fargate, Lambda, Node.js and Jaeger

This is a simple quick start on how to configure a Node.js lambda function API to use OpenTelemetry instrumentation with traces and Tracetest for enhancing your E2E and integration tests with trace-based testing. The infrastructure will use Jaeger as the trace data store and Terraform to provision the required AWS infrastructure (Fargate/Lambda).

## Steps

1. [Install the tracetest CLI](https://github.com/kubeshop/tracetest/blob/main/docs/installing.md#cli-installation)
2. Run `terraform init`, `terraform apply` and accept the changes
3. From the terraform outputs, grab the `tracetes_url` and run `tracetest configure --server-url <tracetest_url>` on a terminal to configure the CLI to send all commands to that address
4. From the terraform outputs, grab the `api_endpoint` and update the `<your_api_endpoint>` section from `tests/test.yaml`
5. Test if it works by running: `tracetest run test -f tests/test.yaml`. This would execute a test against the Node.js API Gateway endpoint that will send spans to Jaeger to be fetched from the Tracetest server.

Feel free to check out the [docs](https://docs.tracetest.io/), and join our [Slack Community](https://dub.sh/tracetest-community) for more info!
