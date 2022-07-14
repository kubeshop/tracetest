import {IPluginComponentMap} from '../../../types/Plugins.types';
import Default from '../Default/Default';
import UploadCollection from './steps/UploadCollection';

export const PluginComponentMap: IPluginComponentMap = {
  SelectPlugin: Default.SelectPlugin,
  BasicDetails: Default.BasicDetails,
  UploadCollection,
};

export default PluginComponentMap;
