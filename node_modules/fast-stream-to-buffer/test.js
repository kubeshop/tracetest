'use strict'

const PassThrough = require('stream').PassThrough
const test = require('tape')
const streamToBuffer = require('./')

test('streamToBuffer() - no buffers', function (t) {
  const stream = new PassThrough()
  streamToBuffer(stream, function (err, buf) {
    t.error(err)
    t.ok(Buffer.isBuffer(buf))
    t.equal(buf.length, 0)
    t.end()
  })

  stream.end()
})

test('streamToBuffer() - single buffer', function (t) {
  const stream = new PassThrough()
  streamToBuffer(stream, function (err, buf) {
    t.error(err)
    t.ok(Buffer.isBuffer(buf))
    t.equal(buf.length, 11)
    t.equal(buf.toString(), 'hello world')
    t.end()
  })

  stream.end('hello world')
})

test('streamToBuffer() - multiple buffers', function (t) {
  const stream = new PassThrough()
  streamToBuffer(stream, function (err, buf) {
    t.error(err)
    t.ok(Buffer.isBuffer(buf))
    t.equal(buf.length, 11)
    t.equal(buf.toString(), 'hello world')
    t.end()
  })

  stream.write('hello')
  stream.end(' world')
})

test('streamToBuffer() - error', function (t) {
  const stream = new PassThrough()
  streamToBuffer(stream, function (err, buf) {
    t.ok(err)
    t.equal(err.message, 'foo')
    t.end()
  })

  stream.emit('error', new Error('foo'))
})

test('streamToBuffer.onStream() - no buffers', function (t) {
  const stream = new PassThrough()
  const onStream = streamToBuffer.onStream(function (err, buf) {
    t.error(err)
    t.ok(Buffer.isBuffer(buf))
    t.equal(buf.length, 0)
    t.end()
  })

  onStream(stream)

  stream.end()
})

test('streamToBuffer.onStream() - single buffer', function (t) {
  const stream = new PassThrough()
  const onStream = streamToBuffer.onStream(function (err, buf) {
    t.error(err)
    t.ok(Buffer.isBuffer(buf))
    t.equal(buf.length, 11)
    t.equal(buf.toString(), 'hello world')
    t.end()
  })

  onStream(stream)

  stream.end('hello world')
})

test('streamToBuffer.onStream() - multiple buffers', function (t) {
  const stream = new PassThrough()
  const onStream = streamToBuffer.onStream(function (err, buf) {
    t.error(err)
    t.ok(Buffer.isBuffer(buf))
    t.equal(buf.length, 11)
    t.equal(buf.toString(), 'hello world')
    t.end()
  })

  onStream(stream)

  stream.write('hello')
  stream.end(' world')
})

test('streamToBuffer.onStream() - error', function (t) {
  const stream = new PassThrough()
  const onStream = streamToBuffer.onStream(function (err, buf) {
    t.ok(err)
    t.equal(err.message, 'foo')
    t.end()
  })

  onStream(stream)

  stream.emit('error', new Error('foo'))
})
