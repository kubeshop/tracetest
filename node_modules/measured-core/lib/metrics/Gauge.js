const { MetricTypes } = require('./Metric');

/**
 * Values that can be read instantly
 * @implements {Metric}
 * @example
 * var Measured = require('measured')
 * var gauge = new Measured.Gauge(function() {
 *     return process.memoryUsage().rss;
 * });
 */
class Gauge {
  /**
   * @param {function} readFn A function that returns the numeric value for this gauge.
   */
  constructor(readFn) {
    this._readFn = readFn;
  }

  /**
   * @return {number} Gauges directly return the value from the callback which should be a number.
   */
  toJSON() {
    return this._readFn();
  }

  /**
   * The type of the Metric Impl. {@link MetricTypes}.
   * @return {string} The type of the Metric Impl.
   */
  getType() {
    return MetricTypes.GAUGE;
  }
}

module.exports = Gauge;
