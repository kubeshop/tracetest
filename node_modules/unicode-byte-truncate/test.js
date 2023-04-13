'use strict'

var test = require('tape')
var trunc = require('./')

test('normal string', function (t) {
  var s = 'foobar'
  t.equal(trunc(s, 0), '')
  t.equal(trunc(s, 3), 'foo')
  t.equal(trunc(s, 6), 'foobar')
  t.equal(trunc(s, 9), 'foobar')
  t.end()
})

test('multibyte string', function (t) {
  var s = 'foo🎉bar'
  t.equal(trunc(s, 3), 'foo')
  t.equal(trunc(s, 4), 'foo')
  t.equal(trunc(s, 5), 'foo')
  t.equal(trunc(s, 6), 'foo')
  t.equal(trunc(s, 7), 'foo🎉')
  t.equal(trunc(s, 8), 'foo🎉b')
  t.end()
})

test('invalid values', function (t) {
  t.throws(function () { trunc() })
  t.throws(function () { trunc(1, 0) })
  t.throws(function () { trunc('') })
  t.throws(function () { trunc('', '') })
  t.throws(function () { trunc('', NaN) })
  t.throws(function () { trunc('', false) })
  t.throws(function () { trunc('', -1) })
  t.throws(function () { trunc('', Infinity) })
  t.end()
})
