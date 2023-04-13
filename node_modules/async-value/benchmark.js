'use strict'

const benchList = require('benchmark-fn-list')
const benchPrettyPrint = require('benchmark-fn-pretty-print')
const AsyncValue = require('./')

benchList([
  {
    name: 'set in constructor',
    iterations: 100000,
    task(cb) {
      var value = new AsyncValue('hello')
      value.get(cb)
    }
  },
  {
    name: 'set then get',
    iterations: 100000,
    task(cb) {
      var value = new AsyncValue()
      value.set('hello')
      value.get(cb)
    }
  },
  {
    name: 'get then set',
    iterations: 100000,
    task(cb) {
      var value = new AsyncValue()
      value.get(cb)
      value.set('hello')
    }
  }
], results => {
  console.log(benchPrettyPrint(results))
})
