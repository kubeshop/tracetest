'use strict'

const decode = Symbol('decode')
const encode = Symbol('encode')

function parsePair (segment) {
  return segment.trim().split('=')
}

function isNotEmpty (value) {
  return value !== ''
}

class SetCookie {
  constructor (input, options) {
    if (Array.isArray(input)) {
      return input.map(item => new SetCookie(item, options))
    }

    this.data = {}
    this.meta = {
      expires: undefined,
      maxAge: undefined,
      domain: undefined,
      path: undefined,
      secure: undefined,
      httpOnly: undefined,
      sameSite: undefined,
    }

    options = options || {}

    // Options
    this[decode] = options.decode || decodeURIComponent
    this[encode] = options.encode || encodeURIComponent

    // Convert strings to objects
    if (typeof input === 'string') {
      const segments = input.split(';')
      const pair = segments.shift()

      // NOTE: `foo=bar`, `=bar` and `bar` are all valid forms.
      // This way of parsing and getting key/value supports it.
      const position = pair.indexOf('=')
      const key = position >= 0 ? pair.slice(0, position) : ''
      const value = pair.slice(position + 1)

      if (!value) {
        throw new Error('Invalid value')
      }

      this.data = {}
      this.data[this[decode](key)] = this[decode](value)

      for (let pair of segments.map(parsePair)) {
        switch (pair[0].toLowerCase()) {
          case 'expires':
            const expires = new Date(pair[1])
            if (isNaN(expires.getTime())) {
              throw new Error('Invalid Expires field')
            }
            this.meta.expires = expires
            break

          case 'max-age':
            const maxAge = parseInt(pair[1], 10)
            if (isNaN(maxAge)) {
              throw new Error('Invalid Max-Age field')
            }
            this.meta.maxAge = maxAge
            break

          case 'domain':
            if (!pair[1]) {
              throw new Error('Invalid Domain field')
            }
            this.meta.domain = pair[1]
            break

          case 'path':
            if (!pair[1]) {
              throw new Error('Invalid Path field')
            }
            this.meta.path = pair[1]
            break

          case 'secure':
            if (pair[1]) {
              throw new Error('Invalid Secure field')
            }
            this.meta.secure = true
            break

          case 'httponly':
            if (pair[1]) {
              throw new Error('Invalid HttpOnly field')
            }
            this.meta.httpOnly = true
            break

          case 'samesite':
            if (!pair[1]) {
              throw new Error('Invalid SameSite field')
            }
            this.meta.sameSite = pair[1]
            break
        }
      }

    // Passthrough objects as-is
    } else if (typeof input === 'object') {
      const data = input.data
      const meta = input.meta

      if (!data || !Object.keys(data).length) {
        throw new Error('Missing data')
      }

      Object.assign(this.data, data)
      Object.assign(this.meta, meta)
    } else {
      throw new Error('Invalid input type')
    }
  }

  toString () {
    const pairs = []
    
    for (let key of Object.keys(this.data)) {
      const pair = [
        this[encode](this.data[key])
      ]
      if (key) {
        pair.unshift('=')
        pair.unshift(this[encode](key))
      }
      pairs.push(pair.join(''))
    }

    if (typeof this.meta.expires !== 'undefined') {
      pairs.push(`Expires=${this.meta.expires.toUTCString()}`)
    }

    if (typeof this.meta.maxAge !== 'undefined') {
      pairs.push(`Max-Age=${this.meta.maxAge}`)
    }

    if (typeof this.meta.domain !== 'undefined') {
      pairs.push(`Domain=${this.meta.domain}`)
    }

    if (typeof this.meta.path !== 'undefined') {
      pairs.push(`Path=${this.meta.path}`)
    }

    if (typeof this.meta.secure !== 'undefined') {
      pairs.push('Secure')
    }

    if (typeof this.meta.httpOnly !== 'undefined') {
      pairs.push('HttpOnly')
    }

    if (typeof this.meta.sameSite !== 'undefined') {
      pairs.push(`SameSite=${this.meta.sameSite}`)
    }

    return pairs.join('; ')
  }
}

module.exports = SetCookie 
