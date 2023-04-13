/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

var util = require('util')

var ObjectIdentityMap = require('object-identity-map')

const constants = require('../constants')
const { DroppedSpanStats } = require('./dropped-span-stats')
var getPathFromRequest = require('./express-utils').getPathFromRequest
var GenericSpan = require('./generic-span')
var parsers = require('../parsers')
var Span = require('./span')
var symbols = require('../symbols')
const {
  TRACE_CONTINUATION_STRATEGY_CONTINUE,
  TRACE_CONTINUATION_STRATEGY_RESTART,
  TRACE_CONTINUATION_STRATEGY_RESTART_EXTERNAL
} = require('../config')
var { TransactionIds } = require('./ids')
const TraceState = require('../tracecontext/tracestate')

module.exports = Transaction
util.inherits(Transaction, GenericSpan)

// Usage:
//    new Transaction(agent, opts?)
//    new Transaction(agent, name, opts?)
//    new Transaction(agent, name, type?, opts?)
//    new Transaction(agent, name, type?, subtype?, opts?)
//    new Transaction(agent, name, type?, subtype?, action?, opts?)
//
// @param {Agent} agent
// @param {string} [name]
// @param {string} [type] - Defaults to 'custom' when serialized.
// @param {string} [subtype] - Deprecated. Unused.
// @param {string} [action] - Deprecated. Unused.
// @param {string} [opts]
//    - opts.childOf - Used to determine the W3C trace-context trace id, parent
//      id, and sampling information for this new transaction. This currently
//      accepts a Transaction instance, Span instance, TraceParent instance, or
//      a traceparent string. (Arguably any but the latter two are non-sensical
//      for a new transaction.)
//    - opts.tracestate - A W3C trace-context tracestate string.
//    - Any other options supported by GenericSpan ...
function Transaction (agent, name, ...args) {
  const opts = typeof args[args.length - 1] === 'object'
    ? (args.pop() || {})
    : {}

  if (opts.timer) {
    process.emitWarning(
      'specifying the `timer` option to `new Transaction()` was never a public API and will be removed',
      'DeprecationWarning',
      'ELASTIC_APM_SPAN_TIMER_OPTION'
    )
  }
  if (opts.tracestate) {
    opts.tracestate = TraceState.fromStringFormatString(opts.tracestate)
  }

  if (opts.childOf) {
    // Possibly restart the trace, depending on `traceContinuationStrategy`.
    // Spec: https://github.com/elastic/apm/blob/main/specs/agents/trace-continuation.md
    let traceContinuationStrategy = agent._conf.traceContinuationStrategy
    if (traceContinuationStrategy === TRACE_CONTINUATION_STRATEGY_RESTART_EXTERNAL) {
      traceContinuationStrategy = TRACE_CONTINUATION_STRATEGY_RESTART
      if (opts.tracestate && opts.tracestate.toMap().has('es')) {
        traceContinuationStrategy = TRACE_CONTINUATION_STRATEGY_CONTINUE
      }
    }
    if (traceContinuationStrategy === TRACE_CONTINUATION_STRATEGY_RESTART) {
      if (!opts.links || !Array.isArray(opts.links)) {
        opts.links = []
      }
      opts.links.push({ context: opts.childOf })
      delete opts.childOf // restart the trace
      delete opts.tracestate
    }
  }

  GenericSpan.call(this, agent, ...args, opts)

  const verb = this.parentId ? 'continue' : 'start'
  agent.logger.debug('%s trace %o', verb, { trans: this.id, parent: this.parentId, trace: this.traceId, name: this.name, type: this.type, subtype: this.subtype, action: this.action })

  this._defaultName = name || ''
  this._customName = ''
  this._user = null
  this._custom = null
  this._result = constants.RESULT_SUCCESS
  this._builtSpans = 0
  this._droppedSpans = 0
  this._breakdownTimings = new ObjectIdentityMap()
  this._faas = undefined
  this._service = undefined
  this._message = undefined
  this._cloud = undefined
  this._droppedSpanStats = new DroppedSpanStats()
  this.outcome = constants.OUTCOME_UNKNOWN
}

Object.defineProperty(Transaction.prototype, 'name', {
  configurable: true,
  enumerable: true,
  get () {
    // Fall back to a somewhat useful name in case no _defaultName is set.
    // This might happen if res.writeHead wasn't called.
    return this._customName ||
      this._defaultName ||
      (this.req ? this.req.method + ' unknown route (unnamed)' : 'unnamed')
  },
  set (name) {
    if (this.ended) {
      this._agent.logger.debug('tried to set transaction.name on already ended transaction %o', { trans: this.id, parent: this.parentId, trace: this.traceId })
      return
    }
    this._agent.logger.debug('setting transaction name %o', { trans: this.id, parent: this.parentId, trace: this.traceId, name: name })
    this._customName = name
  }
})

Object.defineProperty(Transaction.prototype, 'result', {
  configurable: true,
  enumerable: true,
  get () {
    return this._result
  },
  set (result) {
    if (this.ended) {
      this._agent.logger.debug('tried to set transaction.result on already ended transaction %o', { trans: this.id, parent: this.parentId, trace: this.traceId })
      return
    }
    this._agent.logger.debug('setting transaction result %o', { trans: this.id, parent: this.parentId, trace: this.traceId, result: result })
    this._result = result
  }
})

Object.defineProperty(Transaction.prototype, 'ids', {
  get () {
    return this._ids === null
      ? (this._ids = new TransactionIds(this))
      : this._ids
  }
})

Transaction.prototype.toString = function () {
  return this.ids.toString()
}

Transaction.prototype.setUserContext = function (context) {
  if (!context) return
  this._user = Object.assign(this._user || {}, context)
}

Transaction.prototype.setServiceContext = function (serviceContext) {
  if (!serviceContext) return
  this._service = Object.assign(this._service || {}, serviceContext)
}

Transaction.prototype.setMessageContext = function (messageContext) {
  if (!messageContext) return
  this._message = Object.assign(this._message || {}, messageContext)
}

Transaction.prototype.setFaas = function (faasFields) {
  if (!faasFields) return
  this._faas = Object.assign(this._faas || {}, faasFields)
}

Transaction.prototype.setCustomContext = function (context) {
  if (!context) return
  this._custom = Object.assign(this._custom || {}, context)
}

Transaction.prototype.setCloudContext = function (cloudContext) {
  if (!cloudContext) return
  this._cloud = Object.assign(this._cloud || {}, cloudContext)
}

// Create a span on this transaction and make it the current span.
Transaction.prototype.startSpan = function (...args) {
  const span = this.createSpan(...args)
  if (span) {
    this._agent._instrumentation.supersedeWithSpanRunContext(span)
  }
  return span
}

// Create a span on this transaction.
//
// This does *not* replace the current run context to make this span the
// "current" one. This allows instrumentations to avoid impacting the run
// context of the calling code. Compare to `startSpan`.
Transaction.prototype.createSpan = function (...args) {
  if (!this.sampled) {
    return null
  }

  // Exit spans must not have child spans (unless of the same type and subtype).
  // https://github.com/elastic/apm/blob/master/specs/agents/tracing-spans.md#child-spans-of-exit-spans
  const opts = typeof args[args.length - 1] === 'object'
    ? (args.pop() || {})
    : {}
  const [_name, type, subtype] = args // eslint-disable-line no-unused-vars
  opts.childOf = opts.childOf || this._agent._instrumentation.currSpan() || this
  const childOf = opts.childOf
  if (childOf instanceof Span && childOf._exitSpan &&
      !(childOf.type === type && childOf.subtype === subtype)) {
    this._agent.logger.trace({ exitSpanId: childOf.id, newSpanArgs: args },
      'createSpan: drop child span of exit span')
    return null
  }

  const span = new Span(this, ...args, opts)

  if (this._builtSpans >= this._agent._conf.transactionMaxSpans) {
    this._droppedSpans++
    span.setRecorded(false)
  }

  this._builtSpans++
  return span
}

// Note that this only returns a complete result when called *during* the call
// to `transaction.end()`.
Transaction.prototype.toJSON = function () {
  var payload = {
    id: this.id,
    trace_id: this.traceId,
    parent_id: this.parentId,
    name: this.name,
    type: this.type || constants.DEFAULT_SPAN_TYPE,
    duration: this._duration,
    timestamp: this.timestamp,
    result: String(this.result),
    sampled: this.sampled,
    context: undefined,
    span_count: {
      started: this._builtSpans - this._droppedSpans
    },
    outcome: this.outcome,
    faas: this._faas
  }

  if (this.sampled) {
    payload.context = {
      user: Object.assign(
        {},
        this.req && parsers.getUserContextFromRequest(this.req),
        this._user
      ),
      tags: this._labels || {},
      custom: this._custom || {},
      service: this._service || {},
      cloud: this._cloud || {},
      message: this._message || {}
    }
    // Only include dropped count when spans have been dropped.
    if (this._droppedSpans > 0) {
      payload.span_count.dropped = this._droppedSpans
    }

    var conf = this._agent._conf
    if (this.req) {
      payload.context.request = parsers.getContextFromRequest(this.req, conf, 'transactions')
    }
    if (this.res) {
      payload.context.response = parsers.getContextFromResponse(this.res, conf)
    }
  }

  // add sample_rate to transaction
  // https://github.com/elastic/apm/blob/main/specs/agents/tracing-sampling.md
  // Only set sample_rate on transaction payload if a valid trace state
  // variable is set.
  //
  // "If there is no tracestate or no valid es entry with an s attribute,
  //  then the agent must omit sample_rate from non-root transactions and
  //  their spans."
  const sampleRate = this.sampleRate
  if (sampleRate !== null) {
    payload.sample_rate = sampleRate
  }

  this._serializeOTel(payload)

  if (this._links.length > 0) {
    payload.links = this._links
  }

  if (this._droppedSpanStats.size() > 0) {
    payload.dropped_spans_stats = this._droppedSpanStats.encode()
  }
  return payload
}

// Note that this only returns a complete result when called *during* the call
// to `transaction.end()`.
Transaction.prototype._encode = function () {
  if (!this.ended) {
    this._agent.logger.error('cannot encode un-ended transaction: %o', { trans: this.id, parent: this.parentId, trace: this.traceId })
    return null
  }

  return this.toJSON()
}

Transaction.prototype.setDefaultName = function (name) {
  this._agent.logger.debug('setting default transaction name: %s %o', name, { trans: this.id, parent: this.parentId, trace: this.traceId })
  this._defaultName = name
}

Transaction.prototype.setDefaultNameFromRequest = function () {
  var req = this.req
  var path = getPathFromRequest(req, false, this._agent._conf.usePathAsTransactionName)

  if (!path) {
    this._agent.logger.debug('could not extract route name from request %o', {
      url: req.url,
      type: typeof path,
      null: path === null, // because typeof null === 'object'
      route: !!req.route,
      regex: req.route ? !!req.route.regexp : false,
      mountstack: req[symbols.expressMountStack] ? req[symbols.expressMountStack].length : false,
      trans: this.id,
      parent: this.parentId,
      trace: this.traceId
    })
    path = 'unknown route'
  }

  this.setDefaultName(req.method + ' ' + path)
}

Transaction.prototype.ensureParentId = function () {
  return this._context.ensureParentId()
}

Transaction.prototype.end = function (result, endTime) {
  if (this.ended) {
    this._agent.logger.debug('tried to call transaction.end() on already ended transaction %o', { trans: this.id, parent: this.parentId, trace: this.traceId })
    return
  }

  if (result !== undefined && result !== null) {
    this.result = result
  }

  if (!this._defaultName && this.req) this.setDefaultNameFromRequest()

  this._timer.end(endTime)
  this._duration = this._timer.duration
  this._captureBreakdown(this)
  this.ended = true

  this._agent._instrumentation.addEndedTransaction(this)
  this._agent.logger.debug({ trans: this.id, name: this.name, parent: this.parentId, trace: this.traceId, type: this.type, result: this.result, duration: this._duration },
    'ended transaction')

  // Reduce this transaction's memory usage by dropping references except to
  // fields required to support `interface Transaction`.
  // Transaction fields:
  this._customName = this.name // Short-circuit the `name` getter.
  this._defaultName = ''
  this.req = null
  this.res = null
  this._user = null
  this._custom = null
  this._breakdownTimings = null
  this._faas = undefined
  this._service = undefined
  this._message = undefined
  this._cloud = undefined
  // GenericSpan fields:
  // - Cannot drop `this._context` because it is used for `traceparent`, `ids`,
  //   and `.sampled` (when capturing breakdown metrics for child spans).
  this._timer = null
  this._labels = null
}

Transaction.prototype.setOutcome = function (outcome) {
  if (!this._isValidOutcome(outcome)) {
    this._agent.logger.trace(
      'Unknown outcome [%s] seen in Transaction.setOutcome, ignoring',
      outcome
    )
    return
  }

  if (this.ended) {
    this._agent.logger.debug(
      'tried to call Transaction.setOutcome() on already ended transaction %o',
      { trans: this.id, parent: this.parentId, trace: this.traceId })
    return
  }

  this._freezeOutcome()
  this.outcome = outcome
}

Transaction.prototype._setOutcomeFromHttpStatusCode = function (statusCode) {
  // if an outcome's been set from the API we
  // honor its value
  if (this._isOutcomeFrozen) {
    return
  }

  if (statusCode >= 500) {
    this.outcome = constants.OUTCOME_FAILURE
  } else {
    this.outcome = constants.OUTCOME_SUCCESS
  }
}

Transaction.prototype._captureBreakdown = function (span) {
  if (this.ended) {
    return
  }

  const agent = this._agent
  const metrics = agent._metrics
  const conf = agent._conf

  // Avoid unneeded breakdown metrics processing if only propagating trace context.
  if (conf.contextPropagationOnly) {
    return
  }

  // Record span data
  if (this.sampled && conf.breakdownMetrics) {
    captureBreakdown(this, {
      transaction: transactionBreakdownDetails(this),
      span: spanBreakdownDetails(span)
    }, span._timer.selfTime)
  }

  // Record transaction data
  if (span instanceof Transaction) {
    for (const { labels, time, count } of this._breakdownTimings.values()) {
      const flattenedLabels = flattenBreakdown(labels)
      metrics.incrementCounter('span.self_time.count', flattenedLabels, count)
      metrics.incrementCounter('span.self_time.sum.us', flattenedLabels, time)
    }
  }
}

Transaction.prototype.captureDroppedSpan = function (span) {
  return this._droppedSpanStats.captureDroppedSpan(span)
}

function transactionBreakdownDetails ({ name, type } = {}) {
  return {
    name,
    type
  }
}

function spanBreakdownDetails (span) {
  if (span instanceof Transaction) {
    return {
      type: 'app'
    }
  }

  const { type, subtype } = span
  return {
    type,
    subtype
  }
}

function captureBreakdown (transaction, labels, time) {
  const build = () => ({ labels, count: 0, time: 0 })
  const counter = transaction._breakdownTimings.ensure(labels, build)
  counter.time += time
  counter.count++
}

function flattenBreakdown (source, target = {}, prefix = '') {
  for (const [key, value] of Object.entries(source)) {
    if (typeof value === 'undefined' || value === null) continue
    if (typeof value === 'object') {
      flattenBreakdown(value, target, `${prefix}${key}::`)
    } else {
      target[`${prefix}${key}`] = value
    }
  }

  return target
}
