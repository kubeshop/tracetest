/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

const { stringify } = require('querystring')

class Ids {
  toString () {
    return stringify(this, ' ', '=')
  }
}

class SpanIds extends Ids {
  constructor (span) {
    super()
    this['trace.id'] = span.traceId
    this['span.id'] = span.id
    Object.freeze(this)
  }
}

class TransactionIds extends Ids {
  constructor (transaction) {
    super()
    this['trace.id'] = transaction.traceId
    this['transaction.id'] = transaction.id
    Object.freeze(this)
  }
}

module.exports = {
  Ids,
  SpanIds,
  TransactionIds
}
