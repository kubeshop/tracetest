/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

var semver = require('semver')

var shimmer = require('../shimmer')

var onPreAuthSym = Symbol('ElasticAPMOnPreAuth')

module.exports = function (hapi, agent, { version, enabled }) {
  if (!enabled) return hapi

  agent.setFramework({ name: 'hapi', version, overwrite: false })

  if (!semver.satisfies(version, '>=9.0.0')) {
    agent.logger.debug('hapi version %s not supported - aborting...', version)
    return hapi
  }
  const isHapiGte17 = semver.satisfies(version, '>=17')

  agent.logger.debug('shimming hapi.Server.prototype.initialize')

  if (isHapiGte17) {
    shimmer.massWrap(hapi, ['Server', 'server'], function (orig) {
      return function (options) {
        var res = orig.apply(this, arguments)
        patchServer(res)
        return res
      }
    })
  } else {
    shimmer.wrap(hapi.Server.prototype, 'initialize', function (orig) {
      return function () {
        patchServer(this)
        return orig.apply(this, arguments)
      }
    })
  }

  function patchServer (server) {
    // Hooks that are always allowed
    if (typeof server.on === 'function') {
      attachEvents(server)
    } else if (typeof server.events.on === 'function') {
      attachEvents(server.events)
    } else {
      agent.logger.debug('unable to enable hapi error tracking')
    }

    // Prior to hapi 17, when the server has no connections we can't make
    // connection lifecycle hooks (in hapi 17+ the server always has
    // connections, though the `server.connections` property doesn't exists,
    // so this if-statement wont fire)
    var conns = server.connections
    if (conns && conns.length === 0) {
      agent.logger.debug('unable to enable hapi instrumentation on connectionless server')
      return
    }

    // Hooks that are only allowed when the hapi server has connections
    // (with hapi 17+ this is always the case)
    if (typeof server.ext === 'function') {
      server.ext('onPreAuth', onPreAuth)
      server.ext('onPreResponse', onPreResponse)
      if (agent._conf.captureBody !== 'off') {
        server.ext('onPostAuth', onPostAuth)
      }
    } else {
      agent.logger.debug('unable to enable automatic hapi transaction naming')
    }
  }

  function attachEvents (emitter) {
    if (!isHapiGte17) {
      emitter.on('request-error', function (request, error) {
        agent.captureError(error, {
          request: request.raw && request.raw.req
        })
      })
    }

    emitter.on('log', function (event, tags) {
      captureError('log', null, event, tags)
    })

    emitter.on('request', function (req, event, tags) {
      captureError('request', req, event, tags)
    })
  }

  function captureError (type, req, event, tags) {
    if (!event || !tags.error || event.channel === 'internal') {
      return
    }

    // TODO: Find better location to put this than custom
    var payload = {
      custom: {
        tags: event.tags,
        internals: event.internals,
        // Moved from data to error in hapi 17
        data: event.data || event.error
      },
      request: req && req.raw && req.raw.req
    }

    var err = payload.custom.data
    if (!(err instanceof Error) && typeof err !== 'string') {
      err = 'hapi server emitted a ' + type + ' event tagged error'
    }

    agent.captureError(err, payload)
  }

  function onPreAuth (request, reply) {
    agent.logger.debug('received hapi onPreAuth event')

    // Record the fact that the preAuth extension have been called. This
    // info is useful later to know if this is a CORS preflight request
    // that is automatically handled by hapi (as those will not trigger
    // the onPreAuth extention)
    request[onPreAuthSym] = true

    if (request.route) {
      // fingerprint was introduced in hapi 11 and is a little more
      // stable in case the param names change
      // - path example: /foo/{bar*2}
      // - fingerprint example: /foo/?/?
      var fingerprint = request.route.fingerprint || request.route.path

      if (fingerprint) {
        var name = (request.raw && request.raw.req && request.raw.req.method) ||
                   (request.route.method && request.route.method.toUpperCase())

        if (typeof name === 'string') {
          name = name + ' ' + fingerprint
        } else {
          name = fingerprint
        }

        agent._instrumentation.setDefaultTransactionName(name)
      }
    }

    return (isHapiGte17 ? reply.continue : reply.continue())
  }

  function onPostAuth (request, reply) {
    if (request.payload && request.raw && request.raw.req) {
      // Save the parsed req body to be picked up by getContextFromRequest().
      request.raw.req.payload = request.payload
    }
    return (isHapiGte17 ? reply.continue : reply.continue())
  }

  function onPreResponse (request, reply) {
    agent.logger.debug('received hapi onPreResponse event')

    // Detection of CORS preflight requests:
    // There is no easy way in hapi to get the matched route for a
    // CORS preflight request that matches any of the autogenerated
    // routes created by hapi when `cors: true`. The best solution is to
    // detect the request "fingerprint" using the magic if-sentence below
    // and group all those requests into on type of transaction
    if (!request[onPreAuthSym] &&
        request.route && request.route.path === '/{p*}' &&
        request.raw && request.raw.req && request.raw.req.method === 'OPTIONS' &&
        request.raw.req.headers['access-control-request-method']) {
      agent._instrumentation.setDefaultTransactionName('CORS preflight')
    }

    return (isHapiGte17 ? reply.continue : reply.continue())
  }

  return hapi
}
