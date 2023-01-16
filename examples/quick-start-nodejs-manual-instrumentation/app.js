const opentelemetry = require('@opentelemetry/api')
const tracer = opentelemetry.trace.getTracer('quick-start-nodejs-manual-instrumentation-tracer')

const express = require('express')
const app = express()

app.get('/', (req, res) => {
  const span = tracer.startSpan('Hello World')
  span.end()
  res.send('Hello World')
})

app.get('/books', booksListHandler)

function booksListHandler(req, res) {
  const span = tracer.startSpan('Books List')
  const books = getBooks()
  span.setAttribute('books.list.count', books.length)
  span.end()

  res.json(books)
}

function getBooks() {
  return [
    {
      id: 1,
      title: 'Harry Potter',
      isAvailable: 0,
    },
    {
      id: 2,
      title: 'Foundation',
      isAvailable: 0,
    },
    {
      id: 3,
      title: 'Moby Dick',
      isAvailable: 0,
    },
  ]
}

app.listen(8080, () => {
  console.log(`Listening for requests on http://localhost:8080`)
})
