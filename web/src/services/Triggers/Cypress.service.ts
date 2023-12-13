import {ITraceIDValues, ITriggerService} from 'types/Test.types';
import Validator from 'utils/Validator';

const CypressTriggerService = (): ITriggerService => ({
  async getRequest() {
    return null;
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

export default CypressTriggerService();
