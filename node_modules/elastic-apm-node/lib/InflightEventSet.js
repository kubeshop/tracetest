/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

// A `Set` object used to track inflight APM events: ended spans and errors
// that are currently being processed, but have not yet been sent to the Agent
// transport.
//
// `setDrainHandler` allows setting a function to be called when the inflight
// events have drained. Agent#flush() uses this to ensure that a flush waits
// for inflight events to be processed, so they are sent to APM Server before
// calling back.
class InflightEventSet extends Set {
  // Set a `fn` to be called *once* when the set size next goes to zero.
  // If the optional `timeoutMs` is given, then `fn(err)` will be called if
  // the set hasn't yet drained.
  setDrainHandler (fn, timeoutMs) {
    this._drainHandler = fn
    if (timeoutMs) {
      this._drainTimeout = setTimeout(() => {
        this._drain(new Error('inflight event set drain timeout'))
      }, timeoutMs).unref()
    }
  }

  // Call the drain handler, if there is one.
  _drain (err) {
    if (this._drainHandler) {
      if (this._drainTimeout) {
        clearTimeout(this._drainTimeout)
        this._drainTimeout = null
      }
      this._drainHandler(err)
      // Remove the handler so it is only called once.
      this._drainHandler = null
    }
  }

  delete (key) {
    super.delete(key)
    if (this.size === 0) {
      this._drain()
    }
  }
}

module.exports = {
  InflightEventSet
}
