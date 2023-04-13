/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

const TYPE = 'messaging'
const SUBTYPE = 'sns'
const ACTION = 'publish'
const PHONE_NUMBER = '<PHONE_NUMBER>'

function getArnOrPhoneNumberFromRequest (request) {
  let arn = request && request.params && request.params.TopicArn
  if (!arn) {
    arn = request && request.params && request.params.TargetArn
  }
  if (!arn) {
    arn = request && request.params && request.params.PhoneNumber
  }
  return arn
}

function getRegionFromRequest (request) {
  return request && request.service &&
        request.service.config && request.service.config.region
}

function getSpanNameFromRequest (request) {
  const topicName = getDestinationNameFromRequest(request)
  return `SNS PUBLISH to ${topicName}`
}

function getMessageContextFromRequest (request) {
  return {
    queue: {
      name: getDestinationNameFromRequest(request)
    }
  }
}

function getAddressFromRequest (request) {
  return request && request.service && request.service.endpoint &&
    request.service.endpoint.hostname
}

function getPortFromRequest (request) {
  return request && request.service && request.service.endpoint &&
    request.service.endpoint.port
}

function getMessageDestinationContextFromRequest (request) {
  return {
    address: getAddressFromRequest(request),
    port: getPortFromRequest(request),
    cloud: {
      region: getRegionFromRequest(request)
    }
  }
}

function getDestinationNameFromRequest (request) {
  const phoneNumber = request && request.params && request.params.PhoneNumber
  if (phoneNumber) {
    return PHONE_NUMBER
  }

  const topicArn = request && request.params && request.params.TopicArn
  const targetArn = request && request.params && request.params.TargetArn

  if (topicArn) {
    const parts = topicArn.split(':')
    const topicName = parts.pop()
    return topicName
  }

  if (targetArn) {
    const fullName = targetArn.split(':').pop()
    if (fullName.lastIndexOf('/') !== -1) {
      return fullName.substring(0, fullName.lastIndexOf('/'))
    } else {
      return fullName
    }
  }
}

function shouldIgnoreRequest (request, agent) {
  if (request.operation !== 'publish') {
    return true
  }

  // is the named topic on our ignore list?
  if (agent._conf && agent._conf.ignoreMessageQueuesRegExp) {
    const queueName = getArnOrPhoneNumberFromRequest(request)
    if (queueName) {
      for (const rule of agent._conf.ignoreMessageQueuesRegExp) {
        if (rule.test(queueName)) {
          return true
        }
      }
    }
  }

  return false
}

function instrumentationSns (orig, origArguments, request, AWS, agent, { version, enabled }) {
  if (shouldIgnoreRequest(request, agent)) {
    return orig.apply(request, origArguments)
  }

  const ins = agent._instrumentation
  const name = getSpanNameFromRequest(request)
  const span = ins.createSpan(name, TYPE, SUBTYPE, ACTION, { exitSpan: true })
  if (!span) {
    return orig.apply(request, origArguments)
  }

  span._setDestinationContext(getMessageDestinationContextFromRequest(request))
  span.setMessageContext(getMessageContextFromRequest(request))

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
  instrumentationSns,

  // exported for testing
  getSpanNameFromRequest,
  getDestinationNameFromRequest,
  getMessageDestinationContextFromRequest,
  getArnOrPhoneNumberFromRequest
}
