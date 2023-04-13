# @elastic/ecs-pino-format Changelog

## v1.3.0

- TypeScript types. ([#82](https://github.com/elastic/ecs-logging-nodejs/pull/82))

## v1.2.0

- Add an *internal testing-only* option (`opts._elasticApm`) to pass in the
  current loaded "elastic-apm-node" module for use in APM tracing integration.
  This option will be used by tests in the APM agent where the current agent
  import name is a local path rather than "elastic-apm-node" that this code
  normally uses.

## v1.1.2

- Fix a circular-require for code that uses both this package and
  'elastic-apm-node'.
  ([#79](https://github.com/elastic/ecs-logging-nodejs/issues/79))


## v1.1.1

- The ecs-logging spec was [updated to allow "message" to be
  optional](https://github.com/elastic/ecs-logging/pull/55). This allows the
  [change to fallback to an empty string message](https://github.com/elastic/ecs-logging-nodejs/pull/64)
  to be removed -- which is cleaner and fixes a
  [side-effect bug](https://github.com/elastic/ecs-logging-nodejs/issues/73)
  where usage of pino's `prettyPrint: true` was broken.
  ([#75](https://github.com/elastic/ecs-logging-nodejs/pull/75))

## v1.1.0

- Fix a "TypeError: Cannot read property 'host' of undefined" crash when using
  `convertReqRes: true` and logging a `req` field that is not an HTTP request
  object.
  ([#71](https://github.com/elastic/ecs-logging-nodejs/pull/71))

- Set the "message" to the empty string for logger calls that provide no
  message, e.g. `log.info({foo: 'bar'})`. In this case pino will not add a
  message field, which breaks ecs-logging spec.
  ([#64](https://github.com/elastic/ecs-logging-nodejs/pull/64))

- Fix handling when the [`base`](https://getpino.io/#/docs/api?id=base-object)
  option is used to the pino constructor.
  ([#63](https://github.com/elastic/ecs-logging-nodejs/pull/63))

  Before this change, using, for example:
        const log = pino({base: {foo: "bar"}, ...ecsFormat()})
  would result in two issues:
  1. The log records would not include the "foo" field.
  2. The log records would include `"process": {}, "host": {}` for the
     expected process.pid and os.hostname.

  Further, if the following is used:
        const log = pino({base: null, ...ecsFormat()})
  pino does not call `formatters.bindings()` at all, resulting in log
  records that were missing "ecs.version" (making them invalid ecs-logging
  records) and part of the APM integration.

- Add `apmIntegration: false` option to all ecs-logging formatters to
  enable explicitly disabling Elastic APM integration.
  ([#62](https://github.com/elastic/ecs-logging-nodejs/pull/62))

- Fix "elasticApm.isStarted is not a function" crash on startup.
  ([#60](https://github.com/elastic/ecs-logging-nodejs/issues/60))

## v1.0.0

- Update to @elastic/ecs-helpers@1.0.0: ecs.version is now "1.6.0",
  http.request.method is no longer lower-cased, improvements to HTTP
  serialization.

- Add error logging feature. By default if an Error instance is passed as the
  `err` field, then it will be converted to
  [ECS Error fields](https://www.elastic.co/guide/en/ecs/current/ecs-error.html),
  e.g.:


  ```js
  log.info({ err: new Error('boom') }, 'oops')
  ```

  yields:

  ```js
  {
    "log.level": "info",
    "@timestamp": "2021-01-26T17:02:23.697Z",
    ...
    "error": {
      "type": "Error",
      "message": "boom",
      "stack_trace": "Error: boom\n    at Object.<anonymous> (..."
    },
    "message": "oops"
  }
  ```

  This special handling of the `err` field can be disabled via the
  `convertErr: false` formatter option.

- Set "service.name" and "event.dataset" log fields if Elastic APM is started.
  This helps to filter for different log streams in the same pod and the
  latter is required for log anomaly detection.
  ([#41](https://github.com/elastic/ecs-logging-nodejs/issues/41))

- Add support for [ECS tracing fields](https://www.elastic.co/guide/en/ecs/current/ecs-tracing.html).
  If it is detected that [Elastic APM](https://www.npmjs.com/package/elastic-apm-node)
  is in use and there is an active trace, then tracing fields will be added to
  log records. This enables linking between traces and log records in Kibana.
  ([#35](https://github.com/elastic/ecs-logging-nodejs/issues/35))

- BREAKING CHANGE: Conversion of HTTP request and response objects is no longer
  done by default. One must use the new `convertReqRes: true` formatter option.
  As well, only the meta keys `req` and `res` will be handled. Before this
  change the meta keys `req`, `res`, `request`, and `response` would all be
  handled. ([#32](https://github.com/elastic/ecs-logging-nodejs/issues/32))

  Before (no longer works):

  ```
  const log = pino({ ...ecsFormat() })

  http.createServer(function handler (request, response) {
    // ...
    log.info({ request, response }, 'handled request')
  })
  ```

  After:

  ```
  const log = pino({ ...ecsFormat({ convertReqRes: true }) }) // <-- specify convertReqRes option

  http.createServer(function handler (req, res) {
    // ...
    log.info({ req, res }, 'handled request') // <-- only `req` and `res` are special
  })
  ```

## v0.2.0

- Serialize "log.level" as a top-level dotted field per
  https://github.com/elastic/ecs-logging/pull/33 and
  set ["log.logger"](https://www.elastic.co/guide/en/ecs/current/ecs-log.html#field-log-logger)
  to the logger ["name"](https://getpino.io/#/docs/api?id=name-string) if given.
  ([#23](https://github.com/elastic/ecs-logging-nodejs/pull/23))

## v0.1.0

Initial release.
