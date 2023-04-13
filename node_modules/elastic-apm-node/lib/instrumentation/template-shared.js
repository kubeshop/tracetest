/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

exports.wrapCompile = function (agent, moduleName) {
  function wrapTemplate (original) {
    return function wrappedTemplate (data) {
      var span = agent.startSpan(moduleName, 'template', moduleName, 'render')
      var id = span && span.transaction.id

      agent.logger.debug('intercepted call to %s render %o', moduleName, {
        id: id,
        data: data
      })

      var ret = original.apply(this, arguments)
      if (span) span.end()

      return ret
    }
  }

  return function wrappedCompile (original) {
    return function wrappedCompile (input) {
      var span = agent.startSpan(moduleName, 'template', moduleName, 'compile')
      var id = span && span.transaction.id

      agent.logger.debug('intercepted call to %s compile %o', moduleName, {
        id: id,
        input: input
      })

      var ret = original.apply(this, arguments)
      if (span) span.end()

      return typeof ret === 'function' ? wrapTemplate(ret) : ret
    }
  }
}
