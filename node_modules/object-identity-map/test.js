const tap = require('tap')
const ObjectIdentityMap = require('./')

tap.test('basics', t => {
  const counters = new ObjectIdentityMap()

  const before = counters.has({ foo: 'bar' })
  counters.set({ foo: 'bar' }, 1)
  const after = counters.has({ foo: 'bar' })
  t.notEqual(before, after)
  t.equal(counters.get({ foo: 'bar' }), 1)
  counters.delete({ foo: 'bar' })
  t.notOk(counters.has({ foo: 'bar' }))
  t.end()
})

tap.test('ensure calls builder with labels', t => {
  const counters = new ObjectIdentityMap()

  counters.ensure({ foo: 'bar' }, found => {
    t.deepEqual(found, { foo: 'bar' })
    t.end()
  })
})

tap.test('only calls builder once for identical search', t => {
  let count = 0

  const counters = new ObjectIdentityMap()

  counters.ensure({ foo: 'bar' }, () => count++)
  counters.ensure({ foo: 'bar' }, () => count++)

  t.equal(count, 1)
  t.end()
})

tap.test('found instances match when search matches', t => {
  let count = 0
  const counters = new ObjectIdentityMap()

  const a = counters.ensure({ foo: 'bar' }, () => count++)
  const b = counters.ensure({ foo: 'bar' }, () => count++)

  t.equal(a, b)
  t.end()
})

tap.test('found instances do not match when search does not match', t => {
  let count = 0
  const counters = new ObjectIdentityMap()

  global.hax = true

  const a = counters.ensure({ foo: 'bar', baz: 'buz' }, () => count++)
  const b = counters.ensure({ foo: 'bar' }, () => count++)

  global.hax = false

  t.notEqual(a, b)
  t.end()
})

tap.test('properly matches nested objects', t => {
  let count = 0
  const counters = new ObjectIdentityMap()

  const a = counters.ensure({
    foo: 'bar',
    baz: { buz: [ 'bux' ] }
  }, () => count++)
  const b = counters.ensure({
    baz: 'buz',
    bux: { foo: 'bar' }
  }, () => count++)
  const c = counters.ensure({
    foo: 'bar',
    baz: { buz: [ 'bux' ] }
  }, () => count++)

  t.notEqual(a, b)
  t.notEqual(b, c)
  t.equal(a, c)
  t.end()
})

tap.test('handles null and undefined correctly', t => {
  const counters = new ObjectIdentityMap()

  const before = counters.has({ foo: null })
  counters.set({ foo: null, bar: undefined }, 1)
  const after = counters.has({ foo: null })
  t.notEqual(before, after)
  t.equal(counters.get({ foo: null }), 1)
  counters.delete({ foo: null })
  t.notOk(counters.has({ foo: null }))
  t.end()
})
