/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'
const { getMetadataAws } = require('./aws')
const { getMetadataAzure } = require('./azure')
const { getMetadataGcp } = require('./gcp')
const { CallbackCoordination } = require('./callback-coordination')

const logging = require('../logging')

// we "ping" (in the colloquial sense, not the ICMP sense) the metadata
// servers by having a listener that expects the underlying socket to
// connect.  CONNECT_TIMEOUT_MS control that time
const CONNECT_TIMEOUT_MS = 100

// some of the metadata servers have a dedicated domain name.  This
// value is added to the above socket ping to allow extra time for
// the DNS to resolve
const DNS_TIMEOUT_MS = 100

// the metadata servers are HTTP services.  This value is
// used as the timeout for the HTTP request that's
// made.
const HTTP_TIMEOUT_MS = 1000

// timeout for the CallbackCoordination object -- this is a fallback to
// account for a catastrophic error in the CallbackCoordination object
const COORDINATION_TIMEOUT_MS = 3000

class CloudMetadata {
  constructor (cloudProvider, logger, serviceName) {
    this.cloudProvider = cloudProvider
    this.logger = logger
    this.serviceName = serviceName
  }

  /**
   * Fetches Cloud Metadata
   *
   * The module's main entry-point/primary method.  The getCloudMetadata function
   * will query the cloud metadata servers and return (via a callback function)
   * to final metadata object.  This function may be called with a single
   * argument.
   *
   * getCloudMetadata(function(err, metadata){
   *     //...
   * })
   *
   * Or with two
   *
   * getCloudMetadata(config, function(err, metadata){
   *     //...
   * })
   *
   * The config parameter is an object that contains information on the
   * metadata servers.
   */
  getCloudMetadata (config, cb) {
    // normalize arguments
    if (!cb) {
      cb = config
      config = {}
    }

    // fill in blanks if any expected keys are missing
    config.aws = config.aws ? config.aws : null
    config.azure = config.azure ? config.azure : null
    config.gcp = config.gcp ? config.gcp : null

    const fetcher = new CallbackCoordination(COORDINATION_TIMEOUT_MS, this.logger)

    if (this.shouldFetchGcp()) {
      fetcher.schedule((fetcher) => {
        const url = config.gcp
        getMetadataGcp(
          CONNECT_TIMEOUT_MS + DNS_TIMEOUT_MS,
          HTTP_TIMEOUT_MS,
          this.logger,
          url,
          (err, result) => {
            fetcher.recordResult(err, result)
          }
        )
      })
    }

    if (this.shouldFetchAws()) {
      fetcher.schedule((fetcher) => {
        const url = config.aws
        getMetadataAws(
          CONNECT_TIMEOUT_MS,
          HTTP_TIMEOUT_MS,
          this.logger,
          url,
          function (err, result) {
            fetcher.recordResult(err, result)
          }
        )
      })
    }

    if (this.shouldFetchAzure()) {
      fetcher.schedule((fetcher) => {
        const url = config.azure
        getMetadataAzure(
          CONNECT_TIMEOUT_MS,
          HTTP_TIMEOUT_MS,
          this.logger,
          url,
          function (err, result) {
            fetcher.recordResult(err, result)
          }
        )
      })
    }

    fetcher.on('result', function (result) {
      cb(null, result)
    })

    fetcher.on('error', function (err) {
      cb(err)
    })

    fetcher.start()
  }

  shouldFetchGcp () {
    return this.cloudProvider === 'auto' || this.cloudProvider === 'gcp'
  }

  shouldFetchAzure () {
    return this.cloudProvider === 'auto' || this.cloudProvider === 'azure'
  }

  shouldFetchAws () {
    return this.cloudProvider === 'auto' || this.cloudProvider === 'aws'
  }
}

/**
 * Simple Command Line interface to fetch metadata
 *
 * $ node lib/cloud-metadata/index.js
 *
 * Will output metadata object or error if no servers are reachable
 */
function main (args) {
  const cloudMetadata = new CloudMetadata('auto', logging.createLogger('off'))
  cloudMetadata.getCloudMetadata(function (err, metadata) {
    if (err) {
      console.log('could not fetch metadata, see error below')
      console.log(err)
      process.exit(1)
    } else {
      console.log('fetched the following metadata')
      console.log(metadata)
      process.exit(0)
    }
  })
}

if (require.main === module) {
  main(process.argv)
}
module.exports = {
  CloudMetadata
}
