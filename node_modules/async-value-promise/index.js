'use strict'

const AsyncValue = require('async-value')
const assert = require('assert')

function isDefined(thing) {
  return thing !== null & thing !== undefined
}

module.exports = class AsyncValuePromise {
  constructor(value, error) {
    this.value = new AsyncValue()
    this.error = new AsyncValue()

    if (isDefined(error)) {
      this.reject(error)
    } else if (isDefined(value)) {
      this.resolve(value)
    }
  }

  then(pass, fail) {
    const promise = new AsyncValuePromise()

    if (this.value) {
      if (!pass) {
        this.value.send(promise.value)
      } else {
        this.value.get(value => {
          promise.resolve(pass(value))
        })
      }
    }

    if (this.error) {
      if (!fail) {
        this.error.send(promise.error)
      } else {
        this.error.get(error => {
          try {
            promise.resolve(fail(error))
          } catch (err) {
            promise.reject(err)
          }
        })
      }
    }

    return promise
  }

  resolve(value) {
    assert(this.value, 'already rejected')
    if (value instanceof AsyncValuePromise) {
      value.value.send(this.value)
      if (this.error) {
        value.error.send(this.error)
      }
    } else {
      this.value.set(value)
      this.error = null
    }
  }

  reject(error) {
    assert(this.error, 'already resolved')
    if (error instanceof AsyncValuePromise) {
      error.error.send(this.error)
      if (this.value) {
        error.value.send(this.value)
      }
    } else {
      this.error.set(error)
      this.value = null
    }
  }

  catch(fail) {
    return this.then(null, fail)
  }

  static resolve(value) {
    return new AsyncValuePromise(value)
  }

  static reject(error) {
    return new AsyncValuePromise(null, error)
  }

  static all(promises) {
    const promise = new AsyncValuePromise()
    const count = promises.length
    const results = new Array(count)
    let remaining = count
    let failed = false

    function fail(error) {
      if (!failed) {
        promise.reject(error)
        failed = true
      }
    }

    if (count === 0) {
      promise.resolve([])
    } else {
      for (let i = 0; i < count; i++) {
        promises[i].then(value => {
          if (failed) return
          results[i] = value
          if (--remaining === 0) {
            promise.resolve(results)
          }
        }).catch(fail)
      }
    }

    return promise
  }
}
