# fast-stream-to-buffer

Consume a stream of data into a binary Buffer as efficiently as
possible.

[![Build status](https://travis-ci.org/watson/fast-stream-to-buffer.svg?branch=master)](https://travis-ci.org/watson/fast-stream-to-buffer)
[![js-standard-style](https://img.shields.io/badge/code%20style-standard-brightgreen.svg?style=flat)](https://github.com/feross/standard)

## Installation

```
npm install fast-stream-to-buffer --save
```

## Usage

Process an abitrary readable stream:

```js
const streamToBuffer = require('fast-stream-to-buffer')

streamToBuffer(stream, function (err, buf) {
  if (err) throw err
  console.log(buf.toString())
})
```

Or use the `onStream()` helper function:

```js
const http = require('http')
const streamToBuffer = require('fast-stream-to-buffer')

// `http.get` expects a callback as the 2nd argument that will be called
// with a readable stream of the response
http.get('http://example.com', streamToBuffer.onStream(function (err, buf) {
  if (err) throw err
  console.log(buf.toString('utf8'))
})
```

## API

### `streamToBuffer(stream, callback)`

Arguments:

- `stream` - Any readable stream
- `callback` - A callback function which will be called with an optional
  error object as the first argument and a buffer containing the content
  of the `stream` as the 2nd

### `fn = streamToBuffer.onStream(callback)`

Returns a function `fn` which expects a readable stream as its only
argument. When called, it will automatically call `streamToBuffer()`
with the stream as the first argument and the `callback` as the second.

## License

MIT
