/**
 * Moved from JSON files to TS files, as JSONs are not working in composite projects.
 * @see {@link https://github.com/TypeStrong/ts-loader/issues/905}
 */
export default {
  info: {
    _postman_id: '910a51e7-5e97-4e2e-ba77-5f6d8c7625dd',
    name: 'Poke Micro API - DEV',
    schema: 'https://schema.getpostman.com/json/collection/v2.1.0/collection.json',
    _exporter_id: '326743',
  },
  item: [
    {
      name: 'List',
      request: {
        method: 'GET',
        header: [],
        url: {
          raw: 'http://demo-pokemon-api.demo.svc.cluster.local/pokemon?take=20&skip=0',
          protocol: 'http',
          host: ['demo-pokemon-api', 'demo', 'svc', 'cluster', 'local'],
          path: ['pokemon'],
          query: [
            {
              key: 'take',
              value: '20',
            },
            {
              key: 'skip',
              value: '0',
            },
          ],
        },
      },
      response: [],
    },
    {
      name: 'Create',
      request: {
        method: 'POST',
        header: [],
        body: {
          mode: 'raw',
          raw: '{\n    "name": "meowth",\n    "type": "normal",\n    "imageUrl": "https://assets.pokemon.com/assets/cms2/img/pokedex/full/052.png",\n    "isFeatured": true\n}',
          options: {
            raw: {
              language: 'json',
            },
          },
        },
        url: {
          raw: 'http://demo-pokemon-api.demo.svc.cluster.local/pokemon',
          protocol: 'http',
          host: ['demo-pokemon-api', 'demo', 'svc', 'cluster', 'local'],
          path: ['pokemon'],
        },
      },
      response: [],
    },
    {
      name: 'Import',
      request: {
        method: 'POST',
        header: [],
        body: {
          mode: 'raw',
          raw: '{\n    "id": 6\n}',
          options: {
            raw: {
              language: 'json',
            },
          },
        },
        url: {
          raw: 'http://demo-pokemon-api.demo.svc.cluster.local/pokemon/import',
          protocol: 'http',
          host: ['demo-pokemon-api', 'demo', 'svc', 'cluster', 'local'],
          path: ['pokemon', 'import'],
        },
      },
      response: [],
    },
  ],
};
