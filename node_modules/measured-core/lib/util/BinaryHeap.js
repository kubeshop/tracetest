/**
 * Based on http://en.wikipedia.org/wiki/Binary_Heap
 * as well as http://eloquentjavascript.net/appendix2.html
 */
class BinaryHeap {
  constructor(options) {
    options = options || {};

    this._elements = options.elements || [];
    this._score = options.score || this._score;
  }

  /**
   * Add elements to the binary heap.
   * @param {any[]} elements
   */
  add(...elements) {
    elements.forEach(element => {
      this._elements.push(element);
      this._bubble(this._elements.length - 1);
    });
  }

  first() {
    return this._elements[0];
  }

  removeFirst() {
    const root = this._elements[0];
    const last = this._elements.pop();

    if (this._elements.length > 0) {
      this._elements[0] = last;
      this._sink(0);
    }

    return root;
  }

  clone() {
    return new BinaryHeap({
      elements: this.toArray(),
      score: this._score
    });
  }

  toSortedArray() {
    const array = [];
    const clone = this.clone();
    let element;

    while (true) {
      element = clone.removeFirst();
      if (element === undefined) {
        break;
      }

      array.push(element);
    }

    return array;
  }

  toArray() {
    return [].concat(this._elements);
  }

  size() {
    return this._elements.length;
  }

  _bubble(bubbleIndex) {
    const bubbleElement = this._elements[bubbleIndex];
    const bubbleScore = this._score(bubbleElement);
    let parentIndex;
    let parentElement;
    let parentScore;

    while (bubbleIndex > 0) {
      parentIndex = this._parentIndex(bubbleIndex);
      parentElement = this._elements[parentIndex];
      parentScore = this._score(parentElement);

      if (bubbleScore <= parentScore) {
        break;
      }

      this._elements[parentIndex] = bubbleElement;
      this._elements[bubbleIndex] = parentElement;
      bubbleIndex = parentIndex;
    }
  }

  _sink(sinkIndex) {
    const sinkElement = this._elements[sinkIndex];
    const sinkScore = this._score(sinkElement);
    const { length } = this._elements;
    let swapIndex;
    let swapScore;
    let swapElement;
    let childIndexes;
    let i;
    let childIndex;
    let childElement;
    let childScore;

    while (true) {
      swapIndex = null;
      swapScore = null;
      swapElement = null;
      childIndexes = this._childIndexes(sinkIndex);

      for (i = 0; i < childIndexes.length; i++) {
        childIndex = childIndexes[i];

        if (childIndex >= length) {
          break;
        }

        childElement = this._elements[childIndex];
        childScore = this._score(childElement);

        if (childScore > sinkScore) {
          if (swapScore === null || swapScore < childScore) {
            swapIndex = childIndex;
            swapScore = childScore;
            swapElement = childElement;
          }
        }
      }

      if (swapIndex === null) {
        break;
      }

      this._elements[swapIndex] = sinkElement;
      this._elements[sinkIndex] = swapElement;
      sinkIndex = swapIndex;
    }
  }

  _parentIndex(index) {
    return Math.floor((index - 1) / 2);
  }

  _childIndexes(index) {
    return [2 * index + 1, 2 * index + 2];
  }

  _score(element) {
    return element.valueOf();
  }
}

module.exports = BinaryHeap;
