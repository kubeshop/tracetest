# Elastic APM Node.js Agent

This is the official Node.js agent for [Elastic APM](https://www.elastic.co/solutions/apm).

If you have any feedback or questions,
please post them on the [Discuss forum](https://discuss.elastic.co/c/apm).

[![npm](https://img.shields.io/npm/v/elastic-apm-node.svg)](https://www.npmjs.com/package/elastic-apm-node)
[![Build Status](https://apm-ci.elastic.co/buildStatus/icon?job=apm-agent-nodejs%2Fapm-agent-nodejs-mbp%2F3.x)](https://apm-ci.elastic.co/job/apm-agent-nodejs/job/apm-agent-nodejs-mbp/job/3.x/)


## Installation

```
npm install elastic-apm-node --save
```

## Quick start

1. To run Elastic APM for your own applications,
   make sure you have the prerequisites in place first.
   This agent is compatible with [APM Server](https://github.com/elastic/apm-server) v6.5 and above.
   For support for previous releases of the APM Server,
   use version [1.x](https://github.com/elastic/apm-agent-nodejs/tree/1.x) of the agent.
   For details see [Getting Started with Elastic APM](https://www.elastic.co/guide/en/apm/get-started)

1. Now follow the documentation links below relevant to your framework or stack to get set up

## Documentation

- [Table of contents](https://www.elastic.co/guide/en/apm/agent/nodejs)
- [Introduction](https://www.elastic.co/guide/en/apm/agent/nodejs/current/intro.html)
- [Get started with Express](https://www.elastic.co/guide/en/apm/agent/nodejs/current/express.html)
- [Get started with hapi](https://www.elastic.co/guide/en/apm/agent/nodejs/current/hapi.html)
- [Get started with Koa](https://www.elastic.co/guide/en/apm/agent/nodejs/current/koa.html)
- [Get started with Restify](https://www.elastic.co/guide/en/apm/agent/nodejs/current/restify.html)
- [Get started with Fastify](https://www.elastic.co/guide/en/apm/agent/nodejs/current/fastify.html)
- [Get started with Lambda](https://www.elastic.co/guide/en/apm/agent/nodejs/current/lambda.html)
- [Get started with a custom Node.js stack](https://www.elastic.co/guide/en/apm/agent/nodejs/current/custom-stack.html)
- [Advanced Setup and Configuration](https://www.elastic.co/guide/en/apm/agent/nodejs/current/advanced-setup.html)
- [API Reference](https://www.elastic.co/guide/en/apm/agent/nodejs/current/api.html)
- [OpenTelemetry Bridge](https://www.elastic.co/guide/en/apm/agent/nodejs/current/opentelemetry-bridge.html)
- [Custom Transactions](https://www.elastic.co/guide/en/apm/agent/nodejs/current/custom-transactions.html)
- [Custom Spans](https://www.elastic.co/guide/en/apm/agent/nodejs/current/custom-spans.html)
- [Metrics](https://www.elastic.co/guide/en/apm/agent/nodejs/current/metrics.html)
- [Performance Tuning](https://www.elastic.co/guide/en/apm/agent/nodejs/current/performance-tuning.html)
- [Source Map Support](https://www.elastic.co/guide/en/apm/agent/nodejs/current/source-maps.html)
- [Supported Technologies](https://www.elastic.co/guide/en/apm/agent/nodejs/current/supported-technologies.html)
- [Upgrading](https://www.elastic.co/guide/en/apm/agent/nodejs/current/upgrading.html)
- [Troubleshooting](https://www.elastic.co/guide/en/apm/agent/nodejs/current/troubleshooting.html)

## Contributing

Contributions are welcome,
but we recommend that you take a moment and read our [contribution guide](CONTRIBUTING.md) first.

To see what data is being sent to the APM Server,
use the environment variable `ELASTIC_APM_PAYLOAD_LOG_FILE` (or the config option `payloadLogFile`) to specify a log file,
e.g:

```
ELASTIC_APM_PAYLOAD_LOG_FILE=/tmp/payload.ndjson
```

Please see [TESTING.md](TESTING.md) for instructions on how to run the test suite.

## License

[BSD-2-Clause](LICENSE)

<br>Made with ♥️ and ☕️ by Elastic and our community.
