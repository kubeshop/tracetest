import { test, expect } from '@playwright/test';
import { getKey } from './auth';
import Tracetest, { Types } from '@tracetest/playwright';

const { TRACETEST_API_TOKEN = '' } = process.env;

let tracetest: Types.TracetestPlaywright | undefined = undefined;

test.beforeAll(async () => {
  // 1: Create a new Tracetest instance
  tracetest = await Tracetest({ apiToken: TRACETEST_API_TOKEN });
});

test.beforeEach(async ({ page, context }, { title }) => {
  const key = await getKey();
  await context.setExtraHTTPHeaders({
    Authorization: `Bearer ${key}`,
  });

  await page.goto('/');
  // 2: Capture the initial page
  await tracetest?.capture(title, page);
});

test.afterAll(async ({}, testInfo) => {
  // 3: Summary of the test (optional, but recommended)
  await tracetest?.summary();
});

test('Playwright: imports a pokemon', async ({ page }) => {
  expect(await page.getByText('Pokeshop')).toBeTruthy();

  await page.click('text=Import');

  await page.getByLabel('ID').fill(Math.floor(Math.random() * 101).toString());
  await page.getByRole('button', { name: 'OK', exact: true }).click();
});
