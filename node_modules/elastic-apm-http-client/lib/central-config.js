'use strict'

// Central config-related utilities for the APM http client.

const INTERVAL_DEFAULT_S = 300 // 5 min
const INTERVAL_MIN_S = 5
const INTERVAL_MAX_S = 86400 // 1d

/**
 * Determine an appropriate delay until the next fetch of Central Config.
 * Default to 5 minutes, minimum 5s, max 1d.
 *
 * The maximum of 1d ensures we don't get surprised by an overflow value to
 * `setTimeout` per https://developer.mozilla.org/en-US/docs/Web/API/setTimeout#maximum_delay_value
 *
 * @param {Number|undefined} seconds - A number of seconds, typically pulled
 *    from a `Cache-Control: max-age=${seconds}` header on a previous central
 *    config request.
 * @returns {Number}
 */
function getCentralConfigIntervalS (seconds) {
  if (typeof seconds !== 'number' || isNaN(seconds) || seconds <= 0) {
    return INTERVAL_DEFAULT_S
  }
  return Math.min(Math.max(seconds, INTERVAL_MIN_S), INTERVAL_MAX_S)
}

module.exports = {
  getCentralConfigIntervalS,

  // These are exported for testing.
  INTERVAL_DEFAULT_S,
  INTERVAL_MIN_S,
  INTERVAL_MAX_S
}
