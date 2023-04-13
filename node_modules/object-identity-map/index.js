const entries = require('object.entries')

function sortEntriesByKey ([ a ], [ b ]) {
  return a > b ? 1 : -1
}

function rebuildReducer (target, [ key, value ]) {
  target[key] = rebuildAsOrdered(value)
  return target
}

function rebuildAsOrdered (source) {
  let target

  if (Array.isArray(source)) {
    target = []
  } else if (source && typeof source === 'object') {
    target = {}
  } else {
    return source
  }

  return entries(source)
    .sort(sortEntriesByKey)
    .reduce(rebuildReducer, target)
}

function labelsToKey (labels) {
  return JSON.stringify(rebuildAsOrdered(labels))
}

class ObjectIdentityMap extends Map {
  has (labels) {
    return super.has(labelsToKey(labels))
  }

  get (labels) {
    return super.get(labelsToKey(labels))
  }

  set (labels, value) {
    return super.set(labelsToKey(labels), value)
  }

  delete (labels) {
    return super.delete(labelsToKey(labels))
  }

  ensure (labels, build) {
    const key = labelsToKey(labels)
    if (!super.has(key)) {
      super.set(key, build(labels))
    }
    return super.get(key)
  }
}

module.exports = ObjectIdentityMap
