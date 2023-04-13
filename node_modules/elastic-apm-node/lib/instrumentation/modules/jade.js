/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

var shimmer = require('../shimmer')
var templateShared = require('../template-shared')

module.exports = function (jade, agent, { enabled }) {
  if (!enabled) return jade
  agent.logger.debug('shimming jade.compile')
  shimmer.wrap(jade, 'compile', templateShared.wrapCompile(agent, 'jade'))

  return jade
}
