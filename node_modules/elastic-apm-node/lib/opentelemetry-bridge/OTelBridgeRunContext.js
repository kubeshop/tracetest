/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

const otel = require('@opentelemetry/api')

const oblog = require('./oblog')
const { OTelBridgeNonRecordingSpan } = require('./OTelBridgeNonRecordingSpan')
const { OTelSpan } = require('./OTelSpan')
const { RunContext } = require('../instrumentation/run-context')
const Span = require('../instrumentation/span')

const OTEL_CONTEXT_KEY = otel.createContextKey('Elastic APM Context Key OTEL CONTEXT')
let SPAN_KEY = null

// `fetchSpanKey()` is called once during OTel SDK setup to get the `SPAN_KEY`
// that will be used by the OTel JS API during tracing -- when
// `otel.trace.setSpan(context, span)` et al are called.
//
// The fetched SPAN_KEY is used later by OTelBridgeRunContext to intercept
// `Context.{get,set,delete}Value` and translate to the agent's internal
// RunContext semantics for controlling the active/current span.
function fetchSpanKey () {
  const capturingContext = {
    spanKey: null,
    setValue (key, _value) {
      this.spanKey = key
    }
  }
  const fakeSpan = {}
  otel.trace.setSpan(capturingContext, fakeSpan)
  SPAN_KEY = capturingContext.spanKey
  if (!SPAN_KEY) {
    throw new Error('could not fetch OTel API "SPAN_KEY"')
  }
}

// This is a subclass of RunContext that is used when the agent's OTel SDK
// is enabled. It bridges between the OTel API's `Context` and the agent's
// `RunContext`.
//
// 1. It bridges between `<Context>.setValue(SPAN_KEY, ...)`, `.getValue(SPAN_KEY)`
//    used by the OTel API and the RunContext methods used to track the current
//    transaction and span.
// 2. It can propagate an OTel API `Context` instance (e.g. the internal
//    `BaseContext` that is exposed by `otel.ROOT_CONTEXT`) and proxy
//    `.getValue(key)` calls to it. See `OTEL_CONTEXT_KEY` below.
class OTelBridgeRunContext extends RunContext {
  setOTelContext (otelContext) {
    // First, save the `Context` instance in case it holds keys other than SPAN_KEY.
    let runContext = this.setValue(OTEL_CONTEXT_KEY, otelContext)
    // Second, if the `Context` holds a span, then pass that to our `setValue`
    // that knows how to translate that to RunContext semantics.
    const span = otel.trace.getSpan(otelContext)
    if (span) {
      runContext = runContext.setValue(SPAN_KEY, span)
    }
    return runContext
  }

  getValue (key) {
    oblog.apicall('OTelBridgeRunContext.getValue(%o)', key)
    if (key === SPAN_KEY) {
      const curr = this.currSpan() || this.currTransaction()
      if (!curr) {
        return undefined
      } else if (curr instanceof OTelBridgeNonRecordingSpan) {
        return curr
      } else {
        return new OTelSpan(curr)
      }
    }
    const value = super.getValue(key)
    if (value !== undefined) {
      return value
    } else {
      // Fallback to possibly-stashed OTel API Context instance.
      const otelContext = super.getValue(OTEL_CONTEXT_KEY)
      if (otelContext) {
        return otelContext.getValue(key)
      }
    }
  }

  setValue (key, value) {
    oblog.apicall('OTelBridgeRunContext.setValue(%o, %s)', key, value)
    if (key === SPAN_KEY) {
      if (value instanceof OTelSpan) {
        if (value._span instanceof Span) {
          return this.enterSpan(value._span)
        } else {
          // assert(value._span instanceof Transaction || value._span instanceof OTelBridgeNonRecordingSpan)
          return this.enterTrans(value._span)
        }
      } else if (typeof value.isRecording === 'function' && !value.isRecording()) {
        return this.enterTrans(new OTelBridgeNonRecordingSpan(value))
      }
    }
    return super.setValue(key, value)
  }

  deleteValue (key) {
    oblog.apicall('OTelBridgeRunContext.deleteValue(%o)', key)
    if (key === SPAN_KEY) {
      return this.leaveTrans()
    }
    // TODO: Should perhaps proxy deleteValue(key) to the possible underlying OTEL_CONTEXT_KEY entry.
    return super.deleteValue(key)
  }
}

module.exports = {
  fetchSpanKey,
  OTelBridgeRunContext
}
