import {IPluginComponentMap} from 'types/Plugins.types';
import {SupportedPlugins} from 'constants/Common.constants';
import Postman from './Postman';
import Rest from './Rest';
import Grpc from './Grpc';
import Curl from './Curl';
import TraceID from './TraceID';
import Kafka from './Kafka';

export const PluginsComponentMap: Record<SupportedPlugins, IPluginComponentMap> = {
  [SupportedPlugins.REST]: Rest,
  [SupportedPlugins.GRPC]: Grpc,
  [SupportedPlugins.Postman]: Postman,
  [SupportedPlugins.CURL]: Curl,
  [SupportedPlugins.Kafka]: Kafka,
  [SupportedPlugins.OpenAPI]: {},
  [SupportedPlugins.TraceID]: TraceID,
};

export default PluginsComponentMap;
