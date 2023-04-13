/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

// Instrument the @elastic/elasticsearch module.
//
// This uses to 'request' and 'response' events from the Client (documented at
// https://www.elastic.co/guide/en/elasticsearch/client/javascript-api/current/observability.html)
// to hook into all ES server interactions.
//
// Limitations:
// - In @elastic/elasticsearch >=7.14 <8, the diagnostic events sent for ES
//   spans started before the product-check is finished will have an incorrect
//   `currentSpan`.
//
//   An Elasticsearch (ES) request typically results in a single HTTP request to
//   the server. For some of the later 7.x versions of @elastic/elasticsearch
//   there is a product-check "GET /" that blocks the *first* request to the
//   server. The handling of ES requests are effectively queued until that
//   product-check is complete. When they *do* run, the async context is that
//   of the initial ES span. This means that `apm.currentSpan` inside an ES
//   client diagnostic event for these queued ES requests will be wrong.
//   Currently the APM agent is not patching for this.
//
// - When using the non-default `contextManager="patch"` APM Agent option with
//   @elastic/elasticsearch >=8 instrumentation, the diagnostic events do not
//   have the async run context of the current span. There are two impacts:
//   1. Elasticsearch tracing spans will not have additional "http" context
//      about the underlying HTTP request.
//   2. Users cannot access `apm.currentSpan` inside a diagnostic event handler.

const semver = require('semver')
const { URL, URLSearchParams } = require('url')

const { getDBDestination } = require('../../context')
const { getElasticsearchDbStatement } = require('../../elasticsearch-shared')
const shimmer = require('../../shimmer')

/**
 * Get the Elasticsearch cluster name, if possible.
 *
 * This is currently a partial implementation of:
 * https://github.com/elastic/apm/blob/main/specs/agents/tracing-instrumentation-db.md#cluster-name
 *
 * @param {import("@elastic/elasticsearch").DiagnosticResult || null}
 * @returns { string || null }
 */
function getESClusterName (diagResult) {
  if (diagResult && diagResult.headers) {
    const clusterNameFromHeader = diagResult.headers['x-found-handling-cluster']
    if (clusterNameFromHeader) {
      return clusterNameFromHeader
    }
  }
}

module.exports = function (elasticsearch, agent, { version, enabled }) {
  if (!enabled) {
    return elasticsearch
  }
  if (!elasticsearch.Client) {
    agent.logger.debug('@elastic/elasticsearch@%s is not supported (no `elasticsearch.Client`) - aborting...', version)
    return elasticsearch
  }

  // Before v7.7.0 the Transport#request() implementation's Promises support
  // would re-call `this.request(...)` inside a Promise.
  const doubleCallsRequestIfNoCb = semver.lt(version, '7.7.0')
  const ins = agent._instrumentation
  const isGteV8 = semver.satisfies(version, '>=8', { includePrerelease: true })
  const elasticsearchCaptureBodyUrlsRegExp = agent._conf.elasticsearchCaptureBodyUrlsRegExp

  agent.logger.debug('shimming elasticsearch.Transport.prototype.{request,getConnection}')
  shimmer.wrap(elasticsearch.Transport && elasticsearch.Transport.prototype, 'request', wrapRequest)
  shimmer.wrap(elasticsearch.Transport && elasticsearch.Transport.prototype, 'getConnection', wrapGetConnection)
  shimmer.wrap(elasticsearch, 'Client', wrapClient)

  // Tracking the ES client Connection object and DiagnosticResult for each
  // active span. Use WeakMap to avoid a leak from possible spans that don't
  // end.
  const connFromSpan = new WeakMap()
  const diagResultFromSpan = new WeakMap()

  return elasticsearch

  function wrapClient (OrigClient) {
    class ClientTraced extends OrigClient {
      constructor (...args) {
        super(...args)
        const diagnostic = isGteV8 ? this.diagnostic : this
        diagnostic.on('response', (_err, result) => {
          if (result) {
            const currSpan = ins.currSpan()
            if (currSpan) {
              diagResultFromSpan.set(currSpan, result)
            }
          }
        })
      }
    }
    return ClientTraced
  }

  // Transport#request() calls Transport#getConnection() when it is ready to
  // make the HTTP request. This returns the actual connection to be used for
  // the request. This is limited, however:
  // - `getConnection()` is not called if the request was aborted early.
  // - If all connections are marked dead, then this returns null.
  // - We are assuming this is called with the correct async context. See
  //   "Limitations" above.
  function wrapGetConnection (origGetConnection) {
    return function wrappedGetConnection (opts) {
      const conn = origGetConnection.apply(this, arguments)
      const currSpan = ins.currSpan()
      if (conn && currSpan) {
        connFromSpan.set(currSpan, conn)
      }
      return conn
    }
  }

  function wrapRequest (origRequest) {
    return function wrappedRequest (params, options, cb) {
      options = options || {}
      if (typeof options === 'function') {
        cb = options
        options = {}
      }

      if (typeof cb !== 'function' && doubleCallsRequestIfNoCb) {
        return origRequest.apply(this, arguments)
      }

      const method = (params && params.method) || '<UnknownMethod>'
      const path = (params && params.path) || '<UnknownPath>'
      agent.logger.debug({ method, path }, 'intercepted call to @elastic/elasticsearch.Transport.prototype.request')
      const span = ins.createSpan(`Elasticsearch: ${method} ${path}`, 'db', 'elasticsearch', 'request', { exitSpan: true })
      if (!span) {
        return origRequest.apply(this, arguments)
      }

      const parentRunContext = ins.currRunContext()
      const spanRunContext = parentRunContext.enterSpan(span)
      const finish = ins.bindFunctionToRunContext(spanRunContext, (err, result) => {
        // Set destination context.
        // Use the connection from wrappedGetConnection() above, if that worked.
        // Otherwise, fallback to using the first connection on
        // `Transport#connectionPool`, if any.  (This is the best parsed
        // representation of connection options passed to the Client ctor.)
        let conn = connFromSpan.get(span)
        if (conn) {
          connFromSpan.delete(span)
        } else if (this.connectionPool && this.connectionPool.connections) {
          conn = this.connectionPool.connections[0]
        }
        const connUrl = conn && conn.url
        span._setDestinationContext(getDBDestination(
          connUrl && connUrl.hostname, connUrl && connUrl.port))

        // Gather some HTTP context.
        // We are *not* including the response headers b/c they are boring:
        //
        //    X-elastic-product: Elasticsearch
        //    content-type: application/json
        //    content-length: ...
        //
        // Getting the ES client request "DiagnosticResult" object has some edge cases:
        // - In v7 using a callback, we always get it as `result`.
        // - In v7 using a Promise, if the promise is rejected, then `result` is
        //   not passed.
        // - In v8, `result` only includes HTTP response info if `options.meta`
        //   is true. We use the diagnostic 'response' event instead.
        // - In v7, see the limitation note above for the rare start case where
        //   the diagnostic 'response' event may have the wrong currentSpan.
        // The result is that with Promise usage of v7, ES client requests that
        // are queued behind the "product-check" and that reject, won't have a
        // `diagResult`.
        const httpContext = {}
        let haveHttpContext = false
        let diagResult = isGteV8 ? null : result
        if (!diagResult) {
          diagResult = diagResultFromSpan.get(span)
          if (diagResult) {
            diagResultFromSpan.delete(span)
          }
        }
        if (diagResult) {
          if (diagResult.statusCode) {
            haveHttpContext = true
            httpContext.status_code = diagResult.statusCode
          }

          if (diagResult.headers && 'content-length' in diagResult.headers) {
            const contentLength = Number(diagResult.headers['content-length'])
            if (!isNaN(contentLength)) {
              haveHttpContext = true
              httpContext.response = { encoded_body_size: contentLength }
            }
          }
        }

        // Reconstruct the full URL (including query params).
        let origin
        if (connUrl) {
          origin = connUrl.origin
        } else if (diagResult && diagResult.meta && diagResult.meta.connection && diagResult.meta.connection.url) {
          try {
            origin = new URL(diagResult.meta.connection.url).origin
          } catch (_ignoredErr) {}
        }
        if (origin && params && params.path) {
          const fullUrl = new URL(origin)
          fullUrl.pathname = params.path
          fullUrl.search = new URLSearchParams(params.querystring).toString()
          httpContext.url = fullUrl.toString()
          haveHttpContext = true
        }

        if (haveHttpContext) {
          span.setHttpContext(httpContext)
        }

        // Set DB context.
        const dbContext = {
          type: 'elasticsearch'
        }
        if (params) {
          const statement = getElasticsearchDbStatement(
            params.path,
            params.body || params.bulkBody,
            elasticsearchCaptureBodyUrlsRegExp
          )
          if (statement) {
            dbContext.statement = statement
          }
        }
        const clusterName = getESClusterName(diagResult)
        if (clusterName) {
          dbContext.instance = clusterName
        }
        span.setDbContext(dbContext)

        if (err) {
          // Error properties are specified here:
          // https://github.com/elastic/elasticsearch-js/blob/master/lib/errors.d.ts
          // - We capture some data from ResponseError, which is for
          //   Elasticsearch API errors:
          //   https://www.elastic.co/guide/en/elasticsearch/reference/current/common-options.html#common-options-error-options
          // - Otherwise we explicitly turn off `captureAttributes` to avoid
          //   grabbing potentially large and sensitive properties like
          //   `err.data` on DeserializationError.
          const errOpts = {
            captureAttributes: false
          }
          const errBody = err.body
          if (err.name === 'ResponseError' && errBody && errBody.error) {
            // Include some data from the Elasticsearch API response body:
            // https://www.elastic.co/guide/en/elasticsearch/reference/current/common-options.html#common-options-error-options
            const errType = errBody.error.type
            if (errType) {
              // Specialize `error.exception.type` for better error grouping.
              errOpts.exceptionType = `ResponseError (${errType})`
            }
            errOpts.custom = {
              type: errType,
              reason: errBody.error.reason,
              status: errBody.status
            }
            if (errBody.error.caused_by) {
              errOpts.custom.caused_by = errBody.error.caused_by
            }
          }
          agent.captureError(err, errOpts)
        }

        span.end()
      })

      if (typeof cb === 'function') {
        const wrappedCb = (err, result) => {
          finish(err, result)
          ins.withRunContext(parentRunContext, cb, this, err, result)
        }
        return ins.withRunContext(spanRunContext, origRequest, this, params, options, wrappedCb)
      } else {
        const origPromise = ins.withRunContext(spanRunContext, origRequest, this, ...arguments)
        origPromise.then(
          function onResolve (result) {
            finish(null, result)
          },
          function onReject (err) {
            finish(err, null)
          }
        )

        return origPromise
      }
    }
  }
}
