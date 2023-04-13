/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

/**
 * This file is extracted from the 'async-listener' project copyright by
 * Forrest L Norvell. It have been modified slightly to be used in the current
 * context and where possible changes have been contributed back to the
 * original project.
 *
 * https://github.com/othiym23/async-listener
 *
 * Original file:
 *
 * https://github.com/othiym23/async-listener/blob/master/index.js
 *
 * License:
 *
 * BSD-2-Clause, http://opensource.org/licenses/BSD-2-Clause
 */

var { promisify } = require('util')

var isNative = require('is-native')
var semver = require('semver')

var shimmer = require('./shimmer')

var wrap = shimmer.wrap
var massWrap = shimmer.massWrap

var v7plus = semver.gte(process.version, '7.0.0')
var v11plus = semver.gte(process.version, '11.0.0')

module.exports = function (ins) {
  var net = require('net')

  // From Node.js v7.0.0, net._normalizeConnectArgs have been renamed net._normalizeArgs
  if (v7plus && !net._normalizeArgs) {
    // a polyfill in our polyfill etc so forth -- taken from node master on 2017/03/09
    net._normalizeArgs = function (args) {
      if (args.length === 0) {
        return [{}, null]
      }

      var arg0 = args[0]
      var options = {}
      if (typeof arg0 === 'object' && arg0 !== null) {
        // (options[...][, cb])
        options = arg0
      } else if (isPipeName(arg0)) {
        // (path[...][, cb])
        options.path = arg0
      } else {
        // ([port][, host][...][, cb])
        options.port = arg0
        if (args.length > 1 && typeof args[1] === 'string') {
          options.host = args[1]
        }
      }

      var cb = args[args.length - 1]
      if (typeof cb !== 'function') {
        return [options, null]
      } else {
        return [options, cb]
      }
    }
  } else if (!v7plus && !net._normalizeConnectArgs) {
    // a polyfill in our polyfill etc so forth -- taken from node master on 2013/10/30
    net._normalizeConnectArgs = function (args) {
      var options = {}

      function toNumber (x) { return (x = Number(x)) >= 0 ? x : false }

      if (typeof args[0] === 'object' && args[0] !== null) {
        // connect(options, [cb])
        options = args[0]
      } else if (typeof args[0] === 'string' && toNumber(args[0]) === false) {
        // connect(path, [cb])
        options.path = args[0]
      } else {
        // connect(port, [host], [cb])
        options.port = args[0]
        if (typeof args[1] === 'string') {
          options.host = args[1]
        }
      }

      var cb = args[args.length - 1]
      return typeof cb === 'function' ? [options, cb] : [options]
    }
  }

  wrap(net.Server.prototype, '_listen2', function (original) {
    return function () {
      this.on('connection', function (socket) {
        if (socket._handle) {
          socket._handle.onread = ins.bindFunction(socket._handle.onread)
        }
      })

      try {
        return original.apply(this, arguments)
      } finally {
        // the handle will only not be set in cases where there has been an error
        if (this._handle && this._handle.onconnection) {
          this._handle.onconnection = ins.bindFunction(this._handle.onconnection)
        }
      }
    }
  })

  function patchOnRead (ctx) {
    if (ctx && ctx._handle) {
      var handle = ctx._handle
      if (!handle._obOriginalOnread) {
        handle._obOriginalOnread = handle.onread
      }
      handle.onread = ins.bindFunction(handle._obOriginalOnread)
    }
  }

  wrap(net.Socket.prototype, 'connect', function (original) {
    return function () {
      // From Node.js v7.0.0, net._normalizeConnectArgs have been renamed net._normalizeArgs
      var args = v7plus
        ? net._normalizeArgs(arguments)
        : net._normalizeConnectArgs(arguments)
      if (args[1]) args[1] = ins.bindFunction(args[1])
      var result = original.apply(this, args)
      patchOnRead(this)
      return result
    }
  })

  var http = require('http')

  // NOTE: A rewrite occurred in 0.11 that changed the addRequest signature
  // from (req, host, port, localAddress) to (req, options)
  // Here, I use the longer signature to maintain 0.10 support, even though
  // the rest of the arguments aren't actually used
  wrap(http.Agent.prototype, 'addRequest', function (original) {
    return function (req) {
      var onSocket = req.onSocket
      req.onSocket = ins.bindFunction(function (socket) {
        patchOnRead(socket)
        return onSocket.apply(this, arguments)
      })
      return original.apply(this, arguments)
    }
  })

  var childProcess = require('child_process')

  function wrapChildProcess (child) {
    if (Array.isArray(child.stdio)) {
      for (const socket of child.stdio) {
        if (socket && socket._handle) {
          socket._handle.onread = ins.bindFunction(socket._handle.onread)
          wrap(socket._handle, 'close', activatorFirst)
        }
      }
    }

    if (child._handle) {
      child._handle.onexit = ins.bindFunction(child._handle.onexit)
    }
  }

  // iojs v2.0.0+
  if (childProcess.ChildProcess) {
    wrap(childProcess.ChildProcess.prototype, 'spawn', function (original) {
      return function () {
        var result = original.apply(this, arguments)
        wrapChildProcess(this)
        return result
      }
    })
  } else {
    massWrap(childProcess, [
      'execFile', // exec is implemented in terms of execFile
      'fork',
      'spawn'
    ], function (original) {
      return function () {
        var result = original.apply(this, arguments)
        wrapChildProcess(result)
        return result
      }
    })
  }

  // need unwrapped nextTick for use within < 0.9 async error handling
  if (!process._fatalException) {
    process._originalNextTick = process.nextTick
  }

  var processors = []
  if (process._nextDomainTick) processors.push('_nextDomainTick')
  if (process._tickDomainCallback) processors.push('_tickDomainCallback')

  massWrap(
    process,
    processors,
    activator
  )
  wrap(process, 'nextTick', activatorFirst)

  var asynchronizers = [
    'setTimeout',
    'setInterval'
  ]
  if (global.setImmediate) asynchronizers.push('setImmediate')

  var timers = require('timers')
  var patchGlobalTimers = global.setTimeout === timers.setTimeout

  massWrap(
    timers,
    asynchronizers,
    activatorFirst
  )

  if (patchGlobalTimers) {
    massWrap(
      global,
      asynchronizers,
      activatorFirst
    )
  }

  var dns = require('dns')
  massWrap(
    dns,
    [
      'lookup',
      'resolve',
      'resolve4',
      'resolve6',
      'resolveCname',
      'resolveMx',
      'resolveNs',
      'resolveTxt',
      'resolveSrv',
      'reverse'
    ],
    activator
  )

  if (dns.resolveNaptr) wrap(dns, 'resolveNaptr', activator)

  var fs = require('fs')
  var wrappedFsRealpathNative
  if (fs.realpath.native) {
    wrappedFsRealpathNative = wrap(fs.realpath, 'native', activator)
  }

  massWrap(
    fs,
    [
      'watch',
      'rename',
      'truncate',
      'chown',
      'fchown',
      'chmod',
      'fchmod',
      'stat',
      'lstat',
      'fstat',
      'link',
      'symlink',
      'readlink',
      'realpath',
      'unlink',
      'rmdir',
      'mkdir',
      'readdir',
      'close',
      'open',
      'utimes',
      'futimes',
      'fsync',
      'write',
      'read',
      'readFile',
      'writeFile',
      'appendFile',
      'watchFile',
      'unwatchFile',
      'exists'
    ],
    activator
  )

  if (wrappedFsRealpathNative) {
    fs.realpath.native = wrappedFsRealpathNative
  }

  // only wrap lchown and lchmod on systems that have them.
  if (fs.lchown) wrap(fs, 'lchown', activator) // eslint-disable-line node/no-deprecated-api
  if (fs.lchmod) wrap(fs, 'lchmod', activator) // eslint-disable-line node/no-deprecated-api

  // only wrap ftruncate in versions of node that have it
  if (fs.ftruncate) wrap(fs, 'ftruncate', activator)

  // Wrap zlib streams
  var zlib
  try { zlib = require('zlib') } catch (err) { }
  if (zlib && zlib.Deflate && zlib.Deflate.prototype) {
    var proto = Object.getPrototypeOf(zlib.Deflate.prototype)
    if (proto._transform) {
      // streams2
      wrap(proto, '_transform', activator)
    } else if (proto.write && proto.flush && proto.end) {
      // plain ol' streams
      massWrap(
        proto,
        [
          'write',
          'flush',
          'end'
        ],
        activator
      )
    }
  }

  // Wrap Crypto
  var crypto
  try { crypto = require('crypto') } catch (err) { }
  if (crypto) {
    var cryptoFunctions = ['pbkdf2', 'randomBytes']
    if (!v11plus) cryptoFunctions.push('pseudoRandomBytes')
    massWrap(
      crypto,
      cryptoFunctions,
      activator
    )
  }

  var instrumentPromise = isNative(global.Promise)

  // In case it's a non-native Promise, but bind have been used so it
  // looks native. There's still a potential false positive if the
  // non-native Promise library have a `name` property set to "Promise".
  // But worst case, the non-native Promise library will be instrumented
  // twice.
  instrumentPromise = instrumentPromise && global.Promise.name === 'Promise'

  /*
   * Native promises use the microtask queue to make all callbacks run
   * asynchronously to avoid Zalgo issues. Since the microtask queue is not
   * exposed externally, promises need to be modified in a fairly invasive and
   * complex way.
   *
   * The async boundary in promises that must be patched is between the
   * fulfillment of the promise and the execution of any callback that is waiting
   * for that fulfillment to happen. This means that we need to trigger a create
   * when resolve or reject is called and trigger before, after and error handlers
   * around the callback execution. There may be multiple callbacks for each
   * fulfilled promise, so handlers will behave similar to setInterval where
   * there may be multiple before after and error calls for each create call.
   *
   * async-listener monkeypatching has one basic entry point: `wrapCallback`.
   * `wrapCallback` should be called when create should be triggered and be
   * passed a function to wrap, which will execute the body of the async work.
   * The resolve and reject calls can be modified fairly easily to call
   * `wrapCallback`, but at the time of resolve and reject all the work to be done
   * on fulfillment may not be defined, since a call to then, chain or fetch can
   * be made even after the promise has been fulfilled. To get around this, we
   * create a placeholder function which will call a function passed into it,
   * since the call to the main work is being made from within the wrapped
   * function, async-listener will work correctly.
   *
   * There is another complication with monkeypatching Promises. Calls to then,
   * chain and catch each create new Promises that are fulfilled internally in
   * different ways depending on the return value of the callback. When the
   * callback return a Promise, the new Promise is resolved asynchronously after
   * the returned Promise has been also been resolved. When something other than
   * a promise is resolved the resolve call for the new Promise is put in the
   * microtask queue and asynchronously resolved.
   *
   * Then must be wrapped so that its returned promise has a wrapper that can be
   * used to invoke further continuations. This wrapper cannot be created until
   * after the callback has run, since the callback may return either a promise
   * or another value. Fortunately we already have a wrapper function around the
   * callback we can use (the wrapper created by resolve or reject).
   *
   * By adding an additional argument to this wrapper, we can pass in the
   * returned promise so it can have its own wrapper appended. the wrapper
   * function can the call the callback, and take action based on the return
   * value. If a promise is returned, the new Promise can proxy the returned
   * Promise's wrapper (this wrapper may not exist yet, but will by the time the
   * wrapper needs to be invoked). Otherwise, a new wrapper can be create the
   * same way as in resolve and reject. Since this wrapper is created
   * synchronously within another wrapper, it will properly appear as a
   * continuation from within the callback.
   */

  if (instrumentPromise) {
    wrapPromise()
  }

  function wrapPromise () {
    var Promise = global.Promise

    wrap(Promise.prototype, 'then', wrapThen)
    // Node.js <v7 only, alias for .then
    if (Promise.prototype.chain) {
      wrap(Promise.prototype, 'chain', wrapThen)
    }

    function wrapThen (original) {
      return function wrappedThen () {
        var promise = this
        var next = original.apply(promise, Array.prototype.map.call(arguments, bind))

        return next

        // wrap callbacks (success, error) so that the callbacks will be called as a
        // continuations of the resolve or reject call using the __ob_wrapper created above.
        function bind (fn) {
          if (typeof fn !== 'function') return fn
          return ins.bindFunction(fn)
        }
      }
    }
  }

  // Shim activator for functions that have callback last
  function activator (fn) {
    var wrapper
    var fallback = function () {
      var args
      var cbIdx = arguments.length - 1
      if (typeof arguments[cbIdx] === 'function') {
        args = Array(arguments.length)
        for (var i = 0; i < arguments.length - 1; i++) {
          args[i] = arguments[i]
        }
        args[cbIdx] = ins.bindFunction(arguments[cbIdx])
      }
      return fn.apply(this, args || arguments)
    }

    // Preserve function length for small arg count functions.
    switch (fn.length) {
      case 1:
        wrapper = function (cb) {
          if (arguments.length !== 1) return fallback.apply(this, arguments)
          if (typeof cb === 'function') cb = ins.bindFunction(cb)
          return fn.call(this, cb)
        }
        break
      case 2:
        wrapper = function (a, cb) {
          if (arguments.length !== 2) return fallback.apply(this, arguments)
          if (typeof cb === 'function') cb = ins.bindFunction(cb)
          return fn.call(this, a, cb)
        }
        break
      case 3:
        wrapper = function (a, b, cb) {
          if (arguments.length !== 3) return fallback.apply(this, arguments)
          if (typeof cb === 'function') cb = ins.bindFunction(cb)
          return fn.call(this, a, b, cb)
        }
        break
      case 4:
        wrapper = function (a, b, c, cb) {
          if (arguments.length !== 4) return fallback.apply(this, arguments)
          if (typeof cb === 'function') cb = ins.bindFunction(cb)
          return fn.call(this, a, b, c, cb)
        }
        break
      case 5:
        wrapper = function (a, b, c, d, cb) {
          if (arguments.length !== 5) return fallback.apply(this, arguments)
          if (typeof cb === 'function') cb = ins.bindFunction(cb)
          return fn.call(this, a, b, c, d, cb)
        }
        break
      case 6:
        wrapper = function (a, b, c, d, e, cb) {
          if (arguments.length !== 6) return fallback.apply(this, arguments)
          if (typeof cb === 'function') cb = ins.bindFunction(cb)
          return fn.call(this, a, b, c, d, e, cb)
        }
        break
      default:
        wrapper = fallback
    }

    if (promisify.custom in fn) {
      wrapper[promisify.custom] = fn[promisify.custom]
    }

    return wrapper
  }

  // Shim activator for functions that have callback first
  function activatorFirst (fn) {
    var wrapper
    var fallback = function () {
      var args
      if (typeof arguments[0] === 'function') {
        args = Array(arguments.length)
        args[0] = ins.bindFunction(arguments[0])
        for (var i = 1; i < arguments.length; i++) {
          args[i] = arguments[i]
        }
      }
      return fn.apply(this, args || arguments)
    }

    // Preserve function length for small arg count functions.
    switch (fn.length) {
      case 1:
        wrapper = function (cb) {
          if (arguments.length !== 1) return fallback.apply(this, arguments)
          if (typeof cb === 'function') cb = ins.bindFunction(cb)
          return fn.call(this, cb)
        }
        break
      case 2:
        wrapper = function (cb, a) {
          if (arguments.length !== 2) return fallback.apply(this, arguments)
          if (typeof cb === 'function') cb = ins.bindFunction(cb)
          return fn.call(this, cb, a)
        }
        break
      case 3:
        wrapper = function (cb, a, b) {
          if (arguments.length !== 3) return fallback.apply(this, arguments)
          if (typeof cb === 'function') cb = ins.bindFunction(cb)
          return fn.call(this, cb, a, b)
        }
        break
      case 4:
        wrapper = function (cb, a, b, c) {
          if (arguments.length !== 4) return fallback.apply(this, arguments)
          if (typeof cb === 'function') cb = ins.bindFunction(cb)
          return fn.call(this, cb, a, b, c)
        }
        break
      case 5:
        wrapper = function (cb, a, b, c, d) {
          if (arguments.length !== 5) return fallback.apply(this, arguments)
          if (typeof cb === 'function') cb = ins.bindFunction(cb)
          return fn.call(this, cb, a, b, c, d)
        }
        break
      case 6:
        wrapper = function (cb, a, b, c, d, e) {
          if (arguments.length !== 6) return fallback.apply(this, arguments)
          if (typeof cb === 'function') cb = ins.bindFunction(cb)
          return fn.call(this, cb, a, b, c, d, e)
        }
        break
      default:
        wrapper = fallback
    }

    if (promisify.custom in fn) {
      wrapper[promisify.custom] = fn[promisify.custom]
    }

    return wrapper
  }
}

// taken from node master on 2017/03/09
function toNumber (x) {
  return (x = Number(x)) >= 0 ? x : false
}

// taken from node master on 2017/03/09
function isPipeName (s) {
  return typeof s === 'string' && toNumber(s) === false
}
