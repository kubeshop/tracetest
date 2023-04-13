const entries = require('object.entries')

function defaultOnArray () { return [] }
function defaultOnObject () { return {} }

function targetFor (source, key, fieldPath, isNew, {
  onArray = defaultOnArray,
  onObject = defaultOnObject
} = {}) {
  if (Array.isArray(source)) {
    return onArray(source, key, fieldPath, isNew)
  } else if (source !== null && typeof source === 'object') {
    return onObject(source, key, fieldPath, isNew)
  }
}

module.exports = function breadthFilter (root, opts = {}) {
  const { onValue } = opts
  const target = targetFor(root, null, [], true, opts)
  if (!target) return root

  const queue = [[ root, target, [] ]]
  const seen = new Set([ root ])
  let item

  while (item = queue.shift()) {
    const [ source, target, path ] = item
    for (const [ key, value ] of entries(source)) {
      const fieldPath = path.concat(key)
      const isNew = !seen.has(value)
      if (isNew) seen.add(value)

      const newTarget = targetFor(value, key, fieldPath, isNew, opts)
      if (newTarget) {
        target[key] = newTarget
        if (isNew) {
          queue.push([ value, target[key], fieldPath ])
        }
      } else {
        target[key] = onValue(value, key, fieldPath)
      }
    }
  }

  return target
}
