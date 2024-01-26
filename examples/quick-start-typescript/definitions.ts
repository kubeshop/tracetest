import { TestResource } from '@tracetest/client/dist/modules/openapi-client';

export const importDefinition: TestResource = {
  type: 'Test',
  spec: {
    id: '99TOHzpSR',
    name: 'Typescript: Import a Pokemon',
    trigger: {
      type: 'http',
      httpRequest: {
        method: 'POST',
        url: '${var:BASE_URL}/import',
        body: '{"id": ${var:POKEMON_ID}}',
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
        selector: 'span[tracetest.span.type="general" name = "validate request"] span[tracetest.span.type="http"]',
        name: 'All HTTP Spans: Status  code is 200',
        assertions: ['attr:http.status_code = 200'],
      },
      {
        selector: 'span[tracetest.span.type="http" name="GET" http.method="GET"]',
        assertions: ['attr:http.route = "/api/v2/pokemon/${var:POKEMON_ID}"'],
      },
      {
        selector: 'span[tracetest.span.type="database"]',
        name: 'All Database Spans: Processing time is less than 1s',
        assertions: ['attr:tracetest.span.duration < 1s'],
      },
    ],
    outputs: [
      {
        name: 'DATABASE_POKEMON_ID',
        selector:
          'span[tracetest.span.type="database" name="create postgres.pokemon" db.system="postgres" db.name="postgres" db.user="postgres" db.operation="create" db.sql.table="pokemon"]',
        value: "attr:db.result | json_path '$.id'",
      },
    ],
  },
};

export const deleteDefinition: TestResource = {
  type: 'Test',
  spec: {
    id: 'C2gwdktIR',
    name: 'Typescript: Delete a Pokemon',
    trigger: {
      type: 'http',
      httpRequest: {
        method: 'DELETE',
        url: '${var:BASE_URL}/${var:POKEMON_ID}',
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
        selector:
          'span[tracetest.span.type="database" db.system="redis" db.operation="del" db.redis.database_index="0"]',
        assertions: ['attr:db.payload = \'{"key":"pokemon-${var:POKEMON_ID}"}\''],
      },
      {
        selector:
          'span[tracetest.span.type="database" name="delete postgres.pokemon" db.system="postgres" db.name="postgres" db.user="postgres" db.operation="delete" db.sql.table="pokemon"]',
        assertions: ['attr:db.result = 1'],
      },
      {
        selector: 'span[tracetest.span.type="database"]',
        name: 'All Database Spans: Processing time is less than 100ms',
        assertions: ['attr:tracetest.span.duration < 100ms'],
      },
    ],
  },
};
