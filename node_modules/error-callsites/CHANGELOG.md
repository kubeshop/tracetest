# error-callsites Changelog

## 2.0.4

- Fix a crash when a user set `Error.prepareStackTrace = undefined`, which
  is a way of saying "give me the JS engine's default stack formatting".
  ([#3](https://github.com/watson/error-callsites/issues/3))
- Internal maintenance changes:
    - @trentm is mostly maintaining this module now
    - Switch from TravisCI to GitHub CI. Add testing for node v14 and v16.


## 2.0.3

(This changelog was started after this release.)

