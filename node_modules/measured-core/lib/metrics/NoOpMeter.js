const { MetricTypes } = require('./Metric');

/**
 * A No-Op Impl of Meter that can be used with a timer, to only create histogram data.
 * This is useful for some time series aggregators that can calculate rates for you just off of sent count.
 *
 * @implements {Metric}
 * @example
 * const { NoOpMeter, Timer } = require('measured')
 * const meter = new NoOpMeter();
 * const timer = new Timer({meter: meter});
 * ...
 * // do some stuff with the timer and stopwatch api
 * ...
 */
// eslint-disable-next-line padded-blocks
class NoOpMeter {
  /**
   * No-Op impl
   * @param {number} n Number of events to mark.
   */
  // eslint-disable-next-line no-unused-vars
  mark(n) {}

  /**
   * No-Op impl
   */
  start() {}

  /**
   * No-Op impl
   */
  end() {}

  /**
   * No-Op impl
   */
  ref() {}

  /**
   * No-Op impl
   */
  unref() {}

  /**
   * No-Op impl
   */
  reset() {}

  /**
   * No-Op impl
   */
  meanRate() {}

  /**
   * No-Op impl
   */
  currentRate() {}

  /**
   * Returns an empty object
   * @return {{}}
   */
  toJSON() {
    return {};
  }

  /**
   * The type of the Metric Impl. {@link MetricTypes}.
   * @return {string} The type of the Metric Impl.
   */
  getType() {
    return MetricTypes.METER;
  }
}

module.exports = NoOpMeter;
