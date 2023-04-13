# stream-chopper

Chop a single stream of data into a series of readable streams.

[![npm](https://img.shields.io/npm/v/stream-chopper.svg)](https://www.npmjs.com/package/stream-chopper)
[![build status](https://travis-ci.org/watson/stream-chopper.svg?branch=master)](https://travis-ci.org/watson/stream-chopper)
[![codecov](https://img.shields.io/codecov/c/github/watson/stream-chopper.svg)](https://codecov.io/gh/watson/stream-chopper)
[![js-standard-style](https://img.shields.io/badge/code%20style-standard-brightgreen.svg?style=flat)](https://github.com/feross/standard)

Stream Chopper is useful in situations where you have a stream of data
you want to chop up into smaller pieces, either based on time or size.
Each piece will be emitted as a readable stream (called output streams).

Possible use-cases include log rotation, splitting up large data sets,
or chopping up a live stream of data into finite chunks that can then be
stored.

## Control how data is split

Sometimes it's important to ensure that a chunk written to the input
stream isn't split up and devided over two output streams. Stream
Chopper allows you to specify the chopping algorithm used (via the
`type` option) when a chunk is too large to fit into the current output
stream.

By default a chunk too large to fit in the current output stream is
split between it and the next. Alternatively you can decide to either
allow the chunk to "overflow" the size limit, in which case it will be
written to the current output stream, or to "underflow" the size limit,
in which case the current output stream will be ended and the chunk
written to the next output stream.

## Installation

```
npm install stream-chopper --save
```

## Usage

Example app:

```js
const StreamChopper = require('stream-chopper')

const chopper = new StreamChopper({
  size: 30,                    // chop stream when it reaches 30 bytes,
  time: 10000,                 // or when it's been open for 10s,
  type: StreamChopper.overflow // but allow stream to exceed size slightly
})

chopper.on('stream', function (stream, next) {
  console.log('>> Got a new stream! <<')
  stream.pipe(process.stdout)
  stream.on('end', next) // call next when you're ready to receive a new stream
})

chopper.write('This write contains more than 30 bytes\n')
chopper.write('This write contains less\n')
chopper.write('This is the last write\n')
```

Output:

```
>> Got a new stream! <<
This write contains more than 30 bytes
>> Got a new stream! <<
This write contains less
This is the last write
```

## API

### `chopper = new StreamChopper([options])`

Instantiate a `StreamChopper` instance. `StreamChopper` is a [writable]
stream.

Takes an optional `options` object which, besides the normal options
accepted by the [`Writable`][writable] class, accepts the following
config options:

- `size` - The maximum number of bytes that can be written to the
  `chopper` stream before a new output stream is emitted (default:
  `Infinity`)
- `time` - The maximum number of milliseconds that an output stream can
  be in use before a new output stream is emitted (default: `-1` which
  means no limit)
- `type` - Change the algoritm used to determine how a written chunk
  that cannot fit into the current output stream should be handled. The
  following values are possible:
  - `StreamChopper.split` - Fit as much data from the chunk as possible
    into the current stream and write the remainder to the next stream
    (default)
  - `StreamChopper.overflow` - Allow the entire chunk to be written to
    the current stream. After writing, the stream is ended
  - `StreamChopper.underflow` - End the current output stream and write
    the entire chunk to the next stream
- `transform` - An optional function that returns a transform stream
  used for transforming the data in some way (e.g. a zlib Gzip stream).
  If used, the `size` option will count towards the size of the output
  chunks. This config option cannot be used together with the
  `StreamChopper.split` type

If `type` is `StreamChopper.underflow` and the size of the chunk to be
written is larger than `size` an error is emitted.

### Event: `stream`

Emitted every time a new output stream is ready. You must listen for
this event.

The listener function is called with two arguments:

- `stream` - A [readable] output stream
- `next` - A function you must call when you're ready to receive a new
  output stream. If called with an error, the `chopper` stream is
  destroyed

### `chopper.size`

The maximum number of bytes that can be written to the `chopper` stream
before a new output stream is emitted.

Use this property to override it with a new value. The new value will
take effect immediately on the current stream.

### `chopper.time`

The maximum number of milliseconds that an output stream can be in use
before a new output stream is emitted.

Use this property to override it with a new value. The new value will
take effect when the next stream is initialized. To change the current
timer, see [`chopper.resetTimer()`](#chopperresettimertime).

Set to `-1` for no time limit.

### `chopper.type`

The algoritm used to determine how a written chunk that cannot fit into
the current output stream should be handled. The following values are
possible:

- `StreamChopper.split` - Fit as much data from the chunk as possible
  into the current stream and write the remainder to the next stream
- `StreamChopper.overflow` - Allow the entire chunk to be written to
  the current stream. After writing, the stream is ended
- `StreamChopper.underflow` - End the current output stream and write
  the entire chunk to the next stream

Use this property to override it with a new value. The new value will
take effect immediately on the current stream.

### `chopper.chop([callback])`

Manually chop the stream. Forces the current output stream to end even
if its `size` limit or `time` timeout hasn't been reached yet.

Arguments:

- `callback` - An optional callback which will be called once the output
  stream have ended

### `chopper.resetTimer([time])`

Use this function to reset the current timer (configured via the `time`
config option). Calling this function will force the current timer to
start over.

If the optional `time` argument is provided, this value is used as the
new time. This is equivilent to calling:

```js
chopper.time = time
chopper.resetTimer()
```

If the function is called with `time` set to `-1`, the current timer is
cancelled and the time limit is disabled for all future streams.

## License

[MIT](https://github.com/watson/stream-chopper/blob/master/LICENSE)

[writable]: https://nodejs.org/api/stream.html#stream_class_stream_writable
[readable]: https://nodejs.org/api/stream.html#stream_class_stream_readable
