/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

// Indirect usage of the singleton `Agent` to log.
function getLogger () {
  return require('../..').logger
}

/**
 * Class for Managing Tracestate
 *
 * Class that creates objects for managing trace state.
 * This class is capable of parsing both tracestate strings
 * and tracestate binary representations, allowing clients
 * to get and set values in a single list-member/namespace
 * while preserving values in the other namespaces.
 *
 * Capable of working with either the binary of string
 * formatted tracestate values.
 *
 * Usage:
 *   const tracestate = TraceState.fromStringFormatString(headerTracestate, 'es')
 *   tracestate.setValue('s',1)
 *   const newHeader = tracestate.toW3cString()
 */
class TraceState {
  constructor (sourceBuffer, listMemberNamespace = 'es', defaultValues = {}) {
    if (!this._validateVendorKey(listMemberNamespace)) {
      throw new Error('Vendor namespace failed validation.')
    }

    // buffer representation of the trace state string.
    // The initial value of this.buffer will keep the
    // values set in the listMemberNamespace list-member,
    // but as soon as an initial value is set (via setValue)
    // then the listMemberNamespace values will be removed
    // from this.buffer and stored in the this.values.  While
    // slightly more complicated, this allows us to maintain
    // the order of list-member keys in an un-mutated tracestate
    // string, per the W3C spec
    this.buffer = sourceBuffer

    this.listMemberNamespace = listMemberNamespace

    // values for our namespace, set via setValue to
    // ensure names conform
    this.values = {}
    for (const key in defaultValues) {
      const value = defaultValues[key]
      this.setValue(key, value)
    }
  }

  setValue (key, value) {
    const strKey = String(key)
    const strValue = String(value)
    if (!this._validateElasicKeyAndValue(strKey, strValue)) {
      getLogger().trace('could not set tracestate key, invalid characters detected')
      return false
    }
    const isFirstSet = (Object.keys(this.values).length === 0)
    const oldValue = this.values[strKey]
    this.values[strKey] = value

    // per: https://github.com/elastic/apm/blob/d5b2c87326548befcfec6731713932a00e430b99/specs/agents/tracing-distributed-tracing.md
    // If adding another key/value pair to the es entry would exceed the limit
    // of 256 chars, that key/value pair MUST be ignored by agents.
    // The key/value and entry separators : and ; have to be considered as well.
    const serializedValue = this._serializeValues(this.values)
    if (serializedValue.length > 256 && (typeof oldValue === 'undefined')) {
      delete this.values[strKey]
      return false
    }
    if (serializedValue.length > 256 && (typeof oldValue !== 'undefined')) {
      this.values[strKey] = oldValue
      return false
    }

    // the first time we set a value, extract the mutable values from the
    // buffer and set this.values appropriately
    if (isFirstSet && Object.keys(this.values).length === 1) {
      const [buffer, values] = TraceState._removeMemberNamespaceFromBuffer(
        this.buffer,
        this.listMemberNamespace
      )
      this.buffer = buffer
      this.values = values
      values[strKey] = value
    }

    return true
  }

  getValue (keyToGet) {
    const allValues = this.toObject()
    const rawValue = allValues[this.listMemberNamespace]
    if (!rawValue) {
      return rawValue
    }
    const values = TraceState._parseValues(rawValue)
    return values[keyToGet]
  }

  toHexString () {
    const newBuffer = Buffer.alloc(this.buffer.length)
    let newBufferOffset = 0
    for (let i = 0; i < this.buffer.length; i++) {
      const byte = this.buffer[i]
      if (byte === 0) {
        const indexOfKeyLength = i + 1
        const indexOfKey = i + 2
        const lengthKey = this.buffer[indexOfKeyLength]

        const indexOfValueLength = indexOfKey + lengthKey
        const indexOfValue = indexOfValueLength + 1
        const lengthValue = this.buffer[indexOfValueLength]

        const key = this.buffer.slice(indexOfKey, indexOfKey + lengthKey).toString()
        // bail out if this is our mutable namespace
        if (key === this.listMemberNamespace) { continue }

        // if this is not our key copy from the `0` byte to the end of the value
        this.buffer.copy(newBuffer, newBufferOffset, i, indexOfValue + lengthValue)
        newBufferOffset += (indexOfValue + lengthValue)

        // skip ahead to first byte after end of value
        i = indexOfValue + lengthValue - 1
        continue
      }
    }

    // now serialize the internal representation
    const ourBytes = []
    if (Object.keys(this.values).length > 0) {
      // the zero byte
      ourBytes.push(0)

      // the length of the vendor namespace
      ourBytes.push(this.listMemberNamespace.length)
      // the chars of the vendor namespace
      for (let i = 0; i < this.listMemberNamespace.length; i++) {
        ourBytes.push(this.listMemberNamespace.charCodeAt(i))
      }

      // add the length of the value
      const serializedValue = this._serializeValues(this.values)
      ourBytes.push(serializedValue.length)

      // add the bytes of the value
      for (let i = 0; i < serializedValue.length; i++) {
        ourBytes.push(serializedValue.charCodeAt(i))
      }
    }
    const ourBuffer = Buffer.from(ourBytes)
    return Buffer.concat(
      [newBuffer, ourBuffer],
      newBuffer.length + ourBuffer.length
    ).toString('hex')
  }

  /**
   * Returns JSON reprenstation of tracestate key/value pairs
   *
   * Does not parse the mutable list namespace
   */
  toObject () {
    const map = this.toMap()
    const obj = {}
    for (const key of map.keys()) {
      obj[key] = map.get(key)
    }
    return obj
  }

  toMap () {
    const map = new Map()

    // first, serialize values from the internal representation. This means
    // The W3C spec dictates that mutated values need to be on
    // the left of the new trace string
    if (Object.keys(this.values).length) {
      map.set(this.listMemberNamespace, this._serializeValues(
        this.values
      ))
    }
    for (let i = 0; i < this.buffer.length; i++) {
      const byte = this.buffer[i]
      if (byte === 0) {
        const indexOfKeyLength = i + 1
        const indexOfKey = i + 2
        const lengthKey = this.buffer[indexOfKeyLength]

        const indexOfValueLength = indexOfKey + lengthKey
        const indexOfValue = indexOfValueLength + 1
        const lengthValue = this.buffer[indexOfValueLength]

        const key = this.buffer.slice(indexOfKey, indexOfKey + lengthKey).toString()
        const value = this.buffer.slice(indexOfValue, indexOfValue + lengthValue).toString()

        map.set(key, value)

        // skip ahead
        i = indexOfValue + lengthValue - 1
        continue
      }
    }

    return map
  }

  toString () {
    return this.toW3cString()
  }

  toW3cString () {
    const json = this.toObject()
    const chars = []
    for (const key in json) {
      const value = json[key]
      if (!value) { continue }
      chars.push(key)
      chars.push('=')
      chars.push(value)
      chars.push(',')
    }
    chars.pop() // remove final comma
    return chars.join('')
  }

  _serializeValues (keyValues) {
    const chars = []
    for (const key in keyValues) {
      const value = keyValues[key]
      chars.push(`${key}:${value}`)
      chars.push(';')
    }
    chars.pop() // last semi-colon
    return chars.join('')
  }

  _validateVendorKey (key) {
    if (key.length > 256 || key.length < 1) {
      return false
    }

    const re = new RegExp(
      '^[abcdefghijklmnopqrstuvwxyz0123456789_\\-\\*/]*$'
    )
    if (!key.match(re)) {
      return false
    }
    return true
  }

  _validateElasicKeyAndValue (key, value) {
    // 0x20` to `0x7E WITHOUT `,` or `=` or `;` or `;`
    const re = /^[ \][!"#$%&'()*+\-./0123456789<>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ^_abcdefghijklmnopqrstuvwxyz{|}~]*$/

    if (!key.match(re) || !value.match(re)) {
      return false
    }

    if (key.length > 256 || value.length > 256) {
      return false
    }

    return true
  }

  static fromBinaryFormatHexString (string, listMemberNamespace = 'es') {
    return new TraceState(Buffer.from(string, 'hex'), listMemberNamespace)
  }

  static fromStringFormatString (string = '', listMemberNamespace = 'es') {
    // converts string format to byte format
    const bytes = []

    const parts = string.split(',')
    for (let part of parts) {
      part = part.trim() // optional whitespace (OWS)
      if (!part) { continue }
      const [listMember, value] = part.split('=')
      if (!listMember || !value) { continue }
      bytes.push(0)
      bytes.push(listMember.length)
      for (let i = 0; i < listMember.length; i++) {
        bytes.push(listMember.charCodeAt(i))
      }
      bytes.push(value.length)
      for (let i = 0; i < value.length; i++) {
        bytes.push(value.charCodeAt(i))
      }
    }

    return new TraceState(Buffer.from(bytes), listMemberNamespace)
  }

  static _parseValues (rawValues) {
    const parsedValues = {}
    const parts = rawValues.split(';')
    for (const keyValue of parts) {
      if (!keyValue) { continue }
      const [key, value] = keyValue.split(':')
      if (!key || !value) { continue }
      parsedValues[key] = value
    }
    return parsedValues
  }

  static _removeMemberNamespaceFromBuffer (buffer, listMemberNamespace) {
    const newBuffer = Buffer.alloc(buffer.length)
    let newBufferOffset = 0
    const values = {}
    for (let i = 0; i < buffer.length; i++) {
      const byte = buffer[i]
      if (byte === 0) {
        const indexOfKeyLength = i + 1
        const indexOfKey = i + 2
        const lengthKey = buffer[indexOfKeyLength]

        const indexOfValueLength = indexOfKey + lengthKey
        const indexOfValue = indexOfValueLength + 1
        const lengthValue = buffer[indexOfValueLength]

        const key = buffer.slice(indexOfKey, indexOfKey + lengthKey).toString()

        // if this is our mutable namespace extract
        // and set the value in values, otherwise
        // copy into new buffer
        if (key === listMemberNamespace) {
          const rawValues = buffer.slice(indexOfValue, indexOfValue + lengthValue).toString()
          const parsedValues = TraceState._parseValues(rawValues)
          for (const key in parsedValues) {
            values[key] = parsedValues[key]
          }
          continue
        } else {
          buffer.copy(newBuffer, newBufferOffset, i, indexOfValue + lengthValue)
          newBufferOffset += (indexOfValue + lengthValue - i)
        }

        // skip ahead to first byte after end of value
        i = indexOfValue + lengthValue - 1
        continue
      }
    }

    // trim off extra 0 bytes
    const trimmedBuffer = newBuffer.slice(0, newBufferOffset)

    return [
      trimmedBuffer,
      values
    ]
  }
}

module.exports = TraceState
