# mapcap

Cap your Maps with mapcap!

## Install

```sh
npm install mapcap
```

## Usage

```js
const assert = require('assert')
const mapcap = require('mapcap')
const CappedMap = mapcap(Map, 100)
const map = new CappedMap()

for (let i = 0; i < 1e8; i++) {
  map.set(i, i)
}

assert.equal(map.size, 100)
```

To drop "least recently used" items rather than "least recently inserted"
items, set the _third_ argument to `true`. Note that this will result in
an additional `delete(...)` and `set(...)` for every `get(...)` which
_may_ impact performance.

---

### Copyright (c) 2019 Stephen Belanger

#### Licensed under MIT License

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
