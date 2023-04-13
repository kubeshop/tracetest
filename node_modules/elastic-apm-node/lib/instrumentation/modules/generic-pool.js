/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

var semver = require('semver')

var shimmer = require('../shimmer')

module.exports = function (generic, agent, { version }) {
  if (semver.satisfies(version, '^2.0.0')) {
    agent.logger.debug('shimming generic-pool.Pool')
    shimmer.wrap(generic, 'Pool', function (orig) {
      return function wrappedPool () {
        var trans = agent._instrumentation.currTransaction()
        var id = trans && trans.id
        agent.logger.debug('intercepted call to generic-pool.Pool %o', { id: id })

        var pool
        if (this instanceof generic.Pool) {
          var args = [].slice.call(arguments)
          args.unshift(null)
          pool = new (Function.prototype.bind.apply(orig, args))()
        } else {
          pool = orig.apply(this, arguments)
        }

        shimmer.wrap(pool, 'acquire', function (orig) {
          return function wrappedAcquire () {
            var trans = agent._instrumentation.currTransaction()
            var id = trans && trans.id
            agent.logger.debug('intercepted call to pool.acquire %o', { id: id })

            var cb = arguments[0]
            if (typeof cb === 'function') {
              arguments[0] = agent._instrumentation.bindFunction(cb)
            }

            return orig.apply(this, arguments)
          }
        })

        return pool
      }
    })
  } else if (semver.satisfies(version, '^3.1.0') && generic.PriorityQueue) {
    // A work-around as an alternative patching the returned promise from the
    // acquire function, we instead patch its resolve and reject functions.
    //
    // We can do that because they are exposed to the PriorityQueue when
    // enqueuing a ResourceRequest:
    //
    // https://github.com/coopernurse/node-pool/blob/58c275c5146977192165f679e86950396be1b9f1/lib/Pool.js#L404
    agent.logger.debug('shimming generic-pool.PriorityQueue.prototype.enqueue')
    shimmer.wrap(generic.PriorityQueue.prototype, 'enqueue', function (orig) {
      return function wrappedEnqueue () {
        var trans = agent._instrumentation.currTransaction()
        var id = trans && trans.id
        agent.logger.debug('intercepted call to generic-pool.PriorityQueue.prototype.enqueue %o', { id: id })

        var obj = arguments[0]
        // Expect obj to of type Deferred
        if (obj._resolve && obj._reject) {
          obj._resolve = agent._instrumentation.bindFunction(obj._resolve)
          obj._reject = agent._instrumentation.bindFunction(obj._reject)
        }

        return orig.apply(this, arguments)
      }
    })
  } else {
    agent.logger.debug('generic-pool version %s not supported - aborting...', version)
  }

  return generic
}
