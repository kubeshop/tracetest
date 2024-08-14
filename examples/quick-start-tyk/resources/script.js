const { expect } = require('@playwright/test');

const URL = 'http://tyk-gateway:8080';
const API_KEY = '28d220fd77974a4facfb07dc1e49c2aa';

const getKey = async () => {
  const params = {
    headers: {
      'Content-Type': 'application/json',
      'x-tyk-authorization': API_KEY,
      'Response-Type': 'application/json',
    },
  };

  const data = {
    alias: 'website',
    expires: -1,
    access_rights: {
      1: {
        api_id: '1',
        api_name: 'pokeshop',
        versions: ['Default'],
      },
    },
  };

  const res = await fetch(`${URL}/tyk/keys/create`, {
    ...params,
    method: 'POST',
    body: JSON.stringify(data),
  });

  const { key } = await res.json();

  return key;
};

async function importPokemon(page) {
  const key = await getKey();

  await page.setExtraHTTPHeaders({
    Authorization: `Bearer ${key}`,
  });

  await page.goto(URL);

  expect(await page.getByText('Pokeshop')).toBeTruthy();

  await page.click('text=Import');
  await page.getByLabel('ID').fill('143');

  await Promise.all([
    page.waitForResponse((resp) => resp.url().includes('/pokemon/import') && resp.status() === 200),
    page.getByRole('button', { name: 'OK', exact: true }).click(),
  ]);
}

module.exports = { importPokemon };
