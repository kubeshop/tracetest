/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

var http = require('http')
var path = require('path')

var isError = require('core-util-is').isError
var Filters = require('object-filter-sequence')

var config = require('./config')
var connect = require('./middleware/connect')
const constants = require('./constants')
var errors = require('./errors')
const { InflightEventSet } = require('./InflightEventSet')
var Instrumentation = require('./instrumentation')
var { elasticApmAwsLambda } = require('./lambda')
var Metrics = require('./metrics')
var parsers = require('./parsers')
var symbols = require('./symbols')
const { frameCacheStats, initStackTraceCollection } = require('./stacktraces')
const Span = require('./instrumentation/span')
const Transaction = require('./instrumentation/transaction')

var IncomingMessage = http.IncomingMessage
var ServerResponse = http.ServerResponse

var version = require('../package').version

module.exports = Agent

function Agent () {
  // Early configuration to ensure `agent.logger` works before `agent.start()`.
  this.logger = config.configLogger()

  // Get an initial pre-.start() configuration of agent defaults. This is a
  // crutch for Agent APIs that depend on `agent._conf`.
  this._conf = config.initialConfig(this.logger)

  this._httpClient = null
  this._uncaughtExceptionListener = null
  this._inflightEvents = new InflightEventSet()
  this._instrumentation = new Instrumentation(this)
  this._metrics = new Metrics(this)
  this._errorFilters = new Filters()
  this._transactionFilters = new Filters()
  this._spanFilters = new Filters()
  this._transport = null

  this.lambda = elasticApmAwsLambda(this)
  this.middleware = { connect: connect.bind(this) }
}

Object.defineProperty(Agent.prototype, 'currentTransaction', {
  get () {
    return this._instrumentation.currTransaction()
  }
})

Object.defineProperty(Agent.prototype, 'currentSpan', {
  get () {
    return this._instrumentation.currSpan()
  }
})

Object.defineProperty(Agent.prototype, 'currentTraceparent', {
  get () {
    const current = this._instrumentation.currSpan() || this._instrumentation.currTransaction()
    return current ? current.traceparent : null
  }
})

Object.defineProperty(Agent.prototype, 'currentTraceIds', {
  get () {
    return this._instrumentation.ids()
  }
})

// Destroy this agent. This prevents any new agent processing, communication
// with APM server, and resets changed global state *as much as is possible*.
//
// In the typical uses case -- a singleton Agent running for the full process
// lifetime -- it is *not* necessary to call `agent.destroy()`.  It is used
// for some testing.
//
// Limitations:
// - Patching/wrapping of functions for instrumentation *is* undone, but
//   references to the wrapped versions can remain.
// - There may be in-flight tasks (in ins.addEndedSpan() and
//   agent.captureError() for example) that will complete after this destroy
//   completes. They should have no impact other than CPU/resource use.
// - The patching of core node functions when `contextManager="patch"` is *not*
//   undone. This means run context tracking for `contextManager="patch"` is
//   broken with in-process multiple-Agent use.
Agent.prototype.destroy = function () {
  if (this._transport && this._transport.destroy) {
    this._transport.destroy()
  }
  // So in-flight tasks in ins.addEndedSpan() and agent.captureError() do
  // not use the destroyed transport.
  this._transport = null

  // So in-flight tasks do not call user-added filters after the agent has
  // been destroyed.
  this._errorFilters = new Filters()
  this._transactionFilters = new Filters()
  this._spanFilters = new Filters()

  if (this._uncaughtExceptionListener) {
    process.removeListener('uncaughtException', this._uncaughtExceptionListener)
  }
  this._metrics.stop()
  this._instrumentation.stop()

  // Allow a new Agent instance to `.start()`. Typically this is only relevant
  // for tests that may use multiple Agent instances in a single test process.
  global[symbols.agentInitialized] = null

  if (this._origStackTraceLimit && Error.stackTraceLimit !== this._origStackTraceLimit) {
    Error.stackTraceLimit = this._origStackTraceLimit
  }
}

// These are metrics about the agent itself -- separate from the metrics
// gathered on behalf of the using app and sent to APM server. Currently these
// are only useful for internal debugging of the APM agent itself.
//
// **These stats are NOT a promised interface.**
Agent.prototype._getStats = function () {
  const stats = {
    frameCache: frameCacheStats
  }
  if (this._instrumentation._runCtxMgr && this._instrumentation._runCtxMgr._runContextFromAsyncId) {
    stats.runContextFromAsyncIdSize = this._instrumentation._runCtxMgr._runContextFromAsyncId.size
  }
  if (this._transport && typeof this._transport._getStats === 'function') {
    stats.apmclient = this._transport._getStats()
  }
  return stats
}

Agent.prototype.addPatch = function (modules, handler) {
  return this._instrumentation.addPatch.apply(this._instrumentation, arguments)
}

Agent.prototype.removePatch = function (modules, handler) {
  return this._instrumentation.removePatch.apply(this._instrumentation, arguments)
}

Agent.prototype.clearPatches = function (modules) {
  return this._instrumentation.clearPatches.apply(this._instrumentation, arguments)
}

Agent.prototype.startTransaction = function (name, type, subtype, action, { startTime, childOf } = {}) {
  return this._instrumentation.startTransaction.apply(this._instrumentation, arguments)
}

Agent.prototype.endTransaction = function (result, endTime) {
  return this._instrumentation.endTransaction.apply(this._instrumentation, arguments)
}

Agent.prototype.setTransactionName = function (name) {
  return this._instrumentation.setTransactionName.apply(this._instrumentation, arguments)
}

/**
 * Sets outcome value for current transaction
 *
 * The setOutcome method allows users to override the default
 * outcome handling in the agent and set their own value.
 *
 * @param {string} outcome must be one of `failure`, `success`, or `unknown`
 */
Agent.prototype.setTransactionOutcome = function (outcome) {
  return this._instrumentation.setTransactionOutcome.apply(this._instrumentation, arguments)
}

Agent.prototype.startSpan = function (name, type, subtype, action, { startTime, childOf, exitSpan } = {}) {
  return this._instrumentation.startSpan.apply(this._instrumentation, arguments)
}

/**
 * Sets outcome value for current active span
 *
 * The setOutcome method allows users to override the default
 * outcome handling in the agent and set their own value.
 *
 * @param {string} outcome must be one of `failure`, `success`, or `unknown`
 */
Agent.prototype.setSpanOutcome = function (outcome) {
  return this._instrumentation.setSpanOutcome.apply(this._instrumentation, arguments)
}

Agent.prototype._config = function (opts) {
  this._conf = config.createConfig(opts, this.logger)
  this.logger = this._conf.logger
}

Agent.prototype.isStarted = function () {
  return global[symbols.agentInitialized]
}

Agent.prototype.start = function (opts) {
  if (this.isStarted()) {
    throw new Error('Do not call .start() more than once')
  }
  global[symbols.agentInitialized] = true

  this._config(opts)

  if (this._conf.filterHttpHeaders) {
    this.addFilter(require('./filters/http-headers'))
  }

  // Check cases where we do *not* start.
  if (!this._conf.active) {
    this.logger.debug('Elastic APM agent disabled (`active` is false)')
    return this
  } else if (!this._conf.serviceName) {
    this.logger.error('Elastic APM is incorrectly configured: Missing serviceName (APM will be disabled)')
    this._conf.active = false
    return this
  }
  // Sanity check the port from `serverUrl`.
  const parsedUrl = parsers.parseUrl(this._conf.serverUrl)
  const serverPort = (parsedUrl.port
    ? Number(parsedUrl.port)
    : (parsedUrl.protocol === 'https:' ? 443 : 80))
  if (!(serverPort >= 1 && serverPort <= 65535)) {
    this.logger.error('Elastic APM is incorrectly configured: serverUrl "%s" contains an invalid port! (allowed: 1-65535)', this._conf.serverUrl)
    this._conf.active = false
    return this
  }

  if (this._conf.logLevel === 'trace') {
    var stackObj = {}
    Error.captureStackTrace(stackObj)

    // Attempt to load package.json from process.argv.
    var pkg = null
    try {
      var basedir = path.dirname(process.argv[1] || '.')
      pkg = require(path.join(basedir, 'package.json'))
    } catch (e) {}

    this.logger.trace({
      pid: process.pid,
      ppid: process.ppid,
      arch: process.arch,
      platform: process.platform,
      node: process.version,
      agent: version,
      startTrace: stackObj.stack.split(/\n */).slice(1),
      main: pkg ? pkg.main : '<could not determine>',
      dependencies: pkg ? pkg.dependencies : '<could not determine>',
      conf: this._conf.toJSON()
    }, 'agent configured correctly')
  }

  initStackTraceCollection()
  this._transport = this._conf.transport(this._conf, this)

  let runContextClass
  if (this._conf.opentelemetryBridgeEnabled) {
    const { setupOTelBridge, OTelBridgeRunContext } = require('./opentelemetry-bridge')
    runContextClass = OTelBridgeRunContext
    setupOTelBridge(this)
  }
  this._instrumentation.start(runContextClass)
  this._metrics.start()

  this._origStackTraceLimit = Error.stackTraceLimit
  Error.stackTraceLimit = this._conf.stackTraceLimit
  if (this._conf.captureExceptions) this.handleUncaughtExceptions()

  return this
}

Agent.prototype.getServiceName = function () {
  return this._conf ? this._conf.serviceName : undefined
}

Agent.prototype.setFramework = function ({ name, version, overwrite = true }) {
  if (!this._transport || !this._conf) {
    return
  }
  const conf = {}
  if (name && (overwrite || !this._conf.frameworkName)) this._conf.frameworkName = conf.frameworkName = name
  if (version && (overwrite || !this._conf.frameworkVersion)) this._conf.frameworkVersion = conf.frameworkVersion = version
  this._transport.config(conf)
}

Agent.prototype.setUserContext = function (context) {
  var trans = this._instrumentation.currTransaction()
  if (!trans) return false
  trans.setUserContext(context)
  return true
}

Agent.prototype.setCustomContext = function (context) {
  var trans = this._instrumentation.currTransaction()
  if (!trans) return false
  trans.setCustomContext(context)
  return true
}

Agent.prototype.setLabel = function (key, value, stringify) {
  var trans = this._instrumentation.currTransaction()
  if (!trans) return false
  return trans.setLabel(key, value, stringify)
}

Agent.prototype.addLabels = function (labels, stringify) {
  var trans = this._instrumentation.currTransaction()
  if (!trans) return false
  return trans.addLabels(labels, stringify)
}

Agent.prototype.addFilter = function (fn) {
  this.addErrorFilter(fn)
  this.addTransactionFilter(fn)
  this.addSpanFilter(fn)
  // Note: This does *not* add to *metadata* filters, partly for backward
  // compat -- the structure of metadata objects is quite different and could
  // break existing filters -- and partly because that different structure
  // means it makes less sense to re-use the same function to filter them.
}

Agent.prototype.addErrorFilter = function (fn) {
  if (typeof fn !== 'function') {
    this.logger.error('Can\'t add filter of type %s', typeof fn)
    return
  }

  this._errorFilters.push(fn)
}

Agent.prototype.addTransactionFilter = function (fn) {
  if (typeof fn !== 'function') {
    this.logger.error('Can\'t add filter of type %s', typeof fn)
    return
  }

  this._transactionFilters.push(fn)
}

Agent.prototype.addSpanFilter = function (fn) {
  if (typeof fn !== 'function') {
    this.logger.error('Can\'t add filter of type %s', typeof fn)
    return
  }

  this._spanFilters.push(fn)
}

Agent.prototype.addMetadataFilter = function (fn) {
  if (typeof fn !== 'function') {
    this.logger.error('Can\'t add filter of type %s', typeof fn)
    return
  } else if (!this._transport) {
    this.logger.error('cannot add metadata filter to inactive or unconfigured agent (agent has no transport)')
    return
  } else if (typeof this._transport.addMetadataFilter !== 'function') {
    // Graceful failure if unexpectedly using a too-old APM client.
    this.logger.error('cannot add metadata filter: transport does not support addMetadataFilter')
    return
  }

  // Metadata filters are handled by the APM client, where metadata is
  // processed.
  this._transport.addMetadataFilter(fn)
}

const EMPTY_OPTS = {}

// Capture an APM server "error" event for the given `err` and send it to APM
// server.
//
// Usage:
//    captureError(err, opts, cb)
//    captureError(err, opts)
//    captureError(err, cb)
//
// where:
// - `err` is an Error instance, or a string message, or a "parameterized string
//   message" object, e.g.:
//      {
//        message: "this is my message template: %d %s"},
//        params: [ 42, "another param" ]
//      }
// - `opts` can include any of the following (all optional):
//   - `opts.timestamp` - Milliseconds since the Unix epoch. Defaults to now.
//   - `opts.user` - Object to add to `error.context.user`.
//   - `opts.tags` - Deprecated, use `opts.labels`. Object to add to
//     `error.context.labels`.
//   - `opts.labels` - Object to add to `error.context.labels`.
//   - `opts.custom` - Object to add to `error.context.custom`.
//   - `opts.message` - If `err` is an Error instance, this string is added to
//     `error.log.message` (unless it matches err.message).
//   - `opts.request` - HTTP request (node `IncomingMessage` instance) to use
//     for `error.context.request`. If not given, a `req` on the parent
//     transaction (see `opts.parent` below) will be used.
//   - `opts.response` - HTTP response (node `ServerResponse` instance) to use
//     for `error.context.response`. If not given, a `res` on the parent
//     transaction (see `opts.parent` below) will be used.
//   - `opts.handled` - Boolean indicating if this exception was handled by
//     application code. Default true. Setting to `false` also results in the
//     error being flushed to APM server as soon as it is processed.
//   - `opts.captureAttributes` - Boolean. Default true. Set to false to *not*
//     include properties of `err` as attributes on the APM error event.
//   - `opts.skipOutcome` - Boolean. Default false. Set to true to not have
//     this captured error set `<currentSpan>.outcome = failure`.
//   - `opts.parent` - A Transaction or Span instance to make the parent of
//     this error. If not given (undefined), then the current span or
//     transaction will be used. If `null` is given, then no span or transaction
//     will be used.
//   - `opts.exceptionType` - A string to use for `error.exception.type`. By
//     default this is `err.name`. This option is only relevant if `err` is an
//     Error instance.
// - `cb` is a callback `function (captureErr, apmErrorIdString)`. If provided,
//   the error will be flushed to APM server as soon as it is processed, and
//   `cb` will be called when that send is complete.
Agent.prototype.captureError = function (err, opts, cb) {
  if (typeof opts === 'function') {
    cb = opts
    opts = EMPTY_OPTS
  } else if (!opts) {
    opts = EMPTY_OPTS
  }

  const id = errors.generateErrorId()

  if (!this.isStarted()) {
    if (cb) {
      cb(new Error('cannot capture error before agent is started'), id)
    }
    return
  }

  // Avoid unneeded error/stack processing if only propagating trace-context.
  if (this._conf.contextPropagationOnly) {
    if (cb) {
      process.nextTick(cb, null, id)
    }
    return
  }

  const agent = this
  let callSiteLoc = null
  const errIsError = isError(err)
  const handled = opts.handled !== false // default true
  const shouldCaptureAttributes = opts.captureAttributes !== false // default true
  const skipOutcome = Boolean(opts.skipOutcome)
  const timestampUs = (opts.timestamp
    ? Math.floor(opts.timestamp * 1000)
    : Date.now() * 1000)

  // Determine transaction/span context to associate with this error.
  let parent
  let span
  let trans
  if (opts.parent === undefined) {
    parent = this._instrumentation.currSpan() || this._instrumentation.currTransaction()
  } else if (opts.parent === null) {
    parent = null
  } else {
    parent = opts.parent
  }
  if (parent instanceof Transaction) {
    span = null
    trans = parent
  } else if (parent instanceof Span) {
    span = parent
    trans = parent.transaction
  }
  const traceContext = (span || trans || {})._context
  const req = (opts.request instanceof IncomingMessage
    ? opts.request
    : trans && trans.req)
  const res = (opts.response instanceof ServerResponse
    ? opts.response
    : trans && trans.res)

  // As an added feature, for *some* cases, we capture a stacktrace at the point
  // this `captureError` was called. This is added to `error.log.stacktrace`.
  if (handled &&
      (agent._conf.captureErrorLogStackTraces === config.CAPTURE_ERROR_LOG_STACK_TRACES_ALWAYS ||
       (!errIsError && agent._conf.captureErrorLogStackTraces === config.CAPTURE_ERROR_LOG_STACK_TRACES_MESSAGES))
  ) {
    callSiteLoc = {}
    Error.captureStackTrace(callSiteLoc, Agent.prototype.captureError)
  }

  if (span && !skipOutcome) {
    span._setOutcomeFromErrorCapture(constants.OUTCOME_FAILURE)
  }

  // Note this error as an "inflight" event. See Agent#flush().
  const inflightEvents = this._inflightEvents
  inflightEvents.add(id)

  // Move the remaining captureError processing to a later tick because:
  // 1. This allows the calling code to continue processing. For example, for
  //    Express instrumentation this can significantly improve latency in
  //    the app's endpoints because the response does not proceed until the
  //    error handlers return.
  // 2. Gathering `error.context.response` in the same tick results in data
  //    for a response that hasn't yet completed (no headers, unset status_code,
  //    etc.).
  setImmediate(() => {
    // Gather `error.context.*`.
    const errorContext = {
      user: Object.assign(
        {},
        req && parsers.getUserContextFromRequest(req),
        trans && trans._user,
        opts.user
      ),
      tags: Object.assign(
        {},
        trans && trans._labels,
        opts.tags,
        opts.labels
      ),
      custom: Object.assign(
        {},
        trans && trans._custom,
        opts.custom
      )
    }
    if (req) {
      errorContext.request = parsers.getContextFromRequest(req, agent._conf, 'errors')
    }
    if (res) {
      errorContext.response = parsers.getContextFromResponse(res, agent._conf, true)
    }

    errors.createAPMError({
      log: agent.logger,
      id: id,
      exception: errIsError ? err : null,
      logMessage: errIsError ? null : err,
      shouldCaptureAttributes,
      timestampUs,
      handled,
      callSiteLoc,
      message: opts.message,
      sourceLinesAppFrames: agent._conf.sourceLinesErrorAppFrames,
      sourceLinesLibraryFrames: agent._conf.sourceLinesErrorLibraryFrames,
      trans,
      traceContext,
      errorContext,
      exceptionType: opts.exceptionType
    }, function filterAndSendError (_err, apmError) {
      // _err is always null from createAPMError.

      apmError = agent._errorFilters.process(apmError)
      if (!apmError) {
        agent.logger.debug('error ignored by filter %o', { id })
        inflightEvents.delete(id)
        if (cb) {
          cb(null, id)
        }
        return
      }

      if (agent._transport) {
        agent.logger.info('Sending error to Elastic APM: %o', { id })
        agent._transport.sendError(apmError)
        inflightEvents.delete(id)
        if (!handled || cb) {
          // Immediately flush *unhandled* errors -- those from
          // `uncaughtException` -- on the assumption that the process may
          // soon crash. Also flush when a `cb` is provided.
          agent.flush(function (flushErr) {
            if (cb) {
              cb(flushErr, id)
            }
          })
        }
      } else {
        inflightEvents.delete(id)
        if (cb) {
          cb(new Error('cannot send error: missing transport'), id)
        }
      }
    })
  })
}

// The optional callback will be called with the error object after the error
// have been sent to the intake API. If no callback have been provided we will
// automatically terminate the process, so if you provide a callback you must
// remember to terminate the process manually.
Agent.prototype.handleUncaughtExceptions = function (cb) {
  var agent = this

  if (this._uncaughtExceptionListener) {
    process.removeListener('uncaughtException', this._uncaughtExceptionListener)
  }

  this._uncaughtExceptionListener = function (err) {
    agent.logger.debug({ err }, 'Elastic APM caught unhandled exception')
    // The stack trace of uncaught exceptions are normally written to STDERR.
    // The `uncaughtException` listener inhibits this behavior, and it's
    // therefore necessary to manually do this to not break expectations.
    if (agent._conf && agent._conf.logUncaughtExceptions === true) {
      console.error(err)
    }

    agent.captureError(err, { handled: false }, function () {
      cb ? cb(err) : process.exit(1)
    })
  }

  process.on('uncaughtException', this._uncaughtExceptionListener)
}

// Flush all ended APM events (transactions, spans, errors, metricsets) to APM
// server as soon as possible. If the optional `cb` is given, it will be called
// `cb(flushErr)` when this is complete.
//
// Encoding and passing event data to the agent's transport is *asynchronous*
// for some event types: spans and errors. This flush will make a *best effort*
// attempt to wait for those "inflight" events to finish processing before
// flushing data to APM server. To avoid `.flush()` hanging, this times out
// after one second.
//
// If flush is called while still waiting for inflight events in an earlier
// flush call, then the more recent flush will only wait for events that were
// newly inflight *since the last .flush()* call. I.e. the second flush does
// *not* wait for the set of events the first flush is waiting on. This makes
// the semantics of flush less than ideal (one cannot blindly call .flush() to
// flush *everything* that has come before). However it handles the common use
// case of flushing synchronously after ending a span or capturing an error:
//    mySpan.end()
//    apm.flush(function () { ... })
// and it simplifies the implementation.
//
// # Dev Notes
//
// To support the implementation, agent code that creates an inflight event
// must do the following:
// - Take a reference to the current set of inflight events:
//      const inflightEvents = agent._inflightEvents
// - Add a unique ID for the event to the set:
//      inflightEvents.add(id)
// - Delete the ID from the set when sent to the transport (`.sendSpan(...)` et
//   al) or when dropping the event (e.g. because of a filter):
//      inflightEvents.delete(id)
Agent.prototype.flush = function (cb) {
  // This 1s timeout is a subjective balance between "long enough for spans
  // and errors to reasonably encode" and "short enough to not block data
  // being reported to APM server".
  const DEFAULT_INFLIGHT_FLUSH_TIMEOUT_MS = 1000

  return this._flush({ inflightTimeoutMs: DEFAULT_INFLIGHT_FLUSH_TIMEOUT_MS }, cb)
}

// The internal-use `.flush()` that supports some options not exposed to the
// public API.
//
// @param {Number} opts.inflightTimeoutMs - Required. The number of ms to wait
//    for inflight events (spans, errors) to finish being send to the transport
//    before flushing.
// @param {Boolean} opts.lambdaEnd - Optional, default false. Set this to true
//    to signal to the transport that this is the flush at the end of a Lambda
//    function invocation.
//    https://github.com/elastic/apm/blob/main/specs/agents/tracing-instrumentation-aws-lambda.md#data-flushing
Agent.prototype._flush = function (opts, cb) {
  const lambdaEnd = !!opts.lambdaEnd

  if (!this._transport) {
    // Log an *err* to provide a stack for the user.
    const err = new Error('cannot flush agent before it is started')
    this.logger.warn({ err }, err.message)
    if (cb) {
      process.nextTick(cb)
    }
    return
  }

  const boundCb = cb && this._instrumentation.bindFunction(cb)

  // If there are no inflight events then avoid creating additional objects.
  if (this._inflightEvents.size === 0) {
    this._transport.flush({ lambdaEnd }, boundCb)
    return
  }

  // Otherwise, there are inflight events to wait for.  Setup a handler to
  // callback when the current set of inflight events complete.
  const flushingInflightEvents = this._inflightEvents
  flushingInflightEvents.setDrainHandler((drainErr) => {
    // The only possible drainErr is a timeout. This is best effort, so we only
    // log this and move on.
    this.logger.debug({
      numRemainingInflightEvents: flushingInflightEvents.size,
      err: drainErr
    }, 'flush: drained inflight events')

    // Then, flush the intake request to APM server.
    this._transport.flush({ lambdaEnd }, boundCb)
  }, opts.inflightTimeoutMs)

  // Create a new empty set to collect subsequent inflight events.
  this._inflightEvents = new InflightEventSet()
}

Agent.prototype.registerMetric = function (name, labelsOrCallback, callback) {
  var labels
  if (typeof labelsOrCallback === 'function') {
    callback = labelsOrCallback
  } else {
    labels = labelsOrCallback
  }

  if (typeof callback !== 'function') {
    this.logger.error('Can\'t add callback of type %s', typeof callback)
    return
  }

  this._metrics.getOrCreateGauge(name, callback, labels)
}
