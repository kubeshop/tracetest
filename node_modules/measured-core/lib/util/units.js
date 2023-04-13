const NANOSECONDS = 1 / (1000 * 1000);
const MICROSECONDS = 1 / 1000;
const MILLISECONDS = 1;
const SECONDS = 1000 * MILLISECONDS;
const MINUTES = 60 * SECONDS;
const HOURS = 60 * MINUTES;
const DAYS = 24 * HOURS;

/**
 * Time units, as found in Java: {@link http://download.oracle.com/javase/6/docs/api/java/util/concurrent/TimeUnit.html}
 * @module timeUnits
 * @example
 * const timeUnit = require('measured-core').unit
 * setTimeout(() => {}, 5 * timeUnit.MINUTES)
 */
module.exports = {
  /**
   * nanoseconds in milliseconds
   * @type {number}
   */
  NANOSECONDS,
  /**
   * microseconds in milliseconds
   * @type {number}
   */
  MICROSECONDS,
  /**
   * milliseconds in milliseconds
   * @type {number}
   */
  MILLISECONDS,
  /**
   * seconds in milliseconds
   * @type {number}
   */
  SECONDS,
  /**
   * minutes in milliseconds
   * @type {number}
   */
  MINUTES,
  /**
   * hours in milliseconds
   * @type {number}
   */
  HOURS,
  /**
   * days in milliseconds
   * @type {number}
   */
  DAYS
};
