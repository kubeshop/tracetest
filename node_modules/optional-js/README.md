[![npm version](https://badge.fury.io/js/optional-js.svg)](https://badge.fury.io/js/optional-js) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Optional.js
===========

> A container object that wraps possibly undefined values in JavaScript - inspired by [Java Optionals](https://docs.oracle.com/javase/9/docs/api/java/util/Optional.html)

``` javascript
Optional.ofNullable(promptForUserName)
        .map(getUserId)
        .filter(verify)
        .ifPresent(login);
```

# Features

- Runs in browser and Node
- Full Java 8 Optional API is supported, and partial Java 9 API implemented (everything minus stream())
- Zero dependencies
- TypeScript type definitions included
- Lightweight (**<1.0 KB minified, gzipped**)

# Installation

Download the [latest release](https://github.com/JasonStorey/Optional.js/releases) from GitHub or from [NPM](https://www.npmjs.com/package/optional-js)

via npm:
``` bash
$ npm install optional-js
```

then just require in node:
``` javascript
const Optional = require('optional-js');
const emptyOptional = Optional.empty();
```

alternatively, use the browser compatible build in the `./dist` directory of the npm package

Not using a module loader? Include the script, and the browser global `Optional` will be added to window.

# Usage

Java docs - [Java 9 Optionals](https://docs.oracle.com/javase/9/docs/api/java/util/Optional.html)  
TSDocs - [index.d.ts](https://github.com/JasonStorey/Optional.js/blob/master/index.d.ts)

JS Example:
``` javascript
// "login.js"

const Optional = require('optional-js');

// Define some simple operations
const getUserId = 
    username => username === 'root' ? 1234 : 0;

const verify = 
    userId => userId === 1234;

const login = 
    userId => console.log('Logging in as : ' + userId);
    
// Declare a potentially undefined value
const username = process.argv[2];

// Wrap username in an Optional, and build a pipeline using our operations
Optional.ofNullable(username)
        .map(getUserId)
        .filter(verify)
        .ifPresent(login);
```
Then, from the terminal...
``` bash
$ node login.js root
"Logging in as : 1234"
````

# Building

download:
``` bash
git clone git@github.com:JasonStorey/Optional.js.git
```

enter the directory, and install dependencies:
```bash
cd Optional.js && npm install
```

build:
```bash
npm run build
```

# Testing

run the tests:
```bash
npm test
```

# Contributing

Found a bug or missing feature? Please open an [issue](https://github.com/JasonStorey/Optional.js/issues)!

Send your feedback. Send your pull requests. All contributions are appreciated!

# License

Optional.js may be freely distributed under the MIT license - [LICENSE](https://github.com/JasonStorey/Optional.js/blob/master/LICENSE)
