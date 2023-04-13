const consoleLogLevel = require('console-log-level');
const { CachedGauge, SettableGauge, Gauge, Timer, Counter, Meter, Histogram } = require('measured-core');
const DimensionAwareMetricsRegistry = require('./DimensionAwareMetricsRegistry');
const {
  validateSelfReportingMetricsRegistryParameters,
  validateRegisterOptions,
  validateGaugeOptions,
  validateCounterOptions,
  validateHistogramOptions,
  validateTimerOptions,
  validateSettableGaugeOptions,
  validateCachedGaugeOptions
} = require('../validators/inputValidators');

function prefix() {
  return `${new Date().toISOString()}: `;
}

/**
 * A dimensional aware self-reporting metrics registry
 */
class SelfReportingMetricsRegistry {
  /**
   * @param {Reporter|Reporter[]} reporters A single {@link Reporter} or an array of reporters that will be used to report metrics on an interval.
   * @param {SelfReportingMetricsRegistryOptions} [options] Configurable options for the Self Reporting Metrics Registry
   */
  constructor(reporters, options) {
    options = options || {};

    if (!Array.isArray(reporters)) {
      reporters = [reporters];
    }

    validateSelfReportingMetricsRegistryParameters(reporters, options);

    /**
     * @type {Reporter}
     * @protected
     */
    this._reporters = reporters;

    /**
     * @type {DimensionAwareMetricsRegistry}
     * @protected
     */
    this._registry = options.registry || new DimensionAwareMetricsRegistry();
    this._reporters.forEach(reporter => reporter.setRegistry(this._registry));

    /**
     * Loggers to use, defaults to a new console logger if nothing is supplied in options
     * @type {Logger}
     * @protected
     */
    this._log =
      options.logger ||
      consoleLogLevel({ name: 'SelfReportingMetricsRegistry', level: options.logLevel || 'info', prefix: prefix });
  }

  /**
   * Registers a manually created Metric.
   *
   * @param {string} name The Metric name
   * @param {Metric} metric The {@link Metric} to register
   * @param {Dimensions} [dimensions] any custom {@link Dimensions} for the Metric
   * @param {number} [publishingIntervalInSeconds] a optional custom publishing interval
   * @example
   * const settableGauge = new SettableGauge(5);
   * // register the gauge and have it report to every 10 seconds
   * registry.register('my-gauge', settableGauge, {}, 10);
   * interval(() => {
   *    // such as cpu % used
   *    determineAValueThatCannotBeSync((value) => {
   *      settableGauge.update(value);
   *    })
   * }, 10000)
   */
  register(name, metric, dimensions, publishingIntervalInSeconds) {
    validateRegisterOptions(name, metric, dimensions, publishingIntervalInSeconds);

    if (this._registry.hasMetric(name, dimensions)) {
      throw new Error(
        `Metric with name: ${name} and dimensions: ${JSON.stringify(dimensions)} has already been registered`
      );
    } else {
      const key = this._registry.putMetric(name, metric, dimensions);
      this._reporters.forEach(reporter => reporter.reportMetricOnInterval(key, publishingIntervalInSeconds));
    }
    return metric;
  }

  /**
   * Creates a {@link Gauge} or gets the existing Gauge for a given name and dimension combo
   *
   * @param {string} name The Metric name
   * @param {function} callback The callback that will return a value to report to signal fx
   * @param {Dimensions} [dimensions] any custom {@link Dimensions} for the Metric
   * @param {number} [publishingIntervalInSeconds] a optional custom publishing interval
   * @return {Gauge}
   * @example
   * // https://nodejs.org/api/process.html#process_process_memoryusage
   * // Report heap total and heap used at the default interval
   * registry.getOrCreateGauge(
   *   'process-memory-heap-total',
   *   () => {
   *     return process.memoryUsage().heapTotal
   *   }
   * );
   * registry.getOrCreateGauge(
   *   'process-memory-heap-used',
   *   () => {
   *     return process.memoryUsage().heapUsed
   *   }
   * )
   */
  getOrCreateGauge(name, callback, dimensions, publishingIntervalInSeconds) {
    validateGaugeOptions(name, callback, dimensions, publishingIntervalInSeconds);
    let gauge;
    if (this._registry.hasMetric(name, dimensions)) {
      gauge = this._registry.getMetric(name, dimensions);
    } else {
      gauge = new Gauge(callback);
      const key = this._registry.putMetric(name, gauge, dimensions);
      this._reporters.forEach(reporter => reporter.reportMetricOnInterval(key, publishingIntervalInSeconds));
    }
    return gauge;
  }

  /**
   * Creates a {@link Histogram} or gets the existing Histogram for a given name and dimension combo
   *
   * @param {string} name The Metric name
   * @param {Dimensions} [dimensions] any custom {@link Dimensions} for the Metric
   * @param {number} [publishingIntervalInSeconds] a optional custom publishing interval
   * @return {Histogram}
   */
  getOrCreateHistogram(name, dimensions, publishingIntervalInSeconds) {
    validateHistogramOptions(name, dimensions, publishingIntervalInSeconds);

    let histogram;
    if (this._registry.hasMetric(name, dimensions)) {
      histogram = this._registry.getMetric(name, dimensions);
    } else {
      histogram = new Histogram();
      const key = this._registry.putMetric(name, histogram, dimensions);
      this._reporters.forEach(reporter => reporter.reportMetricOnInterval(key, publishingIntervalInSeconds));
    }

    return histogram;
  }

  /**
   * Creates a {@link Meter} or gets the existing Meter for a given name and dimension combo
   *
   * @param {string} name The Metric name
   * @param {Dimensions} [dimensions] any custom {@link Dimensions} for the Metric
   * @param {number} [publishingIntervalInSeconds] a optional custom publishing interval
   * @return {Meter}
   */
  getOrCreateMeter(name, dimensions, publishingIntervalInSeconds) {
    // todo validate options
    let meter;
    if (this._registry.hasMetric(name, dimensions)) {
      meter = this._registry.getMetric(name, dimensions);
    } else {
      meter = new Meter();
      const key = this._registry.putMetric(name, meter, dimensions);
      this._reporters.forEach(reporter => reporter.reportMetricOnInterval(key, publishingIntervalInSeconds));
    }

    return meter;
  }

  /**
   * Creates a {@link Counter} or gets the existing Counter for a given name and dimension combo
   *
   * @param {string} name The Metric name
   * @param {Dimensions} [dimensions] any custom {@link Dimensions} for the Metric
   * @param {number} [publishingIntervalInSeconds] a optional custom publishing interval
   * @return {Counter}
   */
  getOrCreateCounter(name, dimensions, publishingIntervalInSeconds) {
    validateCounterOptions(name, dimensions, publishingIntervalInSeconds);

    let counter;
    if (this._registry.hasMetric(name, dimensions)) {
      counter = this._registry.getMetric(name, dimensions);
    } else {
      counter = new Counter();
      const key = this._registry.putMetric(name, counter, dimensions);
      this._reporters.forEach(reporter => reporter.reportMetricOnInterval(key, publishingIntervalInSeconds));
    }

    return counter;
  }

  /**
   * Creates a {@link Timer} or gets the existing Timer for a given name and dimension combo.
   *
   * @param {string} name The Metric name
   * @param {Dimensions} [dimensions] any custom {@link Dimensions} for the Metric
   * @param {number} [publishingIntervalInSeconds] a optional custom publishing interval
   * @return {Timer}
   */
  getOrCreateTimer(name, dimensions, publishingIntervalInSeconds) {
    validateTimerOptions(name, dimensions, publishingIntervalInSeconds);

    let timer;
    if (this._registry.hasMetric(name, dimensions)) {
      timer = this._registry.getMetric(name, dimensions);
    } else {
      timer = new Timer();
      const key = this._registry.putMetric(name, timer, dimensions);
      this._reporters.forEach(reporter => reporter.reportMetricOnInterval(key, publishingIntervalInSeconds));
    }

    return timer;
  }

  /**
   * Creates a {@link SettableGauge} or gets the existing SettableGauge for a given name and dimension combo.
   *
   * @param {string} name The Metric name
   * @param {Dimensions} [dimensions] any custom {@link Dimensions} for the Metric
   * @param {number} [publishingIntervalInSeconds] a optional custom publishing interval
   * @return {SettableGauge}
   */
  getOrCreateSettableGauge(name, dimensions, publishingIntervalInSeconds) {
    validateSettableGaugeOptions(name, dimensions, publishingIntervalInSeconds);

    let settableGauge;
    if (this._registry.hasMetric(name, dimensions)) {
      settableGauge = this._registry.getMetric(name, dimensions);
    } else {
      settableGauge = new SettableGauge();
      const key = this._registry.putMetric(name, settableGauge, dimensions);
      this._reporters.forEach(reporter => reporter.reportMetricOnInterval(key, publishingIntervalInSeconds));
    }

    return settableGauge;
  }

  /**
   * Creates a {@link CachedGauge} or gets the existing CachedGauge for a given name and dimension combo.
   *
   * @param {string} name The Metric name.
   * @param {function} valueProducingPromiseCallback.
   * @param {number} cachedGaugeUpdateIntervalInSeconds.
   * @param {Dimensions} [dimensions] any custom {@link Dimensions} for the Metric.
   * @param {number} [publishingIntervalInSeconds] a optional custom publishing interval.
   * @return {CachedGauge}
   */
  getOrCreateCachedGauge(
    name,
    valueProducingPromiseCallback,
    cachedGaugeUpdateIntervalInSeconds,
    dimensions,
    publishingIntervalInSeconds
  ) {
    validateCachedGaugeOptions(name, valueProducingPromiseCallback, dimensions, publishingIntervalInSeconds);

    let cachedGauge;
    if (this._registry.hasMetric(name, dimensions)) {
      cachedGauge = this._registry.getMetric(name, dimensions);
    } else {
      cachedGauge = new CachedGauge(valueProducingPromiseCallback, cachedGaugeUpdateIntervalInSeconds);
      const key = this._registry.putMetric(name, cachedGauge, dimensions);
      this._reporters.forEach(reporter => reporter.reportMetricOnInterval(key, publishingIntervalInSeconds));
    }

    return cachedGauge;
  }

  /**
   * Calls end on all metrics in the registry that support end() and calls end on the reporter
   */
  shutdown() {
    // shutdown the reporter
    this._reporters.forEach(reporter => reporter.shutdown());
    // shutdown any metrics that have an end method
    this._registry.allKeys().forEach(key => {
      const metricWrapper = this._registry.getMetricWrapperByKey(key);
      if (metricWrapper.metricImpl.end) {
        metricWrapper.metricImpl.end();
      }
    });
  }
}

module.exports = SelfReportingMetricsRegistry;

/**
 * Configurable options for the Self Reporting Metrics Registry
 *
 * @interface SelfReportingMetricsRegistryOptions
 * @typedef SelfReportingMetricsRegistryOptions
 * @property {Logger} logger the Logger to use
 * @property {string} logLevel The Log level to use if defaulting to included logger
 * @property {DimensionAwareMetricsRegistry} registry The registry to use, defaults to new DimensionAwareMetricsRegistry
 */
