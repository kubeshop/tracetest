import {IHttpValues, ITriggerService, TRawHTTPRequest} from 'types/Test.types';
import Validator from 'utils/Validator';

const HttpTriggerService = (): ITriggerService => ({
  async getRequest(values): Promise<TRawHTTPRequest> {
    const {url, method, auth, headers, body} = values as IHttpValues;

    return {url, method, auth, headers, body};
  },

  async validateDraft(draft): Promise<boolean> {
    const {url, method} = draft as IHttpValues;

    const isValid = Validator.required(url) && Validator.required(method) && Validator.url(url);

    return isValid;
  },
});

export default HttpTriggerService();
