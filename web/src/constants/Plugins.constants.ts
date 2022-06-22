import {IPlugin} from '../types/Plugins.types';

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
  description: 'Test and debug your messaging queue',
  isActive: false,
  stepList: [],
};

const RPC: IPlugin = {
  name: SupportedPlugins.RPC,
  title: 'RPC Request',
  description: 'Test and debug your RPC request',
  isActive: false,
  stepList: [],
};

const Postman: IPlugin = {
  name: SupportedPlugins.RPC,
  title: 'Postman Collection',
  description: 'Define your HTTP Request via a Postman Collection',
  isActive: false,
  stepList: [],
};

const OpenAPI: IPlugin = {
  name: SupportedPlugins.RPC,
  title: 'OpenAPI',
  description: 'Define your HTTP request via an OpenAPI definition',
  isActive: false,
  stepList: [],
};

export const Plugins: Record<SupportedPlugins, IPlugin> = {
  [SupportedPlugins.REST]: Rest,
  [SupportedPlugins.Messaging]: Messaging,
  [SupportedPlugins.RPC]: RPC,
  [SupportedPlugins.Postman]: Postman,
  [SupportedPlugins.OpenAPI]: OpenAPI,
};
