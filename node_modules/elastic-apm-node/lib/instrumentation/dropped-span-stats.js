/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

const MAX_DROPPED_SPAN_STATS = 128

class DroppedSpanStats {
  constructor () {
    this.statsMap = new Map()
  }

  /**
   * Record this span in dropped span stats.
   *
   * @param {Span} span
   * @returns {boolean} True iff this span was added to stats. This return value
   *    is only used for testing.
   */
  captureDroppedSpan (span) {
    if (!span) {
      return false
    }

    const serviceTargetType = span._serviceTarget && span._serviceTarget.type
    const serviceTargetName = span._serviceTarget && span._serviceTarget.name
    const resource = span._destination && span._destination.service && span._destination.service.resource
    if (!span._exitSpan || !(serviceTargetType || serviceTargetName) || !resource) {
      return false
    }

    const stats = this._getOrCreateStats(serviceTargetType, serviceTargetName, resource, span.outcome)
    if (!stats) {
      return false
    }
    stats.duration.count++
    stats.duration.sum.us += (span._duration * 1000)
    return true
  }

  _getOrCreateStats (serviceTargetType, serviceTargetName, resource, outcome) {
    const key = [serviceTargetType, serviceTargetName, resource, outcome].join('')
    let stats = this.statsMap.get(key)
    if (stats) {
      return stats
    }

    if (this.statsMap.size >= MAX_DROPPED_SPAN_STATS) {
      return
    }

    stats = {
      duration: {
        count: 0,
        sum: {
          us: 0
        }
      },
      destination_service_resource: resource,
      outcome: outcome
    }
    if (serviceTargetType) {
      stats.service_target_type = serviceTargetType
    }
    if (serviceTargetName) {
      stats.service_target_name = serviceTargetName
    }
    this.statsMap.set(key, stats)
    return stats
  }

  encode () {
    return Array.from(this.statsMap.values())
  }

  size () {
    return this.statsMap.size
  }
}

module.exports = {
  DroppedSpanStats,

  // Exported for testing-only.
  MAX_DROPPED_SPAN_STATS
}
