/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

var shimmer = require('../shimmer')
var templateShared = require('../template-shared')

module.exports = function (pug, agent, { enabled }) {
  if (!enabled) return pug
  agent.logger.debug('shimming pug.compile')
  shimmer.wrap(pug, 'compile', templateShared.wrapCompile(agent, 'pug'))

  return pug
}
