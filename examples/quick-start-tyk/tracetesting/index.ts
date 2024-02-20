import Tracetest from '@tracetest/client';
import fetch from 'node-fetch';
import { TestResource } from '@tracetest/client/dist/modules/openapi-client';
import { config } from 'dotenv';

config();

const { TRACETEST_API_TOKEN = '', POKESHOP_DEMO_URL = '' } = process.env;

const params = {
  headers: {
    'Content-Type': 'application/json',
    'x-tyk-authorization': '28d220fd77974a4facfb07dc1e49c2aa',
    'Response-Type': 'application/json',
  },
};

const setup = async () => {
  const alias = 'website';

  const data = {
    alias,
    expires: -1,
    access_rights: {
      1: {
        api_id: '1',
        api_name: 'pokeshop',
        versions: ['Default'],
      },
    },
  };

  const res = await fetch('http://tyk-gateway:8080/tyk/keys/create', {
    ...params,
    method: 'POST',
    body: JSON.stringify(data),
  });

  const { key } = (await res.json()) as { key: string };

  return key;
};

const definition: TestResource = {
  type: 'Test',
  spec: {
    id: 'ZV1G3v2IR',
    name: 'Import Pokemon',
    trigger: {
      type: 'http',
      httpRequest: {
        method: 'POST',
        url: '${var:ENDPOINT}/pokeshop/import',
        body: '{"id": ${var:POKEMON_ID}}',
        headers: [
          {
            key: 'Content-Type',
            value: 'application/json',
          },
          {
            key: 'Authorization',
            value: '${var:AUTHORIZATION}',
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
        selector:
          'span[name="tracetest-serverless-dev-api"] span[tracetest.span.type="http" name="GET" http.method="GET"]',
        name: 'The request matches the pokemon Id',
        assertions: ['attr:http.url  =  "https://pokeapi.co/api/v2/pokemon/${var:POKEMON_ID}"'],
      },
    ],
  },
};

const main = async () => {
  const tracetest = await Tracetest(TRACETEST_API_TOKEN);

  const key = await setup();

  const test = await tracetest.newTest(definition);
  await tracetest.runTest(test, {
    variables: [
      {
        key: 'ENDPOINT',
        value: POKESHOP_DEMO_URL.trim(),
      },
      {
        key: 'POKEMON_ID',
        value: `${Math.floor(Math.random() * 100) + 1}`,
      },
      {
        key: 'AUTHORIZATION',
        value: key,
      },
    ],
  });
  console.log(await tracetest.getSummary());
};

main();
