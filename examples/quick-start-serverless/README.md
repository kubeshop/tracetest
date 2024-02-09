# Quick start Tracetest + Serverless 

This repository's objective is to show how you can configure trace-based tests using Tracetest in your AWS Serverless infrastructure.

## Steps

1. Copy the `.env.template` file to `.env`.
2. Log into **the** [Tracetest app](https://app.tracetest.io/).
3. This example is configured to use the OpenTelemetry Collector. Ensure the environment you will be utilizing to run this example is also configured to use the OpenTelemetry Tracing Backend by clicking on Settings, Tracing Backend, OpenTelemetry, and Save.
4. Configure your environment to use [the cloud agent](https://docs.tracetest.io/concepts/cloud-agent), click the Click the Settings link and from the Agent tab select the "Run Agent in tracetest cloud" option.
5. Fill out the [token](https://docs.tracetest.io/concepts/environment-tokens) and [agent url](https://docs.tracetest.io/concepts/cloud-agent) details by editing your .env file. You can find these values in the Settings area for your environment.
6. Run `npm install`.
7. Run `npm start`. The serverless deployment will be triggered first, creating the CloudFormation stack with the necessary resources, and then the trace-based tests included in the `tracetest.ts` script will be executed.
8. When the tests are done, in the terminal you can find the results from the trace-based tests that are triggered.
9. Follow the links in the log to view the test runs programmatically created by your Typescript script.
