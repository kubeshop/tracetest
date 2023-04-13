# original-url

Reconstruct the original URL used in an HTTP request based on the HTTP
request headers.

The module takes into account potential URL rewrites made by proxies,
load balancers, etc along the way (as long as these append special HTTP
headers to the request).

Supported HTTP headers:

- `Host`
- `Forwarded`
- `X-Forwarded-Proto`
- `X-Forwarded-Protocol`
- `X-Url-Scheme`
- `Front-End-Https`
- `X-Forwarded-Ssl`
- `X-Forwarded-Host`
- `X-Forwarded-Port`

If the protocol (http vs https) cannot be determined based on the above
headers, the `encrypted` flag on the TLS connection is used.

[![npm](https://img.shields.io/npm/v/original-url.svg)](https://www.npmjs.com/package/original-url)
[![Build status](https://travis-ci.org/watson/original-url.svg?branch=master)](https://travis-ci.org/watson/original-url)
[![js-standard-style](https://img.shields.io/badge/code%20style-standard-brightgreen.svg?style=flat)](https://github.com/feross/standard)

## Installation

```
npm install original-url --save
```

## Usage

Server example:

```js
const http = require('http')
const originalUrl = require('original-url')

const server = http.createServer(function (req, res) {
  const url = originalUrl(req)
  if (url.full) {
    res.end(`Original URL: ${url.full}\n`)
  } else {
    res.end('Original URL could not be determined\n')
  }
})

server.listen(1337)
```

Request examples:

```
$ curl localhost:1337
Original URL: http://localhost:1337/

$ curl -H 'Host: example.com' localhost:1337
Original URL: http://example.com/

$ curl -H 'Host: example.com:1234' localhost:1337
Original URL: http://example.com:1234/

$ curl -H 'Forwarded: proto=https; host=example.com; for="10.0.0.1:1234"' localhost:1337/sub/path?key=value
Original URL: https://example.com:1234/sub/path?key=value

$ curl -H 'X-Forwarded-Host: example.com' -H 'X-Forwarded-Host: proxy.local' localhost:1337
Original URL: http://example.com/
```

## API

### `result = originalUrl(req)`

This module exposes a single function which takes an HTTP request object
in the form of
[`http.IncomingMessage`](https://nodejs.org/api/http.html#http_class_http_incomingmessage).

When called, the function returns a `result` object with a `full`
property containing the fully resolved URL. The `result` object will
also contain any other property normally returned by the Node.js core
[`url.parse`](https://nodejs.org/api/url.html#url_url_parse_urlstring_parsequerystring_slashesdenotehost)
function.

If the hostname for some reason cannot be determined, `result.full` will
not be present.

## License

[MIT](LICENSE)
