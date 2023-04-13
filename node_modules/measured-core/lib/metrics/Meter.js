const { MetricTypes } = require('./Metric');
const units = require('../util/units');
const EWMA = require('../util/ExponentiallyMovingWeightedAverage');

const RATE_UNIT = units.SECONDS;
const TICK_INTERVAL = 5 * units.SECONDS;

/**
 * Things that are measured as events / interval.
 * @implements {Metric}
 * @example
 * var Measured = require('measured')
 * var meter = new Measured.Meter();
 * http.createServer(function(req, res) {
 *     meter.mark();
 * });
 */
class Meter {
  /**
   * @param {MeterProperties} [properties] see {@link MeterProperties}.
   */
  constructor(properties) {
    this._properties = properties || {};
    this._initializeState();

    if (!this._properties.keepAlive) {
      this.unref();
    }
  }

  /**
   * Initializes the state of this Metric
   * @private
   */
  _initializeState() {
    this._rateUnit = this._properties.rateUnit || RATE_UNIT;
    this._tickInterval = this._properties.tickInterval || TICK_INTERVAL;
    if (this._properties.getTime) {
      this._getTime = this._properties.getTime;
    }

    this._m1Rate = this._properties.m1Rate || new EWMA(units.MINUTES, this._tickInterval);
    this._m5Rate = this._properties.m5Rate || new EWMA(5 * units.MINUTES, this._tickInterval);
    this._m15Rate = this._properties.m15Rate || new EWMA(15 * units.MINUTES, this._tickInterval);
    this._count = 0;
    this._currentSum = 0;
    this._startTime = this._getTime();
    this._lastToJSON = this._getTime();
    this._interval = setInterval(this._tick.bind(this), TICK_INTERVAL);
  }

  /**
   * Register n events as having just occured. Defaults to 1.
   * @param {number} [n]
   */
  mark(n) {
    if (!this._interval) {
      this.start();
    }

    n = n || 1;

    this._count += n;
    this._currentSum += n;
    this._m1Rate.update(n);
    this._m5Rate.update(n);
    this._m15Rate.update(n);
  }

  start() {}

  end() {
    clearInterval(this._interval);
    this._interval = null;
  }

  /**
   * Refs the backing timer again. Idempotent.
   */
  ref() {
    if (this._interval && this._interval.ref) {
      this._interval.ref();
    }
  }

  /**
   * Unrefs the backing timer. The meter will not keep the event loop alive. Idempotent.
   */
  unref() {
    if (this._interval && this._interval.unref) {
      this._interval.unref();
    }
  }

  _tick() {
    this._m1Rate.tick();
    this._m5Rate.tick();
    this._m15Rate.tick();
  }

  /**
   * Resets all values. Meters initialized with custom options will be reset to the default settings (patch welcome).
   */
  reset() {
    this.end();
    this._initializeState();
  }

  meanRate() {
    if (this._count === 0) {
      return 0;
    }

    const elapsed = this._getTime() - this._startTime;
    return this._count / elapsed * this._rateUnit;
  }

  currentRate() {
    const currentSum = this._currentSum;
    const duration = this._getTime() - this._lastToJSON;
    const currentRate = currentSum / duration * this._rateUnit;

    this._currentSum = 0;
    this._lastToJSON = this._getTime();

    // currentRate could be NaN if duration was 0, so fix that
    return currentRate || 0;
  }

  /**
   * @return {MeterData}
   */
  toJSON() {
    return {
      mean: this.meanRate(),
      count: this._count,
      currentRate: this.currentRate(),
      '1MinuteRate': this._m1Rate.rate(this._rateUnit),
      '5MinuteRate': this._m5Rate.rate(this._rateUnit),
      '15MinuteRate': this._m15Rate.rate(this._rateUnit)
    };
  }

  _getTime() {
    if (!process.hrtime) {
      return new Date().getTime();
    }

    const hrtime = process.hrtime();
    return hrtime[0] * 1000 + hrtime[1] / (1000 * 1000);
  }

  /**
   * The type of the Metric Impl. {@link MetricTypes}.
   * @return {string} The type of the Metric Impl.
   */
  getType() {
    return MetricTypes.METER;
  }
}

module.exports = Meter;

/**
 *
 * @interface MeterProperties
 * @typedef MeterProperties
 * @type {Object}
 * @property {number} rateUnit The rate unit. Defaults to 1000 (1 sec).
 * @property {number} tickInterval The interval in which the averages are updated. Defaults to 5000 (5 sec).
 * @property {boolean} keepAlive Optional flag to unref the associated timer. Defaults to `false`.
 * @example
 * const meter = new Meter({ rateUnit: 1000, tickInterval: 5000})
 */

/**
 * The data returned from Meter::toJSON()
 * @interface MeterData
 * @typedef MeterData
 * @typedef {object}
 * @property {number} mean The average rate since the meter was started.
 * @property {number} count The total of all values added to the meter.
 * @property {number} currentRate The rate of the meter since the last toJSON() call.
 * @property {number} 1MinuteRate The rate of the meter biased towards the last 1 minute.
 * @property {number} 5MinuteRate The rate of the meter biased towards the last 5 minutes.
 * @property {number} 15MinuteRate The rate of the meter biased towards the last 15 minutes.
 */
