/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'
const { EventEmitter } = require('events')

/**
 * Coordinates fetching of metadata from multiple providers
 *
 * Implements event based coordination for fetching metadata
 * from multiple providers.  The first provider to return
 * a non-error result "wins".  When this happens the CallbackCoordination
 * object will emit a `result` event.
 *
 * If all the metadata providers fail to return a result, then the
 * object will emit an error event indicating failure to collect metadata
 * from any event.
 *
 * Since a scheduled callback may (accidently) fail to call its own
 * callback function, the CallbackCoordination object includes its
 * own timeout timer to avoid deadlock situation.
 */
class CallbackCoordination extends EventEmitter {
  constructor (maxWaitMS = -1, logger) {
    super()

    this.logger = logger
    // how many results have we seen
    this.resultCount = 0
    this.expectedResults = 0
    this.errors = []
    this.scheduled = []
    this.done = false
    this.started = false
    this.timeout = null
    if (maxWaitMS !== -1) {
      this.timeout = setTimeout(() => {
        if (!this.done) {
          this.complete()
          this.logger.warn('cloud metadata requests timed out, using default values instead')
          const error = new CallbackCoordinationError(
            'callback coordination reached timeout',
            this.errors
          )
          this.emit('error', error)
        }
      }, maxWaitMS)
    }
  }

  /**
   * Finishes coordination
   *
   * Marks as `done`, cleans up timeout (if necessary)
   */
  complete () {
    this.done = true
    if (this.timeout) {
      clearTimeout(this.timeout)
    }
  }

  /**
   * Accepts and schedules a callback function
   *
   * Callback will be in the form
   *     function(fetcher) {
   *         //... work to fetch data ...
   *
   *         // this callback calls the recordResult method
   *         // method of the fetcher
   *         fetcher.recordResult(error, result)
   *     }
   *
   * Method also increments expectedResults counter to keep track
   * of how many callbacks we've scheduled
   */
  schedule (fetcherCallback) {
    if (this.started) {
      this.logger.error('Can not schedule callback, already started')
      return
    }

    this.expectedResults++
    this.scheduled.push(fetcherCallback)
  }

  /**
   * Starts processing of the callbacks scheduled by the `schedule` method
   */
  start () {
    this.started = true
    // if called with nothing, send an error through so we don't hang
    if (this.scheduled.length === 0) {
      const error = new CallbackCoordinationError('no callbacks to run')
      this.recordResult(error)
    }

    for (const cb of this.scheduled) {
      process.nextTick(cb.bind(null, this))
    }
  }

  /**
   * Receives calls from scheduled callbacks.
   *
   * If called with a non-error, the method will emit a `result` event
   * and include the results as an argument.  Only a single result
   * is emitted -- if other callbacks respond with a result this method
   * will ignore them.
   *
   * If called by _all_ scheduled callbacks without a non-error, this method
   * will issue the error event.
   */
  recordResult (error, result) {
    // console.log('.')
    this.resultCount++
    if (error) {
      this.errors.push(error)
      if (this.resultCount >= this.expectedResults && !this.done) {
        this.complete()
        // we've made every request without success, signal an error
        const error = new CallbackCoordinationError(
          'no response from any callback, no cloud metadata will be set (normal outside of cloud env.)',
          this.errors
        )
        this.logger.debug('no cloud metadata servers responded')
        this.emit('error', error)
      }
    }

    if (!error && result && !this.done) {
      this.complete()
      this.emit('result', result)
    }
  }
}

/**
 * Error for CallbackCoordination class
 *
 * Includes the individual errors from each callback
 * of the CallbackCoordination object
 */
class CallbackCoordinationError extends Error {
  constructor (message, allErrors = []) {
    super(message)
    this.allErrors = allErrors
  }
}

module.exports = {
  CallbackCoordination
}
