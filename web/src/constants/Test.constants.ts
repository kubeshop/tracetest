import {HTTP_METHOD} from './Common.constants';

export interface IDemoTestExample {
  name: string;
  url: string;
  method: HTTP_METHOD;
  body: string;
  description: string;
}

export const DemoTestExampleList: IDemoTestExample[] = [
  {
    name: 'Pokemon - List',
    url: 'http://demo-pokemon-api.demo.svc.cluster.local/pokemon?take=20&skip=0',
    method: HTTP_METHOD.GET,
    body: '',
    description: 'Get a Pokemon',
  },
  {
    name: 'Pokemon - Add',
    url: 'http://demo-pokemon-api.demo.svc.cluster.local/pokemon',
    method: HTTP_METHOD.POST,
    body: '{"name":"meowth","type":"normal","imageUrl":"https://assets.pokemon.com/assets/cms2/img/pokedex/full/052.png","isFeatured":true}',
    description: 'Add a Pokemon',
  },
  {
    name: 'Pokemon - Import',
    url: 'http://demo-pokemon-api.demo.svc.cluster.local/pokemon/import',
    method: HTTP_METHOD.POST,
    body: '{"id":52}',
    description: 'Import a Pokemon',
  },
];

export const DEFAULT_HEADERS = [{key: 'Content-Type', value: 'application/json'}];

export enum ResultViewModes {
  Advanced = 'advanced',
  Wizard = 'wizard',
}

export enum TriggerTypes {
  http = 'http',
  grpc = 'grpc',
}
