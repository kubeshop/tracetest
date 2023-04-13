'use strict'

const stringify = require('fast-safe-stringify')

exports.serialize = function serialize (obj) {
  const str = tryJSONStringify(obj) || stringify(obj)
  return str + '\n'
}

function tryJSONStringify (obj) {
  try {
    return JSON.stringify(obj)
  } catch (e) {}
}
