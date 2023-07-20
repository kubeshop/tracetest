// import {CURLParser} from 'parse-curl-js';
import {ICurlValues, ITriggerService} from 'types/Test.types';
import Validator from 'utils/Validator';
import {HTTP_METHOD} from 'constants/Common.constants';
import HttpRequest from 'models/HttpRequest.model';

class CURLParser {
  constructor(command: string) {}

  parse(): any {}
}

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
    const parser = new CURLParser(command);
    const {
      url = '',
      method,
      headers = [],
      body: {data: body = {}},
    } = parser.parse();

    return {
      url: url.split('')[0] === "'" ? url.slice(1, -1) : url,
      auth: undefined,
      command,
      method: method as HTTP_METHOD,
      headers: Object.entries(headers).map(([key, value]) => ({key, value})),
      body: body === 'data' ? '' : JSON.stringify(body),
      sslVerification: false,
    };
  },

  getIsValidCommand(command) {
    try {
      const parser = new CURLParser(command);
      const {
        url = '',
        method,
        body: {data: body = {}},
      } = parser.parse();

      return Boolean((method.toUpperCase() === HTTP_METHOD.POST && body) || (url && method));
    } catch (e) {
      return false;
    }
  },
});

export default CurlTriggerService();
