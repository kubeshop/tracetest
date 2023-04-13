# relative-microtime

This module returns the number of microseconds elapsed since January 1,
1970 00:00:00 UTC.

Note: The returned number of microseconds are locked relative to a
specific time and is therefore not subject to clock drift. This is
contrary to how `Date.now()` operates, whos return value (in
milliseconds) is always based on the system clock.

[![npm](https://img.shields.io/npm/v/relative-microtime.svg)](https://www.npmjs.com/package/relative-microtime)
[![build status](https://travis-ci.org/watson/relative-microtime.svg?branch=master)](https://travis-ci.org/watson/relative-microtime)
[![js-standard-style](https://img.shields.io/badge/code%20style-standard-brightgreen.svg?style=flat)](https://github.com/feross/standard)

## Installation

```
npm install relative-microtime --save
```

## Usage

```js
const timer = require('relative-microtime')()

const milli = Date.now()
const micro = timer()

console.log(milli) // 1537446194859
console.log(micro) // 1537446194859132
```

## Gotcha

## API

### `timer = microtime()`

Sync the microtime clock to the UTC wall clock.

Returns a `timer` function synced with the UTC time when `microtime()`
was called. Any clock drift that happens after the call to `microtime()`
will not influence the returned `timer`.

### `timer()`

Returns a `Number` representing the microseconds elapsed since the UNIX
epoch.

## License

[MIT](LICENSE)
