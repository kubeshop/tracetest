/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

/*
 * This instrumentation exists to work around an issue in mimic-response@1.0.0
 * that was fixed in mimic-response@1.0.1.
 * See https://github.com/elastic/apm-agent-nodejs/issues/423.
 */

var semver = require('semver')

module.exports = function (mimicResponse, agent, { version, enabled }) {
  if (!enabled) return mimicResponse

  if (semver.gte(version, '1.0.1')) {
    agent.logger.debug('mimic-response version %s doesn\'t need to be patched - ignoring...', version)
    return mimicResponse
  }

  var ins = agent._instrumentation

  return function wrappedMimicResponse (fromStream, toStream) {
    // If we bound the `fromStream` emitter, but not the `toStream` emitter, we
    // need to do so as else the `on`, `addListener`, and `prependListener`
    // functions of the `fromStream` will be copied over to the `toStream` but
    // run in the context of the `fromStream`.
    if (fromStream && toStream &&
        ins.isEventEmitterBound(fromStream) &&
        !ins.isEventEmitterBound(toStream)) {
      ins.bindEmitter(toStream)
    }
    return mimicResponse.apply(null, arguments)
  }
}
