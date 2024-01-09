export function CypressCodeSnippet(testName: string) {
  return `import Tracetest from '@tracetest/cypress';

const TRACETEST_API_TOKEN = Cypress.env('TRACETEST_API_TOKEN') || '';
const tracetest = Tracetest();

describe('Cypress Test', () => {
  before(done => {
    tracetest.configure(TRACETEST_API_TOKEN).then(() => done());
  });

  beforeEach(() => {
    cy.visit('/', {
      onBeforeLoad: win => tracetest.capture(win.document),
    });
  });

  afterEach(done => {
    tracetest.runTest('').then(() => done());
  });

  it('${testName}', () => {
    // ...cy commands
  });
});`;
}

export function PlaywrightCodeSnippet(testName: string) {
  return `import { test, expect } from '@playwright/test';
import Tracetest from '@tracetest/playwright';

const { TRACETEST_API_TOKEN = '' } = process.env;

const tracetest = Tracetest();

test.describe.configure({ mode: 'serial' });

test.beforeAll(async () => {
  await tracetest.configure(TRACETEST_API_TOKEN);
});

test.beforeEach(async ({ page }, { title }) => {
  await page.goto('/');
  await tracetest.capture(title, page);
});

test.afterEach(async ({}, { title, config }) => {
  await tracetest.runTest(title, config.metadata.definition ?? '');
});

// optional step to break the playwright script in case a Tracetest test fails
test.afterAll(async ({}, testInfo) => {
  testInfo.setTimeout(60000);
  await tracetest.summary();
});

test('${testName}', () => {
  // ...playwright commands
});`;
}
