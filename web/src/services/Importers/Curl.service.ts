import parseCurl from 'parse-curl';
import {ICurlValues, IImportService} from 'types/Test.types';
import Validator from 'utils/Validator';
import {HTTP_METHOD} from 'constants/Common.constants';
import {Plugins} from 'constants/Plugins.constants';

interface ICurlTriggerService extends IImportService {
  getRequestFromCommand(command: string): ICurlValues;
  getIsValidCommand(command: string): boolean;
}

const CurlTriggerService = (): ICurlTriggerService => ({
  async getRequest(values) {
    const {command} = values as ICurlValues;
    const draft = this.getRequestFromCommand(command);

    return {
      draft: {
        ...draft,
        name: draft.url,
      },
      plugin: Plugins.REST,
    };
  },

  async validateDraft(draft) {
    const {command} = draft as ICurlValues;
    if (!this.getIsValidCommand(command)) return false;

    const {url, method} = this.getRequestFromCommand(command);
    return Validator.required(url) && Validator.required(method);
  },

  getRequestFromCommand(command) {
    const {url = '', method, header = {}, body} = parseCurl(command);

    return {
      url: url.split('')[0] === "'" ? url.slice(1, -1) : url,
      auth: undefined,
      command,
      method: method as HTTP_METHOD,
      headers: Object.entries(header).map(([key, value]) => ({key, value})),
      body,
      sslVerification: false,
    };
  },

  getIsValidCommand(command) {
    try {
      const {url, method, body} = parseCurl(command);

      return Boolean((method.toUpperCase() === HTTP_METHOD.POST && body) || (url && method));
    } catch (e) {
      return false;
    }
  },
});

export default CurlTriggerService();
