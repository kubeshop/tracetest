'use strict'

var tap = require('tap')
var AsyncValue = require('./')

tap.test('set in constructor', t => {
  var value = new AsyncValue('hello')

  value.get(value => {
    t.equal(value, 'hello')
    t.end()
  })
})

tap.test('get then sync set', t => {
  var value = new AsyncValue()

  value.get(value => {
    t.equal(value, 'hello')
    t.end()
  })

  value.set('hello')
})

tap.test('get then async set', t => {
  var value = new AsyncValue()

  value.get(value => {
    t.equal(value, 'hello')
    t.end()
  })

  setImmediate(() => {
    value.set('hello')
  })
})

tap.test('throw on non-function arguments to get', t => {
  var value = new AsyncValue()
  var types = [
    1,
    null,
    undefined,
    'string',
    /regex/,
    {},
    []
  ]
  types.forEach(type => {
    t.throws(() => value.get(type), {
      message: 'callback must be a function'
    })
  })
  t.end()
})

tap.test('throw on multiple calls to set', t => {
  var value = new AsyncValue()
  value.set('first')
  t.throws(() => value.set('second'), {
    message: 'value can only be set once'
  })
  t.end()
})

tap.test('send', t => {
  var a = new AsyncValue()
  var b = new AsyncValue()

  b.get(value => {
    t.equal(value, 'hello')
    t.end()
  })

  a.set('hello')
  a.send(b)
})

tap.test('set sends async values', t => {
  var a = new AsyncValue()
  var b = new AsyncValue()

  b.get(value => {
    t.equal(value, 'hello')
    t.end()
  })

  a.set('hello')
  b.set(a)
})

tap.test('throw on send to non-async-value', t => {
  var value = new AsyncValue()
  var types = [
    1,
    null,
    undefined,
    () => {},
    'string',
    /regex/,
    {},
    []
  ]
  types.forEach(type => {
    t.throws(() => value.send(type), {
      message: 'send target must be an async value'
    })
  })
  t.end()
})
