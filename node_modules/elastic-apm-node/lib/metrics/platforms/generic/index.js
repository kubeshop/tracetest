/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

const os = require('os')

const Stats = require('./stats')

module.exports = function createSystemMetrics (registry) {
  // Base system metrics
  registry.getOrCreateGauge(
    'system.cpu.total.norm.pct',
    require('./system-cpu')
  )
  registry.getOrCreateGauge(
    'system.memory.total',
    () => os.totalmem()
  )
  registry.getOrCreateGauge(
    'system.memory.actual.free',
    () => os.freemem()
  )

  // Process metrics
  const stats = new Stats()
  registry.registerCollector(stats)

  const metrics = [
    'system.process.cpu.total.norm.pct',
    'system.process.cpu.system.norm.pct',
    'system.process.cpu.user.norm.pct'
  ]

  for (const metric of metrics) {
    registry.getOrCreateGauge(metric, () => stats.toJSON()[metric])
  }

  registry.getOrCreateGauge(
    'system.process.memory.rss.bytes',
    () => process.memoryUsage().rss
  )
}
