const consoleLogLevel = require('console-log-level');
const Optional = require('optional-js');
const { validateReporterParameters } = require('../validators/inputValidators');

const DEFAULT_REPORTING_INTERVAL_IN_SECONDS = 10;

function prefix() {
  return `${new Date().toISOString()}: `;
}

/**
 * The abstract reporter that specific implementations can extend to create a Self Reporting Metrics Registry Reporter.
 *
 * {@link SelfReportingMetricsRegistry}
 *
 * @example
 * const os = require('os');
 * const process = require('process');
 * const { SelfReportingMetricsRegistry, Reporter } = require('measured-reporting');
 *
 * // Create a self reporting registry with a named anonymous reporter instance;
 * const registry = new SelfReportingMetricsRegistry(
 *   new class ConsoleReporter extends Reporter {
 *     constructor() {
 *       super({
 *         defaultDimensions: {
 *           hostname: os.hostname(),
 *           env: process.env['NODE_ENV'] ? process.env['NODE_ENV'] : 'unset'
 *         }
 *        })
 *     }
 *
 *     _reportMetrics(metrics) {
 *        metrics.forEach(metric => {
 *          console.log(JSON.stringify({
 *            metricName: metric.name,
 *            dimensions: this._getDimensions(metric),
 *            data: metric.metricImpl.toJSON()
 *           }))
 *       });
 *     }
 *  }()
 * );
 *
 * @example
 * // Create a regular class that extends Reporter
 * class LoggingReporter extends Reporter {
 *   _reportMetrics(metrics) {
 *     metrics.forEach(metric => {
 *       this._log.info(JSON.stringify({
 *        metricName: metric.name,
 *        dimensions: this._getDimensions(metric),
 *        data: metric.metricImpl.toJSON()
 *       }))
 *     });
 *   }
 * }
 *
 * @abstract
 */
class Reporter {
  /**
   * @param {ReporterOptions} [options] The optional params to supply when creating a reporter.
   */
  constructor(options) {
    if (this.constructor === Reporter) {
      throw new TypeError("Can't instantiate abstract class!");
    }

    options = options || {};
    validateReporterParameters(options);

    /**
     * Map of intervals to metric keys, this will be used to look up what metrics should be reported at a given interval.
     *
     * @type {Object.<number, Set<string>>}
     * @private
     */
    this._intervalToMetric = {};
    this._intervals = [];

    /**
     * Map of default dimensions, that should be sent with every metric.
     *
     * @type {Dimensions}
     * @protected
     */
    this._defaultDimensions = options.defaultDimensions || {};

    /**
     * Loggers to use, defaults to a new console logger if nothing is supplied in options
     * @type {Logger}
     * @protected
     */
    this._log =
      options.logger || consoleLogLevel({ name: 'Reporter', level: options.logLevel || 'info', prefix: prefix });

    /**
     * The default reporting interval, a number in seconds.
     * If not overridden via the {@see ReporterOptions}, defaults to 10 seconds.
     *
     * @type {number}
     * @protected
     */
    this._defaultReportingIntervalInSeconds =
      options.defaultReportingIntervalInSeconds || DEFAULT_REPORTING_INTERVAL_IN_SECONDS;

    /**
     * Flag to indicate if reporting timers should be unref'd.
     * If not overridden via the {@see ReporterOptions}, defaults to false.
     *
     * @type {boolean}
     * @protected
     */
    this._unrefTimers = !!options.unrefTimers;

    /**
     * Flag to indicate if metrics should be reset on each reporting interval.
     * If not overridden via the {@see ReporterOptions}, defaults to false.
     *
     * @type {boolean}
     * @protected
     */
    this._resetMetricsOnInterval = !!options.resetMetricsOnInterval;
  }

  /**
   * Sets the registry, this must be called before reportMetricOnInterval.
   *
   * @param {DimensionAwareMetricsRegistry} registry
   */
  setRegistry(registry) {
    this._registry = registry;
  }

  /**
   * Informs the reporter to report a metric on a given interval in seconds.
   *
   * @param {string} metricKey The metric key for the metric in the metric registry.
   * @param {number} intervalInSeconds The interval in seconds to report the metric on.
   */
  reportMetricOnInterval(metricKey, intervalInSeconds) {
    intervalInSeconds = intervalInSeconds || this._defaultReportingIntervalInSeconds;

    if (!this._registry) {
      throw new Error(
        'You must call setRegistry(registry) before telling a Reporter to report a metric on an interval.'
      );
    }

    if (Object.prototype.hasOwnProperty.call(this._intervalToMetric, intervalInSeconds)) {
      this._intervalToMetric[intervalInSeconds].add(metricKey);
    } else {
      this._intervalToMetric[intervalInSeconds] = new Set([metricKey]);
      this._createIntervalCallback(intervalInSeconds);
      setImmediate(() => {
        this._reportMetricsWithInterval(intervalInSeconds);
      });
    }
  }

  /**
   * Creates the timed callback loop for the given interval.
   *
   * @param {number} intervalInSeconds the interval in seconds for the timeout callback
   * @private
   */
  _createIntervalCallback(intervalInSeconds) {
    this._log.debug(`_createIntervalCallback() called with intervalInSeconds: ${intervalInSeconds}`);

    const timer = setInterval(() => {
      this._reportMetricsWithInterval(intervalInSeconds);
    }, intervalInSeconds * 1000);

    if (this._unrefTimers) {
      timer.unref();
    }

    this._intervals.push(timer);
  }

  /**
   * Gathers all the metrics that have been registered to report on the given interval.
   *
   * @param {number} interval The interval to look up what metrics to report
   * @private
   */
  _reportMetricsWithInterval(interval) {
    this._log.debug(`_reportMetricsWithInterval() called with intervalInSeconds: ${interval}`);
    try {
      Optional.of(this._intervalToMetric[interval]).ifPresent(metrics => {
        const metricsToSend = [];
        metrics.forEach(metricKey => {
          metricsToSend.push(this._registry.getMetricWrapperByKey(metricKey));
        });
        this._reportMetrics(metricsToSend);

        if (this._resetMetricsOnInterval) {
          metricsToSend.forEach(({ name, metricImpl }) => {
            if (metricImpl && metricImpl.reset) {
              this._log.debug('Resetting metric', name);
              metricImpl.reset();
            }
          });
        }
      });
    } catch (error) {
      this._log.error('Failed to send metrics to signal fx', error);
    }
  }

  /**
   * This method gets called with an array of {@link MetricWrapper} on an interval, when metrics should be reported.
   *
   * This is the main method that needs to get implemented when created an aggregator specific reporter.
   *
   * @param {MetricWrapper[]} metrics The array of metrics to report.
   * @protected
   * @abstract
   */
  _reportMetrics(metrics) {
    throw new TypeError('Abstract method _reportMetrics(metrics) must be implemented in implementation class');
  }

  /**
   *
   * @param {MetricWrapper} metric The Wrapped Metric Object.
   * @return {Dimensions} The left merged default dimensions with the metric specific dimensions
   * @protected
   */
  _getDimensions(metric) {
    return Object.assign({}, this._defaultDimensions, metric.dimensions);
  }

  /**
   * Clears the intervals that are running to report metrics at an interval, and resets the state.
   */
  shutdown() {
    this._intervals.forEach(interval => clearInterval(interval));
    this._intervals = [];
    this._intervalToMetric = {};
  }
}

/**
 * Options for creating a {@link Reporter}
 * @interface ReporterOptions
 * @typedef ReporterOptions
 * @type {Object}
 * @property {Dimensions} defaultDimensions A dictionary of dimensions to include with every metric reported
 * @property {Logger} logger The logger to use, if not supplied a new Buynan logger will be created
 * @property {string} logLevel The log level to use with the created console logger if you didn't supply your own logger.
 * @property {number} defaultReportingIntervalInSeconds The default reporting interval to use if non is supplied when registering a metric, defaults to 10 seconds.
 * @property {boolean} unrefTimers Indicate if reporting timers should be unref'd, defaults to false.
 * @property {boolean} resetMetricsOnInterval Indicate if metrics should be reset on each reporting interval, defaults to false.
 */

module.exports = Reporter;
