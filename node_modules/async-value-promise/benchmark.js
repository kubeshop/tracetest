'use strict'

const benchList = require('benchmark-fn-list')
const benchPrettyPrint = require('benchmark-fn-pretty-print')
const AsyncValuePromise = require('./')

benchList([
  {
    name: 'then before resolve',
    iterations: 100000,
    task(cb) {
      var result = new AsyncValuePromise()
      result.then(cb, cb)
      result.resolve('hello')
    }
  },
  {
    name: 'resolve before then',
    iterations: 100000,
    task(cb) {
      var result = new AsyncValuePromise()
      result.resolve('hello')
      result.then(cb, cb)
    }
  },
  {
    name: 'then before reject',
    iterations: 100000,
    task(cb) {
      var result = new AsyncValuePromise()
      result.then(cb, cb)
      result.reject('hello')
    }
  },
  {
    name: 'reject before then',
    iterations: 100000,
    task(cb) {
      var result = new AsyncValuePromise()
      result.reject('hello')
      result.then(cb, cb)
    }
  },
  {
    name: 'resolve in constructor',
    iterations: 100000,
    task(cb) {
      var result = new AsyncValuePromise('hello')
      result.then(cb, cb)
    }
  },
  {
    name: 'reject in constructor',
    iterations: 100000,
    task(cb) {
      var result = new AsyncValuePromise(null, 'hello')
      result.then(cb, cb)
    }
  },
  {
    name: 'promise',
    iterations: 100000,
    task(cb) {
      var promise = Promise.resolve('hello')
      promise.then(cb)
    }
  }
], results => {
  console.log(benchPrettyPrint(results))
})
