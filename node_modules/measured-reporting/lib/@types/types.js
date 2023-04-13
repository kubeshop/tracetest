/**
 * A wrapper object around a {@link Metric}, {@link Dimensions} and the metric name
 *
 * @interface MetricWrapper
 * @typedef MetricWrapper
 * @type {Object}
 * @property {string} name The supplied name of the Metric
 * @property {Metric} metricImpl The {@link Metric} object
 * @property {Dimensions} dimensions The {@link Dimensions} for the given {@link Metric}
 */

/**
 * A Dictionary of string, string key value pairs
 *
 * @interface Dimensions
 * @typedef Dimensions
 * @type {Object.<string, string>}
 *
 * @example
 * {
 *   path: "/api/foo"
 *   method: "GET"
 *   statusCode: "200"
 * }
 */
