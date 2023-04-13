'use strict'

class NotAFunctionError extends Error {
  constructor () {
    super('filter must be a function')
  }
}

function validateItems (items, orNumber) {
  if (orNumber && items.length === 1 && typeof items[0] === 'number') {
    return
  }

  if (!Array.isArray(items)) {
    throw new Error('items is undefined')
  }

  for (let item of items) {
    if (typeof item !== 'function') {
      throw new NotAFunctionError()
    }
  }
}

class Filters extends Array {
  constructor (...items) {
    validateItems(items, true)
    super(...items)
  }

  static from (items) {
    validateItems(items)
    return super.from(items)
  }

  concat (...args) {
    const items = args.length > 1
      ? args
      : Array.isArray(args[0]) ? args[0] : [ args[0] ]

    validateItems(items)
    return super.concat(...args)
  }

  push (...items) {
    validateItems(items)
    return super.push(...items)
  }

  unshift (...items) {
    validateItems(items)
    return super.unshift(...items)
  }

  process (payload) {
    let result = payload

    // abort if a filter function doesn't return an object
    this.some(filter => {
      result = filter(result)
      return !result
    })

    return result
  }
}

module.exports = Filters
