/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

/**
 * This file is extracted from the 'process-top' project copyright by
 * Mathias Buus. It has been modified slightly to be used in the current
 * context and where possible changes have been contributed back to the
 * original project.
 *
 * https://github.com/mafintosh/process-top
 *
 * Original file:
 *
 * https://github.com/mafintosh/process-top/blob/master/index.js
 *
 * License:
 *
 * The MIT License (MIT)
 *
 * Copyright (c) 2018 Mathias Buus
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 *
 */

// const os = require('os')
// const p = require('prettier-bytes')
// const eld = require('event-loop-delay')

module.exports = top

// function top (opts) {
function top () { // NOTE: Replacing line above
  // if (!opts) opts = {}

  // const tick = opts.interval || 1000
  // const started = Date.now()
  // const interval = setInterval(perSecond, tick)
  // const win = [{ time: process.hrtime(), cpu: process.cpuUsage(), delay: 0 }, null, null, null]
  const win = [{ time: process.hrtime(), cpu: process.cpuUsage() }, null, null, null] // NOTE: Replacing line above
  // const loopSampler = eld()

  let sec = 1

  // interval.unref()

  return {
    // pid: process.pid,
    // command: process.argv.join(' '),
    // started,
    // time () {
    //   return Date.now() - started
    // },
    // delay () {
    //   const btm = oldest()
    //   const timeDelta = process.hrtime(win[btm].time)
    //   const ms = timeDelta[0] * 1e3 + Math.floor(timeDelta[1] / 1e6)
    //   return Math.floor((loopSampler.delay - win[btm].delay) / (ms / tick))
    // },
    cpu () {
      const btm = oldest()
      const cpuDelta = process.cpuUsage(win[btm].cpu)
      const timeDelta = process.hrtime(win[btm].time)
      const us = timeDelta[0] * 1e6 + timeDelta[1] / 1e3
      perSecond() // NOTE: Added to skip timer and update at each check
      return {
        time: us,
        percent: (cpuDelta.system + cpuDelta.user) / us,
        system: cpuDelta.system,
        user: cpuDelta.user
      }
    }
    // memory () {
    //   const mem = process.memoryUsage()
    //   const total = os.totalmem()
    //   return {
    //     percent: mem.rss / total,
    //     rss: mem.rss,
    //     total,
    //     heapPercent: mem.heapUsed / mem.heapTotal,
    //     heapUsed: mem.heapUsed,
    //     heapTotal: mem.heapTotal,
    //     external: mem.external
    //   }
    // },
    // loadavg () {
    //   return os.loadavg()
    // },
    // destroy () {
    //   clearInterval(interval)
    // },
    // toString () {
    //   const mem = this.memory()
    //   return `cpu: ${pct(this.cpu().percent)} | rss: ${p(mem.rss)} (${pct(mem.percent)}) | heap: ${p(mem.heapUsed)} / ${p(mem.heapTotal)} (${pct(mem.heapPercent)}) | ext: ${p(mem.external)} | delay: ${this.delay()} ms | ${time(this.time())} | loadavg: ${os.loadavg().map(fixed2).join(', ')}`
    // }
  }

  function oldest () {
    let btm = (sec - 4) & 3
    while (!win[btm]) btm = (btm + 1) & 3
    return btm
  }

  function perSecond () {
    const ptr = sec++ & 3
    // win[ptr] = { time: process.hrtime(), cpu: process.cpuUsage(), delay: loopSampler.delay }
    win[ptr] = { time: process.hrtime(), cpu: process.cpuUsage() } // NOTE: Replacing line above
  }
}

// function pct (n) {
//   return (100 * n).toFixed(1) + '%'
// }

// function fixed2 (n) {
//   return n.toFixed(2)
// }

// function time (n) {
//   let secs = Math.floor(n / 1000)
//   let hours = Math.floor(secs / 3600)
//   secs -= hours * 3600
//   let mins = Math.floor(secs / 60)
//   secs -= mins * 60
//   return pad(hours) + ':' + pad(mins) + ':' + pad(secs)
// }

// function pad (n) {
//   return n < 10 ? '0' + n : '' + n
// }
