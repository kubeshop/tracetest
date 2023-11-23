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

export const Plugins = {
  [SupportedPlugins.REST]: Rest,
  [SupportedPlugins.GRPC]: GRPC,
  [SupportedPlugins.Kafka]: Kafka,
  [SupportedPlugins.TraceID]: TraceID,
} as const;

export const TriggerTypeToPlugin = {
  [TriggerTypes.http]: Plugins.REST,
  [TriggerTypes.grpc]: Plugins.GRPC,
  [TriggerTypes.kafka]: Plugins.Kafka,
  [TriggerTypes.traceid]: Plugins.TraceID,
} as const;
