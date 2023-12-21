import {ITraceIDValues, ITriggerService} from 'types/Test.types';
import Validator from 'utils/Validator';
import TraceIDRequest from 'models/TraceIDRequest.model';

const TraceIDTriggerService = (): ITriggerService => ({
  async getRequest(values) {
    const {id} = values as ITraceIDValues;

    return TraceIDRequest({id: id.includes('env:') ||  id.includes('var:') ? id : `\${env:${id}}`});
  },

  async validateDraft(draft) {
    const {id} = draft as ITraceIDValues;
    return Validator.required(id);
  },

  getInitialValues(request) {
    const {id} = request as ITraceIDValues;

    return {id};
  },
});

export default TraceIDTriggerService();
