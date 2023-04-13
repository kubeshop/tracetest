'use strict'

var isInteger = require('is-integer')
var slice = require('unicode-substring')

module.exports = function (str, len) {
  if (typeof str !== 'string') throw new Error('Expected first argument to be a string')
  if (!isInteger(len) || len < 0) throw new Error('Expected second argument be an integer greater than or equal to 0')

  var origLen = len
  while (Buffer.byteLength(str) > origLen) {
    str = slice(str, 0, len--)
  }

  return str
}
