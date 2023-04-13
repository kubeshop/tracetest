/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

/**
 * This file is extracted from the 'shimmer' project copyright by Forrest L
 * Norvell. It have been modified slightly to be used in the current context.
 *
 * https://github.com/othiym23/shimmer
 *
 * Original file:
 *
 * https://github.com/othiym23/shimmer/blob/master/index.js
 *
 * License:
 *
 * BSD-2-Clause, http://opensource.org/licenses/BSD-2-Clause
 */

var symbols = require('../symbols')

var isWrappedSym = Symbol('elasticAPMIsWrapped')

exports.wrap = wrap
exports.massWrap = massWrap
exports.unwrap = unwrap

// Do not load agent until used to avoid circular dependency issues.
var _agent
function logger () {
  if (!_agent) _agent = require('../../')
  return _agent.logger
}

function isFunction (funktion) {
  return typeof funktion === 'function'
}

function wrap (nodule, name, wrapper) {
  if (!nodule || !nodule[name]) {
    logger().debug('no original function %s to wrap', name)
    return
  }

  if (!wrapper) {
    logger().debug({ stack: new Error().stack }, 'no wrapper function')
    return
  }

  if (!isFunction(nodule[name]) || !isFunction(wrapper)) {
    logger().debug('original object and wrapper must be functions')
    return
  }

  if (nodule[name][isWrappedSym]) {
    logger().debug('function %s already wrapped', name)
    return
  }

  var desc = Object.getOwnPropertyDescriptor(nodule, name)
  if (desc && !desc.writable) {
    logger().debug('function %s is not writable', name)
    return
  }

  var original = nodule[name]
  var wrapped = wrapper(original, name)

  wrapped[isWrappedSym] = true
  wrapped[symbols.unwrap] = function elasticAPMUnwrap () {
    if (nodule[name] === wrapped) {
      nodule[name] = original
      wrapped[isWrappedSym] = false
    }
  }

  nodule[name] = wrapped

  return wrapped
}

function massWrap (nodules, names, wrapper) {
  if (!nodules) {
    logger().debug({ stack: new Error().stack }, 'must provide one or more modules to patch')
    return
  } else if (!Array.isArray(nodules)) {
    nodules = [nodules]
  }

  if (!(names && Array.isArray(names))) {
    logger().debug('must provide one or more functions to wrap on modules')
    return
  }

  for (const nodule of nodules) {
    for (const name of names) {
      wrap(nodule, name, wrapper)
    }
  }
}

function unwrap (nodule, name) {
  if (!nodule || !nodule[name]) {
    logger().debug({ stack: new Error().stack }, 'no function to unwrap')
    return
  }

  if (!nodule[name][symbols.unwrap]) {
    logger().debug('no original to unwrap to -- has %s already been unwrapped?', name)
  } else {
    return nodule[name][symbols.unwrap]()
  }
}
