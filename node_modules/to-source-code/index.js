'use strict';

var isNil = require('is-nil');

module.exports = function (func) {

  if (!isNil(func)) {

    try {
      return Function.prototype.toString.call(func);
    } catch (e) {}

    try {
      return (func + '');
    } catch (e) {}
  }

  return '';
};
