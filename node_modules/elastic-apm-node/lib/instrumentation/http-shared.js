/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

var url = require('url')
var endOfStream = require('end-of-stream')

var { parseUrl } = require('../parsers')
var { getHTTPDestination } = require('./context')

const transactionForResponse = new WeakMap()
exports.transactionForResponse = transactionForResponse

exports.instrumentRequest = function (agent, moduleName) {
  var ins = agent._instrumentation
  return function (orig) {
    return function (event, req, res) {
      if (event === 'request') {
        agent.logger.debug('intercepted request event call to %s.Server.prototype.emit for %s', moduleName, req.url)

        if (shouldIgnoreRequest(agent, req)) {
          agent.logger.debug('ignoring request to %s', req.url)
          // Don't leak previous transaction.
          agent._instrumentation.supersedeWithEmptyRunContext()
        } else {
          // Decide whether to use trace-context headers, if any, for a
          // distributed trace.
          const traceparent = req.headers.traceparent || req.headers['elastic-apm-traceparent']
          const tracestate = req.headers.tracestate
          const trans = agent.startTransaction(null, null, {
            childOf: traceparent,
            tracestate: tracestate
          })
          trans.type = 'request'
          trans.req = req
          trans.res = res

          transactionForResponse.set(res, trans)

          ins.bindEmitter(req)
          ins.bindEmitter(res)

          endOfStream(res, function (err) {
            if (trans.ended) return
            if (!err) return trans.end()

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
      }

      return orig.apply(this, arguments)
    }
  }
}

function shouldIgnoreRequest (agent, req) {
  var i

  for (i = 0; i < agent._conf.ignoreUrlStr.length; i++) {
    if (agent._conf.ignoreUrlStr[i] === req.url) return true
  }
  for (i = 0; i < agent._conf.ignoreUrlRegExp.length; i++) {
    if (agent._conf.ignoreUrlRegExp[i].test(req.url)) return true
  }
  for (i = 0; i < agent._conf.transactionIgnoreUrlRegExp.length; i++) {
    if (agent._conf.transactionIgnoreUrlRegExp[i].test(req.url)) return true
  }

  var ua = req.headers['user-agent']
  if (!ua) return false

  for (i = 0; i < agent._conf.ignoreUserAgentStr.length; i++) {
    if (ua.indexOf(agent._conf.ignoreUserAgentStr[i]) === 0) return true
  }
  for (i = 0; i < agent._conf.ignoreUserAgentRegExp.length; i++) {
    if (agent._conf.ignoreUserAgentRegExp[i].test(ua)) return true
  }

  return false
}

function formatURL (item) {
  return {
    href: item.href,
    pathname: item.pathname,
    path: item.pathname + (item.search || ''),
    protocol: item.protocol,
    host: item.host,
    port: item.port,
    hostname: item.hostname,
    hash: item.hash,
    search: item.search
  }
}

// NOTE: This will also stringify and parse URL instances
// to a format which can be mixed into the options object.
function ensureUrl (v) {
  if (typeof v === 'string') {
    return formatURL(parseUrl(v))
  } else if (url.URL && v instanceof url.URL) {
    return formatURL(v)
  } else {
    return v
  }
}

function getSafeHost (res) {
  return res.getHeader ? res.getHeader('Host') : res._headers.host
}

exports.traceOutgoingRequest = function (agent, moduleName, method) {
  var ins = agent._instrumentation

  return function (orig) {
    return function (...args) {
      const parentRunContext = ins.currRunContext()
      var span = ins.createSpan(null, 'external', 'http', { exitSpan: true })
      var id = span && span.transaction.id

      agent.logger.debug('intercepted call to %s.%s %o', moduleName, method, { id: id })

      var options = {}
      var newArgs = [options]
      for (const arg of args) {
        if (typeof arg === 'function') {
          newArgs.push(ins.bindFunctionToRunContext(parentRunContext, arg))
        } else {
          Object.assign(options, ensureUrl(arg))
        }
      }

      if (!options.headers) options.headers = {}

      // W3C trace-context propagation.
      // There are a number of reasons why `span` might be null: child of an
      // exit span, `transactionMaxSpans` was hit, unsampled transaction, etc.
      // If so, then fallback to the current run context's span or transaction,
      // if any.
      var parent = span || parentRunContext.currSpan() || parentRunContext.currTransaction()

      const newHeaders = Object.assign({}, options.headers)
      if (parent) {
        parent.propagateTraceContextHeaders(newHeaders, function (carrier, name, value) {
          carrier[name] = value
        })
      }
      options.headers = newHeaders
      if (!span) {
        return orig.apply(this, newArgs)
      }

      const spanRunContext = parentRunContext.enterSpan(span)
      var req = ins.withRunContext(spanRunContext, orig, this, ...newArgs)

      var protocol = req.agent && req.agent.protocol
      agent.logger.debug('request details: %o', { protocol: protocol, host: getSafeHost(req), id: id })

      ins.bindEmitter(req)

      span.action = req.method
      span.name = req.method + ' ' + getSafeHost(req)

      // TODO: Research if it's possible to add this to the prototype instead.
      // Or if it's somehow preferable to listen for when a `response` listener
      // is added instead of when `response` is emitted.
      const emit = req.emit
      req.emit = function wrappedEmit (type, res) {
        if (type === 'response') onResponse(res)
        if (type === 'abort') onAbort(type)
        return emit.apply(req, arguments)
      }

      const url = getUrlFromRequestAndOptions(req, options, moduleName + ':')
      if (!url) {
        agent.logger.warn('unable to identify http.ClientRequest url %o', { id: id })
      }
      let statusCode
      return req

      // In case the request is ended prematurely
      function onAbort (type) {
        if (span.ended) return
        agent.logger.debug('intercepted http.ClientRequest abort event %o', { id: id })

        onEnd()
      }

      function onEnd () {
        span.setHttpContext({
          method: req.method,
          status_code: statusCode,
          url
        })

        // Add destination info only when socket conn is established
        if (url) {
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
        }

        span._setOutcomeFromHttpStatusCode(statusCode)
        span.end()
      }

      function onResponse (res) {
        agent.logger.debug('intercepted http.ClientRequest response event %o', { id: id })
        ins.bindEmitterToRunContext(parentRunContext, res)
        statusCode = res.statusCode
        res.prependListener('end', function () {
          agent.logger.debug('intercepted http.IncomingMessage end event %o', { id: id })
          onEnd()
        })
      }
    }
  }
}

// Creates a sanitized URL suitable for the span's HTTP context
//
// This function reconstructs a URL using the request object's properties
// where it can (node versions v14.5.0, v12.19.0 and later) , and falling
// back to the options where it can not.  This function also strips any
// authentication information provided with the hostname.  In other words
//
// http://username:password@example.com/foo
//
// becomes http://example.com/foo
//
// NOTE: The options argument may not be the same options that are passed
// to http.request if the caller uses the the http.request(url,options,...)
// method signature. The agent normalizes the url and options into a single
// options array. This function expects those pre-normalized options.
//
// @param {ClientRequest} req
// @param {object} options
// @param {string} fallbackProtocol
// @return string|undefined
function getUrlFromRequestAndOptions (req, options, fallbackProtocol) {
  if (!req) {
    return undefined
  }
  options = options || {}
  req = req || {}
  req.agent = req.agent || {}
  if (isProxiedRequest(req)) {
    return req.path
  }

  const port = options.port ? `:${options.port}` : ''
  // req.host and req.protocol are node versions v14.5.0/v12.19.0 and later
  const host = req.host || options.hostname || options.host || 'localhost'
  const protocol = req.protocol || req.agent.protocol || fallbackProtocol

  return `${protocol}//${host}${port}${req.path}`
}

function isProxiedRequest (req) {
  return req.path.indexOf('https:') === 0 || req.path.indexOf('http:') === 0
}

exports.getUrlFromRequestAndOptions = getUrlFromRequestAndOptions
