import VariableSet, {TRawVariableSet} from 'models/VariableSet.model';

const VariableSetService = () => ({
  getRequest(variableSet: VariableSet): TRawVariableSet {
    return {
      type: 'VariableSet',
      spec: variableSet,
    };
  },

  validateDraft({name = '', values = []}: VariableSet) {
    return !!name && !!values.length;
  },
});

export default VariableSetService();
