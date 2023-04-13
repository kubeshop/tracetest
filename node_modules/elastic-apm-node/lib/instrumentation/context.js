/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

var { parseUrl } = require('../parsers')

// Get the port number including the default port for a protocols
function getPortNumber (port, protocol) {
  if (port === '') {
    port = protocol === 'http:' ? '80' : protocol === 'https:' ? '443' : ''
  }
  return port
}

exports.getHTTPDestination = function (url) {
  const { port, protocol, hostname } = parseUrl(url)
  const portNumber = getPortNumber(port, protocol)

  // If hostname begins with [ and ends with ], assume that it's an IPv6 address.
  // since address and port are recorded separately, we are recording the
  // info in canonical form without square brackets
  const ipv6Hostname =
    hostname[0] === '[' &&
    hostname[hostname.length - 1] === ']'

  const address = ipv6Hostname ? hostname.slice(1, -1) : hostname

  return {
    address,
    port: Number(portNumber)
  }
}

exports.getDBDestination = function (host, port) {
  const destination = {}
  let haveValues = false

  if (host) {
    destination.address = host
    haveValues = true
  }
  port = Number(port)
  if (port) {
    destination.port = port
    haveValues = true
  }

  if (haveValues) {
    return destination
  } else {
    return null
  }
}
