'use strict'

function toNano (t) {
  return (t[0] * 1e9) + t[1]
}

function sumReducer (last, next) {
  return last + next
}

function sortNumber (a, b) {
  a = Number.isNaN(a) ? Number.NEGATIVE_INFINITY : a
  b = Number.isNaN(b) ? Number.NEGATIVE_INFINITY : b

  if (a > b) return 1
  if (a < b) return -1

  return 0
}

function descriptor (target, property, data) {
  Object.defineProperty(target, property, data)
}

function define (target, property, value) {
  descriptor(target, property, { value })
}

function mean (list) {
  return list.reduce(sumReducer, 0) / list.length
}

class EventLoopDelayHistogram {
  constructor ({ resolution = 10 } = {}) {
    define(this, 'resolution', resolution)
    descriptor(this, 'timer', { writable: true })
    descriptor(this, 'samples', { writable: true, value: [] })
  }

  get stddev () {
    const avg = mean(this.samples)

    const squareDiffs = this.samples.map(value => {
      const diff = value - avg
      const sqrDiff = diff * diff
      return sqrDiff
    })

    return Math.sqrt(mean(squareDiffs))
  }

  get mean () {
    return mean(this.samples)
  }

  get min () {
    // This giant number is only present to emulate edge case behaviour.
    // When not yet enabled, this will be the value of min.
    return Math.min(9223372036854776000, ...this.samples)
  }

  get max () {
    return Math.max(0, ...this.samples)
  }

  get percentiles () {
    const map = new Map()
    let last = 0
    if (this.samples.length) {
      map.set(0, this.percentile(0))
      for (let percent = 50; percent < 100; percent += (100 - percent) / 2) {
        const next = this.percentile(percent)
        if (last === next) break
        map.set(percent, next)
        last = next
      }
    }
    map.set(100, this.percentile(100))
    return map
  }

  percentile (percent) {
    percent = Number(percent)

    if (isNaN(percent) || percent < 0 || percent > 100) {
      throw new TypeError('Percent must be a floating point number between 0 and 100')
    }

    const list = this.samples.sort(sortNumber)
    if (percent === 0) return list[0]

    return list[Math.ceil(list.length * (percent / 100)) - 1] || 0
  }

  enable () {
    if (this.timer) return false

    let last = process.hrtime()

    this.timer = setInterval(() => {
      const next = process.hrtime(last)
      this.samples.push(Math.max(0, toNano(next)))
      last = process.hrtime()
    }, this.resolution)

    this.timer.unref()

    return true
  }

  disable () {
    if (!this.timer) return false
    clearInterval(this.timer)
    this.timer = null
    return true
  }

  reset () {
    this.samples = []
  }
}

const proto = EventLoopDelayHistogram.prototype
descriptor(proto, 'stddev', { enumerable: true })
descriptor(proto, 'mean', { enumerable: true })
descriptor(proto, 'min', { enumerable: true })
descriptor(proto, 'max', { enumerable: true })

module.exports = EventLoopDelayHistogram
