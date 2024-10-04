const { expect } = require("@playwright/test");

async function bookstore(page) {
  expect(await page.locator("h1")).toContainText("Bookstore")  
  expect(await page.getByRole('listitem')).toHaveCount(3)  
  expect(await page.getByRole('listitem').filter({ hasText: '❌' })).toHaveCount(1)
}

module.exports = { bookstore };
