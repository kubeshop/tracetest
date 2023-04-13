/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

const truncate = require('unicode-byte-truncate')

const config = require('../config')
const constants = require('../constants')
const { SpanCompression } = require('./span-compression')
const Timer = require('./timer')
const TraceContext = require('../tracecontext')
const { TraceParent } = require('../tracecontext/traceparent')

module.exports = GenericSpan

function GenericSpan (agent, ...args) {
  const opts = typeof args[args.length - 1] === 'object'
    ? (args.pop() || {})
    : {}

  this._timer = new Timer(opts.timer, opts.startTime)

  this._context = TraceContext.startOrResume(opts.childOf, agent._conf, opts.tracestate)
  this._hasPropagatedTraceContext = false

  this._parentSpan = null
  if (opts.childOf instanceof GenericSpan) {
    this.setParentSpan(opts.childOf)
  }
  this._compression = new SpanCompression(agent)
  this._compression.setBufferedSpan(null)

  this._agent = agent
  this._labels = null
  this._ids = null // Populated by sub-types of GenericSpan
  this._otelKind = null
  this._otelAttributes = null

  this._links = []
  if (opts.links) {
    for (let i = 0; i < opts.links.length; i++) {
      const link = linkFromLinkArg(opts.links[i])
      if (link) {
        this._links.push(link)
      }
    }
  }

  this.timestamp = this._timer.start
  this.ended = false
  this._duration = null // Duration in milliseconds. Set on `.end()`.
  this._endTimestamp = null

  this.outcome = constants.OUTCOME_UNKNOWN

  // Freezing the outcome allows us to prefer a value set from
  // from the API and allows a span to keep its unknown status
  // even if it succesfully ends.
  this._isOutcomeFrozen = false

  this.type = null
  this.subtype = null
  this.action = null
  this.setType(...args)
}

Object.defineProperty(GenericSpan.prototype, 'id', {
  enumerable: true,
  get () {
    return this._context.traceparent.id
  }
})

Object.defineProperty(GenericSpan.prototype, 'traceId', {
  enumerable: true,
  get () {
    return this._context.traceparent.traceId
  }
})

Object.defineProperty(GenericSpan.prototype, 'parentId', {
  enumerable: true,
  get () {
    return this._context.traceparent.parentId
  }
})

Object.defineProperty(GenericSpan.prototype, 'sampled', {
  enumerable: true,
  get () {
    return this._context.traceparent.recorded
  }
})

Object.defineProperty(GenericSpan.prototype, 'sampleRate', {
  enumerable: true,
  get () {
    const rate = parseFloat(this._context.tracestate.getValue('s'))
    if (rate >= 0 && rate <= 1) {
      return rate
    }
    return null
  }
})

Object.defineProperty(GenericSpan.prototype, 'traceparent', {
  enumerable: true,
  get () {
    return this._context.toString()
  }
})

// The duration of the span, in milliseconds.
GenericSpan.prototype.duration = function () {
  if (!this.ended) {
    this._agent.logger.debug('tried to call duration() on un-ended transaction/span %o', { id: this.id, parent: this.parentId, trace: this.traceId, name: this.name, type: this.type })
    return null
  }

  return this._duration
}

// The 'stringify' option is for backward compatibility and will likely be
// removed in the next major version.
GenericSpan.prototype.setLabel = function (key, value, stringify = true) {
  const makeLabelValue = () => {
    if (!stringify && (typeof value === 'boolean' || typeof value === 'number')) {
      return value
    }

    return truncate(String(value), config.INTAKE_STRING_MAX_SIZE)
  }

  if (!key) return false
  if (!this._labels) this._labels = {}
  var skey = key.replace(/[.*"]/g, '_')
  if (key !== skey) {
    this._agent.logger.warn('Illegal characters used in tag key: %s', key)
  }
  this._labels[skey] = makeLabelValue()
  return true
}

GenericSpan.prototype.addLabels = function (labels, stringify) {
  if (!labels) return false
  var keys = Object.keys(labels)
  for (const key of keys) {
    if (!this.setLabel(key, labels[key], stringify)) {
      return false
    }
  }
  return true
}

// This method is private because the APM agents spec says that (for OTel
// compat), adding links after span creation should not be allowed.
// https://github.com/elastic/apm/blob/main/specs/agents/span-links.md
//
// To support adding span links for SQS ReceiveMessage and equivalent, the
// message data isn't known until the *response*, after the span has been
// created.
//
// @param {Array} links - An array of objects with a `context` property that is
//    a Transaction, Span, or TraceParent instance, or a W3C trace-context
//    'traceparent' string.
GenericSpan.prototype._addLinks = function (links) {
  if (links) {
    for (let i = 0; i < links.length; i++) {
      const link = linkFromLinkArg(links[i])
      if (link) {
        this._links.push(link)
      }
    }
  }
}

GenericSpan.prototype.setType = function (type = null, subtype = null, action = null) {
  this.type = type || constants.DEFAULT_SPAN_TYPE
  this.subtype = subtype
  this.action = action
}

GenericSpan.prototype._freezeOutcome = function () {
  this._isOutcomeFrozen = true
}

GenericSpan.prototype._isValidOutcome = function (outcome) {
  return outcome === constants.OUTCOME_FAILURE ||
    outcome === constants.OUTCOME_SUCCESS ||
    outcome === constants.OUTCOME_UNKNOWN
}

GenericSpan.prototype.propagateTraceContextHeaders = function (carrier, setter) {
  this._hasPropagatedTraceContext = true
  return this._context.propagateTraceContextHeaders(carrier, setter)
}
GenericSpan.prototype.setParentSpan = function (span) {
  this._parentSpan = span
}

GenericSpan.prototype.getParentSpan = function (span) {
  return this._parentSpan
}

GenericSpan.prototype.getBufferedSpan = function () {
  return this._compression.getBufferedSpan()
}

GenericSpan.prototype.setBufferedSpan = function (span) {
  return this._compression.setBufferedSpan(span)
}

GenericSpan.prototype.isCompositeSameKind = function () {
  return this._compression.isCompositeSameKind()
}

GenericSpan.prototype.isComposite = function () {
  return this._compression.isComposite()
}

GenericSpan.prototype.getCompositeSum = function () {
  return this._compression.composite.sum
}

// https://github.com/elastic/apm/blob/main/specs/agents/tracing-api-otel.md#span-kind
// @param {String} kind
GenericSpan.prototype._setOTelKind = function (kind) {
  this._otelKind = kind
}

// This returns the internal OTel attributes object so it can be mutated.
GenericSpan.prototype._getOTelAttributes = function () {
  if (!this._otelAttributes) {
    this._otelAttributes = {}
  }
  return this._otelAttributes
}

// Serialize OTel-related fields into the given payload, if any.
GenericSpan.prototype._serializeOTel = function (payload) {
  if (this._otelKind) {
    payload.otel = {
      span_kind: this._otelKind
    }
    if (this._otelAttributes) {
      // Though the spec allows it ("MAY"), we are opting *not* to report OTel
      // span attributes as labels for older (<7.16) versions of APM server.
      // This is to avoid the added complexity of guarding allowed attribute
      // value types to those supported by the APM server intake API.
      payload.otel.attributes = this._otelAttributes
    }
  }
}

// Translate a `opts.links` entry (see the `Link` type in "index.d.ts") to a
// span link as it will be serialized and sent to APM server. If the linkArg is
// invalid, this will return null.
//
// @param {Object} linkArg - An object with a `context` property that is a
//    Transaction, Span, or TraceParent instance, or a W3C trace-context
//    'traceparent' string.
function linkFromLinkArg (linkArg) {
  if (!linkArg || !linkArg.context) {
    return null
  }

  const ctx = linkArg.context
  let traceId
  let spanId

  if (ctx._context instanceof TraceContext) { // Transaction or Span
    traceId = ctx._context.traceparent.traceId
    spanId = ctx._context.traceparent.id
  } else if (ctx instanceof TraceParent) {
    traceId = ctx.traceId
    spanId = ctx.id
  } else if (typeof (ctx) === 'string') {
    // Note: Unfortunately TraceParent.fromString doesn't validate the string.
    const traceparent = TraceParent.fromString(ctx)
    traceId = traceparent.traceId
    spanId = traceparent.id
  } else {
    return null
  }

  return {
    trace_id: traceId,
    span_id: spanId
  }
}
