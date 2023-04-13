/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

const { executionAsyncId } = require('async_hooks')
const { URL } = require('url')
var util = require('util')

var Value = require('async-value-promise')

const constants = require('../constants')
var GenericSpan = require('./generic-span')
var { SpanIds } = require('./ids')
const { gatherStackTrace } = require('../stacktraces')

const TEST = process.env.ELASTIC_APM_TEST

module.exports = Span

util.inherits(Span, GenericSpan)

// new Span(transaction)
// new Span(transaction, name?, opts?)
// new Span(transaction, name?, type?, opts?)
// new Span(transaction, name?, type?, subtype?, opts?)
// new Span(transaction, name?, type?, subtype?, action?, opts?)
function Span (transaction, ...args) {
  const opts = typeof args[args.length - 1] === 'object'
    ? (args.pop() || {})
    : {}
  const [name, ...tsaArgs] = args // "tsa" === Type, Subtype, Action

  if (opts.timer) {
    process.emitWarning(
      'specifying the `timer` option to `new Span()` was never a public API and will be removed',
      'DeprecationWarning',
      'ELASTIC_APM_SPAN_TIMER_OPTION'
    )
  }
  if (!opts.childOf) {
    const defaultChildOf = transaction._agent._instrumentation.currSpan() || transaction
    opts.childOf = defaultChildOf
    opts.timer = defaultChildOf._timer
  } else if (opts.childOf._timer) {
    opts.timer = opts.childOf._timer
  }

  this._exitSpan = !!opts.exitSpan
  this.discardable = this._exitSpan

  delete opts.exitSpan

  GenericSpan.call(this, transaction._agent, ...tsaArgs, opts)

  this._db = null
  this._http = null
  this._destination = null
  this._serviceTarget = null
  this._excludeServiceTarget = false
  this._message = null
  this._stackObj = null
  this._capturedStackTrace = null
  this.sync = true
  this._startXid = executionAsyncId()

  this.transaction = transaction
  this.name = name || 'unnamed'

  if (this._agent._conf.spanStackTraceMinDuration >= 0) {
    this._recordStackTrace()
  }

  this._agent.logger.debug('start span %o', { span: this.id, parent: this.parentId, trace: this.traceId, name: this.name, type: this.type, subtype: this.subtype, action: this.action })
}

Object.defineProperty(Span.prototype, 'ids', {
  get () {
    return this._ids === null
      ? (this._ids = new SpanIds(this))
      : this._ids
  }
})

Span.prototype.toString = function () {
  return this.ids.toString()
}

Span.prototype.customStackTrace = function (stackObj) {
  this._agent.logger.debug('applying custom stack trace to span %o', { span: this.id, parent: this.parentId, trace: this.traceId })
  this._recordStackTrace(stackObj)
}

Span.prototype.end = function (endTime) {
  if (this.ended) {
    this._agent.logger.debug('tried to call span.end() on already ended span %o', { span: this.id, parent: this.parentId, trace: this.traceId, name: this.name, type: this.type, subtype: this.subtype, action: this.action })
    return
  }

  this._timer.end(endTime)
  this._endTimestamp = this._timer.endTimestamp
  this._duration = this._timer.duration
  if (executionAsyncId() !== this._startXid) {
    this.sync = false
  }

  this._setOutcomeFromSpanEnd()

  this._inferServiceTargetAndDestinationService()

  this.ended = true
  this._agent.logger.debug('ended span %o', { span: this.id, parent: this.parentId, trace: this.traceId, name: this.name, type: this.type, subtype: this.subtype, action: this.action })

  if (this._capturedStackTrace !== null &&
      this._agent._conf.spanStackTraceMinDuration >= 0 &&
      this._duration / 1000 >= this._agent._conf.spanStackTraceMinDuration) {
    // NOTE: This uses a promise-like thing and not a *real* promise
    // because passing error stacks into a promise context makes it
    // uncollectable by the garbage collector.
    this._stackObj = new Value()
    var self = this
    gatherStackTrace(
      this._agent.logger,
      this._capturedStackTrace,
      this._agent._conf.sourceLinesSpanAppFrames,
      this._agent._conf.sourceLinesSpanLibraryFrames,
      TEST ? null : filterCallSite,
      function (_err, stacktrace) {
        // _err from gatherStackTrace is always null.
        self._stackObj.resolve(stacktrace)
      }
    )
  }

  this._agent._instrumentation.addEndedSpan(this)
  this.transaction._captureBreakdown(this)
}

Span.prototype._inferServiceTargetAndDestinationService = function () {
  // `context.service.target.*` must be set for exit spans. There is a public
  // `span.setServiceTarget(...)` for users to manually set this, but typically
  // it is inferred from other span fields here.
  // https://github.com/elastic/apm/blob/main/specs/agents/tracing-spans-service-target.md#field-values
  if (this._excludeServiceTarget || !this._exitSpan) {
    this._serviceTarget = null
  } else {
    if (!this._serviceTarget) {
      this._serviceTarget = {}
    }
    if (!('type' in this._serviceTarget)) {
      this._serviceTarget.type = this.subtype || this.type || constants.DEFAULT_SPAN_TYPE
    }
    if (!('name' in this._serviceTarget)) {
      if (this._db) {
        if (this._db.instance) {
          this._serviceTarget.name = this._db.instance
        }
      } else if (this._message) {
        if (this._message.queue && this._message.queue.name) {
          this._serviceTarget.name = this._message.queue.name
        }
      } else if (this._http && this._http.url) {
        try {
          this._serviceTarget.name = new URL(this._http.url).host
        } catch (invalidUrlErr) {
          this._agent.logger.debug('cannot set "service.target.name": %s (ignoring)', invalidUrlErr)
        }
      }
    }

    // `destination.service.*` is deprecated, but still required for older
    // APM servers.
    if (!this._destination) { this._destination = {} }
    if (!this._destination.service) { this._destination.service = {} }
    // - `destination.service.{type,name}` could be skipped if the upstream APM server is known to be >=7.14.
    this._destination.service.type = ''
    this._destination.service.name = ''
    // - Infer the now deprecated `context.destination.service.resource` value.
    //   https://github.com/elastic/apm/blob/main/specs/agents/tracing-spans-destination.md#destination-resource
    if (!this._destination.service.resource) {
      if (!this._serviceTarget.name) {
        // If we only have `.type`, then use that.
        this._destination.service.resource = this._serviceTarget.type
      } else if (!this._serviceTarget.type) {
        // If we only have `.name`, then use that.
        this._destination.service.resource = this._serviceTarget.name
      } else if (this.type === 'external') {
        // Typically the "resource" value would now be "$type/$name", e.g.
        // "mysql/customers". However, we want a special case for some spans (to
        // have the same value as historically?) where we do NOT use the
        // "$type/" prefix. One example is HTTP spans. Another is gRPC spans
        // and, I infer from otel_bridge.feature, any OTel "rpc.system"-usage
        // spans as well
        // (https://opentelemetry.io/docs/reference/specification/trace/semantic_conventions/rpc/).
        // Options to infer this from other span data:
        // - Use the presence of "http" context, but without "db" and "message"
        //   context. This is a little brittle, and requires more complete OTel
        //   bridge compatibility mapping of OTel attributes than is currently
        //   being done.
        //      } else if (!this._db && !this._message && this._http && this._http.url) {
        // - Use `span.subtype`: "http", "grpc", ... add others if/when they are
        //   used.
        // - Use `span.type === "external"`. This, at least currently corresponds.
        //   Let's use this one.
        this._destination.service.resource = this._serviceTarget.name
      } else {
        this._destination.service.resource = `${this._serviceTarget.type}/${this._serviceTarget.name}`
      }
    }
  }
}

Span.prototype.setDbContext = function (context) {
  if (!context) return
  this._db = Object.assign(this._db || {}, context)
}

Span.prototype.setHttpContext = function (context) {
  if (!context) return
  this._http = Object.assign(this._http || {}, context)
}

/**
 * This is deprecated and will be dropped in a future version. This was always
 * an internal method, but possibly used by enterprising users of manual
 * instrumentation.
 *
 * @deprecated Users should use the public `setServiceTarget()`.
 *    Internal APM agent code should use `_setDestinationContext()`.
 */
Span.prototype.setDestinationContext = function (destCtx) {
  process.emitWarning(
    '<span>.setDestinationContext() was never a public API and will be removed, use <span>.setServiceTarget().',
    'DeprecationWarning',
    'ELASTIC_APM_SET_DESTINATION_CONTEXT'
  )

  if (destCtx.service && destCtx.service.resource) {
    this.setServiceTarget('', destCtx.service.resource)
  }
  const destCtxWithoutService = Object.assign({}, destCtx)
  delete destCtxWithoutService.service
  this._setDestinationContext(destCtxWithoutService)
}

/**
 * The internal method for setting "destination" context.
 *
 * "destination.service.resource" should only ever be included for special
 * cases. It is typically inferred from other fields via a general algorithm.
 */
Span.prototype._setDestinationContext = function (destCtx) {
  this._destination = Object.assign(this._destination || {}, destCtx)
}

/**
 * Manually set the `service.target.type` and `service.target.name` fields that
 * are used for service maps and the identification of downstream services. The
 * values are only used for "exit" spans -- spans representing outgoing
 * communication, marked with `exitSpan: true` at span creation.
 *
 * If false-y values (e.g. `null`) are given for both `type` and `name`, then
 * `service.target` will explicitly be excluded from this span. This may impact
 * Service Maps and other Kibana APM app reporting for this service.
 *
 * If this method is not called, values are inferred from other span fields per
 * https://github.com/elastic/apm/blob/main/specs/agents/tracing-spans-service-target.md#field-values
 *
 * `service.target.*` fields are ignored for APM Server before v8.3.
 *
 * @param {string | null} type - service target type, usually same value as
 *    `span.subtype`
 * @param {string | null} name - service target name: value depends on type,
 *    for databases it's usually the database name
 */
Span.prototype.setServiceTarget = function (type, name) {
  if (!type && !name) {
    this._excludeServiceTarget = true
    this._serviceTarget = null
    return
  }

  if (typeof type === 'string') {
    this._excludeServiceTarget = false
    if (this._serviceTarget === null) {
      this._serviceTarget = { type }
    } else {
      this._serviceTarget.type = type
    }
  } else {
    this._agent.logger.warn('"type" argument to Span#setServiceTarget must be of type "string", got type "%s": ignoring', typeof type)
  }
  if (typeof name === 'string') {
    this._excludeServiceTarget = false
    if (this._serviceTarget === null) {
      this._serviceTarget = { name }
    } else {
      this._serviceTarget.name = name
    }
  } else {
    this._agent.logger.warn('"name" argument to Span#setServiceTarget must be of type "string", got type "%s": ignoring', typeof name)
  }
}

Span.prototype.setMessageContext = function (context) {
  this._message = Object.assign(this._message || {}, context)
}

Span.prototype.setOutcome = function (outcome) {
  if (!this._isValidOutcome(outcome)) {
    this._agent.logger.trace(
      'Unknown outcome [%s] seen in Span.setOutcome, ignoring',
      outcome
    )
    return
  }

  if (this.ended) {
    this._agent.logger.debug(
      'tried to call Span.setOutcome() on already ended span %o',
      { span: this.id, parent: this.parentId, trace: this.traceId, name: this.name, type: this.type, subtype: this.subtype, action: this.action }
    )
    return
  }
  this._freezeOutcome()
  this._setOutcome(outcome)
}

Span.prototype._setOutcomeFromErrorCapture = function (outcome) {
  if (this._isOutcomeFrozen) {
    return
  }
  this._setOutcome(outcome)
}

Span.prototype._setOutcomeFromHttpStatusCode = function (statusCode) {
  if (this._isOutcomeFrozen) {
    return
  }
  /**
   * The statusCode could be undefined for example if,
   * the request is aborted before socket, in that case
   * we keep the default 'unknown' value.
   */
  if (typeof statusCode !== 'undefined') {
    if (statusCode >= 400) {
      this._setOutcome(constants.OUTCOME_FAILURE)
    } else {
      this._setOutcome(constants.OUTCOME_SUCCESS)
    }
  }

  this._freezeOutcome()
}

Span.prototype._setOutcomeFromSpanEnd = function () {
  if (this.outcome === constants.OUTCOME_UNKNOWN && !this._isOutcomeFrozen) {
    this._setOutcome(constants.OUTCOME_SUCCESS)
  }
}

/**
 * Central setting for outcome
 *
 * Enables "when outcome does X, Y should also happen" behaviors
 */
Span.prototype._setOutcome = function (outcome) {
  this.outcome = outcome
  if (outcome !== constants.OUTCOME_SUCCESS) {
    this.discardable = false
  }
}

Span.prototype._recordStackTrace = function (obj) {
  if (!obj) {
    obj = {}
    Error.captureStackTrace(obj, Span)
  }
  this._capturedStackTrace = obj
}

Span.prototype._encode = function (cb) {
  var self = this

  if (!this.ended) {
    return cb(new Error('cannot encode un-ended span'))
  }

  const payload = {
    id: self.id,
    transaction_id: self.transaction.id,
    parent_id: self.parentId,
    trace_id: self.traceId,
    name: self.name,
    type: self.type || constants.DEFAULT_SPAN_TYPE,
    subtype: self.subtype,
    action: self.action,
    timestamp: self.timestamp,
    duration: self._duration,
    context: undefined,
    stacktrace: undefined,
    sync: self.sync,
    outcome: self.outcome
  }

  // if a valid sample rate is set (truthy or zero), set the property
  const sampleRate = self.sampleRate
  if (sampleRate !== null) {
    payload.sample_rate = sampleRate
  }

  let haveContext = false
  const context = {}
  if (self._serviceTarget) {
    context.service = { target: self._serviceTarget }
    haveContext = true
  }
  if (self._destination) {
    context.destination = self._destination
    haveContext = true
  }
  if (self._db) {
    context.db = self._db
    haveContext = true
  }
  if (self._message) {
    context.message = self._message
    haveContext = true
  }
  if (self._http) {
    context.http = self._http
    haveContext = true
  }
  if (self._labels) {
    context.tags = self._labels
    haveContext = true
  }
  if (haveContext) {
    payload.context = context
  }

  if (self.isComposite()) {
    payload.composite = self._compression.encode()
    payload.timestamp = self._compression.timestamp
    payload.duration = self._compression.duration
  }

  this._serializeOTel(payload)

  if (this._links.length > 0) {
    payload.links = this._links
  }

  if (this._stackObj) {
    this._stackObj.then(
      value => done(null, value),
      error => done(error)
    )
  } else {
    process.nextTick(done)
  }

  function done (err, frames) {
    if (err) {
      self._agent.logger.debug('could not capture stack trace for span %o', { span: self.id, parent: self.parentId, trace: self.traceId, name: self.name, type: self.type, subtype: self.subtype, action: self.action, err: err.message })
    } else if (frames) {
      payload.stacktrace = frames
    }

    // Reduce this span's memory usage by dropping references once they're
    // no longer needed.  We also keep fields required to support
    // `interface Span`.
    // Span fields:
    self._db = null
    self._http = null
    self._message = null
    self._capturedStackTrace = null
    // GenericSpan fields:
    // - Cannot drop `this._context` because it is used for traceparent and ids.
    self._timer = null
    self._labels = null

    cb(null, payload)
  }
}

Span.prototype.isCompressionEligible = function () {
  if (!this.getParentSpan()) {
    return false
  }

  if (this.outcome !== constants.OUTCOME_UNKNOWN &&
      this.outcome !== constants.OUTCOME_SUCCESS
  ) {
    return false
  }

  if (!this._exitSpan) {
    return false
  }

  if (this._hasPropagatedTraceContext) {
    return false
  }

  return true
}

Span.prototype.tryToCompress = function (spanToCompress) {
  return this._compression.tryToCompress(this, spanToCompress)
}

Span.prototype.isRecorded = function () {
  return this._context.isRecorded()
}

Span.prototype.setRecorded = function (value) {
  return this._context.setRecorded(value)
}

Span.prototype.propagateTraceContextHeaders = function (carrier, setter) {
  this.discardable = false
  return GenericSpan.prototype.propagateTraceContextHeaders.call(this, carrier, setter)
}

function filterCallSite (callsite) {
  var filename = callsite.getFileName()
  return filename ? filename.indexOf('/node_modules/elastic-apm-node/') === -1 : true
}
