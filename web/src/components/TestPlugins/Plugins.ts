import {SupportedPlugins} from 'constants/Common.constants';
import Rest from './Forms/Rest';
import Grpc from './Forms/Grpc';
import TraceID from './Forms/TraceID';
import Kafka from './Forms/Kafka';

export const PluginsComponentMap = {
  [SupportedPlugins.REST]: Rest,
  [SupportedPlugins.GRPC]: Grpc,
  [SupportedPlugins.Kafka]: Kafka,
  [SupportedPlugins.TraceID]: TraceID,
} as const;

export default PluginsComponentMap;
