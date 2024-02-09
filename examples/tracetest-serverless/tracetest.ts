import Tracetest from '@tracetest/client';
import { TestResource } from '@tracetest/client/dist/modules/openapi-client';
import { config } from 'dotenv';

config();

const { TRACETEST_API_TOKEN = '', ENDPOINT = '' } = process.env;

const definition: TestResource = {
  type: 'Test',
  spec: {
    id: 'ZV1G3v2IR',
    name: 'Serverless: Import Pokemon',
    trigger: {
      type: 'http',
      httpRequest: {
        method: 'POST',
        url: '${var:ENDPOINT}/import',
        body: '{"id": "${var:POKEMON_ID}"}\n',
        headers: [
          {
            key: 'Content-Type',
            value: 'application/json',
          },
        ],
      },
    },
    specs: [
      {
        selector: 'span[tracetest.span.type="database"]',
        name: 'All Database Spans: Processing time is less than 100ms',
        assertions: ['attr:tracetest.span.duration < 100ms'],
      },
      {
        selector: 'span[tracetest.span.type="http"]',
        name: 'All HTTP Spans: Status  code is 200',
        assertions: ['attr:http.status_code = 200'],
      },
      {
        selector: 'span[name="tracetest-serverless-dev-api"] span[tracetest.span.type="http" name="GET" http.method="GET"]',
        name: 'The request matches the pokemon Id',
        assertions: ['attr:http.url  =  "https://pokeapi.co/api/v2/pokemon/${var:POKEMON_ID}"'],
      },
    ],
  },
};

const main = async () => {
  const tracetest = await Tracetest(TRACETEST_API_TOKEN);

  const test = await tracetest.newTest(definition);
  await tracetest.runTest(test, {
    variables: [
      {
        key: 'ENDPOINT',
        value: ENDPOINT.trim(),
      },
      {
        key: 'POKEMON_ID',
        value: `${Math.floor(Math.random() * 100) + 1}`,
      },
    ],
  });
  console.log(await tracetest.getSummary());
};

main();
