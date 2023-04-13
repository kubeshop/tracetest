/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

// This is a copy of load-source-map@2.0.0. It is inlined in elastic-apm-node
// solely to update its `source-map` dependency to fix
//   https://github.com/elastic/apm-agent-nodejs/issues/2589
// If/when there is a new release of load-source-map with
//   https://github.com/rexxars/load-source-map/pull/6
// then we could move back to using load-source-map as a dependency.

var fs = require('fs')
var path = require('path')
var SourceMapConsumer = require('source-map').SourceMapConsumer

var INLINE_SOURCEMAP_REGEX = /^data:application\/json[^,]+base64,/
var SOURCEMAP_REGEX = /(?:\/\/[@#][ \t]+sourceMappingURL=([^\s'"]+?)[ \t]*$)|(?:\/\*[@#][ \t]+sourceMappingURL=([^*]+?)[ \t]*(?:\*\/)[ \t]*$)/
var READ_FILE_OPTS = { encoding: 'utf8' }

module.exports = function readSourceMap (filename, cb) {
  fs.readFile(filename, READ_FILE_OPTS, function (err, sourceFile) {
    if (err) {
      return cb(err)
    }

    // Look for a sourcemap URL
    var sourceMapUrl = resolveSourceMapUrl(sourceFile, path.dirname(filename))
    if (!sourceMapUrl) {
      return cb()
    }

    // If it's an inline map, decode it and pass it through the same consumer factory
    if (isInlineMap(sourceMapUrl)) {
      return onMapRead(null, decodeInlineMap(sourceMapUrl))
    }

    // Load actual source map from given path
    fs.readFile(sourceMapUrl, READ_FILE_OPTS, onMapRead)

    function onMapRead (readErr, sourceMap) {
      if (readErr) {
        readErr.message = 'Error reading sourcemap for file "' + filename + '":\n' + readErr.message
        return cb(readErr)
      }

      try {
        (new SourceMapConsumer(sourceMap))
          .then(function onConsumerReady (consumer) {
            return cb(null, consumer)
          }, onConsumerError)
      } catch (parseErr) {
        onConsumerError(parseErr)
      }
    }

    function onConsumerError (parseErr) {
      parseErr.message = 'Error parsing sourcemap for file "' + filename + '":\n' + parseErr.message
      return cb(parseErr)
    }
  })
}

function resolveSourceMapUrl (sourceFile, sourcePath) {
  var lines = sourceFile.split(/\r?\n/)
  var sourceMapUrl = null
  for (var i = lines.length - 1; i >= 0 && !sourceMapUrl; i--) {
    sourceMapUrl = lines[i].match(SOURCEMAP_REGEX)
  }

  if (!sourceMapUrl) {
    return null
  }

  return isInlineMap(sourceMapUrl[1])
    ? sourceMapUrl[1]
    : path.resolve(sourcePath, sourceMapUrl[1])
}

function isInlineMap (url) {
  return INLINE_SOURCEMAP_REGEX.test(url)
}

function decodeInlineMap (data) {
  var rawData = data.slice(data.indexOf(',') + 1)
  return Buffer.from(rawData, 'base64').toString()
}
