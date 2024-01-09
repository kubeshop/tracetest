import {IPlugin} from 'types/Plugins.types';
import {SupportedPlugins} from './Common.constants';
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

const Rest: IPlugin = {
  name: SupportedPlugins.REST,
  title: 'HTTP',
  description: 'Test your HTTP service with an HTTP request',
  isActive: true,
  type: TriggerTypes.http,
  demoList: [],
};

const GRPC: IPlugin = {
  name: SupportedPlugins.GRPC,
  title: 'gRPC',
  description: 'Test your gRPC service with a gRPC request',
  isActive: true,
  type: TriggerTypes.grpc,
  demoList: [],
};

const Kafka: IPlugin = {
  name: SupportedPlugins.Kafka,
  title: 'Kafka',
  description: 'Test Kafka based services with a Kafka request',
  isActive: true,
  demoList: [],
  type: TriggerTypes.kafka,
};

const TraceID: IPlugin = {
  name: SupportedPlugins.TraceID,
  title: 'TraceID',
  description: 'Define your test via a TraceID',
  isActive: true,
  demoList: [],
  type: TriggerTypes.traceid,
};

const Cypress: IPlugin = {
  name: SupportedPlugins.Cypress,
  title: 'Cypress',
  description: 'Define your test via Cypress',
  isActive: true,
  demoList: [],
  type: TriggerTypes.cypress,
  requestType: TriggerTypes.traceid,
};

const Playwright: IPlugin = {
  name: SupportedPlugins.Playwright,
  title: 'Playwright',
  description: 'Define your test via Playwright',
  isActive: true,
  demoList: [],
  type: TriggerTypes.playwright,
  requestType: TriggerTypes.traceid,
};

export const Plugins = {
  [SupportedPlugins.REST]: Rest,
  [SupportedPlugins.GRPC]: GRPC,
  [SupportedPlugins.Kafka]: Kafka,
  [SupportedPlugins.TraceID]: TraceID,
  [SupportedPlugins.Cypress]: Cypress,
  [SupportedPlugins.Playwright]: Playwright,
} as const;

export const TriggerTypeToPlugin = {
  [TriggerTypes.http]: Plugins.REST,
  [TriggerTypes.grpc]: Plugins.GRPC,
  [TriggerTypes.kafka]: Plugins.Kafka,
  [TriggerTypes.traceid]: Plugins.TraceID,
  [TriggerTypes.cypress]: Plugins.Cypress,
  [TriggerTypes.playwright]: Plugins.Playwright,
} as const;

export const CreateTriggerTypeToPlugin = {
  [TriggerTypes.http]: Plugins.REST,
  [TriggerTypes.grpc]: Plugins.GRPC,
  [TriggerTypes.kafka]: Plugins.Kafka,
  [TriggerTypes.traceid]: Plugins.TraceID,
} as const;
