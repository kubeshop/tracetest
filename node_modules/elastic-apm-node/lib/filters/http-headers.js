/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

const REDACTED = require('../constants').REDACTED

const cookie = require('cookie')
const redact = require('../redact-secrets')(REDACTED)
const SetCookie = require('set-cookie-serde')

module.exports = httpHeaders

function httpHeaders (obj) {
  const requestHeaders = obj.context && obj.context.request && obj.context.request.headers
  const responseHeaders = obj.context && obj.context.response && obj.context.response.headers

  if (requestHeaders) filterCookieHeaders(requestHeaders)
  if (responseHeaders) filterCookieHeaders(responseHeaders)

  return obj
}

/**
 * Filters cookie _values_ in http headers
 *
 * The filterCookieHeaders method filters individual
 * cookie values
 */
function filterCookieHeaders (headers) {
  for (const key in headers) {
    // undefined headers will be dropped when serialized as
    // json, and don't need to be redacted.  If a cookie
    // header is already redacted there's no need to parse its
    // values.
    if (headers[key] === undefined || REDACTED === headers[key]) {
      continue
    }
    switch (key.toLowerCase()) {
      case 'cookie':
        if (typeof headers[key] === 'string') {
          const cookies = cookie.parse(headers[key])
          redact.forEach(cookies)
          headers[key] = Object.keys(cookies)
            .map(k => k + '=' + cookies[k])
            .join('; ')
        }
        break
      case 'set-cookie':
        // if the sanitize_field_names module has redacted this there's
        // no need to attempt the sanitizaion of individual cookies
        try {
          const setCookies = new SetCookie(headers[key])
          redact.forEach(setCookies)
          headers[key] = stringify(setCookies)
        } catch (err) {
          // Ignore error
          headers[key] = '[malformed set-cookie header]'
        }

        break
    }
  }
}

function stringify (value) {
  return Array.isArray(value)
    ? value.map(value => value.toString())
    : value.toString()
}
