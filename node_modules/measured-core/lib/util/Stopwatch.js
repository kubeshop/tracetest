const { EventEmitter } = require('events');

/**
 * A simple object for tracking elapsed time
 *
 * @extends {EventEmitter}
 */
class Stopwatch extends EventEmitter {
  /**
   * Creates a started Stopwatch
   * @param {StopwatchProperties} [options] See {@link StopwatchProperties}
   */
  constructor(options) {
    super();
    options = options || {};
    EventEmitter.call(this);

    if (options.getTime) {
      this._getTime = options.getTime;
    }
    this._start = this._getTime();
    this._ended = false;
  }

  /**
   * Called to mark the end of the timer task
   * @return {number} the total execution time
   */
  end() {
    if (this._ended) {
      return null;
    }

    this._ended = true;
    const elapsed = this._getTime() - this._start;

    this.emit('end', elapsed);
    return elapsed;
  }

  _getTime() {
    if (!process.hrtime) {
      return Date.now();
    }

    const hrtime = process.hrtime();
    return hrtime[0] * 1000 + hrtime[1] / (1000 * 1000);
  }
}

module.exports = Stopwatch;

/**
 * @interface StopwatchProperties
 * @typedef StopwatchProperties
 * @type {Object}
 * @property {function} getTime optional function override for supplying time., defaults to new Date() / process.hrt()
 */
