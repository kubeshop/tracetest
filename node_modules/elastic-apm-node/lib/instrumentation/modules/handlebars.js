/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

var shimmer = require('../shimmer')
var templateShared = require('../template-shared')

module.exports = function (handlebars, agent, { enabled }) {
  if (!enabled) return handlebars
  agent.logger.debug('shimming handlebars.compile')
  shimmer.wrap(handlebars, 'compile', templateShared.wrapCompile(agent, 'handlebars'))

  return handlebars
}
