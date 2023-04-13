'use strict'

// Logging utilities for the APM http client.

// A logger that does nothing and supports enough of the pino API
// (https://getpino.io/#/docs/api?id=logger) for use as a fallback in
// this package.
class NoopLogger {
  trace () {}
  debug () {}
  info () {}
  warn () {}
  error () {}
  fatal () {}
  child () { return this }
  isLevelEnabled (_level) { return false }
}

module.exports = {
  NoopLogger
}
