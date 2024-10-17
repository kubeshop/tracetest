const { expect } = require("@playwright/test");

async function pokemon(page) {
  await expect(await page.locator("h1")).toContainText("Pokemon");
  await expect(await page.getByRole("listitem")).toHaveCount(5);
}

module.exports = { pokemon };
