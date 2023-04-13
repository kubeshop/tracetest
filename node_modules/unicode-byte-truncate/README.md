# unicode-byte-truncate

Truncate a string to a given byte size by removing bytes from the right
while making sure not to slice in the middle of a multi-byte unicode
character.

[![Build status](https://travis-ci.org/watson/unicode-byte-truncate.svg?branch=master)](https://travis-ci.org/watson/unicode-byte-truncate)
[![js-standard-style](https://img.shields.io/badge/code%20style-standard-brightgreen.svg?style=flat)](https://github.com/feross/standard)

## Installation

```
npm install unicode-byte-truncate --save
```

## Usage

```js
var trunc = require('unicode-byte-truncate')

var str = 'foo🎉bar' // 10 byte string - byte 4 to 7 is a single character

console.log(trunc(str, 4)) // `foo` == 0x666F6F (3 bytes)
console.log(trunc(str, 5)) // `foo` == 0x666F6F (3 bytes)
console.log(trunc(str, 6)) // `foo` == 0x666F6F (3 bytes)
console.log(trunc(str, 7)) // `foo🎉` == 0x666F6FF09F8E89 (7 bytes)
```

## API

The unicode-byte-truncate module exposes a single `trunc` function.

```js
result = trunc(string, maxBytes)
```

Given a `string` and a `maxBytes` integer greater than or equal to zero,
the `trunc` function will slice characters off the end of the string to
ensure that it doesn't contain more bytes than specified by the
`maxBytes` argument.

The truncated string will be returned as the `result`.

The `trunc` function is multi-byte unicode aware and will never cut up
surrogate pairs. This means that the `result` _may_ contain fewer bytes
than specified by the `maxBytes` argument.

## License

MIT
