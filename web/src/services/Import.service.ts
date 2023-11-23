import {ImportTypes} from 'constants/Test.constants';
import {TDraftTest} from 'types/Test.types';
import {IPlugin} from 'types/Plugins.types';
import CurlService from './Importers/Curl.service';
import PostmanService from './Importers/Postman.service';
import DefinitionService from './Importers/Definition.service';

const ImportServiceMap = {
  [ImportTypes.curl]: CurlService,
  [ImportTypes.postman]: PostmanService,
  [ImportTypes.definition]: DefinitionService,
} as const;

const ImportService = () => ({
  async getRequest(type: ImportTypes, draft: TDraftTest): Promise<TDraftTest> {
    return ImportServiceMap[type].getRequest(draft);
  },

  async validateDraft(type: ImportTypes, draft: TDraftTest): Promise<boolean> {
    return ImportServiceMap[type].validateDraft(draft);
  },

  async getPlugin(type: ImportTypes, draft: TDraftTest): Promise<IPlugin> {
    return ImportServiceMap[type].getPlugin(draft);
  },
});

export default ImportService();
