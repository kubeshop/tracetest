const SelfReportingMetricsRegistry = require('./registries/SelfReportingMetricsRegistry');
const Reporter = require('./reporters/Reporter');
const LoggingReporter = require('./reporters/LoggingReporter');
const inputValidators = require('./validators/inputValidators');

/**
 * The main measured module that is referenced when require('measured-reporting') is used.
 * @module measured-reporting
 */
module.exports = {
  /**
   * The Self Reporting Metrics Registry Class.
   *
   * @type {SelfReportingMetricsRegistry}
   */
  SelfReportingMetricsRegistry,

  /**
   * The abstract / base Reporter class.
   *
   * @type {Reporter}
   */
  Reporter,

  /**
   * The basic included reference reporter, simply logs the metrics.
   * See {ReporterOptions} for options.
   *
   * @type {LoggingReporter}
   */
  LoggingReporter,

  /**
   * Various Input Validation functions.
   *
   * @type {inputValidators}
   */
  inputValidators
};
