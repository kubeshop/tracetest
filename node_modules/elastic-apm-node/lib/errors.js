/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

// Handle creating error event objects to be sent to APM server.
// https://github.com/elastic/apm-server/blob/master/docs/spec/v2/error.json

const crypto = require('crypto')
var path = require('path')
const util = require('util')

const { gatherStackTrace } = require('./stacktraces')

const MYSQL_ERROR_MSG_RE = /(ER_[A-Z_]+): /

// ---- internal support functions

// Default `culprit` to the top of the stack or the highest non `library_frame`
// frame if such exists
function culpritFromStacktrace (frames) {
  if (frames.length === 0) return

  var filename = frames[0].filename
  var fnName = frames[0].function
  for (var n = 0; n < frames.length; n++) {
    if (!frames[n].library_frame) {
      filename = frames[n].filename
      fnName = frames[n].function
      break
    }
  }

  return filename ? fnName + ' (' + filename + ')' : fnName
}

// Infer the node.js module name from the top frame filename, if possible.
// Here `frames` is a data structure as returned by `parseStackTrace`.
//
// Examples:
//    node_modules/mymodule/index.js
//                 ^^^^^^^^
//    node_modules/@myorg/mymodule/index.js
//                 ^^^^^^^^^^^^^^^
// or on Windows:
//    node_modules\@myorg\mymodule\lib\subpath\index.js
//                 ^^^^^^^^^^^^^^^
let SEP = path.sep
if (SEP === '\\') {
  SEP = '\\' + SEP // Escape this for use in a regex.
}
const MODULE_NAME_REGEX = new RegExp(`node_modules${SEP}([^${SEP}]*)(${SEP}([^${SEP}]*))?`)
function _moduleNameFromFrames (frames) {
  if (frames.length === 0) return
  var frame = frames[0]
  if (!frame.library_frame) return
  var match = frame.filename.match(MODULE_NAME_REGEX)
  if (!match) return
  var moduleName = match[1]
  if (moduleName && moduleName[0] === '@' && match[3]) {
    // Normalize the module name separator to '/', even on Windows.
    moduleName += '/' + match[3]
  }
  return moduleName
}

// Gather properties from `err` to be used for `error.exception.attributes`.
// If there are no properties to include, then it returns undefined.
function attributesFromErr (err) {
  let n = 0
  const attrs = {}
  const keys = Object.keys(err)
  for (let i = 0; i < keys.length; i++) {
    const key = keys[i]
    if (key === 'stack') {
      continue // 'stack' seems to be enumerable in Node 0.11
    }
    if (key === 'code') {
      continue // 'code' is already used for `error.exception.code`
    }

    let val = err[key]
    if (val === null) {
      continue // null is typeof object and well break the switch below
    }
    switch (typeof val) {
      case 'function':
        continue
      case 'object':
        // Ignore all objects except Dates.
        if (typeof val.toISOString !== 'function' || typeof val.getTime !== 'function') {
          continue
        } else if (Number.isNaN(val.getTime())) {
          val = 'Invalid Date' // calling toISOString() on invalid dates throws
        } else {
          val = val.toISOString()
        }
    }
    attrs[key] = val
    n++
  }
  return n ? attrs : undefined
}

// ---- exports

function generateErrorId () {
  return crypto.randomBytes(16).toString('hex')
}

// Create an "error" APM event object to be sent to APM server.
//
// Required args:
// - Exactly one of `args.exception` or `args.logMessage` must be set.
//   `args.exception` is an Error instance. `args.logMessage` is a log message
//   string, or an object of the form `{ message: 'template', params: [ ... ]}`
//   which will be formated with `util.format()`.
// - `args.id` - An ID for the error. It should be created with
//   `errors.generateErrorId()`.
// - `args.log` {Logger}
// - `args.shouldCaptureAttributes` {Boolean}
// - `args.timestampUs` {Integer} - Timestamp of the error in microseconds.
// - `args.handled` {Boolean}
// - `args.sourceLinesAppFrames` {Integer} - Number of lines of source context
//   to include in stack traces.
// - `args.sourceLinesLibraryFrames` {Integer} - Number of lines of source
//   context to include in stack traces. This and the previous arg are typically
//   select from the agent `sourceLines{Error,Span}{App,Library}Frames` config
//   vars.
//
// Optional args:
// - `args.callSiteLoc` - A `Error.captureStackTrace`d object with a stack to
//   include as `error.log.stacktrace`.
// - `args.traceContext` - The current TraceContext, if any.
// - `args.trans` - The current transaction, if any.
// - `args.errorContext` - An object to be included as `error.context`.
// - `args.message` - A message string that will be included as `error.log.message`
//   if `args.exception` is given. Ignored if `args.logMessage` is given.
// - `args.exceptionType` - A string to use for `error.exception.type`. By
//   default `args.exception.name` is used. This argument is only relevant if
//   `args.exception` was provided.
//
// This always calls back with `cb(null, apmError)`, i.e. it doesn't fail.
function createAPMError (args, cb) {
  let numAsyncStepsRemaining = 0 // finish() will call cb() only when this is 0.

  const error = {
    id: args.id,
    timestamp: args.timestampUs
  }
  if (args.traceContext) {
    error.parent_id = args.traceContext.traceparent.id
    error.trace_id = args.traceContext.traceparent.traceId
  }
  if (args.trans) {
    error.transaction_id = args.trans.id
    error.transaction = {
      name: args.trans.name,
      type: args.trans.type,
      sampled: args.trans.sampled
    }
  }
  if (args.errorContext) {
    error.context = args.errorContext
  }

  if (args.exception) {
    // Handle an exception, i.e. `captureError(<an Error instance>, ...)`.
    const err = args.exception
    const errMsg = String(err.message)
    error.exception = {
      message: errMsg,
      type: args.exceptionType || String(err.name),
      handled: args.handled
    }

    if ('code' in err) {
      error.exception.code = String(err.code)
    } else {
      // To provide better grouping of mysql errors that happens after the async
      // boundery, we modify to exception type to include the custom mysql error
      // type (e.g. ER_PARSE_ERROR)
      var match = errMsg.match(MYSQL_ERROR_MSG_RE)
      if (match) {
        error.exception.code = match[1]
      }
    }

    // Optional add an alternative error message as well as the exception message.
    if (args.message && typeof args.message === 'string') {
      error.log = { message: args.message }
    }

    if (args.shouldCaptureAttributes) {
      const attrs = attributesFromErr(err)
      if (attrs) {
        error.exception.attributes = attrs
      }
    }

    numAsyncStepsRemaining++
    gatherStackTrace(
      args.log,
      args.exception,
      args.sourceLinesAppFrames,
      args.sourceLinesLibraryFrames,
      null, // filterCallSite
      function (_err, stacktrace) {
        // _err from gatherStackTrace is always null.

        const culprit = culpritFromStacktrace(stacktrace)
        if (culprit) {
          error.culprit = culprit
        }
        const moduleName = _moduleNameFromFrames(stacktrace)
        if (moduleName) {
          // TODO: consider if we should include this as it's not originally what module was intended for
          error.exception.module = moduleName
        }
        error.exception.stacktrace = stacktrace
        finish()
      }
    )
  } else {
    // Handle a logMessage, i.e. `captureError(<not an Error instance>, ...)`.
    error.log = {}
    const msg = args.logMessage
    if (typeof msg === 'string') {
      error.log.message = msg
    } else if (typeof msg === 'object' && msg !== null) {
      if (msg.message) {
        error.log.message = util.format.apply(this, [msg.message].concat(msg.params))
        error.log.param_message = msg.message
      } else {
        error.log.message = util.inspect(msg)
      }
    } else {
      error.log.message = String(msg)
    }
  }

  if (args.callSiteLoc) {
    numAsyncStepsRemaining++
    gatherStackTrace(
      args.log,
      args.callSiteLoc,
      args.sourceLinesAppFrames,
      args.sourceLinesLibraryFrames,
      null, // filterCallSite
      function (_err, stacktrace) {
        // _err from gatherStackTrace is always null.

        if (stacktrace) {
          // In case there isn't any log object, we'll make a dummy message
          // as the APM Server requires a message to be present if a
          // stacktrace also present
          if (!error.log) {
            error.log = { message: error.exception.message }
          }
          error.log.stacktrace = stacktrace
          finish()
        }
      }
    )
  } else {
    numAsyncStepsRemaining++
    setImmediate(finish)
  }

  function finish () {
    numAsyncStepsRemaining--
    if (numAsyncStepsRemaining === 0) {
      cb(null, error)
    }
  }
}

module.exports = {
  generateErrorId,
  createAPMError,

  // Exported for testing.
  attributesFromErr,
  _moduleNameFromFrames
}
