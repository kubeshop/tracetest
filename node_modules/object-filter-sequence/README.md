# object-filter-sequence

[![npm](https://img.shields.io/npm/v/object-filter-sequence.svg)](https://www.npmjs.com/package/object-filter-sequence)
[![Build status](https://travis-ci.org/elastic/object-filter-sequence.svg?branch=master)](https://travis-ci.org/elastic/object-filter-sequence)
[![codecov](https://img.shields.io/codecov/c/github/elastic/object-filter-sequence.svg)](https://codecov.io/gh/elastic/object-filter-sequence)
[![Standard - JavaScript Style Guide](https://img.shields.io/badge/code%20style-standard-brightgreen.svg?style=flat)](https://github.com/feross/standard)

This module provides an interface to apply a sequence of filters to an object. It is a subclass of Array, so any array method can be used on it.

## Installation

```
npm install object-filter-sequence
```

## Example Usage

```js
const Filters = require('object-filter-sequence')

const filters = new Filters()

filters.push(previous => {
  const next = {}
  next.key = previous.key.toUpperCase()
  return next
})

filters.push(previous => {
  const next = {}
  next.key = previous.key.reverse()
  return next
})

filters.process({ key: 'value' }) // { key: 'EULAV' }
```

## API

### `filters.process(object)`

This is the only unique method from the Array base class. It is used to apply the filters in the array to the provided object.

## License

[MIT](LICENSE)
