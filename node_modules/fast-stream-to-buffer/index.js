'use strict'

const eos = require('end-of-stream')

module.exports = streamToBuffer

streamToBuffer.onStream = onStream

function streamToBuffer (stream, cb) {
  const buffers = []

  stream.on('data', buffers.push.bind(buffers))

  eos(stream, function (err) {
    switch (buffers.length) {
      case 0:
        cb(err, Buffer.allocUnsafe(0), stream)
        break
      case 1:
        cb(err, buffers[0], stream)
        break
      default:
        cb(err, Buffer.concat(buffers), stream)
    }
  })
}

function onStream (cb) {
  return function (stream) {
    streamToBuffer(stream, cb)
  }
}
