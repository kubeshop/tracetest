/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

const EventEmitter = require('events')
const oblog = require('./oblog')
const { OTelBridgeRunContext } = require('./OTelBridgeRunContext')

// Implements interface ContextManager from:
// https://github.com/open-telemetry/opentelemetry-js-api/blob/v1.0.4/src/context/types.ts#L43
class OTelContextManager {
  constructor (agent) {
    this._agent = agent
    this._ins = agent._instrumentation
  }

  active () {
    oblog.apicall('OTelContextManager.active()')
    return this._ins.currRunContext()
  }

  _runContextFromOTelContext (otelContext) {
    let runContext
    if (otelContext instanceof OTelBridgeRunContext) {
      runContext = otelContext
    } else {
      // `otelContext` is some object implementing OTel's `interface Context`
      // (typically a `BaseContext` from @opentelemetry/api). We derive a new
      // OTelBridgeRunContext from the root run-context that properly uses
      // the Context.
      runContext = this._ins._runCtxMgr.root().setOTelContext(otelContext)
    }
    return runContext
  }

  with (otelContext, fn, thisArg, ...args) {
    oblog.apicall('OTelContextManager.with(%s<...>, function %s, ...)', otelContext.constructor.name, fn.name || '<anonymous>')
    const runContext = this._runContextFromOTelContext(otelContext)
    return this._ins._runCtxMgr.with(runContext, fn, thisArg, ...args)
  }

  bind (otelContext, target) {
    oblog.apicall('OTelContextManager.bind(%s, %s type)', otelContext, typeof target)
    if (target instanceof EventEmitter) {
      const runContext = this._runContextFromOTelContext(otelContext)
      return this._ins._runCtxMgr.bindEE(runContext, target)
    }
    if (typeof target === 'function') {
      const runContext = this._runContextFromOTelContext(otelContext)
      return this._ins._runCtxMgr.bindFn(runContext, target)
    }
    return target
  }

  enable () {
    oblog.apicall('OTelContextManager.enable()')
    this._ins._runCtxMgr.enable()
    return this
  }

  disable () {
    oblog.apicall('OTelContextManager.disable()')
    this._ins._runCtxMgr.disable()
    return this
  }
}

module.exports = {
  OTelContextManager
}
