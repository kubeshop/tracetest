'use strict'

var whitespace = /[\s;]+/g
var borderChars = /([;()])/g

module.exports = function (sql) {
  if (!sql) return ''
  var tokens = tokenize(sql)
  return stipTokens(tokens).join(' ')
}

function stipTokens (tokens) {
  var verb = tokens[0].toUpperCase()
  return [verb]
    .concat(afterVerb(tokens))
    .filter(function (token) { return !!token })
}

function afterVerb (tokens) {
  switch (tokens[0].toUpperCase()) {
    case 'SELECT': return afterToken('FROM', tokens)
    case 'INSERT': return afterToken('INTO', tokens)
    case 'UPDATE': return tokens[1]
    case 'DELETE': return afterToken('FROM', tokens)
    case 'CREATE': return afterToken(['DATABASE', 'TABLE', 'INDEX'], tokens)
    case 'DROP': return afterToken(['DATABASE', 'TABLE'], tokens)
    case 'ALTER': return afterToken('TABLE', tokens)
    case 'DESC': return tokens[1]
    case 'TRUNCATE': return afterToken('TABLE', tokens)
    case 'USE': return tokens[1]
  }
}

function afterToken (find, tokens) {
  var index

  if (!Array.isArray(find)) find = [find]
  find = find.map(function (find) { return find.toUpperCase() })

  for (var n = 0, l = tokens.length - 1; n < l; n++) {
    index = find.indexOf(tokens[n].toUpperCase())
    if (index !== -1) return [find[index], tokens[n + 1]]
  }
}

function tokenize (sql) {
  return normalize(sql).split(' ')
}

function normalize (sql) {
  return sql.replace(borderChars, ' $1 ').replace(whitespace, ' ').trim()
}
