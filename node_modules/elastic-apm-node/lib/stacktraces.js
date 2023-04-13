/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

// Handling for translating `<Error instance>.stack` to the
// `{error.exception,span}.stacktrace` object passed to APM server.
// https://github.com/elastic/apm-server/blob/master/docs/spec/v2/error.json#L633
//
// A lot of this derived from https://github.com/watson/stackman but was
// moved internal to allow breaking compat for perf work.

var fs = require('fs')
var path = require('path')

const asyncCache = require('async-cache')
const afterAllResults = require('after-all-results')

// avoid loading error-callsites until needed to avoid
// Error.prepareStackTrace side-effects
// https://github.com/elastic/apm-agent-nodejs/issues/2833
let errorCallsites
function initStackTraceCollection () {
  errorCallsites = require('error-callsites')
}
const errorStackParser = require('error-stack-parser')
const loadSourceMap = require('./load-source-map')
const LRU = require('lru-cache')

const fileCache = asyncCache({
  max: 500, // fileCacheMax
  load: function (file, cb) {
    fs.readFile(file, { encoding: 'utf8' }, function (err, data) {
      if (err) {
        cb(err)
        return
      }
      cb(null, data.split(/\r?\n/))
    })
  }
})

const sourceMapCache = asyncCache({
  max: 500, // sourceMapCacheMax
  load: function (filename, cb) {
    if (!filename) {
      cb(null, null)
      return
    }
    loadSourceMap(filename, function onSourceMap (err, sourcemap) {
      // Translate sourcemap===undefined to null, because 'async-cache'
      // treats `undefined` as a cache miss. Without this there is no
      // caching for files that have no sourcemap (the common case).
      cb(err, sourcemap || null)
    })
  },
  dispose: function (_filename, sourcemap) {
    if (sourcemap) {
      sourcemap.destroy()
    }
  }
})

const frameCache = new LRU({ max: 1000 })
const frameCacheStats = {
  hits: 0,
  misses: 0
}

let lastCwd = process.cwd()

// ---- internal support functions

// process.cwd() can throw EMFILE if hitting the `ulimit -Hn`. A process's
// cwd should not change regularly, so falling back to the previous one is fine.
// At worst the stacktrace frame.filename values will be relative to an old
// directory.
function getCwd (log) {
  let cwd
  try {
    cwd = process.cwd()
  } catch (ex) {
    log.trace(ex, 'error getting cwd: fallback back to %s', lastCwd)
    return lastCwd
  }
  lastCwd = cwd
  return cwd
}

// If gathering a stacktrace from the structured CallSites fails, this is
// used as a fallback: parsing the `err.stack` *string*.
function stackTraceFromErrStackString (log, err) {
  const stacktrace = []

  // graphqljs adds the `originalError` property which represents the original
  // error thrown within the resolver
  err = err.originalError || err

  if (err.stack === null) {
    return []
  }

  // frames is an array of StackFrame (https://github.com/stacktracejs/stackframe).
  let frames = null
  try {
    frames = errorStackParser.parse(err)
  } catch (parseErr) {
    log.debug('could not parse err.stack string: %s', parseErr)
  }

  if (frames) {
    const cwd = getCwd(log)
    for (var i = 0; i < frames.length; i++) {
      const frame = frames[i]
      const filename = frame.getFileName() || ''
      stacktrace.push({
        filename: getRelativeFileName(filename, cwd),
        function: frame.getFunctionName(),
        lineno: frame.getLineNumber(),
        library_frame: !isStackFrameApp(frame),
        abs_path: filename
      })
    }
  }

  return stacktrace
}

const NODE_MODULES_PATH_SEG = 'node_modules' + path.sep

// Return true iff the given StackFrame
// (https://github.com/stacktracejs/stackframe) is an application frame.
function isStackFrameApp (stackframe) {
  if (isStackFrameNode(stackframe)) {
    return false
  } else {
    const fileName = stackframe.getFileName()
    if (!fileName) {
      return true
    } else if (fileName.indexOf(NODE_MODULES_PATH_SEG) === -1) {
      return true
    } else {
      return false
    }
  }
}

// Return true iff the given StackFrame
// (https://github.com/stacktracejs/stackframe) is a Node frame.
function isStackFrameNode (stackframe) {
  if (stackframe.isNative) {
    return true
  } else {
    const fileName = stackframe.getFileName()
    if (!fileName) {
      return true
    } else {
      return (!path.isAbsolute(fileName) && fileName[0] !== '.')
    }
  }
}

function isCallSiteApp (callsite) {
  if (isCallSiteNode(callsite)) {
    return false
  } else {
    const fileName = callsite.getFileName()
    if (!fileName) {
      return true
    } else if (fileName.indexOf(NODE_MODULES_PATH_SEG) === -1) {
      return true
    } else {
      return false
    }
  }
}

function isCallSiteNode (callsite) {
  if (callsite.isNative()) {
    return true
  } else {
    const fileName = callsite.getFileName()
    if (!fileName) {
      return true
    } else {
      return (!path.isAbsolute(fileName) && fileName[0] !== '.')
    }
  }
}

function isValidCallsites (callsites) {
  return Array.isArray(callsites) &&
    callsites.length > 0 &&
    typeof callsites[0] === 'object' &&
    typeof callsites[0].getFileName === 'function'
}

// From stackman getTypeNameSafely().
function getCallSiteTypeNameSafely (callsite) {
  try {
    return callsite.getTypeName()
  } catch (e) {
    // This seems to happen sometimes when using 'use strict',
    // stemming from `getTypeName`.
    // [TypeError: Cannot read property 'constructor' of undefined]
    return null
  }
}

// From stackman getFunctionNameSanitized().
function getCallSiteFunctionNameSanitized (callsite) {
  var fnName = callsite.getFunctionName()
  if (fnName) return fnName
  var typeName = getCallSiteTypeNameSafely(callsite)
  if (typeName) return typeName + '.' + (callsite.getMethodName() || '<anonymous>')
  return '<anonymous>'
}

function addSourceContextToFrame (frame, lines, lineNum, n) {
  var index = lineNum - 1 // lines 1-based -> index 0-based
  var nBefore = Math.ceil((n - 1) / 2)
  var nAfter = Math.floor((n - 1) / 2)
  frame.pre_context = lines.slice(Math.max(0, index - nBefore), index)
  frame.context_line = lines[index]
  frame.post_context = lines.slice(index + 1, index + 1 + nAfter)
}

// Get the path of `filename` relative to `relTo` -- which should be a directory
// path *without* a trailing path separator.
function getRelativeFileName (filename, relTo) {
  if (filename.startsWith(relTo + path.sep)) {
    return filename.slice(relTo.length + 1)
  } else {
    return filename
  }
}

function getSourceMapConsumer (callsite, cb) {
  if (isCallSiteNode(callsite)) {
    process.nextTick(cb, null, null)
  } else {
    var filename = callsite.getFileName()
    if (!filename) {
      process.nextTick(cb, null, null)
    } else {
      sourceMapCache.get(filename, cb)
    }
  }
}

// Put together an APM stacktrace frame object:
//     {
//       "filename": "...",
//       "lineno": 65,
//       "function": "...",
//       "library_frame": <bool>,
//       "abs_path": "...",
//       "pre_context": [ ... ],
//       "context_line": "...",
//       "post_context": [ ... ],
//     },
// from a v8 CallSite object
// (https://v8.dev/docs/stack-trace-api#customizing-stack-traces).
//
// This calls back with `cb(null, frame)` -- the first err arg is always null.
function frameFromCallSite (log, callsite, cwd, sourceLinesAppFrames, sourceLinesLibraryFrames, cb) {
  // getFileName can return null, e.g. with a `at Generator.next (<anonymous>)` frame.
  const filename = callsite.getFileName() || ''
  const lineno = callsite.getLineNumber()
  const colno = callsite.getColumnNumber()

  // Caching
  const cacheKey = [filename, lineno, colno, sourceLinesAppFrames, sourceLinesLibraryFrames].join(':')
  const cachedFrame = frameCache.get(cacheKey)
  if (cachedFrame !== undefined) {
    frameCacheStats.hits++

    // Guard against later JSON serialization mistakenly changing duplicate
    // frames in a stacktrace (a totally legal thing) into '[Circular]' as a
    // guard against serializing an object with circular references.
    const clonedFrame = Object.assign({}, cachedFrame)
    if (clonedFrame.pre_context) {
      clonedFrame.pre_context = clonedFrame.pre_context.slice()
    }
    if (clonedFrame.post_context) {
      clonedFrame.post_context = clonedFrame.post_context.slice()
    }

    setImmediate(cb, null, clonedFrame)
    return
  }
  function cacheAndCb (frame) {
    frameCacheStats.misses++
    frameCache.set(cacheKey, frame)
    cb(null, frame)
  }

  // If the file has a sourcemap we use that for: filename, lineno, source
  // context.
  getSourceMapConsumer(callsite, function onSourceMapConsumer (sourceMapErr, sourceMapConsumer) {
    let mappedFilename = null
    let absMappedFilename = null
    let mappedLineno = null

    if (sourceMapErr) {
      log.debug({ filename: filename, err: sourceMapErr },
        'could not process file source map')
    } else if (sourceMapConsumer) {
      let pos
      try {
        pos = sourceMapConsumer.originalPositionFor({
          line: lineno,
          column: callsite.getColumnNumber()
        })
      } catch (posErr) {
        log.debug({ filename: filename, line: lineno, err: sourceMapErr },
          'could not get position from sourcemap')
        pos = {
          source: null,
          line: null,
          column: null,
          name: null
        }
      }
      if (pos.source !== null) {
        mappedFilename = pos.source
        absMappedFilename = path.resolve(path.dirname(filename), mappedFilename)
      }
      if (pos.line !== null) {
        mappedLineno = pos.line
      }
      // TODO: Is `pos.name` relevant for `frame.function` if minifying?
    }

    const isApp = isCallSiteApp(callsite)
    const frame = {
      filename: getRelativeFileName(absMappedFilename || filename, cwd),
      lineno: mappedLineno || lineno,
      function: getCallSiteFunctionNameSanitized(callsite),
      library_frame: !isApp
    }
    if (!Number.isFinite(frame.lineno)) {
      // An early comment in stackman suggested this is "sometimes not" an int.
      frame.lineno = 0
    }
    if (filename) {
      frame.abs_path = absMappedFilename || filename
    }

    // Finish early if we do not need to collect source lines of context.
    var linesOfContext = (isApp ? sourceLinesAppFrames : sourceLinesLibraryFrames)
    if (linesOfContext === 0 || !filename || isCallSiteNode(callsite)) {
      cacheAndCb(frame)
      return
    }

    // Attempt to use "sourcesContent" in a sourcemap, if available.
    if (sourceMapConsumer && mappedFilename && mappedLineno) {
      // To use `sourceMapConsumer.sourceContentFor` we need the filename as
      // it is in the "sources" field of the source map. `mappedFilename`,
      // from `sourceMapConsume.originalPositionFor` above, was made relative
      // to "sourceRoot" -- sourceFilename undoes that.
      const sourceFilename = sourceMapConsumer.sourceRoot
        ? path.relative(sourceMapConsumer.sourceRoot, mappedFilename)
        : mappedFilename
      var source = sourceMapConsumer.sourceContentFor(sourceFilename, true)
      log.trace({
        sourceRoot: sourceMapConsumer.sourceRoot,
        mappedFilename,
        sourceFilename,
        haveSourceContent: !!source
      }, 'sourcemap sourceContent lookup')
      if (source) {
        addSourceContextToFrame(frame, source.split(/\r?\n/g), mappedLineno, linesOfContext)
        cacheAndCb(frame)
        return
      }
    }

    // If the file looks like it minimized (as we didn't have a source-map in
    // the processing above), then skip adding source context because it
    // is mostly useless and the typically 500-char lines result in over-large
    // APM error objects.
    if (filename.endsWith('.min.js')) {
      cacheAndCb(frame)
      return
    }

    // Otherwise load the file from disk, if available.
    fileCache.get(frame.abs_path, function onFileCacheGet (fileErr, lines) {
      if (fileErr) {
        log.debug({ filename: frame.abs_path, err: fileErr },
          'could not read file for source context')
      } else {
        addSourceContextToFrame(frame, lines, frame.lineno, linesOfContext)
      }
      cacheAndCb(frame)
    })
  })
}

// ---- exports

// Gather an APM `stacktrace` object from the given `err`.
// This stacktrace object is used for `error.exception.stacktrace`,
// `error.log.stacktrace`, or `span.stacktrace`.
//
// This is a best effort, so it never fails. It will log.debug any parsing
// failures.
//
// @param {Logger} log
// @param {Error|Object} err - Typically an Error instance, but can also be a
//    plain object on which `Error.captureStackTrace(...)` has been called.
// @param {Integer} sourceLinesAppFrames - The number of source lines of
//    context to include for frames in application code.
// @param {Integer} sourceLinesLibraryFrames - The number of source lines of
//    context to include for frames in library code (i.e. under node_modules/)
// @param {Function} filterCallSite - Optional. A function to filter the
//    CallSites to include in the stacktrace.
// @param {Function} cb - `cb(null, stacktrace)`
function gatherStackTrace (log, err, sourceLinesAppFrames, sourceLinesLibraryFrames, filterCallSite, cb) {
  // errorCallsites returns an array of v8 CallSite objects.
  // https://v8.dev/docs/stack-trace-api#customizing-stack-traces
  let callsites = errorCallsites ? errorCallsites(err) : null

  const next = afterAllResults(function finish (_err, stacktrace) {
    // _err is always null from frameFromCallSite.

    // If we are not able to extract callsite information from err, then
    // fallback to parsing the err.stack string.
    if (stacktrace.length === 0) {
      stacktrace = stackTraceFromErrStackString(log, err)
    }

    cb(null, stacktrace)
  })

  if (!isValidCallsites(callsites)) {
    // When can this happen? Another competing Error.prepareStackTrace breaking
    // error-callsites? Also initStackTraceCollection not having been called.
    log.debug('could not get valid callsites from error "%s"', err)
  } else if (callsites) {
    if (filterCallSite) {
      callsites = callsites.filter(filterCallSite)
    }
    const cwd = getCwd(log)
    for (let i = 0; i < callsites.length; i++) {
      frameFromCallSite(
        log,
        callsites[i],
        cwd,
        sourceLinesAppFrames,
        sourceLinesLibraryFrames,
        next()
      )
    }
  }
}

module.exports = {
  gatherStackTrace,
  frameCacheStats,
  initStackTraceCollection,

  // Exported for testing only.
  stackTraceFromErrStackString
}
