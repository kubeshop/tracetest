# async-value-promise

An `AsyncValuePromise` represents a value that will exist in the future. It is
similar to a `Promise` but does not include error boundary zoning or forced
asynchrony, for performance reasons.

This module is internally synchronous. What this means is that if a value
promise has already had `resolve` or `reject` called, any calls to `then` will
trigger the corresponding callback immediately. Conversely, if any callbacks
have already been attached with `then` before a `resolve` or `reject` occurs
they too will run immediately. This is intentional for performance reasons.

## Install

```sh
npm install async-value-promise
```

## Usage

```js
var pass = new AsyncValuePromise()

pass.then(name => {
  console.log(`hello, ${name}!`)
})

pass.resolve('world')

var fail = new AsyncValuePromise()

fail.then(null, error => {
  console.error(err.stack)
})

fail.reject(new Error('oops'))
```

---

### Copyright (c) 2018 Stephen Belanger
#### Licensed under MIT License

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
