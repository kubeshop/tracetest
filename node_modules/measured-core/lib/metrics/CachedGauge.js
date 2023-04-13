const { MetricTypes } = require('./Metric');
const TimeUnits = require('../util/units');

/**
 * A Cached Gauge takes a function that returns a promise that resolves a
 * value that should be cached and updated on a given interval.
 *
 * toJSON() will return the currently cached value.
 *
 * @example
 * const cpuAverageCachedGauge = new CachedGauge(() => {
 *     return new Promise(resolve => {
 *       //Grab first CPU Measure
 *       const startMeasure = cpuAverage();
 *       setTimeout(() => {
 *         //Grab second Measure
 *         const endMeasure = cpuAverage();
 *         const percentageCPU = calculateCpuUsagePercent(startMeasure, endMeasure);
 *         resolve(percentageCPU);
 *       }, sampleTimeInSeconds);
 *     });
 *   }, updateIntervalInSeconds);
 *
 * @implements {Metric}
 */
class CachedGauge {
  /**
   * @param {function} valueProducingPromiseCallback A function that returns a promise than when
   * resolved supplies the value that should be cached in this gauge.
   * @param {number} updateIntervalInSeconds How often the cached gauge should update it's value.
   * @param {number} [timeUnitOverride] by default this function takes updateIntervalInSeconds and multiplies it by TimeUnits.SECONDS (1000),
   * You can override it here.
   */
  constructor(valueProducingPromiseCallback, updateIntervalInSeconds, timeUnitOverride) {
    const timeUnit = timeUnitOverride || TimeUnits.SECONDS;

    this._valueProducingPromiseCallback = valueProducingPromiseCallback;
    this._value = 0;
    this._updateValue();
    this._interval = setInterval(() => {
      this._updateValue();
    }, updateIntervalInSeconds * timeUnit);
  }

  /**
   * Calls the promise producing callback and sets the value when it gets resolved.
   * @private
   */
  _updateValue() {
    this._valueProducingPromiseCallback().then(value => {
      this._value = value;
    });
  }

  /**
   * @return {number} Gauges directly return the value which should be a number.
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

  /**
   * Clears the interval, so that it doesn't keep any processes alive.
   */
  end() {
    clearInterval(this._interval);
    this._interval = null;
  }
}

module.exports = CachedGauge;
