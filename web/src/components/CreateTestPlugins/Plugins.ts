import {IPluginComponentMap} from 'types/Plugins.types';
import {SupportedPlugins} from 'constants/Plugins.constants';
import Rest from './Rest';
import Rpc from './Rpc';

export const PluginsComponentMap: Record<SupportedPlugins, IPluginComponentMap> = {
  [SupportedPlugins.REST]: Rest,
  [SupportedPlugins.Messaging]: {},
  [SupportedPlugins.RPC]: Rpc,
  [SupportedPlugins.Postman]: {},
  [SupportedPlugins.OpenAPI]: {},
};

export default PluginsComponentMap;
