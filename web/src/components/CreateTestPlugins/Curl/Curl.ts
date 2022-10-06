import {IPluginComponentMap} from 'types/Plugins.types';
import RequestDetails from '../Rest/steps/RequestDetails';
import ImportCommand from './steps/ImportCommand';
import Default from '../Default';

export const PluginComponentMap: IPluginComponentMap = {
  SelectPlugin: Default.SelectPlugin,
  BasicDetails: Default.BasicDetails,
  ImportCommand,
  RequestDetails,
};

export default PluginComponentMap;
