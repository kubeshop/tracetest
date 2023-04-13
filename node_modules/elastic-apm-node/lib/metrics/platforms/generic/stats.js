/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

const processCpu = require('./process-cpu')

class Stats {
  constructor () {
    this.stats = {
      'system.process.cpu.total.norm.pct': 0,
      'system.process.cpu.system.norm.pct': 0,
      'system.process.cpu.user.norm.pct': 0
    }
  }

  toJSON () {
    return this.stats
  }

  collect (cb) {
    const cpu = processCpu()
    this.stats['system.process.cpu.total.norm.pct'] = cpu.total
    this.stats['system.process.cpu.system.norm.pct'] = cpu.system
    this.stats['system.process.cpu.user.norm.pct'] = cpu.user

    if (cb) process.nextTick(cb)
  }
}

module.exports = Stats
