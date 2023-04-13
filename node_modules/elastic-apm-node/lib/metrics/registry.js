/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

const os = require('os')

const { SelfReportingMetricsRegistry } = require('measured-reporting')
const DimensionAwareMetricsRegistry = require('measured-reporting/lib/registries/DimensionAwareMetricsRegistry')

const MetricsReporter = require('./reporter')
const createRuntimeMetrics = require('./runtime')
const createSystemMetrics = process.platform === 'linux'
  ? require('./platforms/linux')
  : require('./platforms/generic')

class MetricsRegistry extends SelfReportingMetricsRegistry {
  constructor (agent, { reporterOptions, registryOptions = {} } = {}) {
    const defaultReporterOptions = {
      defaultDimensions: {
        hostname: agent._conf.hostname || os.hostname(),
        env: agent._conf.environment || ''
      }
    }

    const options = Object.assign({}, defaultReporterOptions, reporterOptions)
    const reporter = new MetricsReporter(agent, options)

    registryOptions.registry = new DimensionAwareMetricsRegistry({
      metricLimit: agent._conf.metricsLimit
    })

    super(reporter, registryOptions)
    this._agent = agent

    this._registry.collectors = []
    if (reporter.enabled) {
      createSystemMetrics(this)
      createRuntimeMetrics(this)
    }
  }

  registerCollector (collector) {
    this._registry.collectors.push(collector)
  }
}

module.exports = MetricsRegistry
