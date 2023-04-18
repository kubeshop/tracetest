const {defineConfig} = require('cypress');

module.exports = defineConfig({
  viewportWidth: 1440,
  viewportHeight: 1080,
  responseTimeout: 30000,
  pageLoadTimeout: 20000,
  projectId: '6dm1if',
  requestTimeout: 30000,
  retries: 2,
  e2e: {
    // We've imported your old cypress plugins here.
    // You may want to clean this up later by importing these.
    setupNodeEvents(on, config) {
      // eslint-disable-next-line global-require
      return require('./cypress/plugins/index.js')(on, config);
    },
    specPattern: 'cypress/e2e/**/*.spec.{js,ts,jsx,tsx}',
    baseUrl: process.env.CYPRESS_BASE_URL ||'http://localhost:3000',
  },
  // @ts-ignore
  component: {
    setupNodeEvents() {},
    specPattern: 'cypress/components/**/*.spec.{js,ts,jsx,tsx}',
  },
});
