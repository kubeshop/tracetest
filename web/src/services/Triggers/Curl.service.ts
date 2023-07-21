import parseCurl from 'parse-curl';
import {ICurlValues, ITriggerService} from 'types/Test.types';
import Validator from 'utils/Validator';
import {HTTP_METHOD} from 'constants/Common.constants';
import HttpRequest from 'models/HttpRequest.model';

interface ICurlTriggerService extends ITriggerService {
  getRequestFromCommand(command: string): any;
  getIsValidCommand(command: string): boolean;
}

const CurlTriggerService = (): ICurlTriggerService => ({
  async getRequest(draft) {
    const {url, method, auth, headers, body} = draft as ICurlValues;

    return HttpRequest({url, method, auth, headers, body});
  },

  async validateDraft(draft) {
    const {url, method} = draft as ICurlValues;
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
