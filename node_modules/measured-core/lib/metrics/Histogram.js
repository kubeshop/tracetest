const { MetricTypes } = require('./Metric');
const binarySearch = require('binary-search');
const EDS = require('../util/ExponentiallyDecayingSample');

/**
 * Keeps a reservoir of statistically relevant values biased towards the last 5 minutes to explore their distribution.
 * @implements {Metric}
 * @example
 * var Measured = require('measured')
 * var histogram = new Measured.Histogram();
 * http.createServer(function(req, res) {
 *   if (req.headers['content-length']) {
 *     histogram.update(parseInt(req.headers['content-length'], 10));
 *   }
 * });
 */
class Histogram {
  /**
   @param {HistogramProperties} [properties] see {@link HistogramProperties}.
   */
  constructor(properties) {
    this._properties = properties || {};
    this._initializeState();
  }

  _initializeState() {
    this._sample = this._properties.sample || new EDS();
    this._percentilesMethod = this._properties.percentilesMethod || this._percentiles;
    this._min = null;
    this._max = null;
    this._count = 0;
    this._sum = 0;

    // These are for the Welford algorithm for calculating running constiance
    // without floating-point doom.
    this._constianceM = 0;
    this._constianceS = 0;
  }

  /**
   * Pushes value into the sample. timestamp defaults to Date.now().
   * @param {number} value
   */
  update(value) {
    this._count++;
    this._sum += value;

    this._sample.update(value);
    this._updateMin(value);
    this._updateMax(value);
    this._updateVariance(value);
  }

  _percentiles(percentiles) {
    const values = this._sample.toArray().sort((a, b) => {
      return a === b ? 0 : a - b;
    });

    const results = {};

    let i, percentile, pos, lower, upper;
    for (i = 0; i < percentiles.length; i++) {
      percentile = percentiles[i];
      if (values.length) {
        pos = percentile * (values.length + 1);
        if (pos < 1) {
          results[percentile] = values[0];
        } else if (pos >= values.length) {
          results[percentile] = values[values.length - 1];
        } else {
          lower = values[Math.floor(pos) - 1];
          upper = values[Math.ceil(pos) - 1];
          results[percentile] = lower + (pos - Math.floor(pos)) * (upper - lower);
        }
      } else {
        results[percentile] = null;
      }
    }

    return results;
  }

  weightedPercentiles(percentiles) {
    const values = this._sample.toArrayWithWeights().sort((a, b) => {
      return a.value === b.value ? 0 : a.value - b.value;
    });

    const sumWeight = values.reduce((sum, sample) => {
      return sum + sample.priority;
    }, 0);

    const normWeights = values.map(value => {
      return value.priority / sumWeight;
    });

    const quantiles = [0];
    let i;
    for (i = 1; i < values.length; i++) {
      quantiles[i] = quantiles[i - 1] + normWeights[i - 1];
    }

    function gt(a, b) {
      return a - b;
    }

    const results = {};
    let percentile, pos;
    for (i = 0; i < percentiles.length; i++) {
      percentile = percentiles[i];
      if (values.length) {
        pos = binarySearch(quantiles, percentile, gt);
        if (pos < 0) {
          results[percentile] = values[-pos - 1 - 1].value;
        } else if (pos < 1) {
          results[percentile] = values[0].value;
        } else if (pos >= values.length) {
          results[percentile] = values[values.length - 1].value;
        }
      } else {
        results[percentile] = null;
      }
    }
    return results;
  }

  /**
   * Resets all values. Histograms initialized with custom options will be reset to the default settings (patch welcome).
   */
  reset() {
    // while this is technically a bug?, copying existing logic to maintain current api,
    // TODO reset should reset the sample, not override it with a new EDS()
    this._properties.sample = new EDS();

    this._initializeState();
  }

  /**
   * Checks whether the histogram contains values.
   * @return {boolean} Whether the histogram contains values.
   */
  hasValues() {
    return this._count > 0;
  }

  /**
   * @return {HistogramData}
   */
  toJSON() {
    const percentiles = this._percentilesMethod([0.5, 0.75, 0.95, 0.99, 0.999]);

    return {
      min: this._min,
      max: this._max,
      sum: this._sum,
      variance: this._calculateVariance(),
      mean: this._calculateMean(),
      stddev: this._calculateStddev(),
      count: this._count,
      median: percentiles[0.5],
      p75: percentiles[0.75],
      p95: percentiles[0.95],
      p99: percentiles[0.99],
      p999: percentiles[0.999]
    };
  }

  _updateMin(value) {
    if (this._min === null || value < this._min) {
      this._min = value;
    }
  }

  _updateMax(value) {
    if (this._max === null || value > this._max) {
      this._max = value;
    }
  }

  _updateVariance(value) {
    if (this._count === 1) {
      this._constianceM = value;
      return value;
    }

    const oldM = this._constianceM;

    this._constianceM += (value - oldM) / this._count;
    this._constianceS += (value - oldM) * (value - this._constianceM);

    // TODO is this right, above it returns in the if statement but does nothing but update internal state for the else case?
    return undefined;
  }

  /**
   *
   * @return {number|null}
   * @private
   */
  _calculateMean() {
    return this._count === 0 ? 0 : this._sum / this._count;
  }

  /**
   * @return {number|null}
   * @private
   */
  _calculateVariance() {
    return this._count <= 1 ? null : this._constianceS / (this._count - 1);
  }

  /**
   * @return {number|null}
   * @private
   */
  _calculateStddev() {
    return this._count < 1 ? null : Math.sqrt(this._calculateVariance());
  }

  /**
   * The type of the Metric Impl. {@link MetricTypes}.
   * @return {string} The type of the Metric Impl.
   */
  getType() {
    return MetricTypes.HISTOGRAM;
  }
}

module.exports = Histogram;

/**
 * Properties to create a {@link Histogram} with.
 *
 * @interface HistogramProperties
 * @typedef HistogramProperties
 * @type {Object}
 * @property {object} sample The sample reservoir to use. Defaults to an ExponentiallyDecayingSample.
 */

/**
 * The data returned from Histogram::toJSON()
 * @interface HistogramData
 * @typedef HistogramData
 * @typedef {object}
 * @property {number|null} min The lowest observed value.
 * @property {number|null} max The highest observed value.
 * @property {number|null} sum The sum of all observed values.
 * @property {number|null} variance The variance of all observed values.
 * @property {number|null} mean The average of all observed values.
 * @property {number|null} stddev The stddev of all observed values.
 * @property {number} count The number of observed values.
 * @property {number} median 50% of all values in the resevoir are at or below this value.
 * @property {number} p75 See median, 75% percentile.
 * @property {number} p95 See median, 95% percentile.
 * @property {number} p99 See median, 99% percentile.
 * @property {number} p999 See median, 99.9% percentile.
 */
