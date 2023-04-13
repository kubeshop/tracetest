/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

const assert = require('assert')
const URL = require('url').URL

const otel = require('@opentelemetry/api')

const GenericSpan = require('../instrumentation/generic-span')
const oblog = require('./oblog')
const { otelSpanContextFromTraceContext, epochMsFromOTelTimeInput } = require('./otelutils')
const { RESULT_SUCCESS, OUTCOME_UNKNOWN, OUTCOME_SUCCESS, RESULT_FAILURE, OUTCOME_FAILURE } = require('../constants')
const Span = require('../instrumentation/span')
const Transaction = require('../instrumentation/transaction')

// Based on `isHomogeneousAttributeValueArray` from
// packages/opentelemetry-core/src/common/attributes.ts
function isHomogeneousArrayOfStrNumBool (arr) {
  const len = arr.length
  let elemType = null
  for (let i = 0; i < len; i++) {
    const elem = arr[i]
    if (elem === undefined || elem === null) {
      continue
    }
    if (!elemType) {
      elemType = typeof elem
      if (!(elemType === 'string' || elemType === 'number' || elemType === 'boolean')) {
        return false
      }
    } else if (typeof elem !== elemType) { // eslint-disable-line valid-typeof
      return false
    }
  }
  return true
}

// Set the given attribute key `k` and `v` according to these OTel rules:
// https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/common/README.md#attribute
function maybeSetOTelAttr (attrs, k, v) {
  if (Array.isArray(v)) {
    // Is it homogeneous? Nulls and undefineds are allowed.
    if (isHomogeneousArrayOfStrNumBool(v)) {
      attrs[k] = v.slice()
    }
  } else {
    switch (typeof v) {
      case 'number':
      case 'boolean':
        attrs[k] = v
        break
      case 'string':
        // Truncation (at 1024 bytes) is done in elastic-apm-http-client.
        attrs[k] = v
        break
    }
  }
}

// This wraps a core Transaction or Span in the OTel API's `inteface Span`.
class OTelSpan {
  constructor (span) {
    assert(span instanceof GenericSpan)
    this._span = span
    this._spanContext = null
  }

  toString () {
    return `OTelSpan<${this._span.constructor.name}<${this._span.id}, "${this._span.name}">>`
  }

  // ---- OTel interface Span
  // https://github.com/open-telemetry/opentelemetry-js-api/blob/v1.0.4/src/trace/span.ts

  spanContext () {
    oblog.apicall('%s.spanContext()', this)
    if (!this._spanContext) {
      this._spanContext = otelSpanContextFromTraceContext(this._span._context)
    }
    return this._spanContext
  }

  setAttribute (key, value) {
    if (this._span.ended || !key || typeof key !== 'string') {
      return this
    }

    const attrs = this._span._getOTelAttributes()
    maybeSetOTelAttr(attrs, key, value)

    return this
  }

  setAttributes (attributes) {
    if (this._span.ended || !attributes || typeof attributes !== 'object') {
      return this
    }

    const attrs = this._span._getOTelAttributes()
    for (const k in attributes) {
      if (k.length === 0) continue
      maybeSetOTelAttr(attrs, k, attributes[k])
    }

    return this
  }

  // Span events are not currently supported.
  addEvent (name, attributesOrStartTime, startTime) {
    return this
  }

  setStatus (otelSpanStatus) {
    if (this._span.ended) {
      return this
    }
    switch (otelSpanStatus) {
      case otel.SpanStatusCode.ERROR:
        this._span.setOutcome(OUTCOME_FAILURE)
        break
      case otel.SpanStatusCode.OK:
        this._span.setOutcome(OUTCOME_SUCCESS)
        break
      case otel.SpanStatusCode.UNSET:
        this._span.setOutcome(OUTCOME_UNKNOWN)
        break
    }
    // Also set transaction.result, similar to the Java APM agent.
    if (this._span instanceof Transaction) {
      switch (otelSpanStatus) {
        case otel.SpanStatusCode.ERROR:
          this._span.result = RESULT_FAILURE
          break
        case otel.SpanStatusCode.OK:
          this._span.result = RESULT_SUCCESS
          break
      }
    }
    return this
  }

  updateName (name) {
    if (this._span.ended) {
      return this
    }
    this._span.name = name
    return this
  }

  end (otelEndTime) {
    oblog.apicall('%s.end(endTime=%s)', this, otelEndTime)
    const endTime = otelEndTime === undefined
      ? undefined
      : epochMsFromOTelTimeInput(otelEndTime)
    if (this._span instanceof Transaction) {
      this._transCompatMapping()
      this._span.end(undefined, endTime)
    } else {
      assert(this._span instanceof Span)
      this._spanCompatMapping()
      this._span.end(endTime)
    }
  }

  // https://github.com/elastic/apm/blob/main/specs/agents/tracing-api-otel.md#compatibility-mapping
  _transCompatMapping () {
    const attrs = this._span._otelAttributes
    let type = 'unknown'
    if (!attrs) {
      this._span.type = type
      return
    }

    const otelKind = otel.SpanKind[this._span._otelKind] // map from string to int form
    const isRpc = attrs['rpc.system'] !== undefined
    const isHttp = attrs['http.url'] !== undefined || attrs['http.scheme'] !== undefined
    const isMessaging = attrs['messaging.system'] !== undefined
    if (otelKind === otel.SpanKind.SERVER && (isRpc || isHttp)) {
      type = 'request'
    } else if (otelKind === otel.SpanKind.CONSUMER && isMessaging) {
      type = 'messaging'
    } else {
      type = 'unknown'
    }
    this._span.type = type
  }

  // https://github.com/elastic/apm/blob/main/specs/agents/tracing-api-otel.md#compatibility-mapping
  _spanCompatMapping () {
    const attrs = this._span._otelAttributes
    const otelKind = otel.SpanKind[this._span._otelKind] // map from string to int form
    if (!attrs) {
      if (otelKind === otel.SpanKind.INTERNAL) {
        this._span.type = 'app'
        this._span.subtype = 'internal'
      } else {
        this._span.type = 'unknown'
      }
      return
    }

    let type
    let subtype
    let serviceTargetType = null
    let serviceTargetName = null

    const httpPortFromScheme = function (scheme) {
      if (scheme === 'http') {
        return 80
      } else if (scheme === 'https') {
        return 443
      }
      return -1
    }

    // Extracts 'host:port' from URL.
    const parseNetName = function (url) {
      let u
      try {
        u = new URL(url) // https://developer.mozilla.org/en-US/docs/Web/API/URL
      } catch (_err) {
        return undefined
      }
      if (u.port !== '') {
        return u.host // host:port already in URL
      } else {
        var port = httpPortFromScheme(u.protocol.substring(0, u.protocol.length - 1))
        return port > 0 ? u.hostname + ':' + port : u.hostname
      }
    }

    let netPort = attrs['net.peer.port'] || -1
    const netPeer = attrs['net.peer.name'] || attrs['net.peer.ip']
    let netName = netPeer // netName includes port, if provided
    if (netName && netPort > 0) {
      netName += ':' + netPort
    }

    if (attrs['db.system']) {
      type = 'db'
      subtype = attrs['db.system']
      serviceTargetType = subtype
      serviceTargetName = attrs['db.name'] || null
    } else if (attrs['messaging.system']) {
      type = 'messaging'
      subtype = attrs['messaging.system']
      if (!netName && attrs['messaging.url']) {
        netName = parseNetName(attrs['messaging.url'])
      }
      serviceTargetType = subtype
      serviceTargetName = attrs['messaging.destination'] || null
    } else if (attrs['rpc.system']) {
      type = 'external'
      subtype = attrs['rpc.system']
      serviceTargetType = subtype
      serviceTargetName = netName || attrs['rpc.service'] || null
    } else if (attrs['http.url'] || attrs['http.scheme']) {
      type = 'external'
      subtype = 'http'
      serviceTargetType = 'http'
      const httpHost = attrs['http.host'] || netPeer
      if (httpHost) {
        if (netPort < 0) {
          netPort = httpPortFromScheme(attrs['http.scheme'])
        }
        serviceTargetName = netPort < 0 ? httpHost : httpHost + ':' + netPort
      } else if (attrs['http.url']) {
        serviceTargetName = parseNetName(attrs['http.url'])
      }
    }

    if (type === undefined) {
      if (otelKind === otel.SpanKind.INTERNAL) {
        type = 'app'
        subtype = 'internal'
      } else {
        type = 'unknown'
      }
    }

    this._span.type = type
    if (subtype) {
      this._span.subtype = subtype
    }
    if (serviceTargetType || serviceTargetName) {
      this._span.setServiceTarget(serviceTargetType, serviceTargetName)
    }
  }

  isRecording () {
    return !this._span.ended && this._span.sampled
  }

  recordException (otelException, otelTime) {
    const errOpts = {
      parent: this._span,
      captureAttributes: false,
      skipOutcome: true
    }
    if (otelTime !== undefined) {
      errOpts.timestamp = epochMsFromOTelTimeInput(otelTime)
    }

    this._span._agent.captureError(otelException, errOpts)
  }
}

module.exports = {
  OTelSpan
}
