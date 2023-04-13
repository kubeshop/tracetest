/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'
class QueueMetricsCollector {
  constructor () {
    this.stats = {
      'queue.latency.min.ms': 0,
      'queue.latency.max.ms': 0,
      'queue.latency.avg.ms': 0
    }

    this.total = 0
    this.count = 0
    this.min = 0
    this.max = 0
  }

  // Updates the data used to generate our stats
  //
  // Unlike the Stats and RuntimeCollector, this function
  // also returns the instantiated collector.
  //
  // @param {number} time A javascript ms timestamp
  updateStats (time) {
    if (this.min === 0 || this.min > time) {
      this.min = time
    }
    if (this.max < time) {
      this.max = time
    }

    this.count++
    this.total += time
  }

  // standard `collect` method for stats collectors
  //
  // called by the MetricReporter object prior to its sending data
  //
  // @param {function} cb callback function
  collect (cb) {
    // update average based on count and total
    this.stats['queue.latency.avg.ms'] = 0
    if (this.count > 0) {
      this.stats['queue.latency.avg.ms'] = this.total / this.count
    }

    this.stats['queue.latency.max.ms'] = this.max
    this.stats['queue.latency.min.ms'] = this.min

    // reset for next run
    this.total = 0
    this.count = 0
    this.max = 0
    this.min = 0
    if (cb) process.nextTick(cb)
  }
}

// Creates and Registers a Metric Collector
//
// Unlike the Stats and RuntimeCollector, this function
// also returns the instantiated collector.
//
// @param {string} queueOrTopicName
// @param {MetricRegistry} registry
// @returns QueueMetricsCollector
function createQueueMetrics (queueOrTopicName, registry) {
  const collector = new QueueMetricsCollector()
  registry.registerCollector(collector)
  for (const metric of Object.keys(collector.stats)) {
    registry.getOrCreateGauge(
      metric,
      function returnCurrentValue () {
        return collector.stats[metric]
      },
      { queue_name: queueOrTopicName }
    )
  }
  return collector
}

module.exports = {
  createQueueMetrics,
  QueueMetricsCollector
}
