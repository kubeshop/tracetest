/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

// Instrument the undici module.
//
// This uses undici's diagnostics_channel support for instrumentation.
//    https://github.com/nodejs/undici/blob/main/docs/api/DiagnosticsChannel.md
// Undici is also used for Node >=v18.0.0's `fetch()` implementation, via
// an esbuild bundle. This instrumentation is enabled if either `global.fetch`
// is present or `require('undici')`.
//
// Limitations:
// - Currently this isn't subscribing to 'undici:client:...' messages.
//   With typical undici usage a connection will only be initiated for a
//   request. However, if a user manually does `client.connect(...)` then it is
//   possible for this instrumentation to miss a connection error from
//   'undici:client:connectError'. It would eventually be nice to heuristically
//   add 'connect' spans as children of request spans.
// - This doesn't instrument HTTP CONNECT, as exposed by `undici.connect(...)`.
//   I don't think the current undici diagnostics_channel messages provide a
//   way to watch the completion of the CONNECT request.
// - This hasn't been tested with `undici.upgrade()`.
//
// Some notes on if/when we want to collect some HTTP client metrics:
// - The time between 'undici:client:connected' and 'undici:client:sendHeaders'
//   could be a measure of client-side latency. I'm not sure if client-side
//   queueing of requests would show a time gap here.
// - The time between 'undici:client:sendHeaders' and 'undici:client:bodySent'
//   might be interesting for large bodies, or perhaps for streaming requests.
// - The time between 'undici:client:bodySent' and 'undici:request:headers'
//   could be a measure of response TTFB latency.

let diagch = null
try {
  diagch = require('diagnostics_channel')
} catch (_importErr) {
  // pass
}

const semver = require('semver')

let isInstrumented = false
let spanFromReq = null
let chans = null

// Get the content-length from undici response headers.
// `headers` is an Array of buffers: [k, v, k, v, ...].
// If the header is not present, or has an invalid value, this returns null.
function contentLengthFromResponseHeaders (headers) {
  const name = 'content-length'
  for (let i = 0; i < headers.length; i += 2) {
    const k = headers[i]
    if (k.length === name.length && k.toString().toLowerCase() === name) {
      const v = Number(headers[i + 1])
      if (!isNaN(v)) {
        return v
      } else {
        return null
      }
    }
  }
  return null
}

function uninstrumentUndici () {
  if (!isInstrumented) {
    return
  }
  isInstrumented = false

  spanFromReq = null
  chans.forEach(({ chan, onMessage }) => {
    chan.unsubscribe(onMessage)
  })
  chans = null
}

/**
 * Setup instrumentation for undici. The instrumentation is based entirely on
 * diagnostics_channel usage, so no reference to the loaded undici module is
 * required.
 */
function instrumentUndici (agent) {
  if (isInstrumented) {
    return
  }
  isInstrumented = true

  const ins = agent._instrumentation
  spanFromReq = new WeakMap()

  // Keep ref to avoid https://github.com/nodejs/node/issues/42170 bug and for
  // unsubscribing.
  chans = []
  function diagchSub (name, onMessage) {
    const chan = diagch.channel(name)
    chan.subscribe(onMessage)
    chans.push({
      name,
      chan,
      onMessage
    })
  }

  diagchSub('undici:request:create', ({ request }) => {
    // We do not handle instrumenting HTTP CONNECT. See limitation notes above.
    if (request.method === 'CONNECT') {
      return
    }

    const url = new URL(request.origin)
    const span = ins.createSpan(`${request.method} ${url.host}`, 'external', 'http', request.method, { exitSpan: true })

    // W3C trace-context propagation.
    // If the span is null (e.g. hit `transactionMaxSpans`, unsampled
    // transaction), then fallback to the current run context's span or
    // transaction, if any.
    const parentRunContext = ins.currRunContext()
    const propSpan = span || parentRunContext.currSpan() || parentRunContext.currTransaction()
    if (propSpan) {
      propSpan.propagateTraceContextHeaders(request, function (req, name, value) {
        req.addHeader(name, value)
      })
    }

    if (span) {
      spanFromReq.set(request, span)

      // Set some initial HTTP context, in case the request errors out before a response.
      span.setHttpContext({
        method: request.method,
        url: request.origin + request.path
      })

      const destContext = {
        address: url.hostname
      }
      const port = Number(url.port) || (url.protocol === 'https:' && 443) || (url.protocol === 'http:' && 80)
      if (port) {
        destContext.port = port
      }
      span._setDestinationContext(destContext)
    }
  })

  diagchSub('undici:request:headers', ({ request, response }) => {
    const span = spanFromReq.get(request)
    if (span !== undefined) {
      // We are currently *not* capturing response headers, even though the
      // intake API does allow it, because none of the other `setHttpContext`
      // uses currently do.

      const httpContext = {
        method: request.method,
        status_code: response.statusCode,
        url: request.origin + request.path
      }
      const cLen = contentLengthFromResponseHeaders(response.headers)
      if (cLen !== null) {
        httpContext.response = { encoded_body_size: cLen }
      }
      span.setHttpContext(httpContext)

      span._setOutcomeFromHttpStatusCode(response.statusCode)
    }
  })

  diagchSub('undici:request:trailers', ({ request }) => {
    const span = spanFromReq.get(request)
    if (span !== undefined) {
      span.end()
      spanFromReq.delete(request)
    }
  })

  diagchSub('undici:request:error', ({ request, error }) => {
    const span = spanFromReq.get(request)
    const errOpts = {}
    if (span !== undefined) {
      errOpts.parent = span
      // Cases where we won't have an undici parent span:
      // - We've hit transactionMaxSpans.
      // - The undici HTTP span was suppressed because it is a child of an
      //   exit span (e.g. when used as the transport for the Elasticsearch
      //   client).
      // It might be debatable whether we want to capture the error in the
      // latter case. This could be revisited later.
    }
    agent.captureError(error, errOpts)
    if (span !== undefined) {
      span.end()
      spanFromReq.delete(request)
    }
  })
}

function shimUndici (undici, agent, { version, enabled }) {
  if (!enabled) {
    return undici
  }
  if (semver.lt(version, '4.7.1')) {
    // Undici added its diagnostics_channel messages in v4.7.0. In v4.7.1 the
    // `request.origin` property, that we need, was added.
    agent.logger.debug('cannot instrument undici: undici version %s is not supported', version)
    return undici
  }
  if (!diagch) {
    agent.logger.debug('cannot instrument undici: there is no "diagnostics_channel" module', process.version)
    return undici
  }

  instrumentUndici(agent)
  return undici
}

module.exports = shimUndici
module.exports.instrumentUndici = instrumentUndici
module.exports.uninstrumentUndici = uninstrumentUndici
