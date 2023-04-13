/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

const util = require('util')
const otel = require('@opentelemetry/api')

const logging = require('../logging')
const { fetchSpanKey } = require('./OTelBridgeRunContext')
const oblog = require('./oblog')
const { OTelContextManager } = require('./OTelContextManager')
const { OTelTracerProvider } = require('./OTelTracerProvider')
const { OTelTracer } = require('./OTelTracer')

function setupOTelBridge (agent) {
  let success

  const log = (logging.isLoggerCustom(agent.logger)
    ? agent.logger
    : agent.logger.child({ 'event.module': 'otelbridge' }))

  // `otel.diag` is "the OpenTelemetry internal diagnostic API". If *trace*
  // level logging is enabled, then hook into diag.
  if (log.isLevelEnabled('trace')) {
    success = otel.diag.setLogger({
      verbose: log.trace.bind(log),
      debug: log.debug.bind(log),
      info: log.info.bind(log),
      warn: log.warn.bind(log),
      error: log.error.bind(log)
    }, otel.DiagLogLevel.ALL)
    if (!success) {
      log.error('could not register OpenTelemetry bridge diag logger')
      return
    }
  }

  // For development/debugging-only, set `LOG_OTEL_API_CALLS = true` to get log
  // (almost) every call into the OpenTelemetry API. See docs in "oblog.js".
  const LOG_OTEL_API_CALLS = false
  if (LOG_OTEL_API_CALLS) {
    // oblog.setApiCallLogFn(log.debug.bind(log)) // Alternative, to use our ecs logger.
    oblog.setApiCallLogFn((...args) => {
      const s = util.format(...args)
      console.log('\x1b[90motelapi:\x1b[39m \x1b[32m' + s + '\x1b[39m')
    })
  }

  success = otel.trace.setGlobalTracerProvider(new OTelTracerProvider(new OTelTracer(agent)))
  if (!success) {
    log.error('could not register OpenTelemetry bridge TracerProvider')
    return
  }

  // The OTelBridgeRunContext class needs to get the SPAN_KEY before it can
  // be used.
  fetchSpanKey()

  success = otel.context.setGlobalContextManager(new OTelContextManager(agent))
  if (!success) {
    log.error('could not register OpenTelemetry bridge ContextManager')
  }
}

module.exports = {
  setupOTelBridge
}
