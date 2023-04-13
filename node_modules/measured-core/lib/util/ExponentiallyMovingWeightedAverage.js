const units = require('./units');

const TICK_INTERVAL = 5 * units.SECONDS;

/**
 * ExponentiallyMovingWeightedAverage
 */
class ExponentiallyMovingWeightedAverage {
  constructor(timePeriod, tickInterval) {
    this._timePeriod = timePeriod || units.MINUTE;
    this._tickInterval = tickInterval || TICK_INTERVAL;
    this._alpha = 1 - Math.exp(-this._tickInterval / this._timePeriod);
    this._count = 0;
    this._rate = 0;
  }

  update(n) {
    this._count += n;
  }

  tick() {
    const instantRate = this._count / this._tickInterval;
    this._count = 0;

    this._rate += this._alpha * (instantRate - this._rate);
  }

  rate(timeUnit) {
    return (this._rate || 0) * timeUnit;
  }
}

module.exports = ExponentiallyMovingWeightedAverage;
