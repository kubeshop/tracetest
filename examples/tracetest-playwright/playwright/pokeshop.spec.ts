import { test, expect } from '@playwright/test';
import Tracetest, { Types } from '@tracetest/playwright';

const { TRACETEST_API_TOKEN = '' } = process.env;

let tracetest: Types.TracetestPlaywright | undefined = undefined;

test.describe.configure({ mode: 'serial' });

const definition = `
type: Test
spec:
  id: UGxheXdyaWdodDogaW1wb3J0cyBhIHBva2Vtb24=
  name: "Playwright: imports a pokemon"
  trigger:
    type: playwright
  specs:
    - selector: span[tracetest.span.type="http"] span[tracetest.span.type="http"]
      name: "All HTTP Spans: Status  code is 200"
      assertions:
      - attr:http.status_code   =   200
    - selector: span[tracetest.span.type="database"]
      name: "All Database Spans: Processing time is less than 100ms"
      assertions:
      - attr:tracetest.span.duration < 2s
  outputs:
    - name: MY_OUTPUT
      selector: span[tracetest.span.type="general" name="Tracetest trigger"]
      value: attr:name
`;

test.beforeAll(async () => {
  tracetest = await Tracetest({ apiToken: TRACETEST_API_TOKEN });

  await tracetest.setOptions({
    'Playwright: imports a pokemon': {
      definition,
    },
  });
});

test.beforeEach(async ({ page }, info) => {
  await page.goto('/');
  await tracetest?.capture(page, info);
});

// optional step to break the playwright script in case a Tracetest test fails
test.afterAll(async ({}, testInfo) => {
  testInfo.setTimeout(80000);
  await tracetest?.summary();
});

test('Playwright: creates a pokemon', async ({ page }) => {
  expect(await page.getByText('Pokeshop')).toBeTruthy();

  await page.click('text=Add');

  await page.getByLabel('Name').fill('Charizard');
  await page.getByLabel('Type').fill('Flying');
  await page
    .getByLabel('Image URL')
    .fill('https://upload.wikimedia.org/wikipedia/en/1/1f/Pok%C3%A9mon_Charizard_art.png');
  await page.getByRole('button', { name: 'OK', exact: true }).click();
});

test('Playwright: imports a pokemon', async ({ page }) => {
  expect(await page.getByText('Pokeshop')).toBeTruthy();

  await page.click('text=Import');

  await page.getByLabel('ID').fill(Math.floor(Math.random() * 101).toString());
  await page.getByRole('button', { name: 'OK', exact: true }).click();
});

test('Playwright: deletes a pokemon', async ({ page }) => {
  await page.locator('[data-cy="pokemon-list"]');

  await page.locator('[data-cy="pokemon-card"]').first().click();
  await page.locator('[data-cy="pokemon-card"] [data-cy="delete-pokemon-button"]').first().click();
});
