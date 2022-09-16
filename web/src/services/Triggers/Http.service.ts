import {IHttpValues, ITriggerService, THTTPRequest, TRawHTTPRequest} from 'types/Test.types';
import Validator from 'utils/Validator';
import {HTTP_METHOD} from '../../constants/Common.constants';

const HttpTriggerService = (): ITriggerService => ({
  async getRequest(values): Promise<TRawHTTPRequest> {
    const {url, method, auth, headers, body} = values as IHttpValues;

    return {url, method, auth, headers, body};
  },

  async validateDraft(draft): Promise<boolean> {
    const {url, method} = draft as IHttpValues;
    return Validator.required(url) && Validator.required(method) && Validator.url(url);
  },

  getInitialValues(request) {
    const {url, method, headers, body, auth} = request as THTTPRequest;

    return {
      url,
      auth,
      method: method as HTTP_METHOD,
      headers,
      body,
    };
  },
});

export default HttpTriggerService();
