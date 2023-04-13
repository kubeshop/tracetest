/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

/**
 * Central location for shared constants
 */

module.exports = {
  // The default span or transaction `type`.
  DEFAULT_SPAN_TYPE: 'custom',

  REDACTED: '[REDACTED]',
  OUTCOME_FAILURE: 'failure',
  OUTCOME_SUCCESS: 'success',
  OUTCOME_UNKNOWN: 'unknown',
  RESULT_SUCCESS: 'success',
  RESULT_FAILURE: 'failure',

  // https://github.com/elastic/apm/blob/main/specs/agents/tracing-instrumentation-messaging.md#receiving-trace-context
  MAX_MESSAGES_PROCESSED_FOR_TRACE_CONTEXT: 1000
}
