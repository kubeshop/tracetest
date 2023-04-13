/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

var fs = require('fs')
var path = require('path')

var hook = require('require-in-the-middle')
const semver = require('semver')

const config = require('../config')
var { Ids } = require('./ids')
var NamedArray = require('./named-array')
var Transaction = require('./transaction')
const {
  BasicRunContextManager,
  AsyncHooksRunContextManager,
  AsyncLocalStorageRunContextManager
} = require('./run-context')
const { getLambdaHandlerInfo } = require('../lambda')
const undiciInstr = require('./modules/undici')

const nodeSupportsAsyncLocalStorage = semver.satisfies(process.versions.node, '>=14.5 || ^12.19.0')
// Node v16.5.0 added fetch support (behind `--experimental-fetch` until
// v18.0.0) based on undici@5.0.0. We can instrument undici >=v4.7.1.
const nodeHasInstrumentableFetch = typeof (global.fetch) === 'function'

var MODULES = [
  ['@elastic/elasticsearch', '@elastic/elasticsearch-canary'],
  '@node-redis/client/dist/lib/client',
  '@node-redis/client/dist/lib/client/commands-queue',
  '@redis/client/dist/lib/client',
  '@redis/client/dist/lib/client/commands-queue',
  'apollo-server-core',
  'aws-sdk',
  'bluebird',
  'cassandra-driver',
  'elasticsearch',
  'express',
  'express-graphql',
  'express-queue',
  'fastify',
  'finalhandler',
  'generic-pool',
  'graphql',
  'handlebars',
  ['hapi', '@hapi/hapi'],
  'http',
  'https',
  'http2',
  'ioredis',
  'jade',
  'knex',
  'koa',
  ['koa-router', '@koa/router'],
  'memcached',
  'mimic-response',
  'mongodb-core',
  'mongodb',
  'mysql',
  'mysql2',
  'next/dist/server/api-utils/node',
  'next/dist/server/dev/next-dev-server',
  'next/dist/server/next',
  'next/dist/server/next-server',
  'pg',
  'pug',
  'redis',
  'restify',
  'tedious',
  'undici',
  'ws'
]

module.exports = Instrumentation

function Instrumentation (agent) {
  this._agent = agent
  this._disableInstrumentationsSet = null
  this._hook = null // this._hook is only exposed for testing purposes
  this._started = false
  this._runCtxMgr = null
  this._log = agent.logger

  // NOTE: we need to track module names for patches
  // in a separate array rather than using Object.keys()
  // because the array is given to the hook(...) call.
  this._patches = new NamedArray()

  for (let mod of MODULES) {
    if (!Array.isArray(mod)) mod = [mod]
    const pathName = mod[0]

    this.addPatch(mod, (...args) => {
      // Lazy require so that we don't have to use `require.resolve` which
      // would fail in combination with Webpack. For more info see:
      // https://github.com/elastic/apm-agent-nodejs/pull/957
      return require(`./modules/${pathName}.js`)(...args)
    })
  }

  // patch for lambda handler needs special handling since its
  // module name will always be different than its handler name
  this._lambdaHandlerInfo = getLambdaHandlerInfo(process.env, MODULES, this._log)
  if (this._lambdaHandlerInfo) {
    this.addPatch(
      this._lambdaHandlerInfo.filePath,
      (...args) => {
        return require('./modules/_lambda-handler')(...args)
      }
    )
  }
}

Instrumentation.prototype.currTransaction = function () {
  if (!this._started) {
    return null
  }
  return this._runCtxMgr.active().currTransaction()
}

Instrumentation.prototype.currSpan = function () {
  if (!this._started) {
    return null
  }
  return this._runCtxMgr.active().currSpan()
}

Instrumentation.prototype.ids = function () {
  if (!this._started) {
    return new Ids()
  }
  const runContext = this._runCtxMgr.active()
  const currSpanOrTrans = runContext.currSpan() || runContext.currTransaction()
  if (currSpanOrTrans) {
    return currSpanOrTrans.ids
  }
  return new Ids()
}

Instrumentation.prototype.addPatch = function (modules, handler) {
  if (!Array.isArray(modules)) {
    modules = [modules]
  }

  for (const mod of modules) {
    const type = typeof handler
    if (type !== 'function' && type !== 'string') {
      this._agent.logger.error('Invalid patch handler type: %s', type)
      return
    }

    this._patches.add(mod, handler)
  }

  this._startHook()
}

Instrumentation.prototype.removePatch = function (modules, handler) {
  if (!Array.isArray(modules)) modules = [modules]

  for (const mod of modules) {
    this._patches.delete(mod, handler)
  }

  this._startHook()
}

Instrumentation.prototype.clearPatches = function (modules) {
  if (!Array.isArray(modules)) modules = [modules]

  for (const mod of modules) {
    this._patches.clear(mod)
  }

  this._startHook()
}

Instrumentation.modules = Object.freeze(MODULES)

// Start the instrumentation system.
//
// @param {RunContext} [runContextClass] - A class to use for the core object
//    that is used to track run context. It defaults to `RunContext`. If given,
//    it must be `RunContext` (the typical case) or a subclass of it. The OTel
//    Bridge uses this to provide a subclass that bridges to OpenTelemetry
//    `Context` usage.
Instrumentation.prototype.start = function (runContextClass) {
  if (this._started) return
  this._started = true

  // Could have changed in Agent.start().
  this._log = this._agent.logger

  // Select the appropriate run-context manager.
  const confContextManager = this._agent._conf.contextManager
  if (confContextManager === config.CONTEXT_MANAGER_PATCH) {
    this._runCtxMgr = new BasicRunContextManager(this._log, runContextClass)
    require('./patch-async')(this)
  } else if (confContextManager === config.CONTEXT_MANAGER_ASYNCHOOKS) {
    this._runCtxMgr = new AsyncHooksRunContextManager(this._log, runContextClass)
  } else if (nodeSupportsAsyncLocalStorage) {
    this._runCtxMgr = new AsyncLocalStorageRunContextManager(this._log, runContextClass)
  } else {
    if (confContextManager === config.CONTEXT_MANAGER_ASYNCLOCALSTORAGE) {
      this._log.warn(`config includes 'contextManager="${confContextManager}"', but node ${process.version} does not support AsyncLocalStorage for run-context management: falling back to using async_hooks`)
    }
    this._runCtxMgr = new AsyncHooksRunContextManager(this._log, runContextClass)
  }

  const patches = this._agent._conf.addPatch
  if (Array.isArray(patches)) {
    for (const [mod, path] of patches) {
      this.addPatch(mod, path)
    }
  }

  this._runCtxMgr.enable()
  this._startHook()

  if (nodeHasInstrumentableFetch && this._isModuleEnabled('undici')) {
    this._log.debug('instrumenting fetch')
    undiciInstr.instrumentUndici(this._agent)
  }
}

// Stop active instrumentation and reset global state *as much as possible*.
//
// Limitations: Removing and re-applying 'require-in-the-middle'-based patches
// has no way to update existing references to patched or unpatched exports from
// those modules.
Instrumentation.prototype.stop = function () {
  this._started = false

  // Reset run context tracking.
  if (this._runCtxMgr) {
    this._runCtxMgr.disable()
    this._runCtxMgr = null
  }

  // Reset patching.
  if (this._hook) {
    this._hook.unhook()
    this._hook = null
  }

  if (nodeHasInstrumentableFetch) {
    undiciInstr.uninstrumentUndici()
  }
}

// Reset internal state for (relatively) clean re-use of this Instrumentation.
// Used for testing, while `resetAgent()` + "test/_agent.js" usage still exists.
//
// This does *not* include redoing monkey patching. It resets context tracking,
// so a subsequent test case can re-use the Instrumentation in the same process.
Instrumentation.prototype.testReset = function () {
  if (this._runCtxMgr) {
    this._runCtxMgr.testReset()
  }
}

Instrumentation.prototype._isModuleEnabled = function (modName) {
  if (!this._disableInstrumentationsSet) {
    this._disableInstrumentationsSet = new Set(this._agent._conf.disableInstrumentations)
  }
  return this._agent._conf.instrument && !this._disableInstrumentationsSet.has(modName)
}

Instrumentation.prototype._startHook = function () {
  if (!this._started) return
  if (this._hook) {
    this._agent.logger.debug('removing hook to Node.js module loader')
    this._hook.unhook()
  }

  var self = this

  this._agent.logger.debug('adding hook to Node.js module loader')

  this._hook = hook(this._patches.keys, function (exports, name, basedir) {
    const enabled = self._isModuleEnabled(name)
    var pkg, version

    const isHandlingLambda = self._lambdaHandlerInfo && self._lambdaHandlerInfo.module === name

    if (!isHandlingLambda && basedir) {
      pkg = path.join(basedir, 'package.json')
      try {
        version = JSON.parse(fs.readFileSync(pkg)).version
      } catch (e) {
        self._agent.logger.debug('could not shim %s module: %s', name, e.message)
        return exports
      }
    } else {
      version = process.versions.node
    }

    return self._patchModule(exports, name, version, enabled)
  })
}

Instrumentation.prototype._patchModule = function (exports, name, version, enabled) {
  this._agent.logger.debug('shimming %s@%s module', name, version)
  const isHandlingLambda = this._lambdaHandlerInfo && this._lambdaHandlerInfo.module === name
  let patches
  if (!isHandlingLambda) {
    patches = this._patches.get(name)
  } else if (name === this._lambdaHandlerInfo.module) {
    patches = this._patches.get(this._lambdaHandlerInfo.filePath)
  }

  if (patches) {
    for (let patch of patches) {
      if (typeof patch === 'string') {
        if (patch[0] === '.') {
          patch = path.resolve(process.cwd(), patch)
        }
        patch = require(patch)
      }

      const type = typeof patch
      if (type !== 'function') {
        this._agent.logger.error('Invalid patch handler type "%s" for module "%s"', type, name)
        continue
      }

      exports = patch(exports, this._agent, { version, enabled })
    }
  }
  return exports
}

Instrumentation.prototype.addEndedTransaction = function (transaction) {
  var agent = this._agent

  if (!this._started) {
    agent.logger.debug('ignoring transaction %o', { trans: transaction.id, trace: transaction.traceId })
    return
  }

  const rc = this._runCtxMgr.active()
  if (rc.currTransaction() === transaction) {
    // Replace the active run context with an empty one. I.e. there is now
    // no active transaction or span (at least in this async task).
    this._runCtxMgr.supersedeRunContext(this._runCtxMgr.root())
    this._log.debug({ ctxmgr: this._runCtxMgr.toString() }, 'addEndedTransaction(%s)', transaction.name)
  }

  // Avoid transaction filtering time if only propagating trace-context.
  if (agent._conf.contextPropagationOnly) {
    // This one log.trace related to contextPropagationOnly is included as a
    // possible log hint to future debugging for why events are not being sent
    // to APM server.
    agent.logger.trace('contextPropagationOnly: skip sendTransaction')
    return
  }

  // https://github.com/elastic/apm/blob/main/specs/agents/tracing-sampling.md#non-sampled-transactions
  if (!transaction.sampled && !agent._transport.supportsKeepingUnsampledTransaction()) {
    return
  }

  // if I have ended and I have something buffered, send that buffered thing
  if (transaction.getBufferedSpan()) {
    this._encodeAndSendSpan(transaction.getBufferedSpan())
  }

  var payload = agent._transactionFilters.process(transaction._encode())
  if (!payload) {
    agent.logger.debug('transaction ignored by filter %o', { trans: transaction.id, trace: transaction.traceId })
    return
  }

  agent.logger.debug('sending transaction %o', { trans: transaction.id, trace: transaction.traceId })
  agent._transport.sendTransaction(payload)
}

Instrumentation.prototype.addEndedSpan = function (span) {
  var agent = this._agent

  if (!this._started) {
    agent.logger.debug('ignoring span %o', { span: span.id, parent: span.parentId, trace: span.traceId, name: span.name, type: span.type })
    return
  }

  // Replace the active run context with this span removed. Typically this
  // span is the top of stack (i.e. is the current span). However, it is
  // possible to have out-of-order span.end(), in which case the ended span
  // might not.
  const newRc = this._runCtxMgr.active().leaveSpan(span)
  if (newRc) {
    this._runCtxMgr.supersedeRunContext(newRc)
  }
  this._log.debug({ ctxmgr: this._runCtxMgr.toString() }, 'addEndedSpan(%s)', span.name)

  // Avoid span encoding time if only propagating trace-context.
  if (agent._conf.contextPropagationOnly) {
    return
  }

  if (!span.isRecorded()) {
    span.transaction.captureDroppedSpan(span)
    return
  }

  if (!this._agent._conf.spanCompressionEnabled) {
    this._encodeAndSendSpan(span)
  } else {
    // if I have ended and I have something buffered, send that buffered thing
    if (span.getBufferedSpan()) {
      this._encodeAndSendSpan(span.getBufferedSpan())
      span.setBufferedSpan(null)
    }

    const parentSpan = span.getParentSpan()
    if ((parentSpan && parentSpan.ended) || !span.isCompressionEligible()) {
      const buffered = parentSpan && parentSpan.getBufferedSpan()
      if (buffered) {
        this._encodeAndSendSpan(buffered)
        parentSpan.setBufferedSpan(null)
      }
      this._encodeAndSendSpan(span)
    } else if (!parentSpan.getBufferedSpan()) {
      // span is compressible and there's nothing buffered
      // add to buffer, move on
      parentSpan.setBufferedSpan(span)
    } else if (!parentSpan.getBufferedSpan().tryToCompress(span)) {
      // we could not compress span so SEND bufferend span
      // and buffer the span we could not compress
      this._encodeAndSendSpan(parentSpan.getBufferedSpan())
      parentSpan.setBufferedSpan(span)
    }
  }
}

Instrumentation.prototype._encodeAndSendSpan = function (span) {
  const duration = span.isComposite() ? span.getCompositeSum() : span.duration()
  if (span.discardable && duration / 1000 < this._agent._conf.exitSpanMinDuration) {
    span.transaction.captureDroppedSpan(span)
    return
  }

  const agent = this._agent
  // Note this error as an "inflight" event. See Agent#flush().
  const inflightEvents = agent._inflightEvents
  inflightEvents.add(span.id)

  agent.logger.debug('encoding span %o', { span: span.id, parent: span.parentId, trace: span.traceId, name: span.name, type: span.type })
  span._encode(function (err, payload) {
    if (err) {
      agent.logger.error('error encoding span %o', { span: span.id, parent: span.parentId, trace: span.traceId, name: span.name, type: span.type, error: err.message })
    } else {
      payload = agent._spanFilters.process(payload)
      if (!payload) {
        agent.logger.debug('span ignored by filter %o', { span: span.id, parent: span.parentId, trace: span.traceId, name: span.name, type: span.type })
      } else {
        agent.logger.debug('sending span %o', { span: span.id, parent: span.parentId, trace: span.traceId, name: span.name, type: span.type })
        if (agent._transport) {
          agent._transport.sendSpan(payload)
        }
      }
    }
    inflightEvents.delete(span.id)
  })
}

// Replace the current run context with one where the given transaction is
// current.
Instrumentation.prototype.supersedeWithTransRunContext = function (trans) {
  if (this._started) {
    const rc = this._runCtxMgr.root().enterTrans(trans)
    this._runCtxMgr.supersedeRunContext(rc)
    this._log.debug({ ctxmgr: this._runCtxMgr.toString() }, 'supersedeWithTransRunContext(<Trans %s>)', trans.id)
  }
}

// Replace the current run context with one where the given span is current.
Instrumentation.prototype.supersedeWithSpanRunContext = function (span) {
  if (this._started) {
    const rc = this._runCtxMgr.active().enterSpan(span)
    this._runCtxMgr.supersedeRunContext(rc)
    this._log.debug({ ctxmgr: this._runCtxMgr.toString() }, 'supersedeWithSpanRunContext(<Span %s>)', span.id)
  }
}

// Set the current run context to have *no* transaction. No spans will be
// created in this run context until a subsequent `startTransaction()`.
Instrumentation.prototype.supersedeWithEmptyRunContext = function () {
  if (this._started) {
    this._runCtxMgr.supersedeRunContext(this._runCtxMgr.root())
    this._log.debug({ ctxmgr: this._runCtxMgr.toString() }, 'supersedeWithEmptyRunContext()')
  }
}

// Create a new transaction, but do *not* replace the current run context to
// make this the "current" transaction. Compare to `startTransaction`.
Instrumentation.prototype.createTransaction = function (name, ...args) {
  return new Transaction(this._agent, name, ...args)
}

Instrumentation.prototype.startTransaction = function (name, ...args) {
  const trans = new Transaction(this._agent, name, ...args)
  this.supersedeWithTransRunContext(trans)
  return trans
}

Instrumentation.prototype.endTransaction = function (result, endTime) {
  const trans = this.currTransaction()
  if (!trans) {
    this._agent.logger.debug('cannot end transaction - no active transaction found')
    return
  }
  trans.end(result, endTime)
}

Instrumentation.prototype.setDefaultTransactionName = function (name) {
  const trans = this.currTransaction()
  if (!trans) {
    this._agent.logger.debug('no active transaction found - cannot set default transaction name')
    return
  }
  trans.setDefaultName(name)
}

Instrumentation.prototype.setTransactionName = function (name) {
  const trans = this.currTransaction()
  if (!trans) {
    this._agent.logger.debug('no active transaction found - cannot set transaction name')
    return
  }
  trans.name = name
}

Instrumentation.prototype.setTransactionOutcome = function (outcome) {
  const trans = this.currTransaction()
  if (!trans) {
    this._agent.logger.debug('no active transaction found - cannot set transaction outcome')
    return
  }
  trans.setOutcome(outcome)
}

// Create a new span in the current transaction, if any, and make it the
// current span. The started span is returned. This will return null if a span
// could not be created -- which could happen for a number of reasons.
Instrumentation.prototype.startSpan = function (name, type, subtype, action, opts) {
  const trans = this.currTransaction()
  if (!trans) {
    this._agent.logger.debug('no active transaction found - cannot build new span')
    return null
  }
  return trans.startSpan.apply(trans, arguments)
}

// Create a new span in the current transaction, if any. The created span is
// returned, or null if the span could not be created.
//
// This does *not* replace the current run context to make this span the
// "current" one. This allows instrumentations to avoid impacting the run
// context of the calling code. Compare to `startSpan`.
Instrumentation.prototype.createSpan = function (name, type, subtype, action, opts) {
  const trans = this.currTransaction()
  if (!trans) {
    this._agent.logger.debug('no active transaction found - cannot build new span')
    return null
  }
  return trans.createSpan.apply(trans, arguments)
}

Instrumentation.prototype.setSpanOutcome = function (outcome) {
  const span = this.currSpan()
  if (!span) {
    this._agent.logger.debug('no active span found - cannot set span outcome')
    return null
  }
  span.setOutcome(outcome)
}

Instrumentation.prototype.currRunContext = function () {
  if (!this._started) {
    return null
  }
  return this._runCtxMgr.active()
}

// Bind the given function to the current run context.
Instrumentation.prototype.bindFunction = function (fn) {
  if (!this._started) {
    return fn
  }
  return this._runCtxMgr.bindFn(this._runCtxMgr.active(), fn)
}

// Bind the given function to a given run context.
Instrumentation.prototype.bindFunctionToRunContext = function (runContext, fn) {
  if (!this._started) {
    return fn
  }
  return this._runCtxMgr.bindFn(runContext, fn)
}

// Bind the given function to an *empty* run context.
// This can be used to ensure `fn` does *not* run in the context of the current
// transaction or span.
Instrumentation.prototype.bindFunctionToEmptyRunContext = function (fn) {
  if (!this._started) {
    return fn
  }
  return this._runCtxMgr.bindFn(this._runCtxMgr.root(), fn)
}

// Bind the given EventEmitter to the current run context.
//
// This wraps the emitter so that any added event handler function is bound
// as if `bindFunction` had been called on it. Note that `ee` need not
// inherit from EventEmitter -- it uses duck typing.
Instrumentation.prototype.bindEmitter = function (ee) {
  if (!this._started) {
    return ee
  }
  return this._runCtxMgr.bindEE(this._runCtxMgr.active(), ee)
}

// Bind the given EventEmitter to a given run context.
Instrumentation.prototype.bindEmitterToRunContext = function (runContext, ee) {
  if (!this._started) {
    return ee
  }
  return this._runCtxMgr.bindEE(runContext, ee)
}

// Return true iff the given EventEmitter is bound to a run context.
Instrumentation.prototype.isEventEmitterBound = function (ee) {
  if (!this._started) {
    return false
  }
  return this._runCtxMgr.isEEBound(ee)
}

// Invoke the given function in the context of `runContext`.
Instrumentation.prototype.withRunContext = function (runContext, fn, thisArg, ...args) {
  if (!this._started) {
    return fn.call(thisArg, ...args)
  }
  return this._runCtxMgr.with(runContext, fn, thisArg, ...args)
}
