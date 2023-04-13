/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

const Stats = require('./stats')

module.exports = function createSystemMetrics (registry) {
  const stats = new Stats()

  registry.registerCollector(stats)

  for (const metric of Object.keys(stats.toJSON())) {
    registry.getOrCreateGauge(metric, () => stats.toJSON()[metric])
  }
}
