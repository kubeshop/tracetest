const { MetricTypes } = require('./Metric');
const Histogram = require('./Histogram');
const Meter = require('./Meter');
const Stopwatch = require('../util/Stopwatch');

/**
 *
 * Timers are a combination of Meters and Histograms. They measure the rate as well as distribution of scalar events.
 * <p>
 * Since they are frequently used for tracking how long certain things take, they expose an API for that: See example 1.
 * <p>
 * But you can also use them as generic histograms that also track the rate of events: See example 2.
 *
 * @example
 * var Measured = require('measured')
 * var timer = new Measured.Timer();
 * http.createServer(function(req, res) {
 *     var stopwatch = timer.start();
 *     req.on('end', function() {
 *         stopwatch.end();
 *     });
 * });
 *
 *
 * @example
 * var Measured = require('measured')
 * var timer = new Measured.Timer();
 * http.createServer(function(req, res) {
 *    if (req.headers['content-length']) {
 *        timer.update(parseInt(req.headers['content-length'], 10));
 *    }
 * });
 *
 * @implements {Metric}
 */
class Timer {
  /**
   * @param {TimerProperties} [properties] See {@link TimerProperties}.
   */
  constructor(properties) {
    properties = properties || {};

    this._meter = properties.meter || new Meter({});
    this._histogram = properties.histogram || new Histogram({});
    this._getTime = properties.getTime;
    this._keepAlive = !!properties.keepAlive;

    if (!properties.keepAlive) {
      this.unref();
    }
  }

  /**
   * @return {Stopwatch} Returns a Stopwatch that has been started.
   */
  start() {
    const self = this;
    const watch = new Stopwatch({ getTime: this._getTime });

    watch.once('end', elapsed => {
      self.update(elapsed);
    });

    return watch;
  }

  /**
   * Updates the internal histogram with value and marks one event on the internal meter.
   * @param {number} value
   */
  update(value) {
    this._meter.mark();
    this._histogram.update(value);
  }

  /**
   * Resets all values. Timers initialized with custom options will be reset to the default settings.
   */
  reset() {
    this._meter.reset();
    this._histogram.reset();
  }

  end() {
    this._meter.end();
  }

  /**
   * Refs the backing timer again. Idempotent.
   */
  ref() {
    this._meter.ref();
  }

  /**
   * Unrefs the backing timer. The meter will not keep the event loop alive. Idempotent.
   */
  unref() {
    this._meter.unref();
  }

  /**
   * toJSON output:
   *
   * <li> meter: See <a href="#meter">Meter</a>#toJSON output docs above.</li>
   * <li> histogram: See <a href="#histogram">Histogram</a>#toJSON output docs above.</a></li>
   *
   * @return {any}
   */
  toJSON() {
    return {
      meter: this._meter.toJSON(),
      histogram: this._histogram.toJSON()
    };
  }

  /**
   * The type of the Metric Impl. {@link MetricTypes}.
   * @return {string} The type of the Metric Impl.
   */
  getType() {
    return MetricTypes.TIMER;
  }
}

module.exports = Timer;

/**
 * @interface TimerProperties
 * @typedef TimerProperties
 * @type {Object}
 * @property {Meter} meter The internal meter to use. Defaults to a new {@link Meter}.
 * @property {Histogram} histogram The internal histogram to use. Defaults to a new {@link Histogram}.
 * @property {function} getTime optional function override for supplying time to the {@link Stopwatch}
 * @property {boolean} keepAlive Optional flag to unref the associated timer. Defaults to `false`.
 */
