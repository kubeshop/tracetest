import {IPluginComponentMap} from 'types/Plugins.types';
import RequestDetails from './steps/RequestDetails/RequestDetails';
import Default from '../Default';

export const PluginComponentMap: IPluginComponentMap = {
  SelectPlugin: Default.SelectPlugin,
  BasicDetails: Default.BasicDetails,
  RequestDetails,
};

export default PluginComponentMap;
