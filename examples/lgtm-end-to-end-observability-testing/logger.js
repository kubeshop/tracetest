// logger.js
const winston = require('winston');
const LokiTransport = require('winston-loki');

// Configure Winston to send logs to Loki
const logger = winston.createLogger({
  level: 'info',
  format: winston.format.json(),
  transports: [
    new LokiTransport({
      host: 'http://loki:3100', // Assuming Loki is accessible at this URL
      labels: { job: 'loki-service' },
      json: true,
      batching: true,
      interval: 5, // Send logs in batches every 5 seconds
    }),
  ],
});

module.exports = logger;

