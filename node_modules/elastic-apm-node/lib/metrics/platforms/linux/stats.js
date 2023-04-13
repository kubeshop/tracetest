/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

const fs = require('fs')

const afterAll = require('after-all-results')

const whitespace = /\s+/

class Stats {
  constructor (opts) {
    opts = opts || {}

    this.files = {
      processFile: opts.processFile || '/proc/self/stat',
      memoryFile: opts.memoryFile || '/proc/meminfo',
      cpuFile: opts.cpuFile || '/proc/stat'
    }

    this.previous = {
      cpuTotal: 0,
      cpuUsage: 0,
      memTotal: 0,
      memAvailable: 0,
      utime: 0,
      stime: 0,
      vsize: 0,
      rss: 0
    }

    this.stats = {
      'system.cpu.total.norm.pct': 0,
      'system.memory.actual.free': 0,
      'system.memory.total': 0,
      'system.process.cpu.total.norm.pct': 0,
      'system.process.cpu.system.norm.pct': 0,
      'system.process.cpu.user.norm.pct': 0,
      'system.process.memory.size': 0,
      'system.process.memory.rss.bytes': 0
    }

    this.inProgress = false
    this.timer = null

    // Do initial load
    const files = [
      this.files.processFile,
      this.files.memoryFile,
      this.files.cpuFile
    ]

    try {
      const datas = files.map(readFileSync)
      this.previous = this.readStats(datas)
      this.update(datas)
    } catch (err) {}
  }

  toJSON () {
    return this.stats
  }

  collect (cb) {
    if (this.inProgress) {
      if (cb) process.nextTick(cb)
      return
    }

    this.inProgress = true

    const files = [
      this.files.processFile,
      this.files.memoryFile,
      this.files.cpuFile
    ]

    const next = afterAll((err, files) => {
      if (!err) this.update(files)
      if (cb) cb()
    })

    for (const file of files) {
      fs.readFile(file, next())
    }
  }

  readStats ([processFile, memoryFile, cpuFile]) {
    // CPU data
    //
    // Example of line we're trying to parse:
    // cpu  13978 30 2511 9257 2248 0 102 0 0 0

    const cpuLine = firstLineOfBufferAsString(cpuFile)
    const cpuTimes = cpuLine.split(whitespace)

    let cpuTotal = 0
    for (let i = 1; i < cpuTimes.length; i++) {
      cpuTotal += Number(cpuTimes[i])
    }

    // We're off-by-one in relation to the expected index, because we include
    // the `cpu` label at the beginning of the line
    const idle = Number(cpuTimes[4])
    const iowait = Number(cpuTimes[5])
    const cpuUsage = cpuTotal - idle - iowait

    // Memory data
    let memAvailable = 0
    let memTotal = 0

    let matches = 0
    for (const line of memoryFile.toString().split('\n')) {
      if (/^MemAvailable:/.test(line)) {
        memAvailable = parseInt(line.split(whitespace)[1], 10) * 1024
        matches++
      } else if (/^MemTotal:/.test(line)) {
        memTotal = parseInt(line.split(whitespace)[1], 10) * 1024
        matches++
      }
      if (matches === 2) break
    }

    // Process data
    //
    // Example of line we're trying to parse:
    //
    // 44 (node /app/node_) R 1 44 44 0 -1 4210688 7948 0 0 0 109 21 0 0 20 0 10 0 133652 954462208 12906 18446744073709551615 4194304 32940036 140735797366336 0 0 0 0 4096 16898 0 0 0 17 0 0 0 0 0 0 35037200 35143856 41115648 140735797369050 140735797369131 140735797369131 140735797370852 0
    //
    // We can't just split on whitespace as the 2nd field might contain
    // whitespace. However, the parentheses will always be there, so we can
    // ignore everything from before the `)` to get rid of the whitespace
    // problem.
    //
    // For details about each field, see:
    // http://man7.org/linux/man-pages/man5/proc.5.html

    const processLine = firstLineOfBufferAsString(processFile)
    const processData = processLine.slice(processLine.lastIndexOf(')')).split(whitespace)

    // all fields are referenced by their index, but are off by one because
    // we're dropping the first field from the line due to the whitespace
    // problem described above
    const utime = parseInt(processData[12], 10) // position in file: 14
    const stime = parseInt(processData[13], 10) // position in file: 15
    const vsize = parseInt(processData[21], 10) // position in file: 23

    return {
      cpuUsage,
      cpuTotal,
      memTotal,
      memAvailable,
      utime,
      stime,
      vsize,
      rss: process.memoryUsage().rss // TODO: Calculate using field 24 (rss) * PAGE_SIZE
    }
  }

  update (files) {
    const prev = this.previous
    const next = this.readStats(files)
    const stats = this.stats

    const cpuTotal = next.cpuTotal - prev.cpuTotal
    const cpuUsage = next.cpuUsage - prev.cpuUsage
    const utime = next.utime - prev.utime
    const stime = next.stime - prev.stime

    stats['system.cpu.total.norm.pct'] = cpuUsage / cpuTotal || 0
    stats['system.memory.actual.free'] = next.memAvailable
    stats['system.memory.total'] = next.memTotal

    // We use Math.min to guard against an edge case where /proc/self/stat
    // reported more clock ticks than /proc/stat, in which case it looks like
    // the process spent more CPU time than was used by the system. In that
    // case we just assume it was a 100% CPU.
    //
    // This might happen because we don't read the process file at the same
    // time as the system file. In between the two reads, the process will
    // spend some time on the CPU and hence the two reads are not 100% synced
    // up.
    const cpuProcessPercent = Math.min((utime + stime) / cpuTotal || 0, 1)
    const cpuProcessUserPercent = Math.min(utime / cpuTotal || 0, 1)
    const cpuProcessSystemPercent = Math.min(stime / cpuTotal || 0, 1)

    stats['system.process.cpu.total.norm.pct'] = cpuProcessPercent
    stats['system.process.cpu.user.norm.pct'] = cpuProcessUserPercent
    stats['system.process.cpu.system.norm.pct'] = cpuProcessSystemPercent
    stats['system.process.memory.size'] = next.vsize
    stats['system.process.memory.rss.bytes'] = next.rss

    this.previous = next
    this.inProgress = false
  }
}

function firstLineOfBufferAsString (buff) {
  const newline = buff.indexOf('\n')
  return buff.toString('utf8', 0, newline === -1 ? buff.length : newline)
}

function readFileSync (file) {
  return fs.readFileSync(file)
}

module.exports = Stats
