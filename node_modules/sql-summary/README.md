# sql-summary

Summarize any SQL query.

This JavaScript module will analyse an SQL query and return a very brief
summary string containing:

- Main verb used (`SELECT`, `INSERT`, `UPDATE` etc.)
- Potentially the type operated on (`TABLE`, `DATABASE` etc.)
- The name of the primary table or database operated on

For example, if given the following SQL query:

```sql
SELECT col1, col2 FROM table_name WHERE id=1
```

The following summary string is produced:

```sql
SELECT FROM table_name
```

[![npm](https://img.shields.io/npm/v/sql-summary.svg)](https://www.npmjs.com/package/sql-summary)
[![Build status](https://travis-ci.org/elastic/sql-summary.svg?branch=master)](https://travis-ci.org/elastic/sql-summary)
[![js-standard-style](https://img.shields.io/badge/code%20style-standard-brightgreen.svg?style=flat)](https://github.com/feross/standard)

## Installation

```
npm install sql-summary
```

## Usage

```js
var sqlSummary = require('sql-summary')

var query = 'UPDATE table_name SET col1=1, col2=2 WHERE id=1'

console.log('Running query:', sqlSummary(query)) // => 'Running query: UPDATE table_name'
db.query(query, function (err, result) {
  // ...
})
```

## Use-cases

- In a web-server log output the type of queries used without going into
  too much details about each query
- Group similar queries operating on the same tables

## License

[MIT](https://github.com/elastic/sql-summary/blob/master/LICENSE)
