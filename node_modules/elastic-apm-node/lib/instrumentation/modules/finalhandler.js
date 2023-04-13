/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

// Instrumentation of https://www.npmjs.com/package/finalhandler
// This will `apm.captureError` any error passed to the finalhandler.
//
// Both 'express' and 'connect' use finalhandler, so this effectively handles
// capturing errors from incoming request handlers for Express and Connect apps.
// The `errorReportedSymbol` symbol is used to coordinate with 'express'
// instrumentation to avoid double-reporting errors.

var isError = require('core-util-is').isError

var symbols = require('../../symbols')

function shouldReport (err) {
  if (typeof err === 'string') return true
  if (isError(err) && !err[symbols.errorReportedSymbol]) {
    err[symbols.errorReportedSymbol] = true
    return true
  }
  return false
}

module.exports = function (finalhandler, agent) {
  return function wrappedFinalhandler (req, res, options) {
    var final = finalhandler.apply(this, arguments)

    return function (err) {
      if (shouldReport(err)) {
        agent.captureError(err, { request: req })
      }
      return final.apply(this, arguments)
    }
  }
}
