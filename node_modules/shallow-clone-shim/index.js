'use strict'

module.exports = function clone (obj, orig, shim = {}) {
  const descriptors = Object.getOwnPropertyDescriptors(orig)

  for (const name of Object.keys(shim)) {
    descriptors[name] = shim[name](descriptors[name])
  }

  return Object.defineProperties(obj, descriptors)
}
