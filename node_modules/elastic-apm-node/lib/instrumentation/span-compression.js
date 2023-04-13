/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'
const STRATEGY_EXACT_MATCH = 'exact_match'
const STRATEGY_SAME_KIND = 'same_kind'

class SpanCompression {
  constructor (agent) {
    this._bufferedSpan = null
    this._agent = agent
    this.timestamp = null
    this.duration = 0
    this.composite = {
      count: 0,
      sum: 0,
      compression_strategy: null
    }
  }

  setBufferedSpan (span) {
    this._bufferedSpan = span
  }

  getBufferedSpan () {
    return this._bufferedSpan
  }

  // Compares two spans and returns which compression strategy to use
  // or returns false if the second span can't be compressed into the
  // first.
  //
  // @param Span compositeSpan
  // @param Span toCompressSpan
  // @returns boolean|String
  _getCompressionStrategy (compositeSpan, toCompressSpan) {
    if (!this._isEnabled() || !compositeSpan._serviceTarget || !toCompressSpan._serviceTarget) {
      return false
    }

    const isSameKind = this.isSameKind(compositeSpan, toCompressSpan)

    if (!isSameKind) {
      return false
    }

    let strategy = STRATEGY_SAME_KIND
    if (compositeSpan.name === toCompressSpan.name) {
      strategy = STRATEGY_EXACT_MATCH
    }

    return strategy
  }

  isSameKind (compositeSpan, toCompressSpan) {
    return compositeSpan.type === toCompressSpan.type &&
      compositeSpan.subtype === toCompressSpan.subtype &&
      compositeSpan._serviceTarget.type === toCompressSpan._serviceTarget.type &&
      compositeSpan._serviceTarget.name === toCompressSpan._serviceTarget.name
  }

  // Sets initial composite values or confirms strategy matches
  //
  // Returns true if spanToCompress can be compressed into compositeSpan,
  // returns false otherwise.
  //
  // @param Span compositeSpan
  // @param Span spanToCompress
  // @returns boolean
  _initCompressionStrategy (compositeSpan, spanToCompress) {
    if (!this.composite.compression_strategy) {
      // If no strategy is set, check if strategizable or not. If so,
      // set initial values.  If not, bail.
      this.composite.compression_strategy = this._getCompressionStrategy(
        compositeSpan,
        spanToCompress
      )
      if (!this.composite.compression_strategy) {
        return false
      }

      // set initial composite context values
      this.timestamp = compositeSpan.timestamp
      this.composite.count = 1
      this.composite.sum = compositeSpan._duration
    } else {
      // if so, compare with the compression strat and bail if mismatch
      const strat = this._getCompressionStrategy(compositeSpan, spanToCompress)
      if (strat !== this.composite.compression_strategy) {
        return false
      }
    }

    return true
  }

  // Attempts to compression the second span into the first
  //
  // @param Span compositeSpan
  // @param Span spanToCompress
  // @return boolean
  tryToCompress (compositeSpan, spanToCompress) {
    if (!this._isEnabled()) {
      return false
    }

    // sets initial compression strategy value. returns
    // false if not compression eligable
    if (!this._initCompressionStrategy(compositeSpan, spanToCompress)) {
      return false
    }

    const isAlreadyComposite = this.isComposite()
    const canBeCompressed = isAlreadyComposite
      ? this.tryToCompressComposite(compositeSpan, spanToCompress) : this.tryToCompressRegular(compositeSpan, spanToCompress)

    if (!canBeCompressed) {
      return false
    }

    if (!isAlreadyComposite) {
      this.composite.count = 1
      this.composite.sum = compositeSpan._duration
    }

    this.composite.count++
    this.composite.sum += spanToCompress._duration
    this.duration = (spanToCompress._endTimestamp - compositeSpan.timestamp) / 1000
    return true
  }

  tryToCompressRegular (compositeSpan, spanToCompress) {
    if (!this.isSameKind(compositeSpan, spanToCompress)) {
      return false
    }

    if (compositeSpan.name === spanToCompress.name) {
      if (this.duration <= (this._agent._conf.spanCompressionExactMatchMaxDuration * 1000) && spanToCompress._duration <= (this._agent._conf.spanCompressionExactMatchMaxDuration * 1000)) {
        this.composite.compression_strategy = STRATEGY_EXACT_MATCH
        return true
      }
      return false
    }

    if (this.duration <= (this._agent._conf.spanCompressionSameKindMaxDuration * 1000) && spanToCompress._duration <= (this._agent._conf.spanCompressionSameKindMaxDuration * 1000)) {
      this.composite.compression_strategy = STRATEGY_SAME_KIND
      compositeSpan.name = this._spanNameFromCompositeSpan(compositeSpan)
      return true
    }

    return false
  }

  tryToCompressComposite (compositeSpan, spanToCompress) {
    switch (this.composite.compression_strategy) {
      case STRATEGY_EXACT_MATCH:
        return this.isSameKind(compositeSpan, spanToCompress) &&
            compositeSpan.name === spanToCompress.name &&
            spanToCompress._duration <= (this._agent._conf.spanCompressionExactMatchMaxDuration * 1000)
      case STRATEGY_SAME_KIND:
        return this.isSameKind(compositeSpan, spanToCompress) &&
          spanToCompress._duration <= (this._agent._conf.spanCompressionSameKindMaxDuration * 1000)
    }
  }

  _spanNameFromCompositeSpan (compositeSpan) {
    const prefix = 'Calls to '
    const serviceTarget = compositeSpan._serviceTarget
    if (!serviceTarget.type) {
      if (!serviceTarget.name) {
        return prefix + 'unknown'
      } else {
        return prefix + serviceTarget.name
      }
    } else if (!serviceTarget.name) {
      return prefix + serviceTarget.type
    } else {
      return prefix + serviceTarget.type + '/' + serviceTarget.name
    }
  }

  _isEnabled () {
    return this._agent._conf.spanCompressionEnabled
  }

  isCompositeSameKind () {
    return this.composite.compression_strategy === STRATEGY_SAME_KIND
  }

  isComposite () {
    return this.composite.count > 1
  }

  // Encodes/Serializes composite span properties
  // @return Object
  encode () {
    return {
      compression_strategy: this.composite.compression_strategy,
      count: this.composite.count,
      sum: this.composite.sum
    }
  }
}

module.exports = {
  SpanCompression,
  constants: {
    STRATEGY_EXACT_MATCH,
    STRATEGY_SAME_KIND
  }
}
