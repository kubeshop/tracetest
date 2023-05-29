import {IHttpValues, ITriggerService} from 'types/Test.types';
import Validator from 'utils/Validator';
import {HTTP_METHOD} from 'constants/Common.constants';
import HttpRequest from 'models/HttpRequest.model';

const HttpTriggerService = (): ITriggerService => ({
  async getRequest(values) {
    const {url, method, auth, headers, body, sslVerification} = values as IHttpValues;

    return HttpRequest({url, method, auth, headers, body, sslVerification});
  },

  async validateDraft(draft): Promise<boolean> {
    const {url, method} = draft as IHttpValues;
    return Validator.required(url) && Validator.required(method);
  },

  getInitialValues(request) {
    const {url, method, headers, body, auth, sslVerification} = request as HttpRequest;

    return {
      url,
      auth,
      method: method as HTTP_METHOD,
      headers,
      body,
      sslVerification,
    };
  },
});

export default HttpTriggerService();
