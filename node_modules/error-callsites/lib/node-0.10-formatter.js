'use strict'

// Lifted from Node.js 0.10.40:
// https://github.com/nodejs/node/blob/0439a28d519fb6efe228074b0588a59452fc1677/deps/v8/src/messages.js#L1053-L1080
module.exports = function FormatStackTrace (error, frames) {
  var lines = []
  try {
    lines.push(error.toString())
  } catch (e) {
    try {
      lines.push('<error: ' + e + '>')
    } catch (ee) {
      lines.push('<error>')
    }
  }
  for (var i = 0; i < frames.length; i++) {
    var frame = frames[i]
    var line
    try {
      line = frame.toString()
    } catch (e) {
      try {
        line = '<error: ' + e + '>'
      } catch (ee) {
        // Any code that reaches this point is seriously nasty!
        line = '<error>'
      }
    }
    lines.push('    at ' + line)
  }
  return lines.join('\n')
}
