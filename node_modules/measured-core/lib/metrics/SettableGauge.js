const { MetricTypes } = require('./Metric');

/**
 * Works like a {@link Gauge}, but rather than getting its value from a callback, the value
 * is set when needed. This can be useful for setting a gauges value for asynchronous operations.
 * @implements {Metric}
 * @example
 * const settableGauge = new SettableGauge();
 * // Update the settable gauge ever 10'ish seconds
 * setInterval(() => {
 *     calculateSomethingAsync().then((value) => {
 *         settableGauge.setValue(value);
 *     });
 * }, 10000);
 */
class SettableGauge {
  /**
   * @param {SettableGaugeProperties} [options] See {@link SettableGaugeProperties}.
   */
  constructor(options) {
    options = options || {};
    this._value = options.initialValue || 0;
  }

  setValue(value) {
    this._value = value;
  }

  /**
   * @return {number} Settable Gauges directly return there current value.
   */
  toJSON() {
    return this._value;
  }

  /**
   * The type of the Metric Impl. {@link MetricTypes}.
   * @return {string} The type of the Metric Impl.
   */
  getType() {
    return MetricTypes.GAUGE;
  }
}

module.exports = SettableGauge;

/**
 * Properties that can be supplied to the constructor of a {@link Counter}
 *
 * @interface SettableGaugeProperties
 * @typedef SettableGaugeProperties
 * @type {Object}
 * @property {number} initialValue An initial value to use for this settable gauge. Defaults to 0.
 * @example
 * // Creates a Gauge that with an initial value of 500.
 * const settableGauge = new SettableGauge({ initialValue: 500 })
 *
 */
