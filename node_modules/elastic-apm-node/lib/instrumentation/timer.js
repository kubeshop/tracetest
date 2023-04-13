/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

var microtime = require('relative-microtime')

/**
 * Return `time`, if given and valid, or the current time (calculated using
 * the given `timer`).
 *
 * @param {Timer} timer
 * @param {Number} [time] - An optional float number of milliseconds since
 *    the epoch (i.e. the same as provided by `Date.now()`).
 * @returns {Number} A time in *microseconds* since the Epoch.
 */
function maybeTime (timer, time) {
  if (time >= 0) {
    return time * 1000
  } else {
    return timer._timer()
  }
}

module.exports = class Timer {
  // `startTime`: millisecond float
  constructor (parentTimer, startTime) {
    this._parent = parentTimer
    this._timer = parentTimer ? parentTimer._timer : microtime()
    this.start = maybeTime(this, startTime) // microsecond integer (note: might not be an integer)
    this.endTimestamp = null
    this.duration = null // millisecond float
    this.selfTime = null // millisecond float

    // Track child timings to produce self-time
    this.activeChildren = 0
    this.childStart = 0
    this.childDuration = 0

    if (this._parent) {
      this._parent.startChild(startTime)
    }
  }

  startChild (startTime) {
    if (++this.activeChildren === 1) {
      this.childStart = maybeTime(this, startTime)
    }
  }

  endChild (endTime) {
    if (--this.activeChildren === 0) {
      this.incrementChildDuration(endTime)
    }
  }

  incrementChildDuration (endTime) {
    this.childDuration += (maybeTime(this, endTime) - this.childStart) / 1000
    this.childStart = 0
  }

  // `endTime`: millisecond float
  end (endTime) {
    if (this.duration !== null) return
    this.duration = this.elapsed(endTime)
    if (this.activeChildren) {
      this.incrementChildDuration(endTime)
    }
    this.selfTime = this.duration - this.childDuration
    if (this._parent) {
      this._parent.endChild(endTime)
    }
    this.endTimestamp = this._timer()
  }

  // `endTime`: millisecond float
  // returns: millisecond float
  elapsed (endTime) {
    return (maybeTime(this, endTime) - this.start) / 1000
  }
}
