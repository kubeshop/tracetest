/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

var httpShared = require('../http-shared')
var shimmer = require('../shimmer')

function getSafeHeaders (res) {
  return res.getHeaders ? res.getHeaders() : res._headers
}

module.exports = function (http, agent, { enabled }) {
  if (agent._conf.instrumentIncomingHTTPRequests) {
    agent.logger.debug('shimming http.Server.prototype.emit function')
    shimmer.wrap(http && http.Server && http.Server.prototype, 'emit', httpShared.instrumentRequest(agent, 'http'))

    agent.logger.debug('shimming http.ServerResponse.prototype.writeHead function')
    shimmer.wrap(http && http.ServerResponse && http.ServerResponse.prototype, 'writeHead', wrapWriteHead)
  }

  if (!enabled) return http

  agent.logger.debug('shimming http.request function')
  shimmer.wrap(http, 'request', httpShared.traceOutgoingRequest(agent, 'http', 'request'))

  agent.logger.debug('shimming http.get function')
  shimmer.wrap(http, 'get', httpShared.traceOutgoingRequest(agent, 'http', 'get'))

  return http

  function wrapWriteHead (original) {
    return function wrappedWriteHead () {
      var headers = arguments.length === 1
        ? getSafeHeaders(this) // might be because of implicit headers.
        : arguments[arguments.length - 1]

      var result = original.apply(this, arguments)

      var trans = httpShared.transactionForResponse.get(this)
      if (trans) {
        httpShared.transactionForResponse.delete(this)

        // It shouldn't be possible for the statusCode to be falsy, but just in
        // case we're in a bad state we should avoid throwing
        trans.result = 'HTTP ' + (this.statusCode || '').toString()[0] + 'xx'
        trans._setOutcomeFromHttpStatusCode(this.statusCode)

        // End transacton early in case of SSE
        if (headers && typeof headers === 'object' && !Array.isArray(headers)) {
          Object.keys(headers).some(function (key) {
            if (key.toLowerCase() !== 'content-type') return false
            if (String(headers[key]).toLowerCase().indexOf('text/event-stream') !== 0) return false
            agent.logger.debug('detected SSE response - ending transaction %o', { id: trans.id })
            agent.endTransaction()
            return true
          })
        }
      }

      return result
    }
  }
}
