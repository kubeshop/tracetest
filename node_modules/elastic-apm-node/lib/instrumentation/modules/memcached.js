/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

var semver = require('semver')

var shimmer = require('../shimmer')
var { getDBDestination } = require('../context')

module.exports = function (memcached, agent, { version, enabled }) {
  if (!enabled) {
    return memcached
  }
  if (!semver.satisfies(version, '>=2.2.0')) {
    agent.logger.debug('Memcached version %s not supported - aborting...', version)
    return memcached
  }

  const ins = agent._instrumentation

  agent.logger.debug('shimming memcached.prototype.command')
  shimmer.wrap(memcached.prototype, 'command', wrapCommand)
  shimmer.wrap(memcached.prototype, 'connect', wrapConnect)
  return memcached

  function wrapConnect (original) {
    return function wrappedConnect () {
      const currentSpan = ins.currSpan()
      const server = arguments[0]
      agent.logger.debug('intercepted call to memcached.prototype.connect %o', { server })

      if (currentSpan) {
        const [host, port = 11211] = server.split(':')
        currentSpan._setDestinationContext(getDBDestination(host, port))
      }
      return original.apply(this, arguments)
    }
  }

  // Wrap the generic command that is used to build touch, get, gets etc
  function wrapCommand (original) {
    return function wrappedCommand (queryCompiler, _server) {
      if (typeof queryCompiler !== 'function') {
        return original.apply(this, arguments)
      }

      var query = queryCompiler()
      // Replace the queryCompiler function so it isn't called a second time.
      arguments[0] = function prerunQueryCompiler () {
        return query
      }

      // If the callback is not a function the user doesn't care about result.
      if (!query && typeof query.callback !== 'function') {
        return original.apply(this, arguments)
      }

      const span = ins.createSpan(`memcached.${query.type}`, 'db', 'memcached', query.type, { exitSpan: true })
      if (!span) {
        return original.apply(this, arguments)
      }

      agent.logger.debug('intercepted call to memcached.prototype.command %o', { id: span.id, type: query.type })
      span.setDbContext({ statement: `${query.type} ${query.key}`, type: 'memcached' })

      const spanRunContext = ins.currRunContext().enterSpan(span)
      const origCallback = query.callback
      query.callback = ins.bindFunctionToRunContext(spanRunContext, function tracedCallback () {
        span.end()
        return origCallback.apply(this, arguments)
      })
      return ins.withRunContext(spanRunContext, original, this, ...arguments)
    }
  }
}
