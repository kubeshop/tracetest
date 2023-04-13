'use strict';

var isNil    = require('is-nil');
var toSource = require('to-source-code');

function isObject(value) {

  var type = typeof value;

  return !!value && (type === 'object' || type === 'function');
}

// Checks if `value` is a host object in IE < 9.
function isHostObject(value) {

  // Many host objects are `Object` objects that can coerce to strings
  // despite having improperly defined `toString` methods.

  var result = false;

  if (!isNil(value) && typeof value.toString !== 'function') {
    try {
      result = ('' + value) !== '';
    } catch (e) {}
  }
  return result;
}

function isFunction(value) {

  var tag = isObject(value) ? Object.prototype.toString.call(value) : '';

  return tag === '[object Function]' || tag === '[object GeneratorFunction]';
}


module.exports = function (value) {

  if (!isObject(value)) {
    return false;
  }

  var pattern;

  if (isFunction(value) || isHostObject(value)) {

    var toString       = Function.prototype.toString;
    var hasOwnProperty = Object.prototype.hasOwnProperty;
    var reRegExpChar   = /[\\^$.*+?()[\]{}|]/g;

    pattern = new RegExp('^' +
      toString.call(hasOwnProperty)
        .replace(reRegExpChar, '\\$&')
        .replace(/hasOwnProperty|(function).*?(?=\\\()| for .+?(?=\\\])/g, '$1.*?') + '$'
    );
  } else {
    // detect host constructors (Safari).
    pattern = /^\[object .+?Constructor\]$/;
  }

  return pattern.test(toSource(value));
};
