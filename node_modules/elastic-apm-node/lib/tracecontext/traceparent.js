/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

const { randomFillSync } = require('crypto')

const SIZES = {
  version: 1,
  traceId: 16,
  id: 8,
  flags: 1,
  parentId: 8,

  // Aggregate sizes
  ids: 24, // traceId + id
  all: 34
}

const OFFSETS = {
  version: 0,
  traceId: SIZES.version,
  id: SIZES.version + SIZES.traceId,
  flags: SIZES.version + SIZES.ids,

  // Additional parentId is stored after the header content
  parentId: SIZES.version + SIZES.ids + SIZES.flags
}

const FLAGS = {
  recorded: 0b00000001
}

function defineLazyProp (obj, prop, fn) {
  Object.defineProperty(obj, prop, {
    configurable: true,
    enumerable: true,
    get () {
      const value = fn()
      if (value !== undefined) {
        Object.defineProperty(obj, prop, {
          configurable: true,
          enumerable: true,
          value
        })
      }
      return value
    }
  })
}

function hexSliceFn (buffer, offset, length) {
  return () => buffer.slice(offset, length).toString('hex')
}

function maybeHexSliceFn (buffer, offset, length) {
  const fn = hexSliceFn(buffer, offset, length)
  return () => {
    const value = fn()
    // Check for any non-zero characters to identify a valid ID
    if (/[1-9a-f]/.test(value)) {
      return value
    }
  }
}

function makeChild (buffer) {
  // Move current id into parentId region
  buffer.copy(buffer, OFFSETS.parentId, OFFSETS.id, OFFSETS.flags)

  // Generate new id
  randomFillSync(buffer, OFFSETS.id, SIZES.id)

  return new TraceParent(buffer)
}

function isValidHeader (header) {
  return /^[\da-f]{2}-[\da-f]{32}-[\da-f]{16}-[\da-f]{2}$/.test(header)
}

// NOTE: The version byte is not fully supported yet, but is not important until
// we use the official header name rather than elastic-apm-traceparent.
// https://w3c.github.io/distributed-tracing/report-trace-context.html#versioning-of-traceparent
function headerToBuffer (header) {
  const buffer = Buffer.alloc(SIZES.all)
  buffer.write(header.replace(/-/g, ''), 'hex')
  return buffer
}

function resume (header) {
  return makeChild(headerToBuffer(header))
}

function start (sampled = false) {
  const buffer = Buffer.alloc(SIZES.all)

  // Generate new ids
  randomFillSync(buffer, OFFSETS.traceId, SIZES.ids)

  if (sampled) {
    buffer[OFFSETS.flags] |= FLAGS.recorded
  }

  return new TraceParent(buffer)
}

const bufferSymbol = Symbol('trace-context-buffer')

class TraceParent {
  constructor (buffer) {
    this[bufferSymbol] = buffer
    Object.defineProperty(this, 'recorded', {
      enumerable: true,
      get: function () {
        return !!(this[bufferSymbol][OFFSETS.flags] & FLAGS.recorded)
      }
    })

    defineLazyProp(this, 'version', hexSliceFn(buffer, OFFSETS.version, OFFSETS.traceId))
    defineLazyProp(this, 'traceId', hexSliceFn(buffer, OFFSETS.traceId, OFFSETS.id))
    defineLazyProp(this, 'id', hexSliceFn(buffer, OFFSETS.id, OFFSETS.flags))
    defineLazyProp(this, 'flags', hexSliceFn(buffer, OFFSETS.flags, OFFSETS.parentId))
    defineLazyProp(this, 'parentId', maybeHexSliceFn(buffer, OFFSETS.parentId))
  }

  static startOrResume (childOf, conf) {
    if (childOf instanceof TraceParent) return childOf.child()
    if (childOf && childOf._context instanceof TraceParent) return childOf._context.child()

    return isValidHeader(childOf)
      ? resume(childOf)
      : start(Math.random() <= conf.transactionSampleRate)
  }

  static fromString (header) {
    return new TraceParent(headerToBuffer(header))
  }

  ensureParentId () {
    let id = this.parentId
    if (!id) {
      randomFillSync(this[bufferSymbol], OFFSETS.parentId, SIZES.id)
      id = this.parentId
    }
    return id
  }

  child () {
    return makeChild(Buffer.from(this[bufferSymbol]))
  }

  setRecorded (value) {
    if (value) {
      this[bufferSymbol][OFFSETS.flags] |= FLAGS.recorded
    } else {
      this[bufferSymbol][OFFSETS.flags] &= ~FLAGS.recorded
    }
  }

  toString () {
    return `${this.version}-${this.traceId}-${this.id}-${this.flags}`
  }
}

TraceParent.FLAGS = FLAGS

module.exports = {
  TraceParent
}
