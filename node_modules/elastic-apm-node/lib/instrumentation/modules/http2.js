/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

var eos = require('end-of-stream')

var shimmer = require('../shimmer')
var symbols = require('../../symbols')
var { parseUrl } = require('../../parsers')
var { getHTTPDestination } = require('../context')

// Return true iff this is an HTTP 1.0 or 1.1 request.
//
// @param {http2.Http2ServerRequest} req
function reqIsHTTP1 (req) {
  return req && typeof req.httpVersion === 'string' && req.httpVersion.startsWith('1.')
}

module.exports = function (http2, agent, { enabled }) {
  if (agent._conf.instrumentIncomingHTTPRequests) {
    agent.logger.debug('shimming http2.createServer function')
    shimmer.wrap(http2, 'createServer', wrapCreateServer)
    shimmer.wrap(http2, 'createSecureServer', wrapCreateServer)
  }

  if (!enabled) return http2
  var ins = agent._instrumentation
  agent.logger.debug('shimming http2.connect function')
  shimmer.wrap(http2, 'connect', wrapConnect)

  return http2

  // The `createServer` function will unpatch itself after patching
  // the first server prototype it patches.
  function wrapCreateServer (original) {
    return function wrappedCreateServer (options, handler) {
      var server = original.apply(this, arguments)
      shimmer.wrap(server.constructor.prototype, 'emit', wrapEmit)
      wrappedCreateServer[symbols.unwrap]()
      return server
    }
  }

  function wrapEmit (original) {
    var patched = false
    return function wrappedEmit (event, stream, headers) {
      if (event === 'stream') {
        if (!patched) {
          patched = true
          var proto = stream.constructor.prototype
          shimmer.wrap(proto, 'pushStream', wrapPushStream)
          shimmer.wrap(proto, 'respondWithFile', wrapRespondWith)
          shimmer.wrap(proto, 'respondWithFD', wrapRespondWith)
          shimmer.wrap(proto, 'respond', wrapHeaders)
          shimmer.wrap(proto, 'end', wrapEnd)
        }

        agent.logger.debug('intercepted stream event call to http2.Server.prototype.emit')

        const trans = agent.startTransaction()
        trans.type = 'request'
        // `trans.req` and `trans.res` are fake representations of Node.js's
        // core `http.IncomingMessage` and `http.ServerResponse` objects,
        // sufficient for `parsers.getContextFromRequest()` and
        // `parsers.getContextFromResponse()`, respectively.
        // `remoteAddress` is fetched now, rather than at stream end, because
        // this `stream.session.socket` is a proxy object that can throw
        // `ERR_HTTP2_SOCKET_UNBOUND` if the Http2Session has been destroyed.
        trans.req = {
          headers,
          socket: {
            remoteAddress: stream.session.socket.remoteAddress
          },
          method: headers[':method'],
          url: headers[':path'],
          httpVersion: '2.0'
        }
        trans.res = {
          statusCode: 200,
          headersSent: false,
          finished: false,
          headers: null
        }
        ins.bindEmitter(stream)

        eos(stream, function () {
          trans.end()
        })
      } else if (event === 'request' && reqIsHTTP1(stream)) {
        // http2.createSecureServer() supports a `allowHTTP1: true` option.
        // When true, an incoming client request that supports HTTP/1.x but not
        // HTTP/2 will be allowed. It will result in a 'request' event being
        // emitted. We wrap that here.
        //
        // Note that, a HTTP/2 request results in a 'stream' event (wrapped
        // above) *and* a 'request' event. We do not want to wrap the
        // compatibility 'request' event in this case. Hence the `reqIsHTTP1`
        // guard.
        const req = stream
        const res = headers

        agent.logger.debug('intercepted request event call to http2.Server.prototype.emit for %s', req.url)

        var traceparent = req.headers.traceparent || req.headers['elastic-apm-traceparent']
        var tracestate = req.headers.tracestate
        const trans = agent.startTransaction(null, null, {
          childOf: traceparent,
          tracestate: tracestate
        })
        trans.type = 'request'
        trans.req = req
        trans.res = res

        ins.bindEmitter(req)
        ins.bindEmitter(res)

        eos(res, function (err) {
          if (trans.ended) return
          if (!err) {
            trans.end()
            return
          }

          if (agent._conf.errorOnAbortedRequests) {
            var duration = trans._timer.elapsed()
            if (duration > (agent._conf.abortedErrorThreshold * 1000)) {
              agent.captureError('Socket closed with active HTTP request (>' + agent._conf.abortedErrorThreshold + ' sec)', {
                request: req,
                extra: { abortTime: duration }
              })
            }
          }

          // Handle case where res.end is called after an error occurred on the
          // stream (e.g. if the underlying socket was prematurely closed)
          const end = res.end
          res.end = function () {
            const result = end.apply(this, arguments)
            trans.end()
            return result
          }
        })
      }

      return original.apply(this, arguments)
    }
  }

  function updateHeaders (headers) {
    var trans = agent._instrumentation.currTransaction()
    if (trans && !trans.ended) {
      var status = headers[':status'] || 200
      trans.result = 'HTTP ' + status.toString()[0] + 'xx'
      trans.res.statusCode = status
      trans._setOutcomeFromHttpStatusCode(status)
      trans.res.headers = mergeHeaders(trans.res.headers, headers)
      trans.res.headersSent = true
    }
  }

  function wrapHeaders (original) {
    return function (headers) {
      updateHeaders(headers)
      return original.apply(this, arguments)
    }
  }

  function wrapRespondWith (original) {
    return function (_, headers) {
      updateHeaders(headers)
      return original.apply(this, arguments)
    }
  }

  function wrapEnd (original) {
    return function (headers) {
      var trans = agent._instrumentation.currTransaction()
      // `trans.res` might be removed, because before
      // https://github.com/nodejs/node/pull/20084 (e.g. in node v10.0.0) the
      // 'end' event could be called multiple times for the same Http2Stream,
      // and the `trans.res` ref is removed when the Transaction is ended.
      if (trans && trans.res) {
        trans.res.finished = true
      }
      return original.apply(this, arguments)
    }
  }

  function wrapPushStream (original) {
    return function wrappedPushStream (...args) {
      // Note: Break the run context so that the wrapped `stream.respond` et al
      // for this pushStream do not overwrite outer transaction state.
      var callback = args.pop()
      args.push(agent._instrumentation.bindFunctionToEmptyRunContext(callback))
      return original.apply(this, args)
    }
  }

  function mergeHeaders (source, target) {
    if (source === null) return target
    var result = Object.assign({}, target)
    var keys = Object.keys(source)
    for (let i = 0; i < keys.length; i++) {
      var key = keys[i]
      if (typeof target[key] === 'undefined') {
        result[key] = source[key]
      } else if (Array.isArray(target[key])) {
        result[key].push(source[key])
      } else {
        result[key] = [source[key]].concat(target[key])
      }
    }
    return result
  }

  function wrapConnect (orig) {
    return function (host) {
      const ret = orig.apply(this, arguments)
      shimmer.wrap(ret, 'request', orig => wrapRequest(orig, host))
      return ret
    }
  }

  function wrapRequest (orig, host) {
    return function (headers) {
      agent.logger.debug('intercepted call to http2.request')
      var method = headers[':method'] || 'GET'
      const span = ins.createSpan(null, 'external', 'http', method, { exitSpan: true })

      const parentRunContext = ins.currRunContext()
      var parent = span || parentRunContext.currSpan() || parentRunContext.currTransaction()
      if (parent) {
        const newHeaders = Object.assign({}, headers)
        parent.propagateTraceContextHeaders(newHeaders, function (carrier, name, value) {
          carrier[name] = value
        })
        arguments[0] = newHeaders
      }
      if (!span) {
        return orig.apply(this, arguments)
      }

      const spanRunContext = parentRunContext.enterSpan(span)
      var req = ins.withRunContext(spanRunContext, orig, this, ...arguments)

      ins.bindEmitterToRunContext(parentRunContext, req)

      var urlObj = parseUrl(headers[':path'])
      var path = urlObj.pathname
      var url = host + path
      span.name = method + ' ' + host

      var statusCode
      req.on('response', (headers) => {
        statusCode = headers[':status']
      })

      req.on('end', () => {
        agent.logger.debug('intercepted http2 client end event')

        span.setHttpContext({
          method,
          status_code: statusCode,
          url
        })
        span._setOutcomeFromHttpStatusCode(statusCode)

        // The `getHTTPDestination` function might throw in case an
        // invalid URL is given to the `URL()` function. Until we can
        // be 100% sure this doesn't happen, we better catch it here.
        // For details, see:
        // https://github.com/elastic/apm-agent-nodejs/issues/1769
        try {
          span._setDestinationContext(getHTTPDestination(url))
        } catch (e) {
          agent.logger.error('Could not set destination context: %s', e.message)
        }

        span.end()
      })

      return req
    }
  }
}
