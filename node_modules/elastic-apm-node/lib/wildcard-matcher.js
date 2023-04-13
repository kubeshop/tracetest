/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'
/**
 * Implements a Wildcard matcher
 *
 * Exports a function that implements a simple wildcard matcher
 * per: https://github.com/elastic/apm/issues/144
 */
const escapeStringRegexp = require('escape-string-regexp')

// Converts elastic-wildcard pattern into a
// a javascript regular expression.
const starMatchToRegex = (pattern) => {
  // case insensative by default
  let regexOpts = ['i']
  if (pattern.startsWith('(?-i)')) {
    regexOpts = []
    pattern = pattern.slice(5)
  }

  const patternLength = pattern.length
  const reChars = ['^']
  for (let i = 0; i < patternLength; i++) {
    const char = pattern[i]
    switch (char) {
      case '*':
        reChars.push('.*')
        break
      default:
        reChars.push(
          escapeStringRegexp(char)
        )
    }
  }
  reChars.push('$')
  return new RegExp(reChars.join(''), regexOpts.join(''))
}

class WildcardMatcher {
  compile (pattern) {
    return starMatchToRegex(pattern)
  }

  match (string, pattern) {
    const re = this.compile(pattern)
    return string.search(re) !== -1
  }
}
module.exports = { WildcardMatcher }
