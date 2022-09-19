import {IPlugin} from 'types/Plugins.types';
import {SupportedPlugins} from './Common.constants';
import {DemoByPluginMap} from './Demo.constants';
import {TriggerTypes} from './Test.constants';

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
      name: 'Select test type',
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
  demoList: DemoByPluginMap.REST,
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

const GRPC: IPlugin = {
  name: SupportedPlugins.GRPC,
  title: 'GRPC Request',
  description: 'Test and debug your GRPC request',
  isActive: true,
  type: TriggerTypes.grpc,
  demoList: DemoByPluginMap.GRPC,
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
  demoList: DemoByPluginMap.Postman,
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
  [SupportedPlugins.GRPC]: GRPC,
  [SupportedPlugins.Postman]: Postman,
  [SupportedPlugins.Messaging]: Messaging,
  [SupportedPlugins.OpenAPI]: OpenAPI,
};

export const TriggerTypeToPlugin = {
  [TriggerTypes.http]: Plugins.REST,
  [TriggerTypes.grpc]: Plugins.GRPC,
} as const;
