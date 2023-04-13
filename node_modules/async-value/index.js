'use strict'

const assert = require('assert')

module.exports = class AsyncValue {
  constructor(value) {
    this.value = null
    this.callbacks = []
    if (value !== null && value !== undefined) {
      this.set(value)
    }
  }

  get(callback) {
    assert(typeof callback === 'function', 'callback must be a function')
    if (this.value !== null) {
      callback(this.value)
    } else {
      this.callbacks.push(callback)
    }
  }

  set(value) {
    assert(this.value === null, 'value can only be set once')
    if (value instanceof AsyncValue) {
      value.send(this)
      return
    }
    for (let callback of this.callbacks) {
      callback(value)
    }
    this.callbacks = null
    this.value = value
  }

  send(target) {
    assert(target instanceof AsyncValue, 'send target must be an async value')
    this.get(target.set.bind(target))
  }
}
