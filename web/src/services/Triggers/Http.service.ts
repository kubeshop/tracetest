import {IHttpValues, ITriggerService, TDraftTest, TRawHTTPRequest} from 'types/Test.types';
import Validator from 'utils/Validator';

const HttpTriggerService = (): ITriggerService => ({
  async getRequest(values: TDraftTest): Promise<TRawHTTPRequest> {
    const {url, method, auth, headers, body} = values as IHttpValues;

    return {url, method, auth, headers, body};
  },

  async validateDraft(draft: TDraftTest): Promise<boolean> {
    const {url, method} = draft as IHttpValues;

    const isValid = Validator.required(url) && Validator.required(method) && Validator.url(url);

    return isValid;
  },
});

export default HttpTriggerService();
