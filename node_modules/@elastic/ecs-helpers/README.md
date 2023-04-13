<img align="right" width="auto" height="auto" src="https://www.elastic.co/static-res/images/elastic-logo-200.png">

# @elastic/ecs-helpers

[![Build Status](https://apm-ci.elastic.co/buildStatus/icon?job=apm-agent-nodejs%2Fecs-logging-nodejs-mbp%2Fmaster)](https://apm-ci.elastic.co/job/apm-agent-nodejs/job/ecs-logging-nodejs-mbp/job/master/)  [![js-standard-style](https://img.shields.io/badge/code%20style-standard-brightgreen.svg?style=flat)](http://standardjs.com/)

A set of helpers for the ECS logging libraries. You should not directly used
this package, but the [ECS logging libraries](../loggers) instead.

## Install

```sh
npm install @elastic/ecs-helpers
```

## API

### `version`

The currently supported version of [Elastic Common Schema](https://www.elastic.co/guide/en/ecs/current/index.html).

### `stringify`

Function that serializes (very quickly!) an ECS-format log record object.

```js
const { stringify } = require('@elastic/ecs-helpers')
const ecs = {
  '@timestamp': new Date().toISOString(),
  'log.level': 'info',
  message: 'hello world',
  log: {
    logger: 'test'
  },
  ecs: {
    version: '1.4.0'
  }
}

console.log(stringify(ecs))
```

Note: This uses [fast-json-stringify](https://github.com/fastify/fast-json-stringify)
for serialization. By design this chooses speed over supporting serialization
of objects with circular references. This generally means that ecs-logging-nodejs
libraries will throw a "Converting circular structure to JSON" exception if an
attempt is made to log an object with circular references.

### `formatError(obj, err) -> bool`

A function that adds [ECS Error fields](https://www.elastic.co/guide/en/ecs/current/ecs-error.html)
for a given `Error` object. It returns true iff the given `err` was an Error
object it could process.

```js
const { formatError } = require('@elastic/ecs-helpers')
const logRecord = { msg: 'oops', /* ... */ }
formatError(logRecord, new Error('boom'))
console.log(logRecord)
```

will show:

```js
{
  msg: 'oops',
  error: {
    type: 'Error',
    message: 'boom',
    stack_trace: 'Error: boom\n' +
      '    at REPL30:1:26\n' +
      '    at Script.runInThisContext (vm.js:133:18)\n' +
      // ...
  }
}
```

The ECS logging libraries typically use this to automatically handle an `err`
metadata field passed to a logging statement. E.g.
`log.warn({err: myErr}, '...')` for pino, `log.warn('...', {err: myErr})`
for winston.

### `formatHttpRequest(obj, req) -> bool`

Function that enhances an ECS object with http request data.
The given request object, `req`, must be one of the following:
- Node.js's core [`http.IncomingMessage`](https://nodejs.org/api/all.html#http_class_http_incomingmessage),
- [Express's request object](https://expressjs.com/en/5x/api.html#req) that extends IncomingMessage, or
- a [hapi request object](https://hapi.dev/api/#request).

The function returns true iff the given `req` was a request object it could
process. Note that currently this notably does not process a
[`http.ClientRequest`](https://nodejs.org/api/all.html#http_class_http_clientrequest)
as returned from `http.request()`.

```js
const http = require('http')
const { formatHttpRequest } = require('@elastic/ecs-helpers')

http.createServer(function (req, res) {
  res.end('hi')

  const obj = {}
  formatHttpRequest(obj, req)
  console.log('obj:', JSON.stringify(obj, null, 4))
}).listen(3000)
```

Running this and making a request via `curl http://localhost:3000/` will
print something close to:

```
obj: {
    "http": {
        "version": "1.1",
        "request": {
            "method": "get",
            "headers": {
                "host": "localhost:3000",
                "accept": "*/*"
            }
        }
    },
    "url": {
        "full": "http://localhost:3000/",
        "path": "/"
    },
    "client": {
        "address": "::1",
        "ip": "::1",
        "port": 61969
    },
    "user_agent": {
        "original": "curl/7.64.1"
    }
}
```

### `formatHttpResponse(obj, res)`

Function that enhances an ECS object with http response data.
The given request object, `req`, must be one of the following:
- Node.js's core [`http.ServerResponse`](https://nodejs.org/api/all.html#http_class_http_serverresponse),
- [Express's response object](https://expressjs.com/en/5x/api.html#res) that extends ServerResponse, or
- a [hapi **request** object](https://hapi.dev/api/#request)

The function returns true iff the given `res` was a response object it could
process. Note that currently this notably does not process a
[`http.IncomingMessage`](https://nodejs.org/api/all.html#http_class_http_incomingmessage)
that is the argument to the
["response" event](https://nodejs.org/api/all.html#http_event_response) of a
[client `http.request()`](https://nodejs.org/api/all.html#http_http_request_options_callback)

```js
const http = require('http')
const { formatHttpRequest } = require('@elastic/ecs-helpers')

http.createServer(function (req, res) {
  res.setHeader('Foo', 'Bar')
  res.end('hi')

  const obj = {}
  formatHttpResponse(obj, res)
  console.log('obj:', JSON.stringify(obj, null, 4))
}).listen(3000)
```

Running this and making a request via `curl http://localhost:3000/` will
print something close to:

```
rec: {
    "http": {
        "response": {
            "status_code": 200,
            "headers": {
                "foo": "Bar"
            }
        }
    }
}
```

## License

This software is licensed under the [Apache 2 license](./LICENSE).
