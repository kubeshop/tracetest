'use strict'

const assert = require('assert')
const crypto = require('crypto')
const fs = require('fs')
const http = require('http')
const https = require('https')
const util = require('util')
const os = require('os')
const { performance } = require('perf_hooks')
const { URL } = require('url')
const zlib = require('zlib')

const HttpAgentKeepAlive = require('agentkeepalive')
const HttpsAgentKeepAlive = HttpAgentKeepAlive.HttpsAgent
const Filters = require('object-filter-sequence')
const querystring = require('querystring')
const Writable = require('readable-stream').Writable
const getContainerInfo = require('./lib/container-info')
const eos = require('end-of-stream')
const semver = require('semver')
const streamToBuffer = require('fast-stream-to-buffer')
const StreamChopper = require('stream-chopper')

const ndjson = require('./lib/ndjson')
const { NoopLogger } = require('./lib/logging')
const truncate = require('./lib/truncate')
const { getCentralConfigIntervalS } = require('./lib/central-config')

module.exports = Client

// These symbols are used as markers in the client stream to indicate special
// flush handling.
const kFlush = Symbol('flush')
const kLambdaEndFlush = Symbol('lambdaEndFlush')
function isFlushMarker (obj) {
  return obj === kFlush || obj === kLambdaEndFlush
}

const hostname = os.hostname()
const requiredOpts = [
  'agentName',
  'agentVersion',
  'serviceName',
  'userAgent'
]

// Get handles on uninstrumented functions for making HTTP(S) requests before
// the APM agent has a chance to wrap them. This allows the Client to make
// requests to APM server without interfering with the APM agent's tracing
// of the user application.
const httpGet = http.get
const httpRequest = http.request
const httpsGet = https.get
const httpsRequest = https.request

const containerInfo = getContainerInfo.sync()

const isLambdaExecutionEnvironment = !!process.env.AWS_LAMBDA_FUNCTION_NAME

// All sockets on the agent are unreffed when they are created. This means that
// when the user process's event loop is done, and these are the only handles
// left, the process 'beforeExit' event will be emitted. By listening for this
// we can make sure to end the requests properly before process exit. This way
// we don't keep the process running until the `time` timeout happens.
//
// An exception to this is AWS Lambda which, in some cases (sync function
// handlers that use a callback), will wait for 'beforeExit' to freeze the
// Lambda instance VM *for later re-use*. This means we never want to shutdown
// the `Client` on 'beforeExit'.
const clientsToAutoEnd = []
if (!isLambdaExecutionEnvironment) {
  process.once('beforeExit', function () {
    clientsToAutoEnd.forEach(function (client) {
      if (!client) {
        // Clients remove themselves from the array when they end.
        return
      }
      client._gracefulExit()
    })
  })
}

util.inherits(Client, Writable)

Client.encoding = Object.freeze({
  METADATA: Symbol('metadata'),
  TRANSACTION: Symbol('transaction'),
  SPAN: Symbol('span'),
  ERROR: Symbol('error'),
  METRICSET: Symbol('metricset')
})

function Client (opts) {
  if (!(this instanceof Client)) return new Client(opts)

  Writable.call(this, { objectMode: true })

  this._corkTimer = null
  this._agent = null
  this._activeIntakeReq = false
  this._onIntakeReqConcluded = null
  this._transport = null
  this._configTimer = null
  this._backoffReconnectCount = 0
  this._intakeRequestGracefulExitFn = null // set in makeIntakeRequest
  this._encodedMetadata = null
  this._cloudMetadata = null
  this._extraMetadata = null
  this._metadataFilters = new Filters()
  // _lambdaActive indicates if a Lambda function invocation is active. It is
  // only meaningful if `isLambdaExecutionEnvironment`.
  this._lambdaActive = false

  // Internal runtime stats for developer debugging/tuning.
  this._numEvents = 0 // number of events given to the client
  this._numEventsDropped = 0 // number of events dropped because overloaded
  this._numEventsEnqueued = 0 // number of events written through to chopper
  this.sent = 0 // number of events sent to APM server (not necessarily accepted)
  this._slowWriteBatch = { // data on slow or the slowest _writeBatch
    numOver10Ms: 0,
    // Data for the slowest _writeBatch:
    encodeTimeMs: 0,
    fullTimeMs: 0,
    numEvents: 0,
    numBytes: 0
  }

  this.config(opts)
  this._log = this._conf.logger || new NoopLogger()

  if (this._conf.cloudMetadataFetcher && this._conf.expectExtraMetadata) {
    throw new Error('it is an error to create a Client with both cloudMetadataFetcher and expectExtraMetadata')
  } else if (this._conf.cloudMetadataFetcher) {
    // Start stream in corked mode, uncork when cloud metadata is fetched and
    // assigned.  Also, the _maybeUncork will not uncork until _encodedMetadata
    // is set.
    this._log.trace('corking (cloudMetadataFetcher)')
    this.cork()
    this._fetchAndEncodeMetadata(() => {
      // _fetchAndEncodeMetadata will have set/memoized the encoded
      // metadata to the _encodedMetadata property.

      // This reverses the cork() call in the constructor above. "Maybe" uncork,
      // in case the client has been destroyed before this callback is called.
      this._maybeUncork()
      this._log.trace('uncorked (cloudMetadataFetcher)')

      // the `cloud-metadata` event allows listeners to know when the
      // agent has finished fetching and encoding its metadata for the
      // first time
      this.emit('cloud-metadata', this._encodedMetadata)
    })
  } else if (this._conf.expectExtraMetadata) {
    // Uncorking will happen in the expected `.setExtraMetadata()` call.
    this._log.trace('corking (expectExtraMetadata)')
    this.cork()
  } else {
    this._resetEncodedMetadata()
  }

  // `_apmServerVersion` is one of:
  // - `undefined`: the version has not yet been fetched
  // - `null`: the APM server version is unknown, could not be determined
  // - a semver.SemVer instance
  this._apmServerVersion = this._conf.apmServerVersion ? semver.SemVer(this._conf.apmServerVersion) : undefined
  if (!this._apmServerVersion) {
    this._fetchApmServerVersion()
  }

  this._chopper = new StreamChopper({
    size: this._conf.size,
    time: this._conf.time,
    type: StreamChopper.overflow,
    transform () {
      return zlib.createGzip({
        level: zlib.constants.Z_BEST_SPEED
      })
    }
  })
  const onIntakeError = (err) => {
    if (this.destroyed === false) {
      this.emit('request-error', err)
    }
  }
  this._chopper.on('stream', getChoppedStreamHandler(this, onIntakeError))

  // We don't expect the chopper stream to end until the client is ending.
  // Make sure to clean up if this does happen unexpectedly.
  const fail = () => {
    if (this._writableState.ending === false) this.destroy()
  }
  eos(this._chopper, fail)

  this._index = clientsToAutoEnd.length
  clientsToAutoEnd.push(this)

  // The 'beforeExit' event is significant in Lambda invocation completion
  // handling, so we log it for debugging.
  if (isLambdaExecutionEnvironment && this._log.isLevelEnabled('trace')) {
    process.prependListener('beforeExit', () => {
      this._log.trace('process "beforeExit"')
    })
  }

  if (this._conf.centralConfig) {
    this._pollConfig()
  }
}

// Return current internal stats.
Client.prototype._getStats = function () {
  return {
    numEvents: this._numEvents,
    numEventsDropped: this._numEventsDropped,
    numEventsEnqueued: this._numEventsEnqueued,
    numEventsSent: this.sent,
    slowWriteBatch: this._slowWriteBatch,
    backoffReconnectCount: this._backoffReconnectCount
  }
}

Client.prototype.config = function (opts) {
  this._conf = Object.assign(this._conf || {}, opts)

  this._conf.globalLabels = normalizeGlobalLabels(this._conf.globalLabels)

  const missing = requiredOpts.filter(name => !this._conf[name])
  if (missing.length > 0) throw new Error('Missing required option(s): ' + missing.join(', '))

  // default values
  if (!this._conf.size && this._conf.size !== 0) this._conf.size = 750 * 1024
  if (!this._conf.time && this._conf.time !== 0) this._conf.time = 10000
  if (!this._conf.serverTimeout && this._conf.serverTimeout !== 0) this._conf.serverTimeout = 15000
  if (!this._conf.serverUrl) this._conf.serverUrl = 'http://127.0.0.1:8200'
  if (!this._conf.hostname) this._conf.hostname = hostname
  if (!this._conf.environment) this._conf.environment = process.env.NODE_ENV || 'development'
  if (!this._conf.truncateKeywordsAt) this._conf.truncateKeywordsAt = 1024
  if (!this._conf.truncateStringsAt) this._conf.truncateStringsAt = 1024
  if (!this._conf.truncateCustomKeysAt) this._conf.truncateCustomKeysAt = 1024
  if (!this._conf.truncateLongFieldsAt) this._conf.truncateLongFieldsAt = 10000
  // The deprecated `truncateErrorMessagesAt` will be honored if specified.
  if (!this._conf.bufferWindowTime) this._conf.bufferWindowTime = 20
  if (!this._conf.bufferWindowSize) this._conf.bufferWindowSize = 50
  if (!this._conf.maxQueueSize) this._conf.maxQueueSize = 1024
  if (!this._conf.intakeResTimeout) this._conf.intakeResTimeout = 10000
  if (!this._conf.intakeResTimeoutOnEnd) this._conf.intakeResTimeoutOnEnd = 1000
  this._conf.keepAlive = this._conf.keepAlive !== false
  this._conf.centralConfig = this._conf.centralConfig || false
  if (!('keepAliveMsecs' in this._conf)) this._conf.keepAliveMsecs = 1000
  if (!('maxSockets' in this._conf)) this._conf.maxSockets = Infinity
  if (!('maxFreeSockets' in this._conf)) this._conf.maxFreeSockets = 256
  if (!('freeSocketTimeout' in this._conf)) this._conf.freeSocketTimeout = 4000

  // processed values
  this._conf.serverUrl = new URL(this._conf.serverUrl)

  if (containerInfo) {
    if (!this._conf.containerId && containerInfo.containerId) {
      this._conf.containerId = containerInfo.containerId
    }
    if (!this._conf.kubernetesPodUID && containerInfo.podId) {
      this._conf.kubernetesPodUID = containerInfo.podId
    }
    if (!this._conf.kubernetesPodName && containerInfo.podId) {
      this._conf.kubernetesPodName = hostname
    }
  }

  let AgentKeepAlive
  switch (this._conf.serverUrl.protocol) {
    case 'http:':
      this._transport = http
      this._transportRequest = httpRequest
      this._transportGet = httpGet
      AgentKeepAlive = HttpAgentKeepAlive
      break
    case 'https:':
      this._transport = https
      this._transportRequest = httpsRequest
      this._transportGet = httpsGet
      AgentKeepAlive = HttpsAgentKeepAlive
      break
    default:
      throw new Error('Unknown protocol ' + this._conf.serverUrl.protocol)
  }

  // Only reset `this._agent` if the serverUrl has changed to avoid
  // unnecessarily abandoning keep-alive connections.
  if (!this._agent || (opts && 'serverUrl' in opts)) {
    if (this._agent) {
      this._agent.destroy()
    }
    this._agent = new AgentKeepAlive({
      keepAlive: this._conf.keepAlive,
      keepAliveMsecs: this._conf.keepAliveMsecs,
      freeSocketTimeout: this._conf.freeSocketTimeout,
      timeout: this._conf.serverTimeout,
      maxSockets: this._conf.maxSockets,
      maxFreeSockets: this._conf.maxFreeSockets
    })
  }

  // http request options
  this._conf.requestIntake = getIntakeRequestOptions(this._conf, this._agent)
  this._conf.requestConfig = getConfigRequestOptions(this._conf, this._agent)
  this._conf.requestSignalLambdaEnd = getSignalLambdaEndRequestOptions(this._conf, this._agent)

  this._conf.metadata = getMetadata(this._conf)

  // fixes bug where cached/memoized _encodedMetadata wouldn't be
  // updated when client was reconfigured
  if (this._encodedMetadata) {
    this._resetEncodedMetadata()
  }
}

/**
 * Set extra additional metadata to be sent to APM Server in intake requests.
 *
 * If the Client was configured with `expectExtraMetadata: true` then will
 * uncork the client to allow intake requests to begin.
 *
 * If this is called multiple times, it is additive.
 */
Client.prototype.setExtraMetadata = function (extraMetadata) {
  if (!this._extraMetadata) {
    this._extraMetadata = extraMetadata
  } else {
    metadataMergeDeep(this._extraMetadata, extraMetadata)
  }
  this._resetEncodedMetadata()

  if (this._conf.expectExtraMetadata) {
    this._log.trace('maybe uncork (expectExtraMetadata)')
    this._maybeUncork()
  }
}

/**
 * Add a filter function used to filter the "metadata" object sent to APM
 * server. See the APM Agent `addMetadataFilter` documentation for details.
 * https://www.elastic.co/guide/en/apm/agent/nodejs/current/agent-api.html#apm-add-metadata-filter
 */
Client.prototype.addMetadataFilter = function (fn) {
  assert.strictEqual(typeof fn, 'function', 'fn arg must be a function')
  this._metadataFilters.push(fn)
  if (this._encodedMetadata) {
    this._resetEncodedMetadata()
  }
}

/**
 * (Re)set `_encodedMetadata` from this._conf.metadata, this._cloudMetadata,
 * this._extraMetadata and possible this._metadataFilters.
 */
Client.prototype._resetEncodedMetadata = function () {
  // Make a deep clone so that the originals are not modified when (a) adding
  // `.cloud` and (b) filtering. This isn't perf-sensitive code, so this JSON
  // cycle for cloning should suffice.
  let metadata = deepClone(this._conf.metadata)
  if (this._cloudMetadata) {
    metadata.cloud = deepClone(this._cloudMetadata)
  }
  if (this._extraMetadata) {
    metadataMergeDeep(metadata, deepClone(this._extraMetadata))
  }

  // Possible filters from APM agent's `apm.addMetadataFilter()`.
  if (this._metadataFilters && this._metadataFilters.length > 0) {
    metadata = this._metadataFilters.process(metadata)
  }

  // This is the only code path that should set `_encodedMetadata`.
  this._encodedMetadata = this._encode({ metadata }, Client.encoding.METADATA)
  this._log.trace({ _encodedMetadata: this._encodedMetadata }, '_resetEncodedMetadata')
}

Client.prototype._pollConfig = function () {
  const opts = this._conf.requestConfig
  if (this._conf.lastConfigEtag) {
    opts.headers['If-None-Match'] = this._conf.lastConfigEtag
  }

  const req = this._transportGet(opts, res => {
    res.on('error', err => {
      // Not sure this event can ever be emitted, but just in case
      res.destroy(err)
    })

    this._scheduleNextConfigPoll(getMaxAge(res))

    if (
      res.statusCode === 304 || // No new config since last time
      res.statusCode === 403 || // Central config not enabled in APM Server
      res.statusCode === 404 // Old APM Server that doesn't support central config
    ) {
      res.resume()
      return
    }

    streamToBuffer(res, (err, buf) => {
      if (err) {
        this.emit('request-error', processConfigErrorResponse(res, buf, err))
        return
      }

      if (res.statusCode === 200) {
        // 200: New config available (or no config for the given service.name / service.environment)
        const etag = res.headers.etag
        if (etag) this._conf.lastConfigEtag = etag

        let config
        try {
          config = JSON.parse(buf)
        } catch (parseErr) {
          this.emit('request-error', processConfigErrorResponse(res, buf, parseErr))
          return
        }
        this.emit('config', config)
      } else {
        this.emit('request-error', processConfigErrorResponse(res, buf))
      }
    })
  })

  req.on('error', err => {
    this._scheduleNextConfigPoll()
    this.emit('request-error', err)
  })
}

Client.prototype._scheduleNextConfigPoll = function (seconds) {
  if (this._configTimer !== null) return

  const delayS = getCentralConfigIntervalS(seconds)
  this._configTimer = setTimeout(() => {
    this._configTimer = null
    this._pollConfig()
  }, delayS * 1000)

  this._configTimer.unref()
}

// re-ref the open socket handles
Client.prototype._ref = function () {
  Object.keys(this._agent.sockets).forEach(remote => {
    this._agent.sockets[remote].forEach(function (socket) {
      socket.ref()
    })
  })
}

Client.prototype._write = function (obj, enc, cb) {
  if (isFlushMarker(obj)) {
    this._writeFlush(obj, cb)
  } else {
    const t = process.hrtime()
    const chunk = this._encode(obj, enc)
    this._numEventsEnqueued++
    this._chopper.write(chunk, cb)
    this._log.trace({
      fullTimeMs: deltaMs(t),
      numEvents: 1,
      numBytes: chunk.length
    }, '_write: encode object')
  }
}

Client.prototype._writev = function (objs, cb) {
  // Limit the size of individual writes to manageable batches, primarily to
  // limit large sync pauses due to `_encode`ing in `_writeBatch`. This value
  // is not particularly well tuned. It was selected to get sync pauses under
  // 10ms on a developer machine.
  const MAX_WRITE_BATCH_SIZE = 32

  let offset = 0

  const processBatch = () => {
    if (this.destroyed) {
      cb()
      return
    }

    let flushIdx = -1
    const limit = Math.min(objs.length, offset + MAX_WRITE_BATCH_SIZE)
    for (let i = offset; i < limit; i++) {
      if (isFlushMarker(objs[i].chunk)) {
        flushIdx = i
        break
      }
    }

    if (offset === 0 && flushIdx === -1 && objs.length <= MAX_WRITE_BATCH_SIZE) {
      // A shortcut if there is no flush marker and the whole `objs` fits in a batch.
      this._writeBatch(objs, cb)
    } else if (flushIdx === -1) {
      // No flush marker in this batch.
      this._writeBatch(objs.slice(offset, limit),
        limit === objs.length ? cb : processBatch)
      offset = limit
    } else if (flushIdx > offset) {
      // There are some events in the queue before a flush marker.
      this._writeBatch(objs.slice(offset, flushIdx), processBatch)
      offset = flushIdx
    } else if (flushIdx === objs.length - 1) {
      // The next item is a flush marker, and it is the *last* item in the queue.
      this._writeFlush(objs[flushIdx].chunk, cb)
    } else {
      // The next item in the queue is a flush.
      this._writeFlush(objs[flushIdx].chunk, processBatch)
      offset++
    }
  }

  processBatch()
}

// Write a batch of events (excluding specially handled "flush" events) to
// the stream chopper.
Client.prototype._writeBatch = function (objs, cb) {
  const t = process.hrtime()
  const chunks = []
  for (var i = 0; i < objs.length; i++) {
    const obj = objs[i]
    chunks.push(this._encode(obj.chunk, obj.encoding))
  }
  const chunk = chunks.join('')
  const encodeTimeMs = deltaMs(t)

  this._numEventsEnqueued += objs.length
  this._chopper.write(chunk, cb)
  const fullTimeMs = deltaMs(t)

  if (fullTimeMs > this._slowWriteBatch.fullTimeMs) {
    this._slowWriteBatch.encodeTimeMs = encodeTimeMs
    this._slowWriteBatch.fullTimeMs = fullTimeMs
    this._slowWriteBatch.numEvents = objs.length
    this._slowWriteBatch.numBytes = chunk.length
  }
  if (fullTimeMs > 10) {
    this._slowWriteBatch.numOver10Ms++
  }
  this._log.trace({
    encodeTimeMs: encodeTimeMs,
    fullTimeMs: fullTimeMs,
    numEvents: objs.length,
    numBytes: chunk.length
  }, '_writeBatch')
}

Client.prototype._writeFlush = function (flushMarker, cb) {
  this._log.trace({ activeIntakeReq: this._activeIntakeReq, lambdaEnd: flushMarker === kLambdaEndFlush }, '_writeFlush')

  let onFlushed = cb
  if (isLambdaExecutionEnvironment && flushMarker === kLambdaEndFlush) {
    onFlushed = () => {
      // Signal the Elastic AWS Lambda extension that it is done passing data
      // for this invocation, then call `cb()` so the wrapped Lambda handler
      // can finish.
      this._signalLambdaEnd(cb)
    }
  }

  if (this._activeIntakeReq) {
    this._onIntakeReqConcluded = onFlushed
    this._chopper.chop()
  } else {
    this._chopper.chop(onFlushed)
  }
}

Client.prototype._maybeCork = function () {
  if (!this._writableState.corked) {
    if (isLambdaExecutionEnvironment && !this._lambdaActive) {
      this.cork()
    } else if (this._conf.bufferWindowTime !== -1) {
      this.cork()
      if (this._corkTimer && this._corkTimer.refresh) {
        // the refresh function was added in Node 10.2.0
        this._corkTimer.refresh()
      } else {
        this._corkTimer = setTimeout(() => {
          this.uncork()
        }, this._conf.bufferWindowTime)
      }
    }
  } else if (this._writableState.length >= this._conf.bufferWindowSize) {
    this._maybeUncork()
  }
}

Client.prototype._maybeUncork = function () {
  if (!this._encodedMetadata) {
    // The client must remain corked until cloud metadata has been
    // fetched-or-skipped.
    return
  } else if (isLambdaExecutionEnvironment && !this._lambdaActive) {
    // In a Lambda env, we must only uncork when an invocation is active,
    // otherwise we could start an intake request just before the VM is frozen.
    return
  }

  if (this._writableState.corked) {
    // Wait till next tick, so that the current write that triggered the call
    // to `_maybeUncork` have time to be added to the queue. If we didn't do
    // this, that last write would trigger a single call to `_write`.
    process.nextTick(() => {
      if (this.destroyed === false && !(isLambdaExecutionEnvironment && !this._lambdaActive)) {
        this.uncork()
      }
    })

    if (this._corkTimer) {
      clearTimeout(this._corkTimer)
      this._corkTimer = null
    }
  }
}

Client.prototype._encode = function (obj, enc) {
  const out = {}
  switch (enc) {
    case Client.encoding.SPAN:
      out.span = truncate.span(obj.span, this._conf)
      break
    case Client.encoding.TRANSACTION:
      out.transaction = truncate.transaction(obj.transaction, this._conf)
      break
    case Client.encoding.METADATA:
      out.metadata = truncate.metadata(obj.metadata, this._conf)
      break
    case Client.encoding.ERROR:
      out.error = truncate.error(obj.error, this._conf)
      break
    case Client.encoding.METRICSET:
      out.metricset = truncate.metricset(obj.metricset, this._conf)
      break
  }
  return ndjson.serialize(out)
}

Client.prototype.lambdaStart = function () {
  this._lambdaActive = true
}

// With the cork/uncork handling on this stream, `this.write`ing on this
// stream when already destroyed will lead to:
//    Error: Cannot call write after a stream was destroyed
// when the `_corkTimer` expires.
Client.prototype._isUnsafeToWrite = function () {
  return this.destroyed
}

Client.prototype._shouldDropEvent = function () {
  this._numEvents++
  const shouldDrop = this._writableState.length >= this._conf.maxQueueSize
  if (shouldDrop) {
    this._numEventsDropped++
  }
  return shouldDrop
}

Client.prototype.sendSpan = function (span, cb) {
  if (this._isUnsafeToWrite() || this._shouldDropEvent()) {
    return
  }
  this._maybeCork()
  return this.write({ span }, Client.encoding.SPAN, cb)
}

Client.prototype.sendTransaction = function (transaction, cb) {
  if (this._isUnsafeToWrite() || this._shouldDropEvent()) {
    return
  }
  this._maybeCork()
  return this.write({ transaction }, Client.encoding.TRANSACTION, cb)
}

Client.prototype.sendError = function (error, cb) {
  if (this._isUnsafeToWrite() || this._shouldDropEvent()) {
    return
  }
  this._maybeCork()
  return this.write({ error }, Client.encoding.ERROR, cb)
}

Client.prototype.sendMetricSet = function (metricset, cb) {
  if (this._isUnsafeToWrite() || this._shouldDropEvent()) {
    return
  }
  this._maybeCork()
  return this.write({ metricset }, Client.encoding.METRICSET, cb)
}

/**
 * If possible, start a flush of currently queued APM events to APM server.
 *
 * "If possible," because there are some guards on uncorking. See `_maybeUncork`.
 *
 * @param {Object} opts - Optional.
 *    - {Boolean} opts.lambdaEnd - Optional. Default false. Setting this true
 *      tells the client to also handle the end of a Lambda function invocation.
 * @param {Function} cb - Optional. `cb()` will be called when the data has
 *    be sent to APM Server (or failed in the attempt).
 */
Client.prototype.flush = function (opts, cb) {
  if (typeof opts === 'function') {
    cb = opts
    opts = {}
  } else if (!opts) {
    opts = {}
  }
  const lambdaEnd = !!opts.lambdaEnd

  // Write the special "flush" signal. We do this so that the order of writes
  // and flushes are kept. If we where to just flush the client right here, the
  // internal Writable buffer might still contain data that hasn't yet been
  // given to the _write function.

  if (lambdaEnd && isLambdaExecutionEnvironment && this._lambdaActive) {
    // To flush the current data and ensure that subsequently sent events *in
    // the same tick* do not start a new intake request, we must uncork
    // synchronously -- rather than the nextTick uncork done in `_maybeUncork()`.
    assert(this._encodedMetadata, 'client.flush({lambdaEnd:true}) must not be called before metadata has been set')
    const rv = this.write(kLambdaEndFlush, cb)
    this.uncork()
    this._lambdaActive = false
    return rv
  } else {
    this._maybeUncork()
    return this.write(kFlush, cb)
  }
}

// A handler that can be called on process "beforeExit" to attempt quick and
// orderly shutdown of the client. It attempts to ensure that the current
// active intake API request to APM server is completed quickly.
Client.prototype._gracefulExit = function () {
  this._log.trace('_gracefulExit')

  if (this._intakeRequestGracefulExitFn) {
    this._intakeRequestGracefulExitFn()
  }

  // Calling _ref here, instead of relying on the _ref call in `_final`,
  // is necessary because `client.end()` does *not* result in the Client's
  // `_final()` being called when the process is exiting.
  this._ref()
  this.end()
}

Client.prototype._final = function (cb) {
  this._log.trace('_final')
  if (this._configTimer) {
    clearTimeout(this._configTimer)
    this._configTimer = null
  }
  clientsToAutoEnd[this._index] = null // remove global reference to ease garbage collection
  this._ref()
  this._chopper.end()
  cb()
}

Client.prototype._destroy = function (err, cb) {
  this._log.trace({ err }, '_destroy')
  if (this._configTimer) {
    clearTimeout(this._configTimer)
    this._configTimer = null
  }
  if (this._corkTimer) {
    clearTimeout(this._corkTimer)
    this._corkTimer = null
  }
  clientsToAutoEnd[this._index] = null // remove global reference to ease garbage collection
  this._chopper.destroy()
  this._agent.destroy()
  cb(err)
}

// Return the appropriate backoff delay (in milliseconds) before a next possible
// request to APM server.
// Spec: https://github.com/elastic/apm/blob/main/specs/agents/transport.md#transport-errors
//
// In a Lambda environment, a backoff delay can be harmful: The backoff
// setTimeout is unref'd, to not hold the process open. A subsequent Lambda
// function invocation during that timer will result in no active handles and
// a process "beforeExit" event. That event is interpreted by the Lambda Runtime
// as "the Lambda function callback was never called", and it terminates the
// function and responds with `null`.  The solution is to never backoff in a
// Lambda environment -- we expect and assume the Lambda extension is working,
// and pass responsibility for backoff to the extension.
Client.prototype._getBackoffDelay = function (isErr) {
  let reconnectCount = this._backoffReconnectCount
  if (isErr && !isLambdaExecutionEnvironment) {
    this._backoffReconnectCount++
  } else {
    this._backoffReconnectCount = 0
    reconnectCount = 0
  }

  // min(reconnectCount++, 6) ** 2 ± 10%
  const delayS = Math.pow(Math.min(reconnectCount, 6), 2)
  const jitterS = delayS * (0.2 * Math.random() - 0.1)
  const delayMs = (delayS + jitterS) * 1000
  return delayMs
}

function getChoppedStreamHandler (client, onerror) {
  // Make a request to the apm-server intake API.
  // https://www.elastic.co/guide/en/apm/server/current/events-api.html
  //
  // In normal operation this works as follows:
  // - The StreamChopper (`this._chopper`) calls this function with a newly
  //   created Gzip stream, to which it writes encoded event data.
  // - It `gzipStream.end()`s the stream when:
  //   (a) approximately `apiRequestSize` of data have been written,
  //   (b) `apiRequestTime` seconds have passed, or
  //   (c) `_chopper.chop()` is explicitly called via `client.flush()`,
  //       e.g. used by the Node.js APM agent after `client.sendError()`.
  // - This function makes the HTTP POST to the apm-server, pipes the gzipStream
  //   to it, and waits for the completion of the request and the apm-server
  //   response.
  // - Then it calls the given `next` callback to signal StreamChopper that
  //   another chopped stream can be created, when there is more the send.
  //
  // Of course, things can go wrong. Here are the known ways this pipeline can
  // conclude.
  // - intake response success - A successful response from the APM server. This
  //   is the normal operation case described above.
  // - gzipStream error - An "error" event on the gzip stream.
  // - intake request error - An "error" event on the intake HTTP request, e.g.
  //   ECONNREFUSED or ECONNRESET.
  // - intakeResTimeout - A timer started *after* we are finished sending data
  //   to the APM server by which we require a response (including its body). By
  //   default this is 10s -- a very long time to allow for a slow or far
  //   apm-server. If we hit this, APM server is problematic anyway, so the
  //   delay doesn't add to the problems.
  // - serverTimeout - An idle timeout value (default 30s) set on the socket.
  //   This is a catch-all fallback for an otherwised wedged connection. If this
  //   is being hit, there is some major issue in the application (possibly a
  //   bug in the APM agent).
  // - process completion - The Client takes pains to always `.unref()` its
  //   handles to never keep a using process open if it is ready to exit. When
  //   the process is ready to exit, the following happens:
  //    - The "beforeExit" handler above will call `client._gracefulExit()` ...
  //    - ... which calls `client._ref()` to *hold the process open* to
  //      complete this request, and `client.end()` to end the `gzipStream` so
  //      this request can complete soon.
  //    - We then expect this request to complete quickly and the process will
  //      then finish exiting. A subtlety is if the APM server is not responding
  //      then we'll wait on the shorter `intakeResTimeoutOnEnd` (by default 1s).
  return function makeIntakeRequest (gzipStream, next) {
    const reqId = crypto.randomBytes(16).toString('hex')
    const log = client._log.child({ reqId })
    const startTime = process.hrtime()
    const timeline = []
    let bytesWritten = 0
    let intakeRes
    let intakeReqSocket = null
    let intakeResTimer = null
    let intakeRequestGracefulExitCalled = false
    const intakeResTimeout = client._conf.intakeResTimeout
    const intakeResTimeoutOnEnd = client._conf.intakeResTimeoutOnEnd

    // `_activeIntakeReq` is used to coordinate the callback to `client.flush(db)`.
    client._activeIntakeReq = true

    // Handle conclusion of this intake request. Each "part" is expected to call
    // `completePart()` at least once -- multiple calls are okay for cases like
    // the "error" and "close" events on a stream being called. When a part
    // errors or all parts are completed, then we can conclude.
    let concluded = false
    const completedFromPart = {
      gzipStream: false,
      intakeReq: false,
      intakeRes: false
    }
    let numToComplete = Object.keys(completedFromPart).length
    const completePart = (part, err) => {
      log.trace({ err, concluded }, 'completePart %s', part)
      timeline.push([deltaMs(startTime), `completePart ${part}`, err && err.message])
      assert(part in completedFromPart, `'${part}' is in completedFromPart`)

      if (concluded) {
        return
      }

      // If this is the final part to complete, then we are ready to conclude.
      let allPartsCompleted = false
      if (!completedFromPart[part]) {
        completedFromPart[part] = true
        numToComplete--
        if (numToComplete === 0) {
          allPartsCompleted = true
        }
      }
      if (!err && !allPartsCompleted) {
        return
      }

      // Conclude.
      concluded = true
      if (err) {
        // There was an error: clean up resources.

        // Note that in Node v8, destroying the gzip stream results in it
        // emitting an "error" event as follows. No harm, however.
        //    Error: gzip stream error: zlib binding closed
        //      at Gzip._transform (zlib.js:369:15)
        //      ...
        destroyStream(gzipStream)
        intakeReq.destroy()
        if (intakeResTimer) {
          log.trace('cancel intakeResTimer')
          clearTimeout(intakeResTimer)
          intakeResTimer = null
        }
      }
      client._intakeRequestGracefulExitFn = null

      client.sent = client._numEventsEnqueued
      client._activeIntakeReq = false
      const backoffDelayMs = client._getBackoffDelay(!!err)
      if (err) {
        log.trace({ timeline, bytesWritten, backoffDelayMs, err },
          'conclude intake request: error')
        onerror(err)
      } else {
        log.trace({ timeline, bytesWritten, backoffDelayMs },
          'conclude intake request: success')
      }
      if (client._onIntakeReqConcluded) {
        client._onIntakeReqConcluded()
        client._onIntakeReqConcluded = null
      }

      if (backoffDelayMs > 0) {
        setTimeout(next, backoffDelayMs).unref()
      } else {
        setImmediate(next)
      }
    }

    // Provide a function on the client for it to signal this intake request
    // to gracefully shutdown, i.e. finish up quickly.
    client._intakeRequestGracefulExitFn = () => {
      intakeRequestGracefulExitCalled = true
      if (intakeReqSocket) {
        log.trace('_intakeRequestGracefulExitFn: re-ref intakeReqSocket')
        intakeReqSocket.ref()
      }
      if (intakeResTimer) {
        log.trace('_intakeRequestGracefulExitFn: reset intakeResTimer to short timeout')
        clearTimeout(intakeResTimer)
        intakeResTimer = setTimeout(() => {
          completePart('intakeRes',
            new Error('intake response timeout: APM server did not respond ' +
              `within ${intakeResTimeoutOnEnd / 1000}s of graceful exit signal`))
        }, intakeResTimeoutOnEnd).unref()
      }
    }

    // Start the request and set its timeout.
    const intakeReq = client._transportRequest(client._conf.requestIntake)
    if (Number.isFinite(client._conf.serverTimeout)) {
      intakeReq.setTimeout(client._conf.serverTimeout)
    }
    // TODO: log intakeReq and intakeRes when
    // https://github.com/elastic/ecs-logging-nodejs/issues/67 is implemented.
    log.trace('intake request start')

    // Handle events on the intake request.
    // https://nodejs.org/api/http.html#http_http_request_options_callback docs
    // emitted events on the req and res objects for different scenarios.
    intakeReq.on('timeout', () => {
      log.trace('intakeReq "timeout"')
      // `.destroy(err)` will result in an "error" event.
      intakeReq.destroy(new Error(`APM Server response timeout (${client._conf.serverTimeout}ms)`))
    })

    intakeReq.on('socket', function (socket) {
      intakeReqSocket = socket
      // Unref the socket for this request so that the Client does not keep
      // the node process running if it otherwise would be done. (This is
      // tested by the "unref-client" test in test/side-effects.js.)
      //
      // The HTTP keep-alive agent will unref sockets when unused, and ref them
      // during a request. Given that the normal makeIntakeRequest behaviour
      // is to keep a request open for up to 10s (`apiRequestTime`), we must
      // manually unref the socket.
      //
      // The exception is when in a Lambda environment, where we *do* want to
      // keep the node process running to complete this intake request.
      // Otherwise a 'beforeExit' event can be sent, which the Lambda runtime
      // interprets as "the Lambda handler callback was never called".
      if (!isLambdaExecutionEnvironment && !intakeRequestGracefulExitCalled) {
        log.trace('intakeReq "socket": unref it')
        intakeReqSocket.unref()
      }
    })

    intakeReq.on('response', (intakeRes_) => {
      intakeRes = intakeRes_
      log.trace({ statusCode: intakeRes.statusCode, reqFinished: intakeReq.finished },
        'intakeReq "response"')
      let err
      const chunks = []

      if (!intakeReq.finished) {
        // Premature response from APM server. Typically this is for errors
        // like "queue is full", for which the response body will be parsed
        // below. However, set an `err` as a fallback for the unexpected case
        // that is with a 2xx response.
        if (intakeRes.statusCode >= 200 && intakeRes.statusCode < 300) {
          err = new Error(`premature apm-server response with statusCode=${intakeRes.statusCode}`)
        }
        // There is no point (though no harm) in sending more data to the APM
        // server. In case reading the error response body takes a while, pause
        // the gzip stream until it is destroyed in `completePart()`.
        gzipStream.pause()
      }

      // Handle events on the intake response.
      intakeRes.on('error', (intakeResErr) => {
        // I am not aware of a way to get an "error" event on the
        // IncomingMessage (see also https://stackoverflow.com/q/53691119), but
        // handling it here is preferable to an uncaughtException.
        intakeResErr = wrapError(intakeResErr, 'intake response error event')
        completePart('intakeRes', intakeResErr)
      })
      intakeRes.on('data', (chunk) => {
        chunks.push(chunk)
      })
      // intakeRes.on('close', () => { log.trace('intakeRes "close"') })
      // intakeRes.on('aborted', () => { log.trace('intakeRes "aborted"') })
      intakeRes.on('end', () => {
        log.trace('intakeRes "end"')
        if (intakeResTimer) {
          clearTimeout(intakeResTimer)
          intakeResTimer = null
        }
        if (intakeRes.statusCode < 200 || intakeRes.statusCode > 299) {
          err = processIntakeErrorResponse(intakeRes, Buffer.concat(chunks))
        }
        completePart('intakeRes', err)
      })
    })

    // intakeReq.on('abort', () => { log.trace('intakeReq "abort"') })
    // intakeReq.on('close', () => { log.trace('intakeReq "close"') })
    intakeReq.on('finish', () => {
      log.trace('intakeReq "finish"')
      completePart('intakeReq')
    })
    intakeReq.on('error', (err) => {
      log.trace('intakeReq "error"')
      completePart('intakeReq', err)
    })

    // Handle events on the gzip stream.
    gzipStream.on('data', (chunk) => {
      bytesWritten += chunk.length
    })
    gzipStream.on('error', (gzipErr) => {
      log.trace('gzipStream "error"')
      gzipErr = wrapError(gzipErr, 'gzip stream error')
      completePart('gzipStream', gzipErr)
    })
    gzipStream.on('finish', () => {
      // If the apm-server is not reading its input and the gzip data is large
      // enough to fill buffers, then the gzip stream will emit "finish", but
      // not "end". Therefore, this "finish" event is the best indicator that
      // the ball is now in the apm-server's court.
      //
      // We now start a timer waiting on the response, provided we still expect
      // one (we don't if the request has already errored out, e.g.
      // ECONNREFUSED) and it hasn't already completed (e.g. if it replied
      // quickly with "queue is full").
      log.trace('gzipStream "finish"')
      if (!completedFromPart.intakeReq && !completedFromPart.intakeRes) {
        const timeout = (client._writableState.ending || intakeRequestGracefulExitCalled
          ? intakeResTimeoutOnEnd : intakeResTimeout)
        log.trace({ timeout }, 'start intakeResTimer')
        intakeResTimer = setTimeout(() => {
          completePart('intakeRes',
            new Error('intake response timeout: APM server did not respond ' +
              `within ${timeout / 1000}s of gzip stream finish`))
        }, timeout).unref()
      }
    })
    // Watch the gzip "end" event for its completion, because the "close" event
    // that we would prefer to use, *does not get emitted* for the
    // `client.sendSpan(callback) + client.flush()` test case with
    // *node v12-only*.
    gzipStream.on('end', () => {
      log.trace('gzipStream "end"')
      completePart('gzipStream')
    })
    // gzipStream.on('close', () => { log.trace('gzipStream "close"') })

    // Hook up writing data to a file (only intended for local debugging).
    // Append the intake data to `payloadLogFile`, if given. This is only
    // intended for local debugging because it can have a significant perf
    // impact.
    if (client._conf.payloadLogFile) {
      const payloadLogStream = fs.createWriteStream(client._conf.payloadLogFile, { flags: 'a' })
      gzipStream.pipe(zlib.createGunzip()).pipe(payloadLogStream)
    }

    // Send the metadata object (always first) and hook up the streams.
    assert(client._encodedMetadata, 'client._encodedMetadata is set')
    gzipStream.write(client._encodedMetadata)
    gzipStream.pipe(intakeReq)
  }
}

/**
 * Some behaviors in the APM depend on the APM Server version. These are
 * exposed as `Client#supports...` boolean methods.
 *
 * These `Client#supports...` method names intentionally match those from the Java agent:
 * https://github.com/elastic/apm-agent-java/blob/master/apm-agent-core/src/main/java/co/elastic/apm/agent/report/ApmServerClient.java#L322-L349
 */
Client.prototype.supportsKeepingUnsampledTransaction = function () {
  // Default to assuming we are using a pre-8.0 APM Server if we haven't
  // yet fetched the version. There is no harm in sending unsampled transactions
  // to APM Server >=v8.0.
  if (!this._apmServerVersion) {
    return true
  } else {
    return this._apmServerVersion.major < 8
  }
}

/**
 * Signal to the Elastic AWS Lambda extension that a lambda function execution
 * is done.
 * https://github.com/elastic/apm/blob/main/specs/agents/tracing-instrumentation-aws-lambda.md#data-flushing
 *
 * @param {Function} cb() is called when finished. There are no arguments.
 */
Client.prototype._signalLambdaEnd = function (cb) {
  this._log.trace('_signalLambdaEnd start')
  const startTime = performance.now()
  const finish = errOrErrMsg => {
    const durationMs = performance.now() - startTime
    if (errOrErrMsg) {
      this._log.error({ err: errOrErrMsg, durationMs }, 'error signaling lambda invocation done')
    } else {
      this._log.trace({ durationMs }, 'signaled lambda invocation done')
    }
    cb()
  }

  // We expect to be talking to the localhost Elastic Lambda extension, so we
  // want a shorter timeout than `_conf.serverTimeout`.
  const TIMEOUT_MS = 5000

  const req = this._transportRequest(this._conf.requestSignalLambdaEnd, res => {
    res.on('error', err => {
      // Not sure this event can ever be emitted, but just in case.
      res.destroy(err)
    })
    res.resume()
    if (res.statusCode !== 202) {
      finish(`unexpected response status code: ${res.statusCode}`)
      return
    }
    res.on('end', function () {
      finish()
    })
  })
  req.setTimeout(TIMEOUT_MS)
  req.on('timeout', () => {
    this._log.trace('_signalLambdaEnd timeout')
    req.destroy(new Error(`timeout (${TIMEOUT_MS}ms) signaling Lambda invocation done`))
  })
  req.on('error', err => {
    finish(err)
  })
  req.end()
}

/**
 * Fetch the APM Server version and set `this._apmServerVersion`.
 * https://www.elastic.co/guide/en/apm/server/current/server-info.html
 *
 * If fetching/parsing fails then the APM server version will be set to `null`
 * to indicate "unknown version".
 */
Client.prototype._fetchApmServerVersion = function () {
  const setVerUnknownAndNotify = (errmsg) => {
    this._apmServerVersion = null // means "unknown version"
    if (isLambdaExecutionEnvironment) {
      // In a Lambda environment, where the process can be frozen, it is not
      // unusual for this request to hit an error. As long as APM Server version
      // fetching is not critical to tracing of Lambda invocations, then it is
      // preferable to not add an error message to the users log.
      this._log.debug('verfetch: ' + errmsg)
    } else {
      this.emit('request-error', new Error(errmsg))
    }
  }
  const headers = getHeaders(this._conf)
  // Explicitly do *not* pass in `this._agent` -- the keep-alive http.Agent
  // used for intake requests -- because the socket.ref() handling in
  // `Client#_ref()` conflicts with the socket.unref() below.
  const reqOpts = getBasicRequestOptions('GET', '/', headers, this._conf)
  reqOpts.timeout = 30000

  const req = this._transportGet(reqOpts, res => {
    res.on('error', err => {
      // Not sure this event can ever be emitted, but just in case
      res.destroy(err)
    })

    if (res.statusCode !== 200) {
      res.resume()
      setVerUnknownAndNotify(`unexpected status from APM Server information endpoint: ${res.statusCode}`)
      return
    }

    const chunks = []
    res.on('data', chunk => {
      chunks.push(chunk)
    })
    res.on('end', () => {
      if (chunks.length === 0) {
        setVerUnknownAndNotify('APM Server information endpoint returned no body, often this indicates authentication ("apiKey" or "secretToken") is incorrect')
        return
      }

      let serverInfo
      try {
        serverInfo = JSON.parse(Buffer.concat(chunks))
      } catch (parseErr) {
        setVerUnknownAndNotify(`could not parse APM Server information endpoint body: ${parseErr.message}`)
        return
      }

      if (serverInfo) {
        // APM Server 7.0.0 dropped the "ok"-level in the info endpoint body.
        const verStr = serverInfo.ok ? serverInfo.ok.version : serverInfo.version
        try {
          this._apmServerVersion = semver.SemVer(verStr)
        } catch (verErr) {
          setVerUnknownAndNotify(`could not parse APM Server version "${verStr}": ${verErr.message}`)
          return
        }
        this._log.debug({ apmServerVersion: verStr }, 'fetched APM Server version')
      } else {
        setVerUnknownAndNotify(`could not determine APM Server version from information endpoint body: ${JSON.stringify(serverInfo)}`)
      }
    })
  })

  req.on('socket', socket => {
    // Unref our socket to ensure this request does not keep the process alive.
    socket.unref()
  })
  req.on('timeout', () => {
    this._log.trace('_fetchApmServerVersion timeout')
    req.destroy(new Error(`timeout (${reqOpts.timeout}ms) fetching APM Server version`))
  })
  req.on('error', err => {
    setVerUnknownAndNotify(`error fetching APM Server version: ${err.message}`)
  })
}

/**
 * Fetches cloud metadata, if any, and encodes metadata (to `_encodedMetadata`).
 *
 * @param {function} cb - Called, with no arguments, when complete.
 */
Client.prototype._fetchAndEncodeMetadata = function (cb) {
  assert(this._conf.cloudMetadataFetcher, '_fetchAndEncodeMetadata should not be called without a configured cloudMetadataFetcher')
  this._conf.cloudMetadataFetcher.getCloudMetadata((err, cloudMetadata) => {
    if (err) {
      // We ignore this error (other than logging it). A common case, when
      // not running on one of the big 3 clouds, is "all callbacks failed",
      // which is *fine*. Because it is a common "error" we don't log the
      // stack trace.
      this._log.trace('getCloudMetadata err: %s', err)
    } else if (cloudMetadata) {
      this._cloudMetadata = cloudMetadata
    }
    this._resetEncodedMetadata()
    cb()
  })
}

function getIntakeRequestOptions (opts, agent) {
  const headers = getHeaders(opts)
  headers['Content-Type'] = 'application/x-ndjson'
  headers['Content-Encoding'] = 'gzip'

  return getBasicRequestOptions('POST', '/intake/v2/events', headers, opts, agent)
}

function getSignalLambdaEndRequestOptions (opts, agent) {
  const headers = getHeaders(opts)
  headers['Content-Length'] = 0

  return getBasicRequestOptions('POST', '/intake/v2/events?flushed=true', headers, opts, agent)
}

function getConfigRequestOptions (opts, agent) {
  const path = '/config/v1/agents?' + querystring.stringify({
    'service.name': opts.serviceName,
    'service.environment': opts.environment
  })

  const headers = getHeaders(opts)

  return getBasicRequestOptions('GET', path, headers, opts, agent)
}

function getBasicRequestOptions (method, defaultPath, headers, opts, agent) {
  return {
    agent: agent,
    rejectUnauthorized: opts.rejectUnauthorized !== false,
    ca: opts.serverCaCert,
    hostname: opts.serverUrl.hostname,
    port: opts.serverUrl.port,
    method,
    path: opts.serverUrl.pathname === '/' ? defaultPath : opts.serverUrl.pathname + defaultPath,
    headers
  }
}

function getHeaders (opts) {
  const headers = {}
  if (opts.secretToken) headers.Authorization = 'Bearer ' + opts.secretToken
  if (opts.apiKey) headers.Authorization = 'ApiKey ' + opts.apiKey
  headers.Accept = 'application/json'
  headers['User-Agent'] = opts.userAgent
  return Object.assign(headers, opts.headers)
}

function getMetadata (opts) {
  var payload = {
    service: {
      name: opts.serviceName,
      environment: opts.environment,
      runtime: {
        name: process.release.name,
        version: process.versions.node
      },
      language: {
        name: 'javascript'
      },
      agent: {
        name: opts.agentName,
        version: opts.agentVersion
      },
      framework: undefined,
      version: undefined,
      node: undefined
    },
    process: {
      pid: process.pid,
      ppid: process.ppid,
      title: process.title,
      argv: process.argv
    },
    system: {
      hostname: opts.hostname,
      architecture: process.arch,
      platform: process.platform,
      container: undefined,
      kubernetes: undefined
    },
    labels: opts.globalLabels
  }

  if (opts.serviceNodeName) {
    payload.service.node = {
      configured_name: opts.serviceNodeName
    }
  }

  if (opts.serviceVersion) payload.service.version = opts.serviceVersion

  if (opts.frameworkName || opts.frameworkVersion) {
    payload.service.framework = {
      name: opts.frameworkName,
      version: opts.frameworkVersion
    }
  }

  if (opts.containerId) {
    payload.system.container = {
      id: opts.containerId
    }
  }

  if (opts.kubernetesNodeName || opts.kubernetesNamespace || opts.kubernetesPodName || opts.kubernetesPodUID) {
    payload.system.kubernetes = {
      namespace: opts.kubernetesNamespace,
      node: opts.kubernetesNodeName
        ? { name: opts.kubernetesNodeName }
        : undefined,
      pod: (opts.kubernetesPodName || opts.kubernetesPodUID)
        ? { name: opts.kubernetesPodName, uid: opts.kubernetesPodUID }
        : undefined
    }
  }

  return payload
}

function destroyStream (stream) {
  if (stream instanceof zlib.Gzip ||
      stream instanceof zlib.Gunzip ||
      stream instanceof zlib.Deflate ||
      stream instanceof zlib.DeflateRaw ||
      stream instanceof zlib.Inflate ||
      stream instanceof zlib.InflateRaw ||
      stream instanceof zlib.Unzip) {
    // Zlib streams doesn't have a destroy function in Node.js 6. On top of
    // that simply calling destroy on a zlib stream in Node.js 8+ will result
    // in a memory leak as the handle isn't closed (an operation normally done
    // by calling close). So until that is fixed, we need to manually close the
    // handle after destroying the stream.
    //
    // PR: https://github.com/nodejs/node/pull/23734
    if (typeof stream.destroy === 'function') {
      // Manually close the stream instead of calling `close()` as that would
      // have emitted 'close' again when calling `destroy()`
      if (stream._handle && typeof stream._handle.close === 'function') {
        stream._handle.close()
        stream._handle = null
      }

      stream.destroy()
    } else if (typeof stream.close === 'function') {
      stream.close()
    }
  } else {
    // For other streams we assume calling destroy is enough
    if (typeof stream.destroy === 'function') stream.destroy()
    // Or if there's no destroy (which Node.js 6 will not have on regular
    // streams), emit `close` as that should trigger almost the same effect
    else if (typeof stream.emit === 'function') stream.emit('close')
  }
}

function oneOf (value, list) {
  return list.indexOf(value) >= 0
}

function normalizeGlobalLabels (labels) {
  if (!labels) return
  const result = {}

  for (const key of Object.keys(labels)) {
    const value = labels[key]
    result[key] = oneOf(typeof value, ['string', 'number', 'boolean'])
      ? value
      : value.toString()
  }

  return result
}

// https://httpwg.org/specs/rfc9111.html#cache-response-directive.max-age
function getMaxAge (res) {
  const header = res.headers['cache-control']
  if (!header) {
    return undefined
  }
  const match = header.match(/max-age=(\d+)/i)
  if (!match) {
    return undefined
  }
  return parseInt(match[1], 10)
}

// Wrap the given Error object, including the given message.
//
// Dev Note: Various techniques exist to wrap `Error`s in node.js and JavaScript
// to provide a cause chain, e.g. see
// https://www.joyent.com/node-js/production/design/errors
// However, I'm not aware of a de facto "winner". Eventually there may be
// https://github.com/tc39/proposal-error-cause
// For now we will simply prefix the existing error object's `message` property.
// This is simple and preserves the root error `stack`.
function wrapError (err, msg) {
  err.message = msg + ': ' + err.message
  return err
}

function processIntakeErrorResponse (res, buf) {
  const err = new Error('Unexpected APM Server response')

  err.code = res.statusCode

  if (buf.length > 0) {
    // https://www.elastic.co/guide/en/apm/server/current/events-api.html#events-api-errors
    const body = buf.toString('utf8')
    const contentType = res.headers['content-type']
    if (contentType && contentType.startsWith('application/json')) {
      try {
        const data = JSON.parse(body)
        err.accepted = data.accepted
        err.errors = data.errors
        if (!err.errors) err.response = body
      } catch (e) {
        err.response = body
      }
    } else {
      err.response = body
    }
  }

  return err
}

// Construct or decorate an Error instance from a failing response from the
// APM server central config endpoint.
//
// @param {IncomingMessage} res
// @param {Buffer|undefined} buf - Optional. A Buffer holding the response body.
// @param {Error|undefined} err - Optional. A cause Error instance.
function processConfigErrorResponse (res, buf, err) {
  // This library doesn't have a pattern for wrapping errors yet, so if
  // we already have an Error instance, we will just decorate it. That preserves
  // the stack of the root cause error.
  const errMsg = 'Unexpected APM Server response when polling config'
  if (!err) {
    err = new Error(errMsg)
  } else {
    err.message = errMsg + ': ' + err.message
  }

  err.code = res.statusCode

  if (buf && buf.length > 0) {
    const body = buf.toString('utf8')
    const contentType = res.headers['content-type']
    if (contentType && contentType.startsWith('application/json')) {
      try {
        const response = JSON.parse(body)
        if (typeof response === 'string') {
          err.response = response
        } else if (typeof response === 'object' && response !== null && typeof response.error === 'string') {
          err.response = response.error
        } else {
          err.response = body
        }
      } catch (e) {
        err.response = body
      }
    } else {
      err.response = body
    }
  }

  return err
}

// Return the time difference (in milliseconds) between the given time `t`
// (a 2-tuple as returned by `process.hrtime()`) and now.
function deltaMs (t) {
  const d = process.hrtime(t)
  return d[0] * 1e3 + d[1] / 1e6
}

/**
 * Performs a deep merge of `source` into `target`.  Mutates `target` only but
 * not its objects. Objects are merged, Arrays are not.
 *
 * @author inspired by [eden](https://gist.github.com/ahtcx/0cd94e62691f539160b32ecda18af3d6#gistcomment-2930530)
 */
function metadataMergeDeep (target, source) {
  const isObject = (obj) => obj && typeof obj === 'object' && !Array.isArray(obj)

  if (!isObject(target) || !isObject(source)) {
    return source
  }

  Object.keys(source).forEach(key => {
    const targetValue = target[key]
    const sourceValue = source[key]

    if (isObject(targetValue) && isObject(sourceValue)) {
      target[key] = metadataMergeDeep(Object.assign({}, targetValue), sourceValue)
    } else {
      target[key] = sourceValue
    }
  })

  return target
}

function deepClone (obj) {
  return JSON.parse(JSON.stringify(obj))
}
