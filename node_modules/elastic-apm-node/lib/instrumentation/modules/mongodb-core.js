/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

var semver = require('semver')

const { getDBDestination } = require('../context')
var shimmer = require('../shimmer')

var SERVER_FNS = ['insert', 'update', 'remove', 'auth']
var CURSOR_FNS_FIRST = ['next', '_getmore']

const firstSpan = Symbol('first-span')

module.exports = function (mongodb, agent, { version, enabled }) {
  if (!enabled) return mongodb
  if (!semver.satisfies(version, '>=1.2.19 <4')) {
    agent.logger.debug('mongodb-core version %s not supported - aborting...', version)
    return mongodb
  }

  const ins = agent._instrumentation

  if (mongodb.Server) {
    agent.logger.debug('shimming mongodb-core.Server.prototype.command')
    shimmer.wrap(mongodb.Server.prototype, 'command', wrapCommand)
    agent.logger.debug('shimming mongodb-core.Server.prototype functions: %j', SERVER_FNS)
    shimmer.massWrap(mongodb.Server.prototype, SERVER_FNS, wrapQuery)
  }

  if (mongodb.Cursor) {
    agent.logger.debug('shimming mongodb-core.Cursor.prototype functions: %j', CURSOR_FNS_FIRST)
    shimmer.massWrap(mongodb.Cursor.prototype, CURSOR_FNS_FIRST, wrapCursor)
  }

  return mongodb

  function wrapCommand (orig) {
    return function wrappedFunction (ns, cmd) {
      var trans = agent._instrumentation.currTransaction()
      var id = trans && trans.id
      var span

      agent.logger.debug('intercepted call to mongodb-core.Server.prototype.command %o', { id: id, ns: ns })

      if (trans && arguments.length > 0) {
        var index = arguments.length - 1
        var cb = arguments[index]
        if (typeof cb === 'function') {
          var type
          if (cmd.findAndModify) type = 'findAndModify'
          else if (cmd.createIndexes) type = 'createIndexes'
          else if (cmd.ismaster) type = 'ismaster'
          else if (cmd.count) type = 'count'
          else type = 'command'

          span = ins.createSpan(ns + '.' + type, 'db', 'mongodb', 'query', { exitSpan: true })
          if (span) {
            span.setDbContext({ type: 'mongodb', instance: ns })
            arguments[index] = ins.bindFunctionToRunContext(ins.currRunContext(), wrappedCallback)
          }
        }
      }

      return orig.apply(this, arguments)

      function wrappedCallback (_err, commandResult) {
        agent.logger.debug('intercepted mongodb-core.Server.prototype.command callback %o', { id: id })
        if (commandResult && commandResult.connection) {
          span._setDestinationContext(getDBDestination(
            commandResult.connection.host, commandResult.connection.port))
        }
        span.end()
        return cb.apply(this, arguments)
      }
    }
  }

  function wrapQuery (orig, name) {
    return function wrappedFunction (ns) {
      var trans = agent._instrumentation.currTransaction()
      var id = trans && trans.id
      var span

      agent.logger.debug('intercepted call to mongodb-core.Server.prototype.%s %o', name, { id: id, ns: ns })

      if (trans && arguments.length > 0) {
        var index = arguments.length - 1
        var cb = arguments[index]
        if (typeof cb === 'function') {
          span = ins.createSpan(ns + '.' + name, 'db', 'mongodb', 'query', { exitSpan: true })
          if (span) {
            span.setDbContext({ type: 'mongodb', instance: ns })
            arguments[index] = ins.bindFunctionToRunContext(ins.currRunContext(), wrappedCallback)
          }
        }
      }

      return orig.apply(this, arguments)

      function wrappedCallback (_err, commandResult) {
        agent.logger.debug('intercepted mongodb-core.Server.prototype.%s callback %o', name, { id: id })
        if (commandResult && commandResult.connection) {
          span._setDestinationContext(getDBDestination(
            commandResult.connection.host, commandResult.connection.port))
        }
        span.end()
        return cb.apply(this, arguments)
      }
    }
  }
  function wrapCursor (orig, name) {
    return function wrappedFunction () {
      var trans = agent._instrumentation.currTransaction()
      var id = trans && trans.id
      var span

      agent.logger.debug('intercepted call to mongodb-core.Cursor.prototype.%s %o', name, { id: id })

      if (trans && arguments.length > 0) {
        var cb = arguments[0]
        if (typeof cb === 'function') {
          if (name !== 'next' || !this[firstSpan]) {
            var spanName = `${this.ns}.${this.cmd.find ? 'find' : name}`
            span = ins.createSpan(spanName, 'db', 'mongodb', 'query', { exitSpan: true })
          }
          if (span) {
            span.setDbContext({ type: 'mongodb', instance: this.ns })
            // Limitation: Currently not getting destination address/port for
            // cursor calls.
            arguments[0] = ins.bindFunctionToRunContext(ins.currRunContext(), wrappedCallback)
            if (name === 'next') {
              this[firstSpan] = true
            }
          }
        }
      }

      return orig.apply(this, arguments)

      function wrappedCallback () {
        agent.logger.debug('intercepted mongodb-core.Cursor.prototype.%s callback %o', name, { id: id })
        span.end()
        return cb.apply(this, arguments)
      }
    }
  }
}
