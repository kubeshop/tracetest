/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

const oblog = require('./oblog')

class OTelTracerProvider {
  // @param {OTelTracer} tracer
  constructor (tracer) {
    this._tracer = tracer
  }

  getTracer (_name, _version, _options) {
    oblog.apicall('OTelTracerProvider.getTracer(...)')
    return this._tracer
  }
}

module.exports = {
  OTelTracerProvider
}
