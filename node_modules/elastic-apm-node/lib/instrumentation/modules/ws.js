/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

var semver = require('semver')

var shimmer = require('../shimmer')

module.exports = function (ws, agent, { version, enabled }) {
  if (!enabled) return ws
  if (!semver.satisfies(version, '>=1 <8')) {
    agent.logger.debug('ws version %s not supported - aborting...', version)
    return ws
  }

  const ins = agent._instrumentation

  agent.logger.debug('shimming ws.prototype.send function')
  shimmer.wrap(ws.prototype, 'send', wrapSend)

  return ws

  function wrapSend (orig) {
    return function wrappedSend () {
      agent.logger.debug('intercepted call to ws.prototype.send')
      const span = ins.createSpan('Send WebSocket Message', 'websocket', 'send', { exitSpan: true })
      if (!span) {
        return orig.apply(this, arguments)
      }

      const args = Array.prototype.slice.call(arguments)
      let cb = args[args.length - 1]
      const onDone = function () {
        span.end()
        if (cb) {
          cb.apply(this, arguments)
        }
      }
      if (typeof cb === 'function') {
        args[args.length - 1] = ins.bindFunction(onDone)
      } else {
        cb = null
        args.push(ins.bindFunction(onDone))
      }

      const spanRunContext = ins.currRunContext().enterSpan(span)
      return ins.withRunContext(spanRunContext, orig, this, ...args)
    }
  }
}
