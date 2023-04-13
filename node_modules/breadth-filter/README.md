# breadth-filter

[![Greenkeeper badge](https://badges.greenkeeper.io/Qard/breadth-filter.svg)](https://greenkeeper.io/)

Apply a deep object filter via breadth traversal. It allows replacing values, including retargeting objects and arrays and even detecting circular references to allow the user to decide how to handle the circle. Note that the filter will _not_ descend _into_ circular objects but will only _reach_ them.

## Install

```sh
npm install breadth-filter
```

## Example

```js
const breadthFilter = require('breadth-filter')

const data = {
  user: {
    name: 'someone',
    password: 'THIS SHOULD BE SECRET!!'
  }
}

function onValue (value, key, path) {
  if (key === 'password') {
    console.log('redacted field at', path.join('.'))
    return '[redacted]'
  }
  return value
}

// build a new filtered object...

const filtered = breadthFilter(data, {
  onValue
})

// or, mutate in-place and break cycles...

breadthFilter(data, {
  onValue,
  onObject (value, key, path, isNew) {
    return isNew ? value : '[Circular]'
  },
  onArray (value, key, path, isNew) {
    return isNew ? value : '[Circular]'
  }
})
```

## Options

### onValue(value, key, path)

This function handles primitive value types such as numbers or strings. It can be used to filter things like sensitive data such as passwords or credit card numbers.

Arguments:

* `value` - any
The encountered value to filter.
* `key` - string | number
The property key of the encountered value, at the current depth level. Will be `null` when encountering the root object.
* `path` - array<string>
Any array of all property keys leading from the root to this encountered value.

### onObject(value, key, path, isNew)

This handles encountered objects. It can be used to filter entire objects out of the result, by returning `undefined`, it can enable in-place mutation by return the value directly, and the `isNew` property can be used to identify and break out of circular references, replacing them with something else such as `'[Circular]'`.

Arguments:

* `value` - any
The encountered value to filter.
* `key` - string | number
The property key of the encountered value, at the current depth level. Will be `null` when encountering the root object.
* `path` - array<string>
Any array of all property keys leading from the root to this encountered value.
* `isNew` - boolean
Indicates if this is the first time encountering this value, a false value indicates a circular reference.

### onArray(value, key, path, isNew)

This handles encountered arrays. It can be used to filter entire arrays out of the result, by returning `undefined`, it can enable in-place mutation by return the value directly, and the `isNew` property can be used to identify and break out of circular references, replacing them with something else such as `'[Circular]'`.

Arguments:

* `value` - any
The encountered value to filter.
* `key` - string | number
The property key of the encountered value, at the current depth level. Will be `null` when encountering the root object.
* `path` - array<string>
Any array of all property keys leading from the root to this encountered value.
* `isNew` - boolean
Indicates if this is the first time encountering this value, a false value indicates a circular reference.

---

### Copyright (c) 2019 Stephen Belanger
#### Licensed under MIT License

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
