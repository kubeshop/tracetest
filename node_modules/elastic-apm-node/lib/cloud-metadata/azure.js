/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'
const URL = require('url').URL
const { httpRequest } = require('../http-request')

const DEFAULT_BASE_URL = new URL('/', 'http://169.254.169.254:80')

/**
 * Checks for metadata server then fetches data
 *
 * The getMetadataAzure method will fetch cloud metadata information
 * from Amazon's IMDSv1 endpoint and return (via callback)
 * the formatted metadata.
 *
 * Before fetching data, the server will be "pinged" by attempting
 * to connect via TCP with a short timeout (`connectTimeoutMs`).
 *
 * https://docs.microsoft.com/en-us/azure/virtual-machines/windows/instance-metadata-service?tabs=windows
 */
function getMetadataAzure (connectTimeoutMs, resTimeoutMs, logger, baseUrlOverride, cb) {
  const baseUrl = baseUrlOverride || DEFAULT_BASE_URL
  const options = {
    method: 'GET',
    timeout: resTimeoutMs,
    connectTimeout: connectTimeoutMs,
    headers: {
      Metadata: 'true'
    }
  }

  const url = baseUrl.toString() + 'metadata/instance?api-version=2020-09-01'

  const req = httpRequest(
    url,
    options,
    function (res) {
      const finalData = []
      res.on('data', function (data) {
        finalData.push(data)
      })

      res.on('end', function (data) {
        let result
        try {
          result = formatMetadataStringIntoObject(finalData.join(''))
        } catch (err) {
          logger.trace('azure metadata server responded, but there was an ' +
            'error parsing the result: %o', err)
          cb(err)
          return
        }
        cb(null, result)
      })
    }
  )

  req.on('timeout', function () {
    req.destroy(new Error('request to azure metadata server timed out'))
  })

  req.on('connectTimeout', function () {
    req.destroy(new Error('could not ping azure metadata server'))
  })

  req.on('error', function (err) {
    cb(err)
  })

  req.end()
}

/**
 * Builds metadata object
 *
 * Takes the response from /metadata/instance?api-version=2020-09-01
 * service request and formats it into the cloud metadata object
 */
function formatMetadataStringIntoObject (string) {
  const metadata = {
    account: {
      id: null
    },
    instance: {
      id: null,
      name: null
    },
    project: {
      name: null
    },
    availability_zone: null,
    machine: {
      type: null
    },
    provider: 'azure',
    region: null
  }
  const parsed = JSON.parse(string)
  if (!parsed.compute) {
    return metadata
  }
  const data = parsed.compute
  metadata.account.id = String(data.subscriptionId)
  metadata.instance.id = String(data.vmId)
  metadata.instance.name = String(data.name)
  metadata.project.name = String(data.resourceGroupName)
  metadata.availability_zone = String(data.zone)
  metadata.machine.type = String(data.vmSize)
  metadata.region = String(data.location)

  return metadata
}

module.exports = { getMetadataAzure }
