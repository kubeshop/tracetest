/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

var semver = require('semver')
var clone = require('shallow-clone-shim')
var sqlSummary = require('sql-summary')

var { getDBDestination } = require('../context')

module.exports = function (tedious, agent, { version, enabled }) {
  if (!enabled) return tedious
  if (!semver.satisfies(version, '^1.9.0 || 2.x || 3.x || ^4.0.1 || 5.x || 6.x || 7.x || 8.x || 9.x || 10.x || 11.x || 12.x || 13.x || 14.x || 15.x')) {
    agent.logger.debug('tedious version %s not supported - aborting...', version)
    return tedious
  }

  const ins = agent._instrumentation

  return clone({}, tedious, {
    Connection (descriptor) {
      const getter = descriptor.get
      if (getter) {
        // tedious v6.5.0+
        descriptor.get = function get () {
          return wrapConnection(getter())
        }
      } else if (typeof descriptor.value === 'function') {
        descriptor.value = wrapConnection(descriptor.value)
      } else {
        agent.logger.debug('could not patch `tedious.Connection` property for tedious version %s - aborting...', version)
      }
      return descriptor
    },
    Request (descriptor) {
      const getter = descriptor.get
      if (getter) {
        // tedious v6.5.0+
        descriptor.get = function get () {
          return wrapRequest(getter())
        }
      } else if (typeof descriptor.value === 'function') {
        descriptor.value = wrapRequest(descriptor.value)
      } else {
        agent.logger.debug('could not patch `tedious.Request` property for tedious version %s - aborting...', version)
      }
      return descriptor
    }
  })

  function wrapRequest (OriginalRequest) {
    class Request extends OriginalRequest {
      constructor () {
        super(...arguments)
        ins.bindEmitter(this)
      }
    }

    return Request
  }

  function wrapConnection (OriginalConnection) {
    class Connection extends OriginalConnection {
      constructor () {
        super(...arguments)
        ins.bindEmitter(this)
      }

      makeRequest (request, _packetType, payload) {
        // if not a Request object (i.e. a BulkLoad), then bail
        if (!request.parametersByName) {
          return super.makeRequest(...arguments)
        }
        const span = ins.createSpan(null, 'db', 'mssql', 'query', { exitSpan: true })
        if (!span) {
          return super.makeRequest(...arguments)
        }

        let host, port, instanceName
        if (typeof this.config === 'object') {
          // http://tediousjs.github.io/tedious/api-connection.html#function_newConnection
          host = this.config.server
          if (this.config.options) {
            port = this.config.options.port
            instanceName = this.config.options.instanceName
          }
        }
        span._setDestinationContext(getDBDestination(host, port))

        let sql
        let preparing
        if (payload.parameters !== undefined) {
          // This looks for tedious instance with `RpcRequestPayload` started
          // since version >=v11.0.10, when RPC parameter handling was refactored
          // (https://github.com/tediousjs/tedious/pull/1275).
          preparing = payload.procedure === 'sp_prepare'
          const stmtParam = (payload.parameters.find(({ name }) => name === 'statement') ||
            payload.parameters.find(({ name }) => name === 'stmt'))
          sql = stmtParam ? stmtParam.value : request.sqlTextOrProcedure
        } else {
          preparing = request.sqlTextOrProcedure === 'sp_prepare'
          const params = request.parametersByName
          sql = (params.statement || params.stmt || {}).value
        }
        span.name = sqlSummary(sql) + (preparing ? ' (prepare)' : '')
        const dbContext = { type: 'sql', statement: sql }
        if (instanceName) {
          dbContext.instance = instanceName
        }
        span.setDbContext(dbContext)

        const origCallback = request.userCallback
        request.userCallback = ins.bindFunction(function tracedCallback () {
          // TODO: captureError and setOutcome on err first arg here
          span.end()
          if (origCallback) {
            return origCallback.apply(this, arguments)
          }
        })

        return super.makeRequest(...arguments)
      }
    }

    return Connection
  }
}
