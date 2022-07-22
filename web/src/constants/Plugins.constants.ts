import pokeshopProtoData from 'assets/pokeshop.proto.json';
import pokeshopPostmanData from 'assets/pokeshop.postman_collection.json';
import {IPlugin} from 'types/Plugins.types';
import {HTTP_METHOD} from './Common.constants';
import {TriggerTypes} from './Test.constants';

const pokeshopProtoFile = new File([pokeshopProtoData.proto], 'pokeshop.proto');
const pokeshopPostmanFile = new File([JSON.stringify(pokeshopPostmanData)], 'pokeshop.postman_collection.json');

export enum SupportedPlugins {
  REST = 'REST',
  Messaging = 'Messaging',
  RPC = 'RPC',
  Postman = 'Postman',
  OpenAPI = 'OpenAPI',
}

const Default: IPlugin = {
  name: SupportedPlugins.REST,
  title: 'Default',
  description: '',
  isActive: false,
  type: TriggerTypes.http,
  demoList: [],
  stepList: [
    {
      id: 'plugin-selection',
      name: 'Select a plugin',
      title: 'Choose the way of creating a test',
      component: 'SelectPlugin',
      isDefaultValid: true,
      status: 'selected',
    },
    {
      id: 'basic-details',
      name: 'Basic Details',
      title: 'Provide needed basic information',
      component: 'BasicDetails',
    },
  ],
};

const Rest: IPlugin = {
  name: SupportedPlugins.REST,
  title: 'HTTP Request',
  description: 'Create a basic HTTP request',
  isActive: true,
  type: TriggerTypes.http,
  demoList: [
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
  ],
  stepList: [
    ...Default.stepList,
    {
      id: 'request-details',
      name: 'Request Details',
      title: 'Provide additional information',
      component: 'RequestDetails',
    },
  ],
};

const RPC: IPlugin = {
  name: SupportedPlugins.RPC,
  title: 'RPC Request',
  description: 'Test and debug your RPC request',
  isActive: true,
  type: TriggerTypes.grpc,
  demoList: [
    {
      name: 'GRPC - Pokemon - List',
      url: 'demo-pokemon-api.demo.svc.cluster.local:8082',
      message: '',
      method: 'pokeshop.Pokeshop.getPokemonList',
      description: 'Get a Pokemon',
      protoFile: pokeshopProtoFile,
    },
    {
      name: 'GRPC - Pokemon - Add',
      url: 'demo-pokemon-api.demo.svc.cluster.local:8082',
      message:
        '{"name":"meowth","type":"normal","imageUrl":"https://assets.pokemon.com/assets/cms2/img/pokedex/full/052.png","isFeatured":true}',
      method: 'pokeshop.Pokeshop.createPokemon',
      protoFile: pokeshopProtoFile,
      description: 'Import a Pokemon',
    },
    {
      name: 'GRPC - Pokemon - Import',
      url: 'demo-pokemon-api.demo.svc.cluster.local:8082',
      message: '{"id":52}',
      method: 'pokeshop.Pokeshop.importPokemon',
      protoFile: pokeshopProtoFile,
      description: 'Import a Pokemon',
    },
  ],
  stepList: [
    ...Default.stepList,
    {
      id: 'request-details',
      name: 'Request Details',
      title: 'Provide additional information',
      component: 'RequestDetails',
    },
  ],
};

const Messaging: IPlugin = {
  name: SupportedPlugins.Messaging,
  title: 'Message Queue',
  description: 'Put a message on a queue to initiate a Tracetest',
  isActive: false,
  stepList: [],
  demoList: [],
  type: TriggerTypes.http,
};

const Postman: IPlugin = {
  name: SupportedPlugins.Postman,
  title: 'Postman Collection',
  description: 'Define your HTTP Request via a Postman Collection',
  type: TriggerTypes.http,
  isActive: true,
  demoList: [
    {
      name: 'Postman - Pokemon - List',
      url: 'http://demo-pokemon-api.demo.svc.cluster.local/pokemon?take=20&skip=0',
      method: HTTP_METHOD.GET,
      body: '',
      description: 'Get a Pokemon',
      collectionTest: 'List',
      collectionFile: pokeshopPostmanFile,
    },
    {
      name: 'Postman - Pokemon - Add',
      url: 'http://demo-pokemon-api.demo.svc.cluster.local/pokemon',
      method: HTTP_METHOD.POST,
      body: '{"name":"meowth","type":"normal","imageUrl":"https://assets.pokemon.com/assets/cms2/img/pokedex/full/052.png","isFeatured":true}',
      description: 'Add a Pokemon',
      collectionTest: 'Create',
      collectionFile: pokeshopPostmanFile,
    },
    {
      name: 'Pokemon - Import',
      url: 'http://demo-pokemon-api.demo.svc.cluster.local/pokemon/import',
      method: HTTP_METHOD.POST,
      body: '{"id":52}',
      description: 'Import a Pokemon',
      collectionTest: 'Import',
      collectionFile: pokeshopPostmanFile,
    },
  ],
  stepList: [
    ...Default.stepList,
    {
      id: 'import-postman-collection',
      name: 'Import Postman collection',
      title: 'Upload Postman collection',
      component: 'UploadCollection',
    },
  ],
};

const OpenAPI: IPlugin = {
  name: SupportedPlugins.OpenAPI,
  title: 'OpenAPI',
  description: 'Define your HTTP request via an OpenAPI definition',
  isActive: false,
  stepList: [],
  demoList: [],
  type: TriggerTypes.http,
};

export const Plugins: Record<SupportedPlugins, IPlugin> = {
  [SupportedPlugins.REST]: Rest,
  [SupportedPlugins.RPC]: RPC,
  [SupportedPlugins.Postman]: Postman,
  [SupportedPlugins.Messaging]: Messaging,
  [SupportedPlugins.OpenAPI]: OpenAPI,
};

export const TriggerTypeToPlugin = {
  [TriggerTypes.http]: Plugins.REST,
  [TriggerTypes.grpc]: Plugins.RPC,
} as const;
