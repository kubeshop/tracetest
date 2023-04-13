/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

var semver = require('semver')
var sqlSummary = require('sql-summary')

var shimmer = require('../shimmer')
var symbols = require('../../symbols')
var { getDBDestination } = require('../context')

module.exports = function (mysql2, agent, { version, enabled }) {
  if (!enabled) {
    return mysql2
  }
  if (!semver.satisfies(version, '>=1 <3')) {
    agent.logger.debug('mysql2 version %s not supported - aborting...', version)
    return mysql2
  }

  var ins = agent._instrumentation

  shimmer.wrap(mysql2.Connection.prototype, 'query', wrapQuery)
  shimmer.wrap(mysql2.Connection.prototype, 'execute', wrapQuery)

  return mysql2

  function wrapQuery (original) {
    return function wrappedQuery (sql, values, cb) {
      agent.logger.debug('intercepted call to mysql2.%s', original.name)

      var span = ins.createSpan(null, 'db', 'mysql', 'query', { exitSpan: true })
      if (!span) {
        return original.apply(this, arguments)
      }

      if (this[symbols.knexStackObj]) {
        span.customStackTrace(this[symbols.knexStackObj])
        this[symbols.knexStackObj] = null
      }

      let hasCallback = false
      const wrapCallback = function (origCallback) {
        hasCallback = true
        return ins.bindFunction(function wrappedCallback (_err) {
          span.end()
          return origCallback.apply(this, arguments)
        })
      }

      let host, port, user, database
      if (typeof this.config === 'object') {
        ({ host, port, user, database } = this.config)
      }
      span._setDestinationContext(getDBDestination(host, port))

      let sqlStr
      switch (typeof sql) {
        case 'string':
          sqlStr = sql
          break
        case 'object':
          // `mysql2.{query,execute}` support the only arg being an instance
          // of the internal mysql2 `Command` object.
          if (typeof sql.onResult === 'function') {
            sql.onResult = wrapCallback(sql.onResult)
          }
          sqlStr = sql.sql
          break
        case 'function':
          arguments[0] = wrapCallback(sql)
          break
      }
      if (sqlStr) {
        span.setDbContext({ type: 'sql', instance: database, user, statement: sqlStr })
        span.name = sqlSummary(sqlStr)
      } else {
        span.setDbContext({ type: 'sql', instance: database, user })
      }

      if (typeof values === 'function') {
        arguments[1] = wrapCallback(values)
      } else if (typeof cb === 'function') {
        arguments[2] = wrapCallback(cb)
      }
      const spanRunContext = ins.currRunContext().enterSpan(span)
      const result = ins.withRunContext(spanRunContext, original, this, ...arguments)

      if (result && !hasCallback) {
        ins.bindEmitter(result)
        shimmer.wrap(result, 'emit', function (origEmit) {
          return function (event) {
            switch (event) {
              case 'error':
              case 'close':
              case 'end':
                span.end()
            }
            return origEmit.apply(this, arguments)
          }
        })
      }

      return result
    }
  }
}
