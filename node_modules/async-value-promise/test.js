'use strict'

var tap = require('tap')
var AsyncValuePromise = require('./')

tap.test('resolve', t => {
  var promise = new AsyncValuePromise()

  promise.then(value => {
    t.equal(value, 'hello')
    t.end()
  })

  setImmediate(() => {
    promise.resolve('hello')
  })
})

tap.test('reject', t => {
  var promise = new AsyncValuePromise()

  promise.then(null, error => {
    t.equal(error, 'hello')
    t.end()
  })

  setImmediate(() => {
    promise.reject('hello')
  })
})

tap.test('throw on reject after resolve', t => {
  var promise = new AsyncValuePromise()
  promise.resolve('pass')
  t.throws(() => promise.reject('fail'))
  t.end()
})

tap.test('throw on resolve after reject', t => {
  var promise = new AsyncValuePromise()
  promise.reject('fail')
  t.throws(() => promise.resolve('pass'))
  t.end()
})

tap.test('catch', t => {
  var promise = new AsyncValuePromise()

  promise.catch(error => {
    t.equal(error, 'hello')
    t.end()
  })

  setImmediate(() => {
    promise.reject('hello')
  })
})

tap.test('all success', t => {
  var promiseA = new AsyncValuePromise()
  var promiseB = new AsyncValuePromise()
  var promise = AsyncValuePromise.all([
    promiseA,
    promiseB
  ])

  promise
    .then((values) => {
      t.equal(values[0], 'first')
      t.equal(values[1], 'second')
      t.end()
    })
    .catch(() => {
      t.fail('should not reject')
    })

  setImmediate(() => {
    promiseB.resolve('second')
    setImmediate(() => {
      promiseA.resolve('first')
    })
  })
})

tap.test('all failure', t => {
  var promiseA = new AsyncValuePromise()
  var promiseB = new AsyncValuePromise()
  var promise = AsyncValuePromise.all([ promiseA, promiseB ])

  promise
    .then(() => t.fail('should not resolve'))
    .catch(error => {
      t.equal(error, 'error')
      t.end()
    })

  setImmediate(() => {
    promiseA.resolve('success')
    promiseB.reject('error')
  })
})

tap.test('all empty', t => {
  var promise = AsyncValuePromise.all([])

  promise
    .then((values) => {
      t.equal(values.length, 0)
      t.end()
    })
    .catch(() => {
      t.fail('should not reject')
    })
})

tap.test('static resolve', t => {
  var promise = AsyncValuePromise.resolve('hello')

  promise
    .then(value => {
      t.equal(value, 'hello')
      t.end()
    })
    .catch(() => t.fail('should not reject'))
})

tap.test('static reject', t => {
  var promise = AsyncValuePromise.reject('hello')

  promise
    .then(() => t.fail('should not resolve'))
    .catch(error => {
      t.equal(error, 'hello')
      t.end()
    })
})

tap.test('catch', t => {
  var promise = AsyncValuePromise.reject('world')

  promise
    .catch(name => 'hello ' + name)
    .then(value => {
      t.equal(value, 'hello world')
      t.end()
    })
    .catch(() => t.fail('should not reject'))
})
