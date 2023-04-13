const BinaryHeap = require('./BinaryHeap');
const units = require('./units');

const RESCALE_INTERVAL = units.HOURS;
const ALPHA = 0.015;
const SIZE = 1028;

/**
 * ExponentiallyDecayingSample
 */
class ExponentiallyDecayingSample {
  constructor(options) {
    options = options || {};

    this._elements = new BinaryHeap({
      score: element => -element.priority
    });

    this._rescaleInterval = options.rescaleInterval || RESCALE_INTERVAL;
    this._alpha = options.alpha || ALPHA;
    this._size = options.size || SIZE;
    this._random = options.random || this._random;
    this._landmark = null;
    this._nextRescale = null;
  }

  update(value, timestamp) {
    const now = Date.now();
    if (!this._landmark) {
      this._landmark = now;
      this._nextRescale = this._landmark + this._rescaleInterval;
    }

    timestamp = timestamp || now;

    const newSize = this._elements.size() + 1;

    const element = {
      priority: this._priority(timestamp - this._landmark),
      value: value
    };

    if (newSize <= this._size) {
      this._elements.add(element);
    } else if (element.priority > this._elements.first().priority) {
      this._elements.removeFirst();
      this._elements.add(element);
    }

    if (now >= this._nextRescale) {
      this._rescale(now);
    }
  }

  toSortedArray() {
    return this._elements.toSortedArray().map(element => element.value);
  }

  toArray() {
    return this._elements.toArray().map(element => element.value);
  }

  toArrayWithWeights() {
    return this._elements.toArray();
  }

  _weight(age) {
    // We divide by 1000 to not run into huge numbers before reaching a
    // rescale event.
    return Math.exp(this._alpha * (age / 1000));
  }

  _priority(age) {
    return this._weight(age) / this._random();
  }

  _random() {
    return Math.random();
  }

  _rescale(now) {
    now = now || Date.now();

    const self = this;
    const oldLandmark = this._landmark;
    this._landmark = now || Date.now();
    this._nextRescale = now + this._rescaleInterval;

    const factor = self._priority(-(self._landmark - oldLandmark));

    this._elements.toArray().forEach(element => {
      element.priority *= factor;
    });
  }
}

module.exports = ExponentiallyDecayingSample;
