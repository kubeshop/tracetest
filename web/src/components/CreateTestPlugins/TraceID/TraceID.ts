import {IPluginComponentMap} from 'types/Plugins.types';
import VariableName from './steps/VariableName';
import Default from '../Default';

export const PluginComponentMap: IPluginComponentMap = {
  SelectPlugin: Default.SelectPlugin,
  BasicDetails: Default.BasicDetails,
  TraceIdVariableName: VariableName,
};

export default PluginComponentMap;
