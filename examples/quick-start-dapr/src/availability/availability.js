const opentelemetry = require('@opentelemetry/api')
const tracer = opentelemetry.trace.getTracer('bookstore')

const express = require('express')
const app = express()

const axios = require('axios')

app.get('/', (req, res) => {
  const span = tracer.startSpan('Hello Availability service!')
  span.end()
  res.send('Hello Availability service!')
})

app.get('/availability/:bookId', availabilityHandler)

async function availabilityHandler(req, res) {
  const span = tracer.startSpan('Availablity check')
  const bookId = req.params.bookId
  span.setAttribute('bookId', bookId)
  const isAvailable = await isBookAvailable(bookId)
  span.setAttribute('isAvailable', isAvailable)
  res.json({ isAvailable })
  span.end()
}

async function isBookAvailable(bookId) {
  const endpoint = `http://stock:8080/stock/${bookId}`
  const { data: { isAvailable } } = await axios.get(`${endpoint}`)
  return isAvailable
}

app.listen(8080, () => {
  console.log(`Listening for requests on http://localhost:8080`)
})
