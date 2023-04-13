/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

const { RunContext } = require('./RunContext')

const ADD_LISTENER_METHODS = [
  'addListener',
  'on',
  'once',
  'prependListener',
  'prependOnceListener'
]

// An abstract base RunContextManager class that implements the following
// methods that all run context manager implementations can share:
//    root()
//    bindFn(runContext, target)
//    bindEE(runContext, eventEmitter)
//    isEEBound(eventEmitter)
// and stubs out the remaining public methods of the RunContextManager
// interface.
//
// (This class has largerly the same API as @opentelemetry/api `ContextManager`.
// The implementation is adapted from
// https://github.com/open-telemetry/opentelemetry-js/blob/main/packages/opentelemetry-context-async-hooks/src/AbstractAsyncHooksContextManager.ts)
class AbstractRunContextManager {
  constructor (log, runContextClass = RunContext) {
    this._log = log
    this._kListeners = Symbol('ElasticListeners')
    this._root = new runContextClass() // eslint-disable-line new-cap
  }

  // Get the root run context. This is always empty (no current trans or span).
  //
  // This is the equivalent of OTel JS API's `ROOT_CONTEXT` constant. Ours
  // is not a top-level constant, because the RunContext class can be
  // overriden.
  root () {
    return this._root
  }

  enable () {
    throw new Error('abstract method not implemented')
  }

  disable () {
    throw new Error('abstract method not implemented')
  }

  // Reset state for re-use of this context manager by tests in the same process.
  testReset () {
    this.disable()
    this.enable()
  }

  active () {
    throw new Error('abstract method not implemented')
  }

  with (runContext, fn, thisArg, ...args) {
    throw new Error('abstract method not implemented')
  }

  supersedeRunContext (runContext) {
    throw new Error('abstract method not implemented')
  }

  // The OTel ContextManager API has a single .bind() like this:
  //
  // bind (runContext, target) {
  //   if (target instanceof EventEmitter) {
  //     return this._bindEventEmitter(runContext, target)
  //   }
  //   if (typeof target === 'function') {
  //     return this._bindFunction(runContext, target)
  //   }
  //   return target
  // }
  //
  // Is there any value in this over our two separate `.bind*` methods?

  bindFn (runContext, target) {
    if (typeof target !== 'function') {
      return target
    }
    // this._log.trace('bind %s to fn "%s"', runContext, target.name)

    const self = this
    const wrapper = function () {
      return self.with(runContext, () => target.apply(this, arguments))
    }
    Object.defineProperty(wrapper, 'length', {
      enumerable: false,
      configurable: true,
      writable: false,
      value: target.length
    })

    return wrapper
  }

  // (This implementation is adapted from OTel's `_bindEventEmitter`.)
  bindEE (runContext, ee) {
    // Explicitly do *not* guard with `ee instanceof EventEmitter`. The
    // `Request` object from the aws-sdk@2 module, for example, has an `on`
    // with the EventEmitter API that we want to bind, but it is not otherwise
    // an EventEmitter.

    const map = this._getPatchMap(ee)
    if (map !== undefined) {
      // No double-binding.
      return ee
    }
    this._createPatchMap(ee)

    // patch methods that add a listener to propagate context
    ADD_LISTENER_METHODS.forEach(methodName => {
      if (ee[methodName] === undefined) return
      ee[methodName] = this._patchAddListener(ee, ee[methodName], runContext)
    })
    // patch methods that remove a listener
    if (typeof ee.removeListener === 'function') {
      ee.removeListener = this._patchRemoveListener(ee, ee.removeListener)
    }
    if (typeof ee.off === 'function') {
      ee.off = this._patchRemoveListener(ee, ee.off)
    }
    // patch method that remove all listeners
    if (typeof ee.removeAllListeners === 'function') {
      ee.removeAllListeners = this._patchRemoveAllListeners(
        ee,
        ee.removeAllListeners
      )
    }
    return ee
  }

  // Return true iff the given EventEmitter is already bound to a run context.
  isEEBound (ee) {
    return (this._getPatchMap(ee) !== undefined)
  }

  // Patch methods that remove a given listener so that we match the "patched"
  // version of that listener (the one that propagate context).
  _patchRemoveListener (ee, original) {
    const contextManager = this
    return function (event, listener) {
      const map = contextManager._getPatchMap(ee)
      const listeners = map && map[event]
      if (listeners === undefined) {
        return original.call(this, event, listener)
      }
      const patchedListener = listeners.get(listener)
      return original.call(this, event, patchedListener || listener)
    }
  }

  // Patch methods that remove all listeners so we remove our internal
  // references for a given event.
  _patchRemoveAllListeners (ee, original) {
    const contextManager = this
    return function (event) {
      const map = contextManager._getPatchMap(ee)
      if (map !== undefined) {
        if (arguments.length === 0) {
          contextManager._createPatchMap(ee)
        } else if (map[event] !== undefined) {
          delete map[event]
        }
      }
      return original.apply(this, arguments)
    }
  }

  // Patch methods on an event emitter instance that can add listeners so we
  // can force them to propagate a given context.
  _patchAddListener (ee, original, runContext) {
    const contextManager = this
    return function (event, listener) {
      let map = contextManager._getPatchMap(ee)
      if (map === undefined) {
        map = contextManager._createPatchMap(ee)
      }
      let listeners = map[event]
      if (listeners === undefined) {
        listeners = new WeakMap()
        map[event] = listeners
      }
      const patchedListener = contextManager.bindFn(runContext, listener)
      // store a weak reference of the user listener to ours
      listeners.set(listener, patchedListener)
      return original.call(this, event, patchedListener)
    }
  }

  _createPatchMap (ee) {
    const map = Object.create(null)
    ee[this._kListeners] = map
    return map
  }

  _getPatchMap (ee) {
    return ee[this._kListeners]
  }
}

module.exports = {
  AbstractRunContextManager
}
