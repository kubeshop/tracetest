/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

// Instrumentation of fastify.
// https://www.fastify.io/docs/latest/LTS/

const semver = require('semver')

module.exports = function (fastify, agent, { version, enabled }) {
  if (!enabled) return fastify

  agent.setFramework({ name: 'fastify', version, overwrite: false })

  agent.logger.debug('wrapping fastify build function')

  if (semver.gte(version, '3.0.0')) {
    wrappedBuild3Plus.default = wrappedBuild3Plus
    wrappedBuild3Plus.fastify = wrappedBuild3Plus
    return wrappedBuild3Plus
  }

  if (semver.gte(version, '2.0.0-rc')) return wrappedBuild2Plus

  return wrappedBuildPre2

  function wrappedBuild3Plus () {
    const _fastify = fastify.apply(null, arguments)

    agent.logger.debug('adding onRequest hook to fastify')
    _fastify.addHook('onRequest', (req, reply, next) => {
      const method = req.routerMethod || req.raw.method // Fallback for fastify >3 <3.3.0
      const url = req.routerPath || reply.context.config.url // Fallback for fastify >3 <3.3.0
      const name = method + ' ' + url
      agent._instrumentation.setDefaultTransactionName(name)
      next()
    })

    agent.logger.debug('adding onError hook to fastify')
    _fastify.addHook('onError', (req, reply, err, next) => {
      agent.captureError(err, { request: req.raw })
      next()
    })

    return _fastify
  }

  function wrappedBuild2Plus () {
    const _fastify = fastify.apply(null, arguments)

    agent.logger.debug('adding onRequest hook to fastify')
    _fastify.addHook('onRequest', (req, reply, next) => {
      const context = reply.context
      const name = req.raw.method + ' ' + context.config.url
      agent._instrumentation.setDefaultTransactionName(name)
      next()
    })

    agent.logger.debug('adding onError hook to fastify')
    _fastify.addHook('onError', (req, reply, err, next) => {
      agent.captureError(err, { request: req.raw })
      next()
    })

    return _fastify
  }

  function wrappedBuildPre2 () {
    const _fastify = fastify.apply(null, arguments)

    agent.logger.debug('adding onRequest hook to fastify')
    _fastify.addHook('onRequest', (req, reply, next) => {
      const context = reply._context
      const name = req.method + ' ' + context.config.url
      agent._instrumentation.setDefaultTransactionName(name)
      next()
    })

    agent.logger.warn('Elastic APM cannot automaticaly capture errors on this verison of Fastify. Upgrade to version 2.0.0 or later.')

    return _fastify
  }
}
