import VariableSet, {TRawVariableSet} from 'models/VariableSet.model';

const VariableSetService = () => ({
  getRequest(variableSet: VariableSet): TRawVariableSet {
    return {
      type: 'VariableSet',
      spec: variableSet,
    };
  },

  validateDraft({name = '', description = '', values = []}: VariableSet) {
    return !!name && !!description && !!values.length;
  },
});

export default VariableSetService();
