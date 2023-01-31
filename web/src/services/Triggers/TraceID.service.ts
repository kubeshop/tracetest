import {ITraceIDValues, ITriggerService, TRawTRACEIDRequest} from 'types/Test.types';
import Validator from 'utils/Validator';

const TraceIDTriggerService = (): ITriggerService => ({
  async getRequest(values): Promise<TRawTRACEIDRequest> {
    const {id} = values as ITraceIDValues;

    return {id: id.includes('env:') ? id : `\${env:${id}}`};
  },

  async validateDraft(draft): Promise<boolean> {
    const {id} = draft as ITraceIDValues;
    return Validator.required(id);
  },

  getInitialValues(request) {
    const {id} = request as ITraceIDValues;

    return {id};
  },
});

export default TraceIDTriggerService();
