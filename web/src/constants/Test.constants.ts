import { Schemas } from "../types/Common.types";

interface DemoTestExample {
  name: string;
  url: string;
  method: Schemas['HTTPRequest']['method'];
  body: string;
  description: string;
}

export const DemoTestExampleList: DemoTestExample[] = [
  {
    name: 'Shopping app',
    url: 'http://shop/buy',
    method: 'GET',
    body: '',
    description: 'Generic get',
  },
  {
    name: 'Pokemon - List',
    url: 'http://demo-pokemon-api.demo.svc.cluster.local/pokemon?take=20&skip=0',
    method: 'GET',
    body: '',
    description: 'Get a Pokemon',
  },
  {
    name: 'Pokemon - Add',
    url: 'http://demo-pokemon-api.demo.svc.cluster.local/pokemon',
    method: 'POST',
    body: '{"name":"meowth","type":"normal","imageUrl":"https://assets.pokemon.com/assets/cms2/img/pokedex/full/052.png","isFeatured":true}',
    description: 'Add a Pokemon',
  },
  {
    name: 'Pokemon - Import',
    url: 'http://demo-pokemon-api.demo.svc.cluster.local/pokemon/import',
    method: 'POST',
    body: '{"id":52}',
    description: 'Import a Pokemon',
  },
];
