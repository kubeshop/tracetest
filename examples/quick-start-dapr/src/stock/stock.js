const opentelemetry = require('@opentelemetry/api')
const tracer = opentelemetry.trace.getTracer('bookstore')

const express = require('express')
const app = express()

app.get('/', (req, res) => {
  const span = tracer.startSpan('Hello Stock service!')
  span.end()
  res.send('Hello Stock service!')
})

app.get('/stock/:bookId', stockHandler)

function stockHandler(req, res) {
  const span = tracer.startSpan('Stock check')
  const bookId = req.params.bookId
  span.setAttribute('bookId', bookId)
  const isAvailable = isInStock(bookId)
  span.setAttribute('isAvailable', isAvailable)
  res.json({ isAvailable })
  span.end()
}

function isInStock(bookId) {
  const { stock } = getStock().find(book => book.id == bookId)
  return stock > 0
}

function getStock() {
  return [
    { id: 1, stock: 6 },
    { id: 2, stock: 8 },
    { id: 3, stock: 0 }
  ]
}

app.listen(8080, () => {
  console.log(`Listening for requests on http://localhost:8080`)
})
