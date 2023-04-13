/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

var semver = require('semver')

const shimmer = require('../shimmer')

module.exports = function (expressGraphql, agent, { version, enabled }) {
  if (!enabled) {
    return expressGraphql
  }

  if (semver.satisfies(version, '>=0.10.0 <0.13.0')) {
    // https://github.com/graphql/express-graphql/pull/626 changed `graphqlHTTP`
    // to no longer be the top-level export:
    //   {
    //     graphqlHTTP: [Function: graphqlHTTP],
    //     getGraphQLParams: [AsyncFunction: getGraphQLParams]
    //   }
    shimmer.wrap(expressGraphql, 'graphqlHTTP', wrapGraphqlHTTP)
    return expressGraphql
  } else if (semver.satisfies(version, '>=0.6.1 <0.10.0') && typeof expressGraphql === 'function') {
    // Up to and including 0.9.x, `require('express-graphql')` is:
    //   [Function: graphqlHTTP] {
    //     getGraphQLParams: [AsyncFunction: getGraphQLParams]
    //   }
    const wrappedGraphqlHTTP = wrapGraphqlHTTP(expressGraphql)
    for (const key of Object.keys(expressGraphql)) {
      wrappedGraphqlHTTP[key] = expressGraphql[key]
    }
    return wrappedGraphqlHTTP
  } else {
    agent.logger.debug('express-graphql@%s not supported: skipping instrumentation', version)
    return expressGraphql
  }

  function wrapGraphqlHTTP (origGraphqlHTTP) {
    return function wrappedGraphqlHTTP () {
      var orig = origGraphqlHTTP.apply(this, arguments)

      if (typeof orig !== 'function') {
        return orig
      }

      // Express is very particular with the number of arguments!
      return function (req, res) {
        var trans = agent._instrumentation.currTransaction()
        if (trans) {
          trans._graphqlRoute = true
        }
        return orig.apply(this, arguments)
      }
    }
  }
}
