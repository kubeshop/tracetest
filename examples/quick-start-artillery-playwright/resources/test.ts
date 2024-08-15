import { Page } from 'playwright';
import { expect, TestInfo } from '@playwright/test';
import Tracetest from '@tracetest/playwright';
import { config } from 'dotenv';

config();

const { TRACETEST_TOKEN = '', TRACETEST_ENVIRONMENT_ID = '' } = process.env;

const definition = `
type: Test
spec:
  id: artillery-playwight-import-pokemon
  name: "Artillery Playwright - Import Pokemon"
  trigger:
    type: playwright
  specs:
    - selector: span[tracetest.span.type="general" name = "validate request"] span[tracetest.span.type="http"]
      name: "All HTTP Spans: Status  code is 200"
      assertions:
        - attr:http.status_code = 200
    - selector: span[tracetest.span.type="http" name="GET" http.method="GET"]
      assertions:
        - attr:http.route = "/api/v2/pokemon/\${var:POKEMON_ID}"
    - selector: span[tracetest.span.type="database"]
      name: "All Database Spans: Processing time is less than 1s"
      assertions:
        - attr:tracetest.span.duration < 1s
  outputs:
    - name: DATABASE_POKEMON_ID
      selector: span[tracetest.span.type="database" name="create postgres.pokemon" db.system="postgres" db.name="postgres" db.user="postgres" db.operation="create" db.sql.table="pokemon"]
      value: attr:db.result | json_path '$.id'
`;

export async function importPokemon(page: Page) {
  const tracetest = await Tracetest({ apiToken: TRACETEST_TOKEN, environmentId: TRACETEST_ENVIRONMENT_ID });
  const title = 'Artillery Playwright - Import Pokemon';
  const pokemonId = Math.floor(Math.random() * 101).toString();
  await page.goto('/');

  await tracetest?.setOptions({
    [title]: {
      definition,
      runInfo: {
        variables: [
          {
            key: 'POKEMON_ID',
            value: pokemonId,
          },
        ],
      },
    },
  });

  await tracetest?.capture(page, { title: 'Artillery Playwright - Import Pokemon', config: {} } as TestInfo);

  expect(await page.getByText('Pokeshop')).toBeTruthy();

  await page.click('text=Import');

  await page.getByLabel('ID').fill(pokemonId);
  await page.getByRole('button', { name: 'OK', exact: true }).click();

  await tracetest?.summary();
}
