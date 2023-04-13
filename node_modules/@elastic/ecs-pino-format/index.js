// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

'use strict'

const {
  version,
  formatError,
  formatHttpRequest,
  formatHttpResponse
} = require('@elastic/ecs-helpers')

const { hasOwnProperty } = Object.prototype
let triedElasticApmImport = false
let elasticApm = null

// Create options for `pino(...)` that configure it for ecs-logging output.
//
// @param {Object} opts - Optional.
//    - {Boolean} opts.convertErr - Whether to convert a logged `err` field
//      to ECS error fields. Default true, to match Pino's default of having
//      an `err` serializer.
//    - {Boolean} opts.convertReqRes - Whether to convert logged `req` and `res`
//      HTTP request and response fields to ECS HTTP, User agent, and URL
//      fields. Default false.
//    - {Boolean} opts.apmIntegration - Whether to automatically integrate with
//      Elastic APM (https://github.com/elastic/apm-agent-nodejs). If a started
//      APM agent is detected, then log records will include the following
//      fields:
//        - "service.name" - the configured serviceName in the agent
//        - "event.dataset" - set to "$serviceName.log" for correlation in Kibana
//        - "trace.id", "transaction.id", and "span.id" - if there is a current
//          active trace when the log call is made
//      Default true.
function createEcsPinoOptions (opts) {
  let convertErr = true
  let convertReqRes = false
  let apmIntegration = true
  if (opts) {
    if (hasOwnProperty.call(opts, 'convertErr')) {
      convertErr = opts.convertErr
    }
    if (hasOwnProperty.call(opts, 'convertReqRes')) {
      convertReqRes = opts.convertReqRes
    }
    if (hasOwnProperty.call(opts, 'apmIntegration')) {
      apmIntegration = opts.apmIntegration
    }
  }

  let apm = null
  let apmServiceName = null
  if (apmIntegration) {
    // istanbul ignore if
    if (opts && opts._elasticApm) {
      // `opts._elasticApm` is an internal/testing-only option to be used
      // for testing in the APM agent where the import is a local path
      // rather than "elastic-apm-node".
      elasticApm = opts._elasticApm
    } else if (!triedElasticApmImport) {
      triedElasticApmImport = true
      // We lazily require this module here instead of at the top-level to
      // avoid a possible circular-require if the user code does
      // `require('@elastic/ecs-pino-format')` and has a "node_modules/"
      // where 'elastic-apm-node' shares the same ecs-pino-format install.
      try {
        elasticApm = require('elastic-apm-node')
      } catch (ex) {
        // Silently ignore.
      }
    }
    if (elasticApm && elasticApm.isStarted && elasticApm.isStarted()) {
      apm = elasticApm
      // Elastic APM v3.11.0 added getServiceName(). Fallback to private `apm._conf`.
      // istanbul ignore next
      apmServiceName = apm.getServiceName
        ? apm.getServiceName()
        : apm._conf.serviceName
    }
  }

  let isServiceNameInBindings = false
  let isEventDatasetInBindings = false

  const ecsPinoOptions = {
    messageKey: 'message',
    timestamp: () => `,"@timestamp":"${new Date().toISOString()}"`,
    formatters: {
      level (label, number) {
        return { 'log.level': label }
      },

      bindings (bindings) {
        const {
          // `pid` and `hostname` are default bindings, unless overriden by
          // a `base: {...}` passed to logger creation.
          pid,
          hostname,
          // name is defined if `log = pino({name: 'my name', ...})`
          name,
          // Warning: silently drop any "ecs" value from `base`. See
          // "ecs.version" comment below.
          ecs,
          ...ecsBindings
        } = bindings

        if (pid !== undefined) {
          // https://www.elastic.co/guide/en/ecs/current/ecs-process.html#field-process-pid
          ecsBindings.process = { pid: pid }
        }
        if (hostname !== undefined) {
          // https://www.elastic.co/guide/en/ecs/current/ecs-host.html#field-host-hostname
          ecsBindings.host = { hostname: hostname }
        }
        if (name !== undefined) {
          // https://www.elastic.co/guide/en/ecs/current/ecs-log.html#field-log-logger
          ecsBindings.log = { logger: name }
        }

        // Note if service.name & event.dataset are set, to not do so again below.
        if (bindings.service && bindings.service.name) {
          isServiceNameInBindings = true
        }
        if (bindings.event && bindings.event.dataset) {
          isEventDatasetInBindings = true
        }

        return ecsBindings
      },

      log (obj) {
        const {
          req,
          res,
          err,
          ...ecsObj
        } = obj

        // https://www.elastic.co/guide/en/ecs/current/ecs-ecs.html
        // For "ecs.version" we take a heavier-handed approach, because it is
        // a require ecs-logging field: overwrite any possible "ecs" value from
        // the log statement. This means we don't need to spend the time
        // guarding against "ecs" being null, Array, Buffer, Date, etc.
        ecsObj.ecs = { version }

        if (apm) {
          // A mis-configured APM Agent can be "started" but not have a
          // "serviceName".
          if (apmServiceName) {
            // Per https://github.com/elastic/ecs-logging/blob/master/spec/spec.json
            // "service.name" and "event.dataset" should be automatically set
            // if not already by the user.
            if (!isServiceNameInBindings) {
              const service = ecsObj.service
              if (service === undefined) {
                ecsObj.service = { name: apmServiceName }
              } else if (!isVanillaObject(service)) {
                // Warning: "service" type conflicts with ECS spec. Overwriting.
                ecsObj.service = { name: apmServiceName }
              } else if (typeof service.name !== 'string') {
                ecsObj.service.name = apmServiceName
              }
            }
            if (!isEventDatasetInBindings) {
              const event = ecsObj.event
              if (event === undefined) {
                ecsObj.event = { dataset: apmServiceName + '.log' }
              } else if (!isVanillaObject(event)) {
                // Warning: "event" type conflicts with ECS spec. Overwriting.
                ecsObj.event = { dataset: apmServiceName + '.log' }
              } else if (typeof event.dataset !== 'string') {
                ecsObj.event.dataset = apmServiceName + '.log'
              }
            }
          }

          // https://www.elastic.co/guide/en/ecs/current/ecs-tracing.html
          const tx = apm.currentTransaction
          if (tx) {
            ecsObj.trace = ecsObj.trace || {}
            ecsObj.trace.id = tx.traceId
            ecsObj.transaction = ecsObj.transaction || {}
            ecsObj.transaction.id = tx.id
            const span = apm.currentSpan
            // istanbul ignore else
            if (span) {
              ecsObj.span = ecsObj.span || {}
              ecsObj.span.id = span.id
            }
          }
        }

        // https://www.elastic.co/guide/en/ecs/current/ecs-http.html
        if (err !== undefined) {
          if (!convertErr) {
            ecsObj.err = err
          } else {
            formatError(ecsObj, err)
          }
        }

        // https://www.elastic.co/guide/en/ecs/current/ecs-http.html
        if (req !== undefined) {
          if (!convertReqRes) {
            ecsObj.req = req
          } else {
            formatHttpRequest(ecsObj, req)
          }
        }
        if (res !== undefined) {
          if (!convertReqRes) {
            ecsObj.res = res
          } else {
            formatHttpResponse(ecsObj, res)
          }
        }

        return ecsObj
      }
    }
  }

  return ecsPinoOptions
}

// Return true if the given arg is a "vanilla" object. Roughly the intent is
// whether this is basic mapping of string keys to values that will serialize
// as a JSON object.
//
// Currently, it excludes Map. The uses above don't really expect a user to:
//     service = new Map([["foo", "bar"]])
//     log.info({ service }, '...')
//
// There are many ways tackle this. See some attempts and benchmarks at:
// https://gist.github.com/trentm/34131a92eede80fd2109f8febaa56f5a
function isVanillaObject (o) {
  return (typeof o === 'object' &&
    (!o.constructor || o.constructor.name === 'Object'))
}

module.exports = createEcsPinoOptions
