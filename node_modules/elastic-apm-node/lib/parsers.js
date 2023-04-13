/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

var url = require('url')

var basicAuth = require('basic-auth')
var getUrlFromRequest = require('original-url')
var parseHttpHeadersFromReqOrRes = require('http-headers')
var stringify = require('fast-safe-stringify')

const {
  redactKeysFromObject,
  redactKeysFromPostedFormVariables
} = require('./filters/sanitize-field-names')

function getContextFromRequest (req, conf, type) {
  var captureBody = conf.captureBody === type || conf.captureBody === 'all'

  var context = {
    http_version: req.httpVersion,
    method: req.method,
    url: getUrlFromRequest(req),
    headers: undefined
  }
  if (req.socket && req.socket.remoteAddress) {
    context.socket = {
      remote_address: req.socket.remoteAddress
    }
  }

  if (conf.captureHeaders) {
    context.headers = redactKeysFromObject(
      Object.assign({}, req.headers),
      conf.sanitizeFieldNamesRegExp
    )
  }

  var contentLength = parseInt(req.headers['content-length'], 10)
  var transferEncoding = req.headers['transfer-encoding']
  var chunked = typeof transferEncoding === 'string' && transferEncoding.toLowerCase() === 'chunked'
  var body = req.json || req.body || req.payload
  var haveBody = body && (chunked || contentLength > 0)

  if (haveBody) {
    if (!captureBody) {
      context.body = '[REDACTED]'
    } else if (Buffer.isBuffer(body)) {
      context.body = '<Buffer>'
    } else {
      body = redactKeysFromPostedFormVariables(body, req.headers, conf.sanitizeFieldNamesRegExp)
      if (typeof body !== 'string') {
        body = tryJsonStringify(body) || stringify(body)
      }
      context.body = body
    }
  }

  // TODO: Tempoary fix for https://github.com/elastic/apm-agent-nodejs/issues/813
  if (context.url && context.url.port) {
    context.url.port = String(context.url.port)
  }

  return context
}

function getContextFromResponse (res, conf, isError) {
  var context = {
    status_code: res.statusCode,
    headers: undefined
  }

  if (conf.captureHeaders) {
    context.headers = res.headers || parseHttpHeadersFromReqOrRes(res, true)
    context.headers = redactKeysFromObject(context.headers, conf.sanitizeFieldNamesRegExp)
  }

  if (isError) {
    context.headers_sent = res.headersSent
    context.finished = res.finished
  }

  return context
}

function getUserContextFromRequest (req) {
  var user = req.user || basicAuth(req) || req.session
  if (!user) return

  var context = {}

  if (typeof user.id === 'string' || typeof user.id === 'number') {
    context.id = user.id
  } else if (typeof user._id === 'string' || typeof user._id === 'number') {
    context.id = user._id
  }

  if (typeof user.username === 'string') {
    context.username = user.username
  } else if (typeof user.name === 'string') {
    context.username = user.name
  }

  if (typeof user.email === 'string') {
    context.email = user.email
  }

  return context
}

function parseUrl (urlStr) {
  return new url.URL(urlStr, 'relative:///')
}

function tryJsonStringify (obj) {
  try {
    return JSON.stringify(obj)
  } catch (e) {}
}

module.exports = {
  getContextFromRequest,
  getContextFromResponse,
  getUserContextFromRequest,
  parseUrl
}
