'use strict'

module.exports = function sync () {
  const nanoStart = process.hrtime()
  const microStart = Date.now() * 1000

  return function microtime () {
    const diff = process.hrtime(nanoStart)
    const microRemainder = diff[1] / 1000 | 0 // Use bitwise OR to remove the decimals
    const microDelta = diff[0] * 1e6 + microRemainder
    return microStart + microDelta
  }
}
