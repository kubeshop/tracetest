/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'
const semver = require('semver')
const shimmer = require('../shimmer')
const { instrumentationS3 } = require('./aws-sdk/s3')
const { instrumentationSqs } = require('./aws-sdk/sqs')
const { instrumentationDynamoDb } = require('./aws-sdk/dynamodb.js')
const { instrumentationSns } = require('./aws-sdk/sns.js')

const instrumentorFromSvcId = {
  s3: instrumentationS3,
  sqs: instrumentationSqs,
  dynamodb: instrumentationDynamoDb,
  sns: instrumentationSns
}

// Called in place of AWS.Request.send and AWS.Request.promise
//
// Determines which amazon service an API request is for
// and then passes call on to an appropriate instrumentation
// function.
function instrumentOperation (orig, origArguments, request, AWS, agent, { version, enabled }) {
  const instrumentor = instrumentorFromSvcId[request.service.serviceIdentifier]
  if (instrumentor) {
    return instrumentor(orig, origArguments, request, AWS, agent, { version, enabled })
  }

  // if we're still here, then we still need to call the original method
  return orig.apply(request, origArguments)
}

// main entry point for aws-sdk instrumentation
module.exports = function (AWS, agent, { version, enabled }) {
  if (!enabled) return AWS
  if (!semver.satisfies(version, '>1 <3')) {
    agent.logger.debug('aws-sdk version %s not supported - aborting...', version)
    return AWS
  }

  shimmer.wrap(AWS.Request.prototype, 'send', function (orig) {
    return function _wrappedAWSRequestSend () {
      return instrumentOperation(orig, arguments, this, AWS, agent, { version, enabled })
    }
  })

  shimmer.wrap(AWS.Request.prototype, 'promise', function (orig) {
    return function _wrappedAWSRequestPromise () {
      return instrumentOperation(orig, arguments, this, AWS, agent, { version, enabled })
    }
  })
  return AWS
}
