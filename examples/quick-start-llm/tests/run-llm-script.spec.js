const { expect } = require("@playwright/test");

async function runGeminiSummarization(page) {
  await page.getByRole("button", { name: "Add example text", exact: true }).click();
  await page.getByRole("button", { name: "Summarize", exact: true }).click();
}

module.exports = { runGeminiSummarization };
