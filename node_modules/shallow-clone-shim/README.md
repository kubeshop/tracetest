# shallow-clone-shim

Shallow clones an object, including non-enumerable properties, while
respecting the original data and accessor descriptors of the properties.
This means that for instances getters and setters are copied faithfully.
Optionally allows for shimming/overwriting properties by redefining or
manipulating existing property descriptors.

[![npm](https://img.shields.io/npm/v/shallow-clone-shim.svg)](https://www.npmjs.com/package/shallow-clone-shim)
[![build status](https://travis-ci.org/watson/shallow-clone-shim.svg?branch=master)](https://travis-ci.org/watson/shallow-clone-shim)
[![js-standard-style](https://img.shields.io/badge/code%20style-standard-brightgreen.svg?style=flat)](https://github.com/feross/standard)

## Installation

```
npm install shallow-clone-shim --save
```

## Usage

```js
const assert = require('assert')
const clone = require('shallow-clone-shim')

const original = Object.defineProperties({}, {
  foo: { // non-writable
    value: 1
  },
  bar: { // non-configurable
    enumerable: true,
    get: function get () {
      return 2
    }
  }
})

assert.strictEqual(original.foo, 1)
assert.strictEqual(original.bar, 2)

const copy = clone({}, original, {
  bar (descriptor) {
    // descriptor == Object.getOwnPropertyDescriptor(original, 'bar')
    const getter = descriptor.get
    descriptor.get = function get () {
      return getter() + 1
    }
    return descriptor
  }
})

assert.strictEqual(original.foo, 1)
assert.strictEqual(original.bar, 3)
```

## API

### `object = clone(object, original[, shim])`

Shallow copies all own properties of the `original` into `object`. Both
enumerable and non-enumerable properties are copied.

The `object` is also returned.

If the optional `shim` argument is supplied, it's expected to be an
object containing functions. The names of the `shim` object propeties is
expected to match the names of properties in the `original` object. Each
shim function is called with the property descriptor for that particular
property in `original`. The function is expected to return a valid
property descriptor as expected by [`Object.defineProperty()`]. The
returned desciptor will replace the original descriptor in the copied
object.

## License

[MIT](LICENSE)

[`Object.defineProperty()`]: https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Object/defineProperty
