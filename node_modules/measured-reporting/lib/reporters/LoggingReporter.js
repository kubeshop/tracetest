const Reporter = require('./Reporter');

/**
 * A reporter impl that simply logs the metrics via the Logger.
 *
 * @example
 * const { SelfReportingMetricsRegistry, LoggingReporter } = require('measured-reporting');
 * const registry = new SelfReportingMetricsRegistry(new LoggingReporter());
 *
 * @extends {Reporter}
 */
class LoggingReporter extends Reporter {
  /**
   * @param {LoggingReporterOptions} [options]
   */
  constructor(options) {
    super(options);
    const level = (options || {}).logLevelToLogAt;
    this._logLevel = (level || 'info').toLowerCase();
  }

  /**
   * Logs the metrics via the inherited logger instance.
   * @param {MetricWrapper[]} metrics
   * @protected
   */
  _reportMetrics(metrics) {
    metrics.forEach(metric => {
      this._log[this._logLevel](
        JSON.stringify({
          metricName: metric.name,
          dimensions: this._getDimensions(metric),
          data: metric.metricImpl.toJSON()
        })
      );
    });
  }
}

module.exports = LoggingReporter;

/**
 * @interface LoggingReporterOptions
 * @typedef LoggingReporterOptions
 * @type {Object}
 * @property {Dimensions} defaultDimensions A dictionary of dimensions to include with every metric reported
 * @property {Logger} [logger] The logger to use, if not supplied a new Buynan logger will be created
 * @property {string} [logLevel] The log level to use with the created console logger if you didn't supply your own logger.
 * @property {number} [defaultReportingIntervalInSeconds] The default reporting interval to use if non is supplied when registering a metric, defaults to 10 seconds.
 * @property {string} [logLevelToLogAt] You can specify the log level ['debug', 'info', 'warn', 'error'] that this reporter will use when logging the metrics via the logger.
 */
