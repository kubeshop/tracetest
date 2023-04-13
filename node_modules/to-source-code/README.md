# to-source-code

> Converts function to its source code.

[![MIT License](https://img.shields.io/badge/license-MIT_License-green.svg?style=flat-square)](https://github.com/gearcase/to-source-code/blob/master/LICENSE)

[![build:?](https://img.shields.io/travis/gearcase/to-source-code/master.svg?style=flat-square)](https://travis-ci.org/gearcase/to-source-code)
[![coverage:?](https://img.shields.io/coveralls/gearcase/to-source-code/master.svg?style=flat-square)](https://coveralls.io/github/gearcase/to-source-code)


## Install

```
$ npm install --save to-source-code 
```


## Usage

> For more use-cases see the [tests](https://github.com/gearcase/to-source-code/blob/master/test/spec/index.js)


```js
var toSourceCode = require('to-source-code');

toSourceCode(Math.min);               
// => 'function min() { [native code] }'

toSourceCode(function () { console.log(123); });              
// => 'function () { console.log(123); }'

toSourceCode();      // => ''
toSourceCode(null);  // => ''
toSourceCode(1);     // => '1'
toSourceCode('1');   // => '1'
toSourceCode(true);  // => 'true'
toSourceCode({});    // => '[object Object]'
```

## Related

- [is-native](https://github.com/gearcase/is-native) - Checks if the given value is a native function.
- [to-num](https://github.com/gearcase/to-num) - Converts the given value to a number.
- [to-integer](https://github.com/gearcase/to-integer) - Converts the given value to an integer.
- [to-length](https://github.com/gearcase/to-length) - Converts value to an integer suitable for use as the length of an array-like object.
- [to-str](https://github.com/gearcase/to-str) - Converts the given value to a string.
- [to-lower](https://github.com/gearcase/to-lower) - Converts string, as a whole, to lower case.
- [to-upper](https://github.com/gearcase/to-upper) - Converts string, as a whole, to upper case.
- [to-path](https://github.com/gearcase/to-path) - Converts value to a property path array. 

## Contributing

Pull requests and stars are highly welcome.

For bugs and feature requests, please [create an issue](https://github.com/gearcase/to-source-code/issues/new).
