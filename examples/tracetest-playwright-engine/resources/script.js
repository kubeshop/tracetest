const { expect } = require("@playwright/test");

async function addPokemon(page) {
  expect(await page.getByText("Pokeshop")).toBeTruthy();

  await page.click("text=Add");

  await page.getByLabel("Name").fill("Charizard");
  await page.getByLabel("Type").fill("Flying");
  await page
    .getByLabel("Image URL")
    .fill("https://upload.wikimedia.org/wikipedia/en/1/1f/Pok%C3%A9mon_Charizard_art.png");
  await page.getByRole("button", { name: "OK", exact: true }).click();
}

async function deletePokemon(page) {
  expect(await page.getByText("Pokeshop")).toBeTruthy();

  await page.locator('[data-cy="pokemon-list"]');
  await page.locator('[data-cy="pokemon-card"]').first().click();
  await page.locator('[data-cy="pokemon-card"] [data-cy="delete-pokemon-button"]').first().click();
}

async function importPokemon(page) {
  expect(await page.getByText("Pokeshop")).toBeTruthy();

  await page.click("text=Import");
  await page.getByLabel("ID").fill("143");

  await Promise.all([
    page.waitForResponse((resp) => resp.url().includes("/pokemon/import") && resp.status() === 200),
    page.getByRole("button", { name: "OK", exact: true }).click(),
  ]);
}

module.exports = { addPokemon, deletePokemon, importPokemon };
