/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

var semver = require('semver')

var shimmer = require('../shimmer')

module.exports = function (Router, agent, { version, enabled }) {
  if (!enabled) return Router
  if (!semver.satisfies(version, '>=5.2.0 <=12')) {
    agent.logger.debug('koa-router version %s not supported - aborting...', version)
    return Router
  }

  agent.logger.debug('shimming koa-router prototype.match function')
  shimmer.wrap(Router.prototype, 'match', function (orig) {
    return function (_, method) {
      var matched = orig.apply(this, arguments)

      if (typeof method !== 'string') {
        agent.logger.debug('unexpected method type in koa-router prototype.match: %s', typeof method)
        return matched
      }

      if (Array.isArray(matched && matched.pathAndMethod)) {
        const layer = matched.pathAndMethod.find(function (layer) {
          return layer && layer.opts && layer.opts.end === true
        })

        var path = layer && layer.path
        if (typeof path === 'string') {
          var name = method + ' ' + path
          agent._instrumentation.setDefaultTransactionName(name)
        } else {
          agent.logger.debug('unexpected path type in koa-router prototype.match: %s', typeof path)
        }
      } else {
        agent.logger.debug('unexpected match result in koa-router prototype.match: %s', typeof matched)
      }

      return matched
    }
  })

  return Router
}
