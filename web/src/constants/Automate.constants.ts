export function CypressCodeSnippet(testName: string) {
  return `import Tracetest, { Types } from '@tracetest/cypress';
const TRACETEST_API_TOKEN = Cypress.env('TRACETEST_API_TOKEN') || '';
let tracetest: Types.TracetestCypress | undefined = undefined;

describe('Home', { defaultCommandTimeout: 60000 }, () => {
  before(done => {
    Tracetest({ apiToken: TRACETEST_API_TOKEN }).then(() => done());
  });

  beforeEach(() => {
    cy.visit('/', {
      onBeforeLoad: win => tracetest.capture(win.document),
    });
  });

  // uncomment to wait for trace tests to be done
  after(done => {
    tracetest.summary().then(() => done());
  });

  it('${testName}', () => {
    // ...cy commands
  });
});`;
}

export function PlaywrightCodeSnippet(testName: string) {
  return `import { test, expect } from '@playwright/test';
import Tracetest, { Types } from '@tracetest/playwright';
const { TRACETEST_API_TOKEN = '' } = process.env;
let tracetest: Types.TracetestPlaywright | undefined = undefined;

test.describe.configure({ mode: 'serial' });
test.beforeAll(async () => {
  tracetest = await Tracetest({ apiToken: TRACETEST_API_TOKEN });
});

test.beforeEach(async ({ page }, { title }) => {
  await page.goto('/');
  await tracetest?.capture(title, page);
});

// optional step to break the playwright script in case a Tracetest test fails
test.afterAll(async ({}, testInfo) => {
  testInfo.setTimeout(60000);
  await tracetest?.summary();
});

test('${testName}', () => {
  // ...playwright commands
});`;
}
