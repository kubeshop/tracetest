const opentelemetry = require('@opentelemetry/api')
const tracer = opentelemetry.trace.getTracer('quick-start-nodejs-manual-instrumentation-tracer')

const express = require('express')
const app = express()

const axios = require('axios')

app.get('/', (req, res) => {
  const span = tracer.startSpan('Hello World')
  span.end()
  res.send('Hello World')
})

app.get('/books', booksListHandler)
app.get('/books/:bookId', async function(req, res){
  const avRes = await axios.get(`http://localhost:8081/availability/${req.params.bookId}`)
  const isAvailable = avRes.data
  res.json(isAvailable)
})

async function booksListHandler(req, res) {
  const span = tracer.startSpan('Books List')
  const books = await getAvailableBooks()
  span.setAttribute('books.list.count', books.length)
  span.end()

  res.json(books)
}

async function getAvailableBooks() {
  const books = getBooks()
  const availableBooks = await Promise.all(
    books.map(async book => {
      const endpoint = `http://availability:8080/availability/${book.id}`
      const { data: { isAvailable } } = await axios.get(`${endpoint}`)
      return { ...book, isAvailable }
    })
  )
  return availableBooks
}

function getBooks() {
  return [
    {
      id: 1,
      title: 'Harry Potter',
    },
    {
      id: 2,
      title: 'Foundation',
    },
    {
      id: 3,
      title: 'Moby Dick',
    },
  ]
}

app.listen(8080, () => {
  console.log(`Listening for requests on http://localhost:8080`)
})
