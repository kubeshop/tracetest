/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

const EventEmitter = require('events')

var semver = require('semver')
var sqlSummary = require('sql-summary')

var shimmer = require('../shimmer')
var symbols = require('../../symbols')
var { getDBDestination } = require('../context')

module.exports = function (pg, agent, { version, enabled }) {
  if (!enabled) {
    return pg
  }
  if (!semver.satisfies(version, '>=4.0.0 <9.0.0')) {
    agent.logger.debug('pg version %s not supported - aborting...', version)
    return pg
  }

  patchClient(pg.Client, 'pg.Client', agent)

  // Trying to access the pg.native getter will trigger and log the warning
  // "Cannot find module 'pg-native'" to STDERR if the module isn't installed.
  // Overwriting the getter we can lazily patch the native client only if the
  // user is acually requesting it.
  var getter = pg.__lookupGetter__('native')
  if (getter) {
    delete pg.native
    // To be as true to the original pg module as possible, we use
    // __defineGetter__ instead of Object.defineProperty.
    pg.__defineGetter__('native', function () {
      var native = getter()
      if (native && native.Client) {
        patchClient(native.Client, 'pg.native.Client', agent)
      }
      return native
    })
  }

  return pg
}

function patchClient (Client, klass, agent) {
  agent.logger.debug('shimming %s.prototype.query', klass)
  shimmer.wrap(Client.prototype, 'query', wrapQuery)

  function wrapQuery (orig, name) {
    return function wrappedFunction (sql) {
      agent.logger.debug('intercepted call to %s.prototype.%s', klass, name)
      const ins = agent._instrumentation
      const span = ins.createSpan('SQL', 'db', 'postgresql', 'query', { exitSpan: true })
      if (!span) {
        return orig.apply(this, arguments)
      }

      // Get connection parameters from Client.
      let host, port, database, user
      if (typeof this.connectionParameters === 'object') {
        ({ host, port, database, user } = this.connectionParameters)
      }
      span._setDestinationContext(getDBDestination(host, port))

      const dbContext = { type: 'sql' }
      let sqlText = sql
      if (sql && typeof sql.text === 'string') {
        sqlText = sql.text
      }
      if (typeof sqlText === 'string') {
        span.name = sqlSummary(sqlText)
        dbContext.statement = sqlText
      } else {
        agent.logger.debug('unable to parse sql form pg module (type: %s)', typeof sqlText)
      }
      if (database) {
        dbContext.instance = database
      }
      if (user) {
        dbContext.user = user
      }
      span.setDbContext(dbContext)

      if (this[symbols.knexStackObj]) {
        span.customStackTrace(this[symbols.knexStackObj])
        this[symbols.knexStackObj] = null
      }

      let index = arguments.length - 1
      let cb = arguments[index]
      if (Array.isArray(cb)) {
        index = cb.length - 1
        cb = cb[index]
      }

      const spanRunContext = ins.currRunContext().enterSpan(span)
      const onQueryEnd = ins.bindFunctionToRunContext(spanRunContext, (_err) => {
        agent.logger.debug('intercepted end of %s.prototype.%s', klass, name)
        span.end()
      })

      if (typeof cb === 'function') {
        arguments[index] = ins.bindFunction((err, res) => {
          onQueryEnd(err)
          return cb(err, res)
        })
        return orig.apply(this, arguments)
      } else {
        var queryOrPromise = orig.apply(this, arguments)

        // It is important to prefer `.on` to `.then` for pg <7 >=6.3.0, because
        // `query.then` is broken in those versions. See
        // https://github.com/brianc/node-postgres/commit/b5b49eb895727e01290e90d08292c0d61ab86322#r23267714
        if (typeof queryOrPromise.on === 'function') {
          queryOrPromise.on('end', onQueryEnd)
          queryOrPromise.on('error', onQueryEnd)
          if (queryOrPromise instanceof EventEmitter) {
            ins.bindEmitter(queryOrPromise)
          }
        } else if (typeof queryOrPromise.then === 'function') {
          queryOrPromise.then(
            () => { onQueryEnd() },
            onQueryEnd
          )
        } else {
          agent.logger.debug('ERROR: unknown pg query type: %s', typeof queryOrPromise)
        }

        return queryOrPromise
      }
    }
  }
}
