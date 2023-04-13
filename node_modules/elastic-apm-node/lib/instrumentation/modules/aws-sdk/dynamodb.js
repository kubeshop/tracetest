/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

const TYPE = 'db'
const SUBTYPE = 'dynamodb'
const ACTION = 'query'

function getRegionFromRequest (request) {
  return request && request.service &&
        request.service.config && request.service.config.region
}

function getPortFromRequest (request) {
  return request && request.service &&
        request.service.endpoint && request.service.endpoint.port
}

function getMethodFromRequest (request) {
  const method = request && request.operation
  if (method) {
    return method[0].toUpperCase() + method.slice(1)
  }
}

function getStatementFromRequest (request) {
  const method = getMethodFromRequest(request)
  if (method === 'Query' && request && request.params && request.params.KeyConditionExpression) {
    return request.params.KeyConditionExpression
  }
  return undefined
}

function getAddressFromRequest (request) {
  return request && request.service && request.service.endpoint &&
        request.service.endpoint.hostname
}

function getTableFromRequest (request) {
  const table = request && request.params && request.params.TableName
  if (!table) {
    return ''
  }
  return ` ${table}`
}

// Creates the span name from request information
function getSpanNameFromRequest (request) {
  const method = getMethodFromRequest(request)
  const table = getTableFromRequest(request)
  const name = `DynamoDB ${method}${table}`
  return name
}

function shouldIgnoreRequest (request, agent) {
  return false
}

// Main entrypoint for SQS instrumentation
//
// Must call (or one of its function calls must call) the
// `orig` function/method
function instrumentationDynamoDb (orig, origArguments, request, AWS, agent, { version, enabled }) {
  if (shouldIgnoreRequest(request, agent)) {
    return orig.apply(request, origArguments)
  }

  const ins = agent._instrumentation
  const name = getSpanNameFromRequest(request)
  const span = ins.createSpan(name, TYPE, SUBTYPE, ACTION, { exitSpan: true })
  if (!span) {
    return orig.apply(request, origArguments)
  }

  const region = getRegionFromRequest(request)
  const dbContext = {
    type: SUBTYPE,
    instance: region
  }
  const dbStatement = getStatementFromRequest(request)
  if (dbStatement) {
    dbContext.statement = dbStatement
  }
  span.setDbContext(dbContext)
  span._setDestinationContext({
    address: getAddressFromRequest(request),
    port: getPortFromRequest(request),
    cloud: {
      region
    }
  })

  const onComplete = function (response) {
    if (response && response.error) {
      agent.captureError(response.error)
    }
    span.end()
  }
  // Bind onComplete to the span's run context so that `captureError` picks
  // up the correct currentSpan.
  const parentRunContext = ins.currRunContext()
  const spanRunContext = parentRunContext.enterSpan(span)
  request.on('complete', ins.bindFunctionToRunContext(spanRunContext, onComplete))

  const cb = origArguments[origArguments.length - 1]
  if (typeof cb === 'function') {
    origArguments[origArguments.length - 1] = ins.bindFunctionToRunContext(parentRunContext, cb)
  }
  return ins.withRunContext(spanRunContext, orig, request, ...origArguments)
}

module.exports = {
  instrumentationDynamoDb,

  // exported for testing
  getRegionFromRequest,
  getPortFromRequest,
  getStatementFromRequest,
  getAddressFromRequest,
  getMethodFromRequest
}
