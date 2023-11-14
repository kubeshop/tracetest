import {ImportTypes} from 'constants/Test.constants';
import {TDraftTest, IImportResult} from 'types/Test.types';
import CurlService from './Importers/Curl.service';
import PostmanService from './Importers/Postman.service';

const ImportServiceMap = {
  [ImportTypes.curl]: CurlService,
  [ImportTypes.postman]: PostmanService,
} as const;

const ImportService = () => ({
  async getRequest(type: ImportTypes, draft: TDraftTest): Promise<IImportResult> {
    return ImportServiceMap[type].getRequest(draft);
  },

  async validateDraft(type: ImportTypes, draft: TDraftTest): Promise<boolean> {
    return ImportServiceMap[type].validateDraft(draft);
  },
});

export default ImportService();
