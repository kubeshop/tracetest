/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

// Use this module via `node --require=elastic-apm-node/start-next.js ...`
// to monitor a Next.js app with Elastic APM.

const apm = require('./').start()

// Flush APM data on server process termination.
// https://nextjs.org/docs/deployment#manual-graceful-shutdowns
// Note: Support for NEXT_MANUAL_SIG_HANDLE was added in next@12.1.7-canary.7,
// so this `apm.flush()` will only happen in that and later versions.
process.env.NEXT_MANUAL_SIG_HANDLE = 1
function flushApmAndExit () {
  apm.flush(() => {
    process.exit(0)
  })
}
process.on('SIGTERM', flushApmAndExit)
process.on('SIGINT', flushApmAndExit)

module.exports = apm
