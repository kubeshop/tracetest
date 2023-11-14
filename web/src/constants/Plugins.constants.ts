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
  title: 'HTTP Request',
  description: 'Create a basic HTTP request',
  isActive: true,
  type: TriggerTypes.http,
  demoList: [],
};

const GRPC: IPlugin = {
  name: SupportedPlugins.GRPC,
  title: 'GRPC Request',
  description: 'Test and debug your GRPC request',
  isActive: true,
  type: TriggerTypes.grpc,
  demoList: [],
};

const Kafka: IPlugin = {
  name: SupportedPlugins.Kafka,
  title: 'Kafka',
  description: 'Test consumers with Kafka messages',
  isActive: true,
  demoList: [],
  type: TriggerTypes.kafka,
};

const TraceID: IPlugin = {
  name: SupportedPlugins.TraceID,
  title: 'TraceID',
  description: 'Define your test via a Trace ID',
  isActive: true,
  demoList: [],
  type: TriggerTypes.traceid,
};

export const Plugins = {
  [SupportedPlugins.REST]: Rest,
  [SupportedPlugins.GRPC]: GRPC,
  [SupportedPlugins.TraceID]: TraceID,
  [SupportedPlugins.Kafka]: Kafka,
} as const;

export const TriggerTypeToPlugin = {
  [TriggerTypes.http]: Plugins.REST,
  [TriggerTypes.grpc]: Plugins.GRPC,
  [TriggerTypes.traceid]: Plugins.TraceID,
  [TriggerTypes.kafka]: Plugins.Kafka,
} as const;
