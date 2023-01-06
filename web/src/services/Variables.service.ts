import {uniqBy} from 'lodash';
import {TTestVariables, TDraftVariables, TTestVariablesMap} from 'types/Variables.types';
import {TEnvironmentValue} from '../types/Environment.types';
import {TTest} from '../types/Test.types';

const VariablesService = () => ({
  getVariableEntries(testsVariables: TTestVariables[]): TTestVariablesMap {
    return testsVariables.reduce((entries, {test, variables: {missing}, variables}, index) => {
      return missing.length
        ? {
            ...entries,
            [`${test.id}-${index}`]: {
              test,
              variables,
            },
          }
        : entries;
    }, {});
  },

  getDraftVariables(variableMap: TTestVariablesMap): TDraftVariables {
    return {
      variables: Object.values(variableMap).reduce(
        (acc, {variables: {missing}}) =>
          missing.length
            ? {
                ...acc,
                ...missing.reduce((missingAcc, {key, defaultValue}) => ({...missingAcc, [key]: defaultValue}), {}),
              }
            : acc,
        {}
      ),
    };
  },

  getFlatVariablesFromDraft({variables}: TDraftVariables): TEnvironmentValue[] {
    return Object.entries(variables).flatMap(([key, value]) => [{key, value}]);
  },

  getMatchingTests(variables: TTestVariablesMap, key: string): TTest[] {
    const list = Object.values(variables).reduce<TTest[]>(
      (acc, {test, variables: {missing}}) =>
        missing.find(missingVariable => missingVariable.key === key) ? acc.concat(test) : acc,
      []
    );

    return uniqBy(list, 'id');
  },
});

export default VariablesService();
