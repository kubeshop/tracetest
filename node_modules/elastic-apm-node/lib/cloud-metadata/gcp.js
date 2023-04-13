/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'
const URL = require('url').URL
const { httpRequest } = require('../http-request')
const DEFAULT_BASE_URL = new URL('/', 'http://metadata.google.internal:80')
/**
 * Checks for metadata server then fetches data
 *
 * The getMetadataGcp function will fetch cloud metadata information
 * from Amazon's IMDSv1 endpoint and return (via callback)
 * the formatted metadata.
 *
 * Before fetching data, the server will be "pinged" by attempting
 * to connect via TCP with a short timeout. (`connectTimeoutMs`)
 *
 * https://cloud.google.com/compute/docs/storing-retrieving-metadata
 */
function getMetadataGcp (connectTimeoutMs, resTimeoutMs, logger, baseUrlOverride, cb) {
  const baseUrl = baseUrlOverride || DEFAULT_BASE_URL
  const options = {
    method: 'GET',
    timeout: resTimeoutMs,
    connectTimeout: connectTimeoutMs,
    headers: {
      'Metadata-Flavor': 'Google'
    }
  }
  const url = baseUrl + 'computeMetadata/v1/?recursive=true'
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
          logger.trace('gcp metadata server responded, but there was an ' +
            'error parsing the result: %o', err)
          cb(err)
          return
        }
        cb(null, result)
      })
    }
  )

  req.on('timeout', function () {
    req.destroy(new Error('request to metadata server timed out'))
  })

  req.on('connectTimeout', function () {
    req.destroy(new Error('could not ping metadata server'))
  })

  req.on('error', function (err) {
    cb(err)
  })

  req.end()
}

/**
 * Builds metadata object
 *
 * Takes the response from a /computeMetadata/v1/?recursive=true
 * service request and formats it into the cloud metadata object
 */
function formatMetadataStringIntoObject (string) {
  const data = JSON.parse(string)
  // cast string manipulation fields as strings "just in case"
  if (data.instance) {
    data.instance.machineType = String(data.instance.machineType)
    data.instance.zone = String(data.instance.zone)
  }

  const metadata = {
    availability_zone: null,
    region: null,
    instance: {
      id: null
    },
    machine: {
      type: null
    },
    provider: null,
    project: {
      id: null,
      name: null
    }
  }

  metadata.availability_zone = null
  metadata.region = null
  if (data.instance && data.instance.zone) {
    // `projects/513326162531/zones/us-west1-b` manipuated into
    // `us-west1-b`, and then `us-west1`
    const regionWithZone = data.instance.zone.split('/').pop()
    const parts = regionWithZone.split('-')
    parts.pop()
    metadata.region = parts.join('-')
    metadata.availability_zone = regionWithZone
  }

  if (data.instance) {
    metadata.instance = {
      id: String(data.instance.id)
    }

    metadata.machine = {
      type: String((data.instance.machineType.split('/').pop()))
    }
  } else {
    metadata.instance = {
      id: null
    }

    metadata.machine = {
      type: null
    }
  }

  metadata.provider = 'gcp'

  if (data.project) {
    metadata.project = {
      id: String(data.project.numericProjectId),
      name: String(data.project.projectId)
    }
  } else {
    metadata.project = {
      id: null,
      name: null
    }
  }
  return metadata
}

module.exports = { getMetadataGcp }
