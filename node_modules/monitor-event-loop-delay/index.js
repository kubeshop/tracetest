try {
  const perfHooks = require('perf_hooks')
  if (typeof perfHooks.monitorEventLoopDelay !== 'function') {
    throw new Error('No builtin event loop monitor')
  }
  module.exports = opts => perfHooks.monitorEventLoopDelay(opts)
} catch (err) {
  const EventLoopDelayHistogram = require('./polyfill')
  module.exports = opts => new EventLoopDelayHistogram(opts)
}
