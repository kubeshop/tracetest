import {IPluginComponentMap} from 'types/Plugins.types';
import Value from './steps/Value';
import Default from '../Default';

export const PluginComponentMap: IPluginComponentMap = {
  SelectPlugin: Default.SelectPlugin,
  BasicDetails: Default.BasicDetails,
  TraceIDValue: Value,
};

export default PluginComponentMap;
