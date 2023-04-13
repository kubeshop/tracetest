/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

module.exports = function (koa, agent, { version, enabled }) {
  if (!enabled) return koa

  agent.setFramework({ name: 'koa', version, overwrite: false })

  return koa
}
