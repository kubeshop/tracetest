/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

// Instrumentation of the 'ioredis' package:
// https://github.com/luin/ioredis
// https://github.com/luin/ioredis/blob/master/API.md

const semver = require('semver')

const constants = require('../../constants')
const { getDBDestination } = require('../context')
const shimmer = require('../shimmer')

const TYPE = 'db'
const SUBTYPE = 'redis'
const ACTION = 'query'
const hasIoredisSpanSym = Symbol('ElasticAPMHasIoredisSpan')

module.exports = function (ioredis, agent, { version, enabled }) {
  if (!enabled) {
    return ioredis
  }
  if (!semver.satisfies(version, '>=2.0.0 <6.0.0')) {
    agent.logger.debug('ioredis version %s not supported - aborting...', version)
    return ioredis
  }

  const ins = agent._instrumentation

  agent.logger.debug('shimming ioredis.prototype.sendCommand')
  shimmer.wrap(ioredis.prototype, 'sendCommand', wrapSendCommand)
  return ioredis

  function wrapSendCommand (origSendCommand) {
    return function wrappedSendCommand (command) {
      if (!command || !command.name || !command.promise) {
        // Doesn't look like an ioredis.Command, skip instrumenting.
        return origSendCommand.apply(this, arguments)
      }
      if (command[hasIoredisSpanSym]) {
        // Avoid double-instrumenting a command when ioredis *re*-calls
        // sendCommand for queued commands when "ready".
        return origSendCommand.apply(this, arguments)
      }

      agent.logger.debug({ command: command.name }, 'intercepted call to ioredis.prototype.sendCommand')
      const span = ins.createSpan(command.name.toUpperCase(), TYPE, SUBTYPE, ACTION, { exitSpan: true })
      if (!span) {
        return origSendCommand.apply(this, arguments)
      }

      command[hasIoredisSpanSym] = true

      const options = this.options || {} // `this` is the `Redis` client.
      span._setDestinationContext(getDBDestination(options.host, options.port))
      span.setDbContext({ type: 'redis' })

      const spanRunContext = ins.currRunContext().enterSpan(span)
      command.promise.then(
        () => {
          span.end()
        },
        ins.bindFunctionToRunContext(spanRunContext, (err) => {
          span._setOutcomeFromErrorCapture(constants.OUTCOME_FAILURE)
          agent.captureError(err, { skipOutcome: true })
          span.end()
        })
      )
      return ins.withRunContext(spanRunContext, origSendCommand, this, ...arguments)
    }
  }
}
