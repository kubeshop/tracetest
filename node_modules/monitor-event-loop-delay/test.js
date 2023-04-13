const tap = require('tap')

let perfHooks
try {
  perfHooks = require('perf_hooks')
} catch (err) {}

const EventLoopDelayHistogram = require('./polyfill')

function nearly (value, expected, range) {
  return value > expected * (1 - range) && value < expected * (1 + range)
}

const emptyPercentiles = new Map([
  [100, 0]
])

function test (makeMonitor) {
  return t => {
    t.test('basic', t => {
      const monitor = makeMonitor({
        resolution: 10
      })

      t.comment('before enable')
      t.ok(isNaN(monitor.stddev), 'stddev is NaN')
      t.ok(isNaN(monitor.mean), 'mean is NaN')
      t.equal(monitor.min, 9223372036854776000, 'min is 9223372036854776000 - due to Node.js core edge case')
      t.equal(monitor.max, 0, 'max is zero')
      t.equal(monitor.percentile(50), 0, '50th percentile is zero')
      t.deepEqual(monitor.percentiles, emptyPercentiles)

      t.ok(monitor.enable(), 'enables when disabled')
      t.notOk(monitor.enable(), 'does not enable when already enabled')

      t.comment('after enable')
      t.ok(isNaN(monitor.stddev), 'stddev is NaN')
      t.ok(isNaN(monitor.mean), 'mean is NaN')
      t.equal(monitor.min, 9223372036854776000, 'min is 9223372036854776000 - due to Node.js core edge case')
      t.equal(monitor.max, 0, 'max is zero')
      t.equal(monitor.percentile(50), 0, '50th percentile is zero')
      t.deepEqual(monitor.percentiles, emptyPercentiles)

      setTimeout(() => {
        const {
          stddev,
          mean,
          min,
          max,
          percentiles
        } = monitor

        t.comment('after timeout')
        t.ok(nearly(stddev, 1.5e7, 1.5), 'stddev is in expected range')
        t.ok(nearly(mean, 1.5e7, 0.7), 'mean is in expected range')
        t.ok(nearly(min, 1.5e7, 0.7), 'min is in expected range')
        t.ok(nearly(max, 1.5e7, 0.7), 'max is in expected range')
        t.ok(nearly(monitor.percentile(50), mean, 0.5), '50th percentile is zero')
        t.ok(percentiles.has(75), 'percentiles map filled out')

        t.ok(monitor.disable(), 'disables when enabled')
        t.notOk(monitor.disable(), 'does not disable when already disabled')

        setTimeout(() => {
          t.comment('after disable')
          t.ok(nearly(monitor.stddev, stddev, 0.2), 'stddev values are equal')
          t.ok(nearly(monitor.mean, mean, 0.2), 'mean values are equal')
          t.ok(nearly(monitor.min, min, 0.2), 'min values are equal')
          t.ok(nearly(monitor.max, max, 0.2), 'max values are equal')
          t.ok(nearly(monitor.percentile(50), mean, 0.2), '50th percentile is zero')
          t.ok(nearly(monitor.percentiles.get(75), percentiles.get(75), 0.2), 'percentiles map filled out')

          monitor.reset()

          t.comment('after reset')
          t.ok(isNaN(monitor.stddev), 'stddev is NaN')
          t.ok(isNaN(monitor.mean), 'mean is NaN')
          t.equal(monitor.min, 9223372036854776000, 'min is 9223372036854776000 - due to Node.js core edge case')
          t.equal(monitor.max, 0, 'max is zero')
          t.equal(monitor.percentile(50), 0, '50th percentile is zero')
          t.deepEqual(monitor.percentiles, emptyPercentiles)

          t.end()
        }, 100)
      }, 100)
    })

    t.test('storm', t => {
      t.plan(10)

      const sampler = makeMonitor({
        resolution: 10
      })
      sampler.enable()
      let last = process.hrtime()

      const check = setInterval(() => {
        const value = nanoToMilli(sampler.mean || 0)
        const total = toMillis(process.hrtime(last))
        last = process.hrtime()
        sampler.reset()
        t.ok(value >= 0 && value < total, `value is realistic (value: ${value}, total: ${total})`)
      }, 100)

      t.on('end', () => {
        clearInterval(check)
        sampler.disable()
        storm.stop()
      })

      const storm = makeStorm(1e4)
    })

    t.end()
  }
}

tap.test('polyfill', test(opts => {
  return new EventLoopDelayHistogram(opts)
}))

if (perfHooks) {
  tap.test('perf_hooks', test(opts => {
    return perfHooks.monitorEventLoopDelay(opts)
  }))
}

class TimeoutCollection {
  constructor () {
    this.pending = new Set()
    this.stopped = false
  }

  create (cb, ...args) {
    if (this.stopped) return

    const timer = setTimeout(() => {
      this.pending.delete(timer)
      cb()
    }, ...args)

    this.pending.add(timer)
    return timer
  }

  stop () {
    this.stopped = true
    for (let timer of this.pending) {
      this.pending.delete(timer)
      clearTimeout(timer)
    }
  }
}

function makeStorm (n) {
  const timers = new TimeoutCollection()

  function chain () {
    timers.create(chain, Math.floor(Math.random() * 50))
  }

  for (let i = 0; i < n; i++) {
    chain()
  }

  return timers
}

function toMillis (t) {
  return (t[0] * 1e3) + (t[1] / 1e6)
}

function nanoToMilli (ms) {
  return ms / 1e6
}
