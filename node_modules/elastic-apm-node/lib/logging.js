/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

// Internal logging for the Elastic APM Node.js Agent.
//
// Promised interface:
// - The amount of logging can be controlled via the `logLevel` config var,
//   and via the `log_level` central config var.
// - A custom logger can be provided via the `logging` config var.
//
// Nothing else about this package's logging (e.g. structure or the particular
// message text) is promised/stable.
//
// Per https://github.com/elastic/apm/blob/main/specs/agents/logging.md
// the valid log levels are:
//  - trace
//  - debug
//  - info (default)
//  - warning
//  - error
//  - critical
//  - off
//
// Before this spec, the supported levels were:
//  - trace
//  - debug
//  - info (default)
//  - warn - both "warn" and "warning" will be supported for backward compat
//  - error
//  - fatal - mapped to "critical" for backward compat

var ecsFormat = require('@elastic/ecs-pino-format')
var pino = require('pino')
var semver = require('semver')

const DEFAULT_LOG_LEVEL = 'info'

// Used to mark loggers created here, for use by `isLoggerCustom()`.
const LOGGER_IS_OURS_SYM = Symbol('ElasticAPMLoggerIsOurs')

const PINO_LEVEL_FROM_LEVEL_NAME = {
  trace: 'trace',
  debug: 'debug',
  info: 'info',
  warning: 'warn',
  warn: 'warn', // Supported for backwards compat
  error: 'error',
  critical: 'fatal',
  fatal: 'fatal', // Supported for backwards compat
  off: 'silent'
}

// SafePinoDestWrapper provides a pino destination that will pass logging calls
// to a given `customLogger`. The custom logger must have the following API:
//
// - `.trace(string)`
// - `.debug(string)`
// - `.info(string)`
// - `.warn(string)`
// - `.error(string)`
// - `.fatal(string)`
//
// The limitation of this wrapping is that structured data fields are *not*
// passed on to the custom logger. I.e. this is a fallback mechanism.
class SafePinoDestWrapper {
  constructor (customLogger) {
    this.customLogger = customLogger
    this.logFnNameFromLastLevel = pino.levels.labels
    this[Symbol.for('pino.metadata')] = true
  }

  write (s) {
    const { lastMsg, lastLevel } = this
    const logFnName = this.logFnNameFromLastLevel[lastLevel]
    this.customLogger[logFnName](lastMsg)
  }
}

// Create a pino logger for the agent.
//
// By default `createLogger()` will return a pino logger that logs to stdout
// in ecs-logging format, set to the "info" level.
//
// @param {String} levelName - Optional, default "info". It is meant to be one
//    of the log levels specified in the top of file comment. For backward
//    compatibility it falls back to "trace".
// @param {Object} customLogger - Optional. A custom logger object to which
//    log messages will be passed. It must provide
//    trace/debug/info/warn/error/fatal methods that take a string argument.
//
//    Internally the agent uses structured logging using the pino API
//    (https://getpino.io/#/docs/api?id=logger). However, with a custom logger,
//    log record fields other than the *message* are dropped, to avoid issues
//    with incompatible logger APIs.
//
//    As a special case, if the provided logger is a *pino logger instance*,
//    then it will be used directly.
function createLogger (levelName, customLogger) {
  let dest
  const serializers = {
    err: pino.stdSerializers.err,
    req: pino.stdSerializers.req,
    res: pino.stdSerializers.res
  }

  if (!levelName) {
    levelName = DEFAULT_LOG_LEVEL
  }
  let pinoLevel = PINO_LEVEL_FROM_LEVEL_NAME[levelName]
  if (!pinoLevel) {
    // For backwards compat, support an earlier bug where an unknown log level
    // was accepted.
    // TODO: Consider being more strict on this for v4.0.0.
    pinoLevel = 'trace'
  }

  if (customLogger) {
    // Is this a pino logger? If so, it supports the API the agent requires and
    // can be used directly. We must add our custom serializers.
    if (Symbol.for('pino.serializers') in customLogger) {
      // Pino added `options` second arg to `logger.child` in 6.12.0.
      if (semver.gte(customLogger.version, '6.12.0')) {
        return customLogger.child({}, { serializers })
      }

      return customLogger.child({
        serializers: serializers
      })
    }

    // Otherwise, we fallback to wrapping the provided logger such that the
    // agent can use the pino logger API without breaking. The limitation is
    // that only the log *message* is logged. Extra structured fields are
    // dropped.
    dest = new SafePinoDestWrapper(customLogger)
    // Our wrapping logger level should be 'trace', to pass through all
    // messages to the wrapped logger.
    pinoLevel = 'trace'
  } else {
    // Log to stdout, the same default as pino itself.
    dest = pino.destination(1)
  }

  const logger = pino({
    name: 'elastic-apm-node',
    base: {}, // Don't want pid and hostname fields.
    level: pinoLevel,
    serializers: serializers,
    ...ecsFormat({ apmIntegration: false })
  }, dest)

  if (!customLogger) {
    logger[LOGGER_IS_OURS_SYM] = true // used for isLoggerCustom()
  }
  return logger
}

function isLoggerCustom (logger) {
  return !logger[LOGGER_IS_OURS_SYM]
}

// Adjust the level on the given logger.
function setLogLevel (logger, levelName) {
  const pinoLevel = PINO_LEVEL_FROM_LEVEL_NAME[levelName]
  if (!pinoLevel) {
    logger.warn('unknown log levelName "%s": cannot setLogLevel', levelName)
  } else {
    logger.level = pinoLevel
  }
}

module.exports = {
  DEFAULT_LOG_LEVEL: DEFAULT_LOG_LEVEL,
  createLogger: createLogger,
  isLoggerCustom: isLoggerCustom,
  setLogLevel: setLogLevel
}
