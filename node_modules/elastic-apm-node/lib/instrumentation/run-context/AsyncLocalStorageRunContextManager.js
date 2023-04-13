/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

const { AsyncLocalStorage } = require('async_hooks')

const { AbstractRunContextManager } = require('./AbstractRunContextManager')

/**
 * A RunContextManager that uses core node `AsyncLocalStorage` as the mechanism
 * for run-context tracking.
 *
 * (Adapted from https://github.com/open-telemetry/opentelemetry-js/blob/main/packages/opentelemetry-context-async-hooks/src/AsyncLocalStorageContextManager.ts)
 */
class AsyncLocalStorageRunContextManager extends AbstractRunContextManager {
  constructor (log, runContextClass) {
    super(log, runContextClass)
    this._asyncLocalStorage = new AsyncLocalStorage()
  }

  // A string representation useful for debug logging. For example,
  //    AsyncLocalStorageRunContextManager( RC(Trans(685ead, manual), [Span(9dd31c, GET httpstat.us, ended)]) )
  toString () {
    return `${this.constructor.name}( ${this.active().toString()} )`
  }

  enable () {
    return this
  }

  disable () {
    this._asyncLocalStorage.disable()
    return this
  }

  // Reset state for re-use of this context manager by tests in the same process.
  testReset () {
    this.disable()
    this.enable()
  }

  active () {
    const store = this._asyncLocalStorage.getStore()
    if (store == null) {
      return this.root()
    } else {
      return store
    }
  }

  with (runContext, fn, thisArg, ...args) {
    const cb = thisArg == null ? fn : fn.bind(thisArg)
    return this._asyncLocalStorage.run(runContext, cb, ...args)
  }

  // This public method is needed to support the semantics of
  // apm.startTransaction() and apm.startSpan() that impact the current run
  // context.
  //
  // Otherwise, all run context changes are via `.with()` -- scoped to a
  // function call -- or via the "before" async hook -- scoped to an async task.
  supersedeRunContext (runContext) {
    this._asyncLocalStorage.enterWith(runContext)
  }
}

module.exports = {
  AsyncLocalStorageRunContextManager
}
