/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

// A no-op (does nothing) Agent transport -- i.e. the APM server client API
// provided by elastic-apm-http-client.
//
// This is used for some configurations (when `disableSend=true` or when
// `contextPropagationOnly=true`) and in some tests.

class NoopTransport {
  config (opts) {}

  addMetadataFilter (fn) {}
  setExtraMetadata (metadata) {}
  lambdaStart () {}

  sendSpan (span, cb) {
    if (cb) {
      process.nextTick(cb)
    }
  }

  sendTransaction (transaction, cb) {
    if (cb) {
      process.nextTick(cb)
    }
  }

  sendError (_error, cb) {
    if (cb) {
      process.nextTick(cb)
    }
  }

  sendMetricSet (metricset, cb) {
    if (cb) {
      process.nextTick(cb)
    }
  }

  flush (opts, cb) {
    if (typeof opts === 'function') {
      cb = opts
      opts = {}
    } else if (!opts) {
      opts = {}
    }
    if (cb) {
      process.nextTick(cb)
    }
  }

  supportsKeepingUnsampledTransaction () {
    return true
  }

  // Inherited from Writable, called in agent.js.
  destroy () {}
}

module.exports = {
  NoopTransport
}
