const express = require('express');
const logger = require('./logger');
const meter = require('./meter');
// const tracer = require('./tracer');

// Create an Express app
const app = express();

// Define a custom metric (e.g., a request counter)
const requestCounter = meter.createCounter('http_requests', {
  description: 'Counts HTTP requests',
});

// Middleware to increment the counter on every request
app.use((req, res, next) => {
  // Increment the request counter
  logger.info(`Received request for ${req.url}`);
  requestCounter.add(1, { method: req.method, route: req.path });
  next();
});

// Define a simple route
app.get('/', (req, res) => {
  // const span = tracer.startSpan('handle_root_request');

  // Simulate some work
  setTimeout(() => {
    res.send('Hello, World!');
    // span.end();
  }, 100);
});

// Start the server
app.listen(8081, () => {
  logger.info('Server is running on http://localhost:5000');
  console.log('Server is running on http://localhost:5000');
});