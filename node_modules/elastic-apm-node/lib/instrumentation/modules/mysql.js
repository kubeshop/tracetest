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

module.exports = function (mysql, agent, { version, enabled }) {
  if (!enabled) {
    return mysql
  }
  if (!semver.satisfies(version, '^2.0.0')) {
    agent.logger.debug('mysql version %s not supported - aborting...', version)
    return mysql
  }

  agent.logger.debug('shimming mysql.createPool')
  shimmer.wrap(mysql, 'createPool', wrapCreatePool)

  agent.logger.debug('shimming mysql.createPoolCluster')
  shimmer.wrap(mysql, 'createPoolCluster', wrapCreatePoolCluster)

  agent.logger.debug('shimming mysql.createConnection')
  shimmer.wrap(mysql, 'createConnection', wrapCreateConnection)

  return mysql

  function wrapCreateConnection (original) {
    return function wrappedCreateConnection () {
      var connection = original.apply(this, arguments)

      wrapQueryable(connection, 'connection', agent)

      return connection
    }
  }

  function wrapCreatePool (original) {
    return function wrappedCreatePool () {
      var pool = original.apply(this, arguments)

      agent.logger.debug('shimming mysql pool.getConnection')
      shimmer.wrap(pool, 'getConnection', wrapGetConnection)

      return pool
    }
  }

  function wrapCreatePoolCluster (original) {
    return function wrappedCreatePoolCluster () {
      var cluster = original.apply(this, arguments)

      agent.logger.debug('shimming mysql cluster.of')
      shimmer.wrap(cluster, 'of', function wrapOf (original) {
        return function wrappedOf () {
          var ofCluster = original.apply(this, arguments)

          agent.logger.debug('shimming mysql cluster of.getConnection')
          shimmer.wrap(ofCluster, 'getConnection', wrapGetConnection)

          return ofCluster
        }
      })

      return cluster
    }
  }

  function wrapGetConnection (original) {
    return function wrappedGetConnection () {
      var cb = arguments[0]

      if (typeof cb === 'function') {
        arguments[0] = agent._instrumentation.bindFunction(function wrapedCallback (err, connection) { // eslint-disable-line handle-callback-err
          if (connection) wrapQueryable(connection, 'getConnection() > connection', agent)
          return cb.apply(this, arguments)
        })
      }

      return original.apply(this, arguments)
    }
  }
}

function wrapQueryable (connection, objType, agent) {
  const ins = agent._instrumentation

  agent.logger.debug('shimming mysql %s.query', objType)
  shimmer.wrap(connection, 'query', wrapQuery)

  let host, port, user, database
  if (typeof connection.config === 'object') {
    ({ host, port, user, database } = connection.config)
  }

  function wrapQuery (original) {
    return function wrappedQuery (sql, values, cb) {
      agent.logger.debug('intercepted call to mysql %s.query', objType)

      var span = ins.createSpan(null, 'db', 'mysql', 'query', { exitSpan: true })
      if (!span) {
        return original.apply(this, arguments)
      }

      var hasCallback = false
      var sqlStr

      if (this[symbols.knexStackObj]) {
        span.customStackTrace(this[symbols.knexStackObj])
        this[symbols.knexStackObj] = null
      }

      const wrapCallback = function (origCallback) {
        hasCallback = true
        return ins.bindFunction(function wrappedCallback (_err) {
          span.end()
          return origCallback.apply(this, arguments)
        })
      }

      switch (typeof sql) {
        case 'string':
          sqlStr = sql
          break
        case 'object':
          if (typeof sql._callback === 'function') {
            sql._callback = wrapCallback(sql._callback)
          }
          sqlStr = sql.sql
          break
        case 'function':
          arguments[0] = wrapCallback(sql)
          break
      }

      if (sqlStr) {
        agent.logger.debug({ sql: sqlStr }, 'extracted sql from mysql query')
        span.setDbContext({ statement: sqlStr, type: 'sql', user, instance: database })
        span.name = sqlSummary(sqlStr)
      }
      span._setDestinationContext(getDBDestination(host, port))

      if (typeof values === 'function') {
        arguments[1] = wrapCallback(values)
      } else if (typeof cb === 'function') {
        arguments[2] = wrapCallback(cb)
      }

      const spanRunContext = ins.currRunContext().enterSpan(span)
      const result = ins.withRunContext(spanRunContext, original, this, ...arguments)

      if (!hasCallback && result instanceof EventEmitter) {
        // Wrap `result.emit` instead of `result.once('error', ...)` to avoid
        // changing app behaviour by possibly setting the only 'error' handler.
        shimmer.wrap(result, 'emit', function (origEmit) {
          return function wrappedEmit (event, data) {
            // The 'mysql' module emits 'end' even after an 'error' event.
            switch (event) {
              case 'error':
                break
              case 'end':
                span.end()
                break
            }
            return origEmit.apply(this, arguments)
          }
        })
        // Ensure event handlers execute in the caller run context.
        ins.bindEmitter(result)
      }

      return result
    }
  }
}
