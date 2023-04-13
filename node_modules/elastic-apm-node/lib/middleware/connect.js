/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

module.exports = function connectMiddleware () {
  var agent = this
  return function (err, req, res, next) {
    agent.captureError(err, { request: req }, function elasticAPMMiddleware () {
      next(err)
    })
  }
}
