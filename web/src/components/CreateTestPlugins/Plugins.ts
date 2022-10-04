import {IPluginComponentMap} from 'types/Plugins.types';
import {SupportedPlugins} from 'constants/Common.constants';
import Postman from './Postman';
import Rest from './Rest';
import Grpc from './Grpc';
import Curl from './Curl';

export const PluginsComponentMap: Record<SupportedPlugins, IPluginComponentMap> = {
  [SupportedPlugins.REST]: Rest,
  [SupportedPlugins.GRPC]: Grpc,
  [SupportedPlugins.Postman]: Postman,
  [SupportedPlugins.CURL]: Curl,
  [SupportedPlugins.Messaging]: {},
  [SupportedPlugins.OpenAPI]: {},
};

export default PluginsComponentMap;
