/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

const asyncHooks = require('async_hooks')

const { BasicRunContextManager } = require('./BasicRunContextManager')

// A run context manager that uses an async hook to automatically track
// run context across async tasks.
//
// (Adapted from https://github.com/open-telemetry/opentelemetry-js/blob/main/packages/opentelemetry-context-async-hooks/src/AsyncHooksContextManager.ts)
class AsyncHooksRunContextManager extends BasicRunContextManager {
  constructor (log, runContextClass) {
    super(log, runContextClass)
    this._runContextFromAsyncId = new Map()
    this._asyncHook = asyncHooks.createHook({
      init: this._init.bind(this),
      before: this._before.bind(this),
      after: this._after.bind(this),
      destroy: this._destroy.bind(this),
      promiseResolve: this._destroy.bind(this)
    })
  }

  enable () {
    super.enable()
    this._asyncHook.enable()
    return this
  }

  disable () {
    super.disable()
    this._asyncHook.disable()
    this._runContextFromAsyncId.clear()
    return this
  }

  // Reset state for re-use of this context manager by tests in the same process.
  testReset () {
    // Absent a core node async_hooks bug, the easy way to implement this method
    // would be: `this.disable(); this.enable()`.
    // However there is a bug in Node.js v12.0.0 - v12.2.0 (inclusive) where
    // disabling the async hook could result in it never getting re-enabled.
    // https://github.com/nodejs/node/issues/27585
    // https://github.com/nodejs/node/pull/27590 (included in node v12.3.0)
    this._runContextFromAsyncId.clear()
    this._stack = []
  }

  /**
   * Init hook will be called when userland create a async context, setting the
   * context as the current one if it exist.
   * @param asyncId id of the async context
   * @param type the resource type
   */
  _init (asyncId, type, triggerAsyncId) {
    // ignore TIMERWRAP as they combine timers with same timeout which can lead to
    // false context propagation. TIMERWRAP has been removed in node 11
    // every timer has it's own `Timeout` resource anyway which is used to propagete
    // context.
    if (type === 'TIMERWRAP') {
      return
    }

    const context = this._stack[this._stack.length - 1]
    if (context !== undefined) {
      this._runContextFromAsyncId.set(asyncId, context)
    }
  }

  /**
   * Destroy hook will be called when a given context is no longer used so we can
   * remove its attached context.
   * @param asyncId id of the async context
   */
  _destroy (asyncId) {
    this._runContextFromAsyncId.delete(asyncId)
  }

  /**
   * Before hook is called just before executing a async context.
   * @param asyncId id of the async context
   */
  _before (asyncId) {
    const context = this._runContextFromAsyncId.get(asyncId)
    if (context !== undefined) {
      this._enterRunContext(context)
    }
  }

  /**
   * After hook is called just after completing the execution of a async context.
   */
  _after () {
    this._exitRunContext()
  }
}

module.exports = {
  AsyncHooksRunContextManager
}
