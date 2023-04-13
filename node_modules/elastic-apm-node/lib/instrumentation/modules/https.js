/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

var semver = require('semver')

var httpShared = require('../http-shared')
var shimmer = require('../shimmer')

module.exports = function (https, agent, { version, enabled }) {
  if (agent._conf.instrumentIncomingHTTPRequests) {
    agent.logger.debug('shimming https.Server.prototype.emit function')
    shimmer.wrap(https && https.Server && https.Server.prototype, 'emit', httpShared.instrumentRequest(agent, 'https'))
  }

  if (!enabled) return https
  // From Node.js v9.0.0 and onwards, https requests no longer just call the
  // http.request function. So to correctly instrument outgoing HTTPS requests
  // in all supported Node.js versions, we'll only only instrument the
  // https.request function if the Node version is v9.0.0 or above.
  //
  // This change was introduced in:
  // https://github.com/nodejs/node/commit/5118f3146643dc55e7e7bd3082d1de4d0e7d5426
  if (semver.gte(version, '9.0.0')) {
    agent.logger.debug('shimming https.request function')
    shimmer.wrap(https, 'request', httpShared.traceOutgoingRequest(agent, 'https', 'request'))

    agent.logger.debug('shimming https.get function')
    shimmer.wrap(https, 'get', httpShared.traceOutgoingRequest(agent, 'https', 'get'))
  } else {
    // We must ensure that the `http` module is instrumented to intercept
    // `http.{request,get}` that `https.{request,get}` are using.
    require('http')
  }

  return https
}
