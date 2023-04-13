const Optional = require('optional-js');
const { validateMetric } = require('measured-core').metricValidators;

/**
 * This module contains various validators to validate publicly exposed input.
 *
 * @module inputValidators
 */
module.exports = {
  /**
   * Validates @{link Gauge} options.
   *
   * @param {string} name The metric name
   * @param {function} callback The callback for the Gauge
   * @param {Dimensions} dimensions The optional custom dimensions
   * @param {number} publishingIntervalInSeconds the optional publishing interval
   */
  validateGaugeOptions: (name, callback, dimensions, publishingIntervalInSeconds) => {
    module.exports.validateCommonMetricParameters(name, dimensions, publishingIntervalInSeconds);
    module.exports.validateNumberReturningCallback(callback);
  },

  /**
   * Validates @{link Gauge} options.
   *
   * @param {string} name The metric name
   * @param {function} callback The callback for the CachedGauge
   * @param {Dimensions} dimensions The optional custom dimensions
   * @param {number} publishingIntervalInSeconds the optional publishing interval
   */
  validateCachedGaugeOptions: (name, callback, dimensions, publishingIntervalInSeconds) => {
    module.exports.validateCommonMetricParameters(name, dimensions, publishingIntervalInSeconds);
    // Should we validate the promise call back, it may be expensive or produce a race condition in some use-cases.
  },

  /**
   * Validates the create histogram Options.
   *
   * @param {string} name The metric name
   * @param {Dimensions} dimensions The optional custom dimensions
   * @param {number} publishingIntervalInSeconds the optional publishing interval
   */
  validateHistogramOptions: (name, dimensions, publishingIntervalInSeconds) => {
    module.exports.validateCommonMetricParameters(name, dimensions, publishingIntervalInSeconds);
  },

  /**
   * Validates the create counter Options.
   *
   * @param {string} name The metric name
   * @param {Dimensions} dimensions The optional custom dimensions
   * @param {number} publishingIntervalInSeconds the optional publishing interval
   */
  validateCounterOptions: (name, dimensions, publishingIntervalInSeconds) => {
    module.exports.validateCommonMetricParameters(name, dimensions, publishingIntervalInSeconds);
  },

  /**
   * Validates the create timer Options.
   *
   * @param {string} name The metric name
   * @param {Dimensions} dimensions The optional custom dimensions
   * @param {number} publishingIntervalInSeconds the optional publishing interval
   */
  validateTimerOptions: (name, dimensions, publishingIntervalInSeconds) => {
    module.exports.validateCommonMetricParameters(name, dimensions, publishingIntervalInSeconds);
  },

  /**
   * Validates the create timer Options.
   *
   * @param {string} name The metric name
   * @param {Metric} metric The metric instance
   * @param {Dimensions} dimensions The optional custom dimensions
   * @param {number} publishingIntervalInSeconds the optional publishing interval
   */
  validateRegisterOptions: (name, metric, dimensions, publishingIntervalInSeconds) => {
    module.exports.validateMetric(metric);
    module.exports.validateCommonMetricParameters(name, dimensions, publishingIntervalInSeconds);
  },

  /**
   * Validates the create settable gauge Options.
   *
   * @param {string} name The metric name
   * @param {Dimensions} dimensions The optional custom dimensions
   * @param {number} publishingIntervalInSeconds the optional publishing interval
   */
  validateSettableGaugeOptions: (name, dimensions, publishingIntervalInSeconds) => {
    module.exports.validateCommonMetricParameters(name, dimensions, publishingIntervalInSeconds);
  },

  /**
   * Validates the options that are common amoung all create metric methods
   *
   * @param {string} name The metric name
   * @param {Dimensions} dimensions The optional custom dimensions
   * @param {number} publishingIntervalInSeconds the optional publishing interval
   */
  validateCommonMetricParameters: (name, dimensions, publishingIntervalInSeconds) => {
    module.exports.validateMetricName(name);
    module.exports.validateOptionalDimensions(dimensions);
    module.exports.validateOptionalPublishingInterval(publishingIntervalInSeconds);
  },

  /**
   * Validates the metric name.
   *
   * @param name The metric name.
   */
  validateMetricName: name => {
    const type = typeof name;
    if (type !== 'string') {
      throw new TypeError(`options.name is a required option and must be of type string, actual type: ${type}`);
    }
  },

  /**
   * Validates that a metric implements the metric interface.
   *
   * @function
   * @name validateMetric
   * @param {Metric} metric The object that is supposed to be a metric.
   */
  validateMetric,

  /**
   * Validates the provided callback.
   *
   * @param callback The provided callback for a gauge.
   */
  validateNumberReturningCallback: callback => {
    const type = typeof callback;
    if (type !== 'function') {
      throw new TypeError(`options.callback is a required option and must be function, actual type: ${type}`);
    }

    const callbackType = typeof callback();
    if (callbackType !== 'number') {
      throw new TypeError(`options.callback must return a number, actual return type: ${callbackType}`);
    }
  },

  /**
   * Validates a set of optional dimensions
   * @param dimensionsOptional
   */
  validateOptionalDimensions: dimensionsOptional => {
    Optional.ofNullable(dimensionsOptional).ifPresent(dimensions => {
      const type = typeof dimensions;
      if (type !== 'object') {
        throw new TypeError(`options.dimensions should be an object, actual type: ${type}`);
      }

      if (Array.isArray(dimensions)) {
        throw new TypeError('dimensions where detected to be an array, expected Object<string, string>');
      }

      Object.keys(dimensions).forEach(key => {
        const valueType = typeof dimensions[key];
        if (valueType !== 'string') {
          throw new TypeError(`options.dimensions.${key} should be of type string, actual type: ${type}`);
        }
      });
    });
  },

  /**
   * Validates that an optional logger instance at least has the methods we expect.
   * @param loggerOptional
   */
  validateOptionalLogger: loggerOptional => {
    Optional.ofNullable(loggerOptional).ifPresent(logger => {
      if (
        typeof logger.debug !== 'function' ||
        typeof logger.info !== 'function' ||
        typeof logger.warn !== 'function' ||
        typeof logger.error !== 'function'
      ) {
        throw new TypeError(
          'The logger that was passed in does not support all required ' +
            'logging methods, expected object to have functions debug, info, warn, and error with ' +
            'method signatures (...msgs) => {}'
        );
      }
    });
  },

  /**
   * Validates the optional publishing interval.
   *
   * @param publishingIntervalInSecondsOptional The optional publishing interval.
   */
  validateOptionalPublishingInterval: publishingIntervalInSecondsOptional => {
    Optional.ofNullable(publishingIntervalInSecondsOptional).ifPresent(publishingIntervalInSeconds => {
      const type = typeof publishingIntervalInSeconds;
      if (type !== 'number') {
        throw new TypeError(`options.publishingIntervalInSeconds must be of type number, actual type: ${type}`);
      }
    });
  },

  /**
   * Validates optional params for a Reporter
   * @param {ReporterOptions} options The optional params
   */
  validateReporterParameters: options => {
    if (options) {
      module.exports.validateOptionalDimensions(options.defaultDimensions);
      module.exports.validateOptionalLogger(options.logger);

      const type = typeof options.unrefTimers;
      if (type !== 'boolean' && type !== 'undefined') {
        throw new TypeError(`options.unrefTimers should be a boolean or undefined, actual type: ${type}`);
      }
    }
  },

  /**
   * Validates that a valid Reporter object has been supplied
   *
   * @param {Reporter} reporter
   */
  validateReporterInstance: reporter => {
    if (!reporter) {
      throw new TypeError('The reporter was undefined, when it was required');
    }
    if (typeof reporter.setRegistry !== 'function') {
      throw new TypeError(
        'A reporter must implement setRegistry(registry), see the abstract Reporter class in the docs.'
      );
    }
    if (typeof reporter.reportMetricOnInterval !== 'function') {
      throw new TypeError(
        'A reporter must implement reportMetricOnInterval(metricKey, intervalInSeconds), see the abstract Reporter class in the docs.'
      );
    }
  },

  /**
   * Validates the input parameters for a {@link SelfReportingMetricsRegistry}
   * @param {Reporter[]} reporters
   * @param {SelfReportingMetricsRegistryOptions} [options]
   */
  validateSelfReportingMetricsRegistryParameters: (reporters, options) => {
    reporters.forEach(reporter => module.exports.validateReporterInstance(reporter));
    if (options) {
      module.exports.validateOptionalLogger(options.logger);
    }
  }
};
