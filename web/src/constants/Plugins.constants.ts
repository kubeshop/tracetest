import {IPlugin} from 'types/Plugins.types';
import {SupportedPlugins} from './Common.constants';
import {DemoByPluginMap} from './Demo.constants';
import {TriggerTypes} from './Test.constants';

export enum ComponentNames {
  SelectPlugin = 'SelectPlugin',
  BasicDetails = 'BasicDetails',
  RequestDetails = 'RequestDetails',
  UploadCollection = 'UploadCollection',
  TestsSelection = 'TestsSelection',
  ImportCommand = 'ImportCommand',
  TraceIdVariableName = 'TraceIdVariableName',
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
      name: 'Select test type',
      title: 'Choose the way of creating a test',
      component: ComponentNames.SelectPlugin,
      isDefaultValid: true,
      status: 'selected',
    },
    {
      id: 'basic-details',
      name: 'Basic Details',
      title: 'Provide needed basic information',
      component: ComponentNames.BasicDetails,
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
      component: ComponentNames.RequestDetails,
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
      component: ComponentNames.RequestDetails,
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
      component: ComponentNames.UploadCollection,
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

const Curl: IPlugin = {
  name: SupportedPlugins.CURL,
  title: 'CURL Command',
  description: 'Define your HTTP request via a CURL command',
  isActive: true,
  stepList: [
    ...Default.stepList,
    {
      id: 'import-curl',
      name: 'Import CURL',
      title: 'Paste the CURL command',
      component: ComponentNames.ImportCommand,
    },
    {
      id: 'request-details',
      name: 'Request Details',
      title: 'Provide additional information',
      component: ComponentNames.RequestDetails,
    },
  ],
  demoList: DemoByPluginMap.CURL,
  type: TriggerTypes.http,
};

const TraceID: IPlugin = {
  name: SupportedPlugins.TraceID,
  title: 'TraceID',
  description: 'Define your test via a Trace ID',
  isActive: true,
  stepList: [
    ...Default.stepList,
    {
      id: 'trace-id-value',
      name: 'Variable Name',
      title: 'Add a Variable Name',
      component: ComponentNames.TraceIdVariableName,
    },
  ],
  demoList: [],
  type: TriggerTypes.traceid,
};

export const Plugins: Record<SupportedPlugins, IPlugin> = {
  [SupportedPlugins.REST]: Rest,
  [SupportedPlugins.GRPC]: GRPC,
  [SupportedPlugins.CURL]: Curl,
  [SupportedPlugins.Postman]: Postman,
  [SupportedPlugins.TraceID]: TraceID,
  [SupportedPlugins.Messaging]: Messaging,
  [SupportedPlugins.OpenAPI]: OpenAPI,
};

export const TriggerTypeToPlugin = {
  [TriggerTypes.http]: Plugins.REST,
  [TriggerTypes.grpc]: Plugins.GRPC,
  [TriggerTypes.traceid]: Plugins.TraceID,
} as const;
