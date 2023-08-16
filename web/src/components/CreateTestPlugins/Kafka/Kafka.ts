import {IPluginComponentMap} from 'types/Plugins.types';
import Default from '../Default';
import RequestDetails from './steps/RequestDetails';

export const PluginComponentMap: IPluginComponentMap = {
  SelectPlugin: Default.SelectPlugin,
  BasicDetails: Default.BasicDetails,
  RequestDetails,
};

export default PluginComponentMap;
