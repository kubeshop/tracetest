# async-value

An `AsyncValue` represents a value that will exist in the future. This is
different from a `Promise` in that it has no concept of errors or how they are
handled.

Note that, unlike a `Promise`, this module is also internally synchronous. What
this means is that if a value has already been `set` a call to `get` will
trigger the callback immediately. Conversely, if any callbacks have already been
assigned with `get` before a `set` occurs they too will run immediately when
`set` is called. This behavior is intentional for performance reasons.

## Install

```sh
npm install async-value
```

## Usage

```js
var value = new AsyncValue()

value.get(name => {
  console.log(`hello, ${name}!`)
})

value.set('world')
```

---

### Copyright (c) 2018 Stephen Belanger
#### Licensed under MIT License

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
