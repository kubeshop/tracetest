import { test, expect } from "@playwright/test"

test.beforeEach(async ({ page }, info) => {
  await page.goto("/")
})

test("should have a h1 heading of Bookstore", async ({ page }) => {
  await expect(await page.locator("h1")).toContainText("Bookstore")  
})

test("should have a list with 3 items", async ({ page }) => {
  await expect(await page.getByRole('listitem')).toHaveCount(3)  
})

test("should have one list item with a red X", async ({ page }) => {
  await expect(await page.getByRole('listitem').filter({ hasText: '❌' })).toHaveCount(1)
})
