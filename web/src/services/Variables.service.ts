import {TDraftVariables, TMissingVariable} from 'types/Variables.types';
import {TEnvironmentValue} from 'types/Environment.types';

const VariablesService = () => ({
  getDraftVariables(missingVariables: TMissingVariable[]): TDraftVariables {
    return {
      variables: missingVariables.reduce(
        (acc, {key, defaultValue}) => ({
          ...acc,
          [key]: defaultValue,
        }),
        {}
      ),
    };
  },

  getSubmitValues({variables}: TDraftVariables): TEnvironmentValue[] {
    return Object.entries(variables).map(([key, value]) => ({
      key,
      value,
    }));
  },
});

export default VariablesService();
