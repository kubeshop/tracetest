/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

const semver = require('semver')
const shimmer = require('../shimmer')
const clone = require('shallow-clone-shim')

module.exports = function (apolloServerCore, agent, { version, enabled }) {
  if (!enabled) return apolloServerCore

  if (!semver.satisfies(version, '^2.0.2 || ^3.0.0')) {
    agent.logger.debug('apollo-server-core version %s not supported - aborting...', version)
    return apolloServerCore
  }

  function wrapRunHttpQuery (orig) {
    return function wrappedRunHttpQuery () {
      var trans = agent._instrumentation.currTransaction()
      if (trans) trans._graphqlRoute = true
      return orig.apply(this, arguments)
    }
  }

  if (semver.satisfies(version, '<2.14')) {
    shimmer.wrap(apolloServerCore, 'runHttpQuery', wrapRunHttpQuery)
    return apolloServerCore
  }

  // apollo-server-core >= 2.14 does not allow overriding the exports object
  return clone({}, apolloServerCore, {
    runHttpQuery (descriptor) {
      const getter = descriptor.get
      if (getter) {
        descriptor.get = function get () {
          return wrapRunHttpQuery(getter())
        }
      }
      return descriptor
    }
  })
}
