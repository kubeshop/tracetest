const tap = require('tap')
const mapcap = require('./')

const cap = 5

tap.test('instance', t => {
  validate(t, mapcap(new Map(), cap))
})

tap.test('instance - least recently used', t => {
  validate(t, mapcap(new Map(), cap, true), true)
})

tap.test('class', t => {
  const CappedMap = mapcap(Map, cap)
  validate(t, new CappedMap())
})

tap.test('class - least recently used', t => {
  const CappedMap = mapcap(Map, cap, true)
  validate(t, new CappedMap(), true)
})

function validate (t, instance, lru) {
  t.ok(instance instanceof Map, 'is an instance of Map')
  t.equal(instance.size, 0, 'is empty')

  // Verify the size is capped
  for (let i = 0; i < cap; i++) {
    instance.set(`key #${i}`, 'data')
    t.equal(instance.size, i + 1, `has size of ${i + 1}`)
  }

  // Walk the map backwards to re-order least-recently-used state
  if (lru) {
    for (let i = cap; i; i--) {
      instance.get(`key #${i}`)
    }
  }

  instance.set('should rotate old key out', 'data')
  t.equal(instance.size, cap, `still has only ${cap} items`)

  if (lru) {
    t.notOk(instance.has(`key #${cap}`), 'does not have last item')
    t.ok(instance.has(`key #${cap - 1}`), 'has second last item')
  } else {
    t.notOk(instance.has('key #0'), 'does not have first item')
    t.ok(instance.has('key #1'), 'has second item')
  }

  t.end()
}
