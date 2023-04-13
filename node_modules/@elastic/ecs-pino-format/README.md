<img align="right" width="auto" height="auto" src="https://www.elastic.co/static-res/images/elastic-logo-200.png">

# @elastic/ecs-pino-format

[![Build Status](https://apm-ci.elastic.co/buildStatus/icon?job=apm-agent-nodejs%2Fecs-logging-nodejs-mbp%2Fmaster)](https://apm-ci.elastic.co/job/apm-agent-nodejs/job/ecs-logging-nodejs-mbp/job/master/)  [![js-standard-style](https://img.shields.io/badge/code%20style-standard-brightgreen.svg?style=flat)](http://standardjs.com/)

This Node.js package provides a formatter for the [pino](https://www.npmjs.com/package/pino)
logger compatible with [Elastic Common Schema (ECS) logging](https://www.elastic.co/guide/en/ecs/current/index.html).<br/>
In combination with the [filebeat](https://www.elastic.co/products/beats/filebeat)
shipper, you can send your logs directly to Elasticsearch and leverage
[Kibana's Logs app](https://www.elastic.co/guide/en/observability/current/monitor-logs.html)
to inspect all logs in one single place.

Please see the [Node.js ECS pino documentation](https://www.elastic.co/guide/en/ecs-logging/nodejs/current/pino.html).


## Install

```sh
npm install @elastic/ecs-pino-format
```

## Usage

This package will configure Pino's `formatters`, `messageKey` and `timestamp` options.

```js
const ecsFormat = require('@elastic/ecs-pino-format')
const pino = require('pino')

const log = pino(ecsFormat())
log.info('Hello world')

const child = log.child({ module: 'foo' })
child.warn('From child')
```

Running this will produce log output similar to the following:

```sh
{"log.level":"info","@timestamp":"2021-01-19T22:51:12.142Z","ecs":{"version":"1.6.0"},"process":{"pid":82240},"host":{"hostname":"pink.local"},"message":"Hello world"}
{"log.level":"warn","@timestamp":"2021-01-19T22:51:12.143Z","ecs":{"version":"1.6.0"},"process":{"pid":82240},"host":{"hostname":"pink.local"},"module":"foo","message":"From child"}
```

Please see the [Node.js ECS pino documentation](https://www.elastic.co/guide/en/ecs-logging/nodejs/current/pino.html) for more.

## License

This software is licensed under the [Apache 2 license](./LICENSE).
