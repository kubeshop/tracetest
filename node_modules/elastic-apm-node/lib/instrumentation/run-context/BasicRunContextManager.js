/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

const { AbstractRunContextManager } = require('./AbstractRunContextManager')

// A basic manager for run context. It handles a stack of run contexts, but does
// no automatic tracking (via async_hooks or otherwise). In combination with
// "patch-async.js" it does a limited job of context tracking for much of the
// core Node.js API.
//
// (Adapted from https://github.com/open-telemetry/opentelemetry-js/blob/main/packages/opentelemetry-context-async-hooks/src/AsyncHooksContextManager.ts)
class BasicRunContextManager extends AbstractRunContextManager {
  constructor (log, runContextClass) {
    super(log, runContextClass)
    this._stack = [] // Top of stack is the current run context.
  }

  // A string representation useful for debug logging. For example,
  //    BasicRunContextManager(
  //      RC(Trans(685ead, manual), [Span(9dd31c, GET httpstat.us, ended)]),
  //      RC(Trans(685ead, manual)) )
  toString () {
    return `${this.constructor.name}( ${this._stack.map(rc => rc.toString()).join(', ')} )`
  }

  enable () {
    return this
  }

  disable () {
    this._stack = []
    return this
  }

  // Reset state for re-use of this context manager by tests in the same process.
  testReset () {
    this.disable()
    this.enable()
  }

  active () {
    return this._stack[this._stack.length - 1] || this.root()
  }

  with (runContext, fn, thisArg, ...args) {
    this._enterRunContext(runContext)
    try {
      return fn.call(thisArg, ...args)
    } finally {
      this._exitRunContext()
    }
  }

  // This public method is needed to support the semantics of
  // apm.startTransaction() and apm.startSpan() that impact the current run
  // context.
  //
  // Otherwise, all run context changes are via `.with()` -- scoped to a
  // function call -- or via the "before" async hook -- scoped to an async task.
  supersedeRunContext (runContext) {
    this._exitRunContext()
    this._enterRunContext(runContext)
  }

  _enterRunContext (runContext) {
    this._stack.push(runContext)
  }

  _exitRunContext () {
    this._stack.pop()
  }
}

module.exports = {
  BasicRunContextManager
}
