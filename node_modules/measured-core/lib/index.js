const Collection = require('./Collection');
const Counter = require('./metrics/Counter');
const Gauge = require('./metrics/Gauge');
const SettableGauge = require('./metrics/SettableGauge');
const CachedGauge = require('./metrics/CachedGauge');
const Histogram = require('./metrics/Histogram');
const Meter = require('./metrics/Meter');
const NoOpMeter = require('./metrics/NoOpMeter');
const Timer = require('./metrics/Timer');
const BinaryHeap = require('./util/BinaryHeap');
const ExponentiallyDecayingSample = require('./util/ExponentiallyDecayingSample');
const ExponentiallyMovingWeightedAverage = require('./util/ExponentiallyMovingWeightedAverage');
const Stopwatch = require('./util/Stopwatch');
const units = require('./util/units');
const { MetricTypes } = require('./metrics/Metric');
const metricValidators = require('./validators/metricValidators');

/**
 * The main measured-core module that is referenced when require('measured-core') is used.
 * @module measured-core
 */
module.exports = {
  /**
   * See {@link Collection}
   * @type {Collection}
   */
  Collection,

  /**
   * See {@link Counter}
   * @type {Counter}
   */
  Counter,

  /**
   * See {@link Gauge}
   * @type {Gauge}
   */
  Gauge,

  /**
   * See {@link SettableGauge}
   * @type {SettableGauge}
   */
  SettableGauge,

  /**
   * See {@link CachedGauge}
   * @type {CachedGauge}
   */
  CachedGauge,

  /**
   * See {@link Histogram}
   * @type {Histogram}
   */
  Histogram,

  /**
   * See {@link Meter}
   * @type {Meter}
   */
  Meter,

  /**
   * See {@link NoOpMeter}
   * @type {NoOpMeter}
   */
  NoOpMeter,

  /**
   * See {@link Timer}
   * @type {Timer}
   */
  Timer,

  /**
   * See {@link BinaryHeap}
   * @type {BinaryHeap}
   */
  BinaryHeap,

  /**
   * See {@link ExponentiallyDecayingSample}
   * @type {ExponentiallyDecayingSample}
   */
  ExponentiallyDecayingSample,

  /**
   * See {@link ExponentiallyMovingWeightedAverage}
   * @type {ExponentiallyMovingWeightedAverage}
   */
  ExponentiallyMovingWeightedAverage,

  /**
   * See {@link Stopwatch}
   * @type {Stopwatch}
   */
  Stopwatch,

  /**
   * See {@link MetricTypes}
   * @type {MetricTypes}
   */
  MetricTypes,

  /**
   * See {@link units}
   * @type {units}
   */
  units,

  /**
   * See {@link units}
   * @type {units}
   */
  TimeUnits: units,

  /**
   * See {@link module:metricValidators}
   * @type {Object.<string, function>}
   */
  metricValidators,

  /**
   * Creates a named collection. See {@link Collection} for more details
   *
   * @param name The name for the collection
   * @return {Collection}
   */
  createCollection: name => {
    return new Collection(name);
  }
};
