/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'
const { TraceParent } = require('./traceparent')
const TraceState = require('./tracestate')

class TraceContext {
  constructor (traceparent, tracestate, conf = {}) {
    this.traceparent = traceparent
    this.tracestate = tracestate
    this._conf = conf
  }

  // Resume (a.k.a. continue) a trace if `childOf` includes trace-context.
  // Otherwise, make a sampling decision and start new trace-context.
  //
  // @param { Transaction | Span | TraceParent | string | undefined } childOf
  //    If this is a Transaction, Span, or TraceParent instance, or a valid
  //    W3C trace-context 'traceparent' string, then the trace is continued.
  //    Otherwise a sampling decision made based on
  //    `conf.transactionSampleRate` and a new trace is started.
  // @param { Object } conf - The agent configuration.
  // @param { TraceState | string } tracestateArg - A TraceState instance or
  //    a W3C trace-context 'tracestate' string.
  static startOrResume (childOf, conf, tracestateArg) {
    if (childOf && childOf._context instanceof TraceContext) return childOf._context.child()
    const traceparent = TraceParent.startOrResume(childOf, conf)
    const tracestate = tracestateArg instanceof TraceState
      ? tracestateArg
      : TraceState.fromStringFormatString(tracestateArg)

    // if a root transaction/span, set the initial sample rate in the tracestate

    if (!childOf && traceparent.recorded) {
      // if this is a sampled/recorded transaction, set the rate
      tracestate.setValue('s', conf.transactionSampleRate)
    } else if (!childOf) {
      // if this is a un-sampled/unreocrded transaction, set the
      // rate to zero, per the spec
      //
      // https://github.com/elastic/apm/blob/main/specs/agents/tracing-sampling.md
      tracestate.setValue('s', 0)
    }

    return new TraceContext(traceparent, tracestate, conf)
  }

  static fromString (header) {
    return TraceParent.fromString(header)
  }

  ensureParentId () {
    return this.traceparent.ensureParentId()
  }

  child () {
    const childTraceParent = this.traceparent.child()
    const childTraceContext = new TraceContext(
      childTraceParent, this.tracestate, this._conf
    )
    return childTraceContext
  }

  /**
   * Returns traceparent string only
   *
   * @todo legacy -- can we remove to avoid confusion?
   */
  toString () {
    return this.traceparent.toString()
  }

  toTraceStateString () {
    return this.tracestate.toW3cString()
  }

  toTraceParentString () {
    return this.traceparent.toString()
  }

  propagateTraceContextHeaders (carrier, setter) {
    if (!carrier || !setter) {
      return
    }
    const traceparentStr = this.toTraceParentString()
    const tracestateStr = this.toTraceStateString()
    if (traceparentStr) {
      setter(carrier, 'traceparent', traceparentStr)
      if (this._conf.useElasticTraceparentHeader) {
        setter(carrier, 'elastic-apm-traceparent', traceparentStr)
      }
    }

    if (tracestateStr) {
      setter(carrier, 'tracestate', tracestateStr)
    }
  }

  setRecorded () {
    return this.traceparent.setRecorded()
  }

  isRecorded () {
    return this.traceparent.recorded
  }
}

TraceContext.FLAGS = TraceParent.FLAGS

module.exports = TraceContext
