/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

const Instrumentation = require('../index')
const { getLambdaHandlerInfo } = require('../../lambda')
const propwrap = require('../../propwrap')

module.exports = function (module, agent, { version, enabled }) {
  if (!enabled) {
    return module
  }

  const { field } = getLambdaHandlerInfo(process.env, Instrumentation.modules, agent.logger)
  try {
    const newMod = propwrap.wrap(module, field, (orig) => {
      return agent.lambda(orig)
    })
    return newMod
  } catch (wrapErr) {
    agent.logger.warn('could not wrap lambda handler: %s', wrapErr)
    return module
  }
}
