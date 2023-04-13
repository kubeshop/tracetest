/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

var shimmer = require('../shimmer')

module.exports = function (expressQueue, agent, { enabled }) {
  if (!enabled) return expressQueue

  var ins = agent._instrumentation

  return function wrappedExpressQueue (config) {
    var result = expressQueue(config)
    shimmer.wrap(result.queue, 'createJob', function (original) {
      return function (job) {
        if (job.next) {
          job.next = ins.bindFunction(job.next)
        }
        return original.apply(this, arguments)
      }
    })
    return result
  }
}
