const tap = require('tap')

const breadthFilter = require('./')

function reverse (input) {
  return input.split('').reverse().join('')
}

function identity (input) {
  return input
}

function whenKeyMatches (expected, fn) {
  return (value, key) => {
    return key === expected ? fn(value) : value
  }
}

function whenPathMatches (expected, fn) {
  return (value, key, path) => {
    return path.join('.') === expected ? fn(value) : value
  }
}

// Mutation handlers
function onObject (source) {
  return source
}
function onArray (source) {
  return source
}

tap.test('basic value mapping', (t) => {
  const input = {
    foo: {
      bar: {
        baz: 'buz'
      },
      bux: 'bax'
    }
  }

  const output = breadthFilter(input, {
    onValue: reverse
  })

  const expected = {
    foo: {
      bar: {
        baz: 'zub'
      },
      bux: 'xab'
    }
  }

  t.deepEqual(output, expected, 'matches expected output')
  t.end()
})

tap.test('receives key', (t) => {
  const input = {
    foo: {
      bar: 'baz'
    },
    bux: {
      bar: 'baz'
    },
    bax: {
      bex: 'baz'
    }
  }

  const output = breadthFilter(input, {
    onValue: whenKeyMatches('bar', reverse)
  })

  const expected = {
    foo: {
      bar: 'zab'
    },
    bux: {
      bar: 'zab'
    },
    bax: {
      bex: 'baz'
    }
  }

  t.deepEqual(output, expected, 'matches expected output')
  t.end()
})

tap.test('receives path', (t) => {
  const input = {
    foo: {
      bar: {
        baz: 'buz'
      },
      bux: 'bax'
    }
  }

  const output = breadthFilter(input, {
    onValue: whenPathMatches('foo.bar.baz', reverse)
  })

  const expected = {
    foo: {
      bar: {
        baz: 'zub'
      },
      bux: 'bax'
    }
  }

  t.deepEqual(output, expected, 'matches expected output')
  t.end()
})

tap.test('supports arrays', (t) => {
  const input = {
    foo: [
      { bar: 'baz' },
      { bax: 'baz' }
    ]
  }

  const output = breadthFilter(input, {
    onValue: whenPathMatches('foo.0.bar', reverse)
  })

  const expected = {
    foo: [
      { bar: 'zab' },
      { bax: 'baz' }
    ]
  }

  t.deepEqual(output, expected, 'matches expected output')
  t.end()
})

tap.test('does not mutate by default', (t) => {
  const input = {
    foo: {
      bar: {
        baz: 'buz'
      },
      bux: 'bax'
    }
  }

  breadthFilter(input, {
    onValue: reverse
  })

  const expected = {
    foo: {
      bar: {
        baz: 'buz'
      },
      bux: 'bax'
    }
  }

  t.deepEqual(input, expected, 'matches expected output')
  t.end()
})

tap.test('mutates in destructive mode', (t) => {
  const input = {
    foo: {
      bar: {
        baz: 'buz'
      },
      bux: ['bax']
    }
  }

  breadthFilter(input, {
    onValue: reverse,
    onObject,
    onArray
  })

  const expected = {
    foo: {
      bar: {
        baz: 'zub'
      },
      bux: ['xab']
    }
  }

  t.deepEqual(input, expected, 'matches expected output')
  t.end()
})

tap.test('gracefully handle circular references', (t) => {
  const input = {
    foo: {
      bar: {
        baz: 'buz'
      },
      bux: 'bax'
    }
  }

  // Form a circular reference
  input.foo.foo = input.foo
  input.foo.bar.root = input

  const output = breadthFilter(input, {
    onValue: reverse,
    onArray,
    onObject(source, key, path, isNew) {
      return isNew ? {} : '[Circular]'
    }
  })

  const expected = {
    foo: {
      bar: {
        baz: 'zub',
        root: '[Circular]'
      },
      bux: 'xab',
      foo: '[Circular]'
    }
  }

  t.deepEqual(output, expected, 'matches expected output')
  t.end()
})


tap.test('supports null and undefined values', (t) => {
  const input = {
    foo: {
      bar: {
        baz: null
      },
      bux: undefined
    }
  }

  breadthFilter(input, {
    onValue: identity,
    onObject,
    onArray
  })

  const expected = {
    foo: {
      bar: {
        baz: null
      },
      bux: undefined
    }
  }

  t.deepEqual(input, expected, 'matches expected output')
  t.end()
})
