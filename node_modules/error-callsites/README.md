# error-callsites

Yet another module for extracting callsites (a.k.a. stack-frames) from
Node.js `Error` objects.

[![npm](https://img.shields.io/npm/v/error-callsites.svg)](https://www.npmjs.com/package/error-callsites)
[![Tests](https://github.com/watson/error-callsites/workflows/Tests/badge.svg)](https://github.com/watson/error-callsites/actions)
[![js-standard-style](https://img.shields.io/badge/code%20style-standard-brightgreen.svg?style=flat)](https://github.com/feross/standard)

## Installation

```
npm install error-callsites
```

## Usage

```js
var callsites = require('error-callsites')

var err = new Error('foo')
var frames = callsites(err)

console.log('Error occurred on line', frames[0].getLineNumber())
```

## API

The module exposes a single function which expects an `Error` object as
the first argument:

```js
var frames = callsites(err)
```

More specifically, it expects an object on which `Error.captureStackTrace(obj)`
has been called, which is the case for `Error` objects. You can do it directly
yourself:

```js
var obj = {}
Error.captureStackTrace(obj)
var frames = callsites(err)
```

The function returns an array of
[v8 `CallSite` objects](https://v8.dev/docs/stack-trace-api#customizing-stack-traces)
-- one for each frame in the stack trace. CallSite objects provide accessor
methods for information on that frame. For example:

- `callsite.getThis()` - returns the value of this
- `callsite.getTypeName()` - returns the type of this as a string. This is the name of the function stored in the constructor field of this, if available, otherwise the object's [[Class]] internal property.
- `callsite.getFunction()` - returns the current function
- `callsite.getFunctionName()` - returns the name of the current function, typically its name property. If a name property is not available an attempt will be made to try to infer a name from the function's context.
- `callsite.getMethodName()` - returns the name of the property of this or one of its prototypes that holds the current function
- `callsite.getFileName()` - if this function was defined in a script returns the name of the script
- `callsite.getLineNumber()` - if this function was defined in a script returns the current line number
- `callsite.getColumnNumber()` - if this function was defined in a script returns the current column number
- `callsite.getEvalOrigin()` - if this function was created using a call to eval returns a CallSite object representing the location where eval was called
- `callsite.isToplevel()` - is this a toplevel invocation, that is, is this the global object?
- `callsite.isEval()` - does this call take place in code defined by a call to eval?
- `callsite.isNative()` - is this call in native V8 code?
- `callsite.isConstructor()` - is this a constructor call?

See the link above for a complete list of the accessor methods. Some have been
added in more recent versions of v8 and Node.js.


## How it works

This package works by setting a [custom `Error.prepareStackTrace`](https://v8.dev/docs/stack-trace-api#customizing-stack-traces)
property that transparently captures the structured stack trace that v8 passes
to this API, while also supporting other Node.js code in the process setting
its own `Error.prepareStackTrace`.


## License

[MIT](LICENSE)
