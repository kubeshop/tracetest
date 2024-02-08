const { defineConfig } = require('cypress');
const { config } = require('dotenv');

config();

module.exports = defineConfig({
  chromeWebSecurity: false,
  e2e: {
    baseUrl: process.env.POKESHOP_DEMO_URL || 'http://localhost:8081',
    env: {
      TRACETEST_API_TOKEN: process.env.TRACETEST_API_TOKEN,
    },
    setupNodeEvents() {
      // implement node event listeners here
    },
  },
});
