/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

var isError = require('core-util-is').isError
var semver = require('semver')

var shimmer = require('../shimmer')
var symbols = require('../../symbols')

module.exports = function (express, agent, { version, enabled }) {
  if (!enabled) return express

  agent.setFramework({ name: 'express', version, overwrite: false })

  if (!semver.satisfies(version, '^4.0.0')) {
    agent.logger.debug('express version %s not supported - aborting...', version)
    return express
  }

  // express 5 moves the router methods onto a prototype
  var routerProto = semver.satisfies(version, '^5')
    ? (express.Router && express.Router.prototype)
    : express.Router

  var layerPatchedSymbol = Symbol('layer-patched')

  function shouldReport (err) {
    if (!agent._conf.captureExceptions) return false
    if (typeof err === 'string') return true
    if (isError(err) && !err[symbols.errorReportedSymbol]) {
      err[symbols.errorReportedSymbol] = true
      return true
    }
    return false
  }

  function safePush (obj, prop, value) {
    if (!obj[prop]) obj[prop] = []
    obj[prop].push(value)
  }

  function patchLayer (layer, layerPath) {
    if (!layer[layerPatchedSymbol]) {
      layer[layerPatchedSymbol] = true
      agent.logger.debug('shimming express.Router.Layer.handle function: %s', layer.name)
      shimmer.wrap(layer, 'handle', function (orig) {
        let handle

        if (orig.length !== 4) {
          handle = function (req, res, next) {
            if (!layer.route && layerPath && typeof next === 'function') {
              safePush(req, symbols.expressMountStack, layerPath)
              arguments[2] = function (nextArg) {
                // https://github.com/expressjs/express/blob/4.18.1/lib/router/route.js#L116-L149
                // The argument to an Express handler's `next()` can be:
                // falsey (call the next handler), 'route' (skip handlers for
                // this route), 'router' (skip handlers for this Router), or
                // any other value is considered an error (it doesn't have to
                // be an instance of Error). For all but the last one, Express
                // will consider other routes, so we want to pop the mount
                // stack.
                if (!nextArg || nextArg === 'route' || nextArg === 'router') {
                  req[symbols.expressMountStack].pop()
                }
                return next.apply(this, arguments)
              }
            }

            return orig.apply(this, arguments)
          }
        } else {
          handle = function (err, req, res, next) {
            if (shouldReport(err)) {
              agent.captureError(err, { request: req })
            }
            return orig.apply(this, arguments)
          }
        }

        for (const prop in orig) {
          if (Object.prototype.hasOwnProperty.call(orig, prop)) {
            handle[prop] = orig[prop]
          }
        }

        return handle
      })
    }
  }

  agent.logger.debug('shimming express.Router.use function')

  shimmer.wrap(routerProto, 'route', orig => {
    return function route (path) {
      var route = orig.apply(this, arguments)
      var layer = this.stack[this.stack.length - 1]
      patchLayer(layer, path)
      return route
    }
  })

  shimmer.wrap(routerProto, 'use', orig => {
    return function use (path) {
      var route = orig.apply(this, arguments)
      var layer = this.stack[this.stack.length - 1]
      patchLayer(layer, typeof path === 'string' && path)
      return route
    }
  })

  agent.logger.debug('shimming express.static function')

  shimmer.wrap(express, 'static', function wrapStatic (orig) {
    // By the time of this writing, Express adds a `mime` property to the
    // `static` function that needs to be copied to the wrapped function.
    // Instead of only copying the `mime` function, let's loop over all
    // properties in case new properties are added in later versions of
    // Express.
    for (const prop of Object.keys(orig)) {
      agent.logger.debug('copying property %s from express.static', prop)
      wrappedStatic[prop] = orig[prop]
    }

    return wrappedStatic

    function wrappedStatic () {
      var origServeStatic = orig.apply(this, arguments)
      return function serveStatic (req, res, next) {
        req[symbols.staticFile] = true

        return origServeStatic(req, res, nextHook)

        function nextHook (err) {
          if (!err) req[symbols.staticFile] = false
          return next.apply(this, arguments)
        }
      }
    }
  })

  return express
}
