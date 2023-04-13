const { MetricTypes } = require('./Metric');

/**
 * Counters are things that increment or decrement
 * @implements {Metric}
 * @example
 * var Measured = require('measured')
 * var activeUploads = new Measured.Counter();
 * http.createServer(function(req, res) {
 *    activeUploads.inc();
 *    req.on('end', function() {
 *         activeUploads.dec();
 *    });
 * });
 */
class Counter {
  /**
   * @param {CounterProperties} [properties] see {@link CounterProperties}
   */
  constructor(properties) {
    properties = properties || {};

    this._count = properties.count || 0;
  }

  /**
   * Counters directly return their currently value.
   * @return {number}
   */
  toJSON() {
    return this._count;
  }

  /**
   * Increments the counter.
   * @param {number} n Increment the counter by n. Defaults to 1.
   */
  inc(n) {
    this._count += arguments.length ? n : 1;
  }

  /**
   * Decrements the counter
   * @param {number} n Decrement the counter by n. Defaults to 1.
   */
  dec(n) {
    this._count -= arguments.length ? n : 1;
  }

  /**
   * Resets the counter back to count Defaults to 0.
   * @param {number} count Resets the counter back to count Defaults to 0.
   */
  reset(count) {
    this._count = count || 0;
  }

  /**
   * The type of the Metric Impl. {@link MetricTypes}.
   * @return {string} The type of the Metric Impl.
   */
  getType() {
    return MetricTypes.COUNTER;
  }
}

module.exports = Counter;

/**
 * Properties that can be supplied to the constructor of a {@link Counter}
 *
 * @interface CounterProperties
 * @typedef CounterProperties
 * @type {Object}
 * @property {number} count An initial count for the counter. Defaults to 0.
 * @example
 * // Creates a counter that starts at 5.
 * const counter = new Counter({ count: 5 })
 */
