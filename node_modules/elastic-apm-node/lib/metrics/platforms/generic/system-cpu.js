/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

const os = require('os')

function cpuAverage () {
  const times = {
    user: 0,
    nice: 0,
    sys: 0,
    idle: 0,
    irq: 0,
    total: 0
  }

  const cpus = os.cpus()
  for (const cpu of cpus) {
    for (const type of Object.keys(cpu.times)) {
      times[type] += cpu.times[type]
      times.total += cpu.times[type]
    }
  }

  // Average over CPU count
  const averages = {}
  for (const type of Object.keys(times)) {
    averages[type] = times[type] / cpus.length
  }

  return averages
}

function cpuPercent (last, next) {
  const idle = next.idle - last.idle
  const total = next.total - last.total
  return 1 - idle / total || 0
}

let last = cpuAverage()

module.exports = function systemCPUUsage () {
  const next = cpuAverage()
  const result = cpuPercent(last, next)
  last = next
  return result
}
