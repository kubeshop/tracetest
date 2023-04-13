/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

// Shared functionality between the instrumentations of:
// - elasticsearch - the legacy Elasticsearch JS client
// - @elastic/elasticsearch - the new Elasticsearch JS client

// Only capture the ES request body if the request path matches the
// `elasticsearchCaptureBodyUrls` config.
function shouldCaptureBody (path, elasticsearchCaptureBodyUrlsRegExp) {
  if (!path) {
    return false
  }
  for (var i = 0; i < elasticsearchCaptureBodyUrlsRegExp.length; i++) {
    const re = elasticsearchCaptureBodyUrlsRegExp[i]
    if (re.test(path)) {
      return true
    }
  }
  return false
}

/**
 * Get an appropriate `span.context.db.statement` for this ES client request, if any.
 * https://github.com/elastic/apm/blob/main/specs/agents/tracing-instrumentation-db.md#elasticsearch_capture_body_urls-configuration
 *
 * @param {string | null} path
 * @param {string | null} body
 * @param {RegExp[]} elasticsearchCaptureBodyUrlsRegExp
 * @return {string | undefined}
 */
function getElasticsearchDbStatement (path, body, elasticsearchCaptureBodyUrlsRegExp) {
  if (body && shouldCaptureBody(path, elasticsearchCaptureBodyUrlsRegExp)) {
    if (typeof (body) === 'string') {
      return body
    } else if (Buffer.isBuffer(body) || typeof body.pipe === 'function') {
      // Never serialize a Buffer or a Readable. These guards mirror
      // `shouldSerialize()` in the ES client, e.g.:
      // https://github.com/elastic/elastic-transport-js/blob/069172506d1fcd544b23747d8c2d497bab053038/src/Transport.ts#L614-L618
    } else if (Array.isArray(body)) {
      try {
        return body.map(JSON.stringify).join('\n') + '\n' // ndjson
      } catch (_ignoredErr) {}
    } else if (typeof (body) === 'object') {
      try {
        return JSON.stringify(body)
      } catch (_ignoredErr) {}
    }
  }
}

module.exports = {
  getElasticsearchDbStatement
}
