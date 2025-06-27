import { test, expect } from "@playwright/test";
import Tracetest from "@tracetest/playwright";
let tracetest;

test.describe.configure({ mode: "serial" });

const { TRACETEST_TOKEN, TRACETEST_SERVER_URL = "https://app.tracetest.io", TRACETEST_ENVIRONMENT_ID } = process.env;

const definition = `
type: Test
spec:
  id: phAZcrT4A
  name: Playwright - Books list with availability
  description: Testing the books list and availability check
  trigger:
    type: playwright
  specs:
  - selector: span[tracetest.span.type="http" name="GET /books" http.target="/books"
      http.method="GET"]
    assertions:
    - attr:tracetest.span.duration  < 500ms
  - selector: span[tracetest.span.type="general" name="Books List"]
    assertions:
    - attr:books.list.count = 3
  - selector: span[tracetest.span.type="http" name="GET /availability/:bookId" http.method="GET"]
    assertions:
    - attr:http.host = "availability:8080"
  - selector: span[tracetest.span.type="general" name="Availablity check"]
    assertions:
    - attr:isAvailable = "true"

`;

test.beforeAll(async () => {
  tracetest = await Tracetest({
    apiToken: TRACETEST_TOKEN,
    serverUrl: TRACETEST_SERVER_URL,
    serverPath: "",
    environmentId: TRACETEST_ENVIRONMENT_ID,
  });

  await tracetest.setOptions({
    "should validate Bookstore": {
      definition,
    },
  });
});

test.beforeEach(async ({ page, context }, info) => {
  await tracetest?.capture({ context, info });
  await page.goto("/");
});

// optional step to break the playwright script in case a Tracetest test fails
test.afterAll(async ({}, testInfo) => {
  testInfo.setTimeout(80000);
  await tracetest?.summary();
});

test("should validate Bookstore", async ({ page }) => {
  await expect(await page.locator("h1")).toContainText("Bookstore");
  await expect(await page.getByRole("listitem")).toHaveCount(3);
  await expect(await page.getByRole("listitem").filter({ hasText: "❌" })).toHaveCount(1);
});
