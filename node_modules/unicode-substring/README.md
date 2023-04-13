# unicode-substring [![Build Status](https://travis-ci.org/lautis/unicode-substring.svg?branch=master)](https://travis-ci.org/lautis/unicode-substring)

Unicode-aware substring for JavaScript. Surrogate pairs are counted as one character.

## What?

Characters in JavaScript strings are exposed as 16-bit code points, also known as UCS-2 encoding. This usually good enough, but since there are more than 2^16 characters in Unicode, 16 bits is not enough to represent all characters. To overcome this limitation, characters with scalar value over `0x10FFFF` need to be encoded as surrogate pairs. This encoding is known as UTF-16.

The purpose of this library is to treat surrogate pairs as one character when extracting substrings from a string. This might be preferable if indices are returned from an Unicode-compatible environment.

## Usage

```javascript

var unicodeSubstring = require('unicode-substring')
// unicodeSubstring(string, start, end)
unicodeSubstring("Hello World!", 0, 5)
// => "Hello"
```

The `start` and `end` parameters behave similarly as [String.prototype.substring](https://developer.mozilla.org/en/docs/Web/JavaScript/Reference/Global_Objects/String/substring).

