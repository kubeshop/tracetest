import {TDraftVariables, TMissingVariable, TTestVariablesMap} from 'types/Variables.types';
import {TEnvironmentValue} from 'types/Environment.types';
import {TTest} from 'types/Test.types';
import {uniqBy} from 'lodash';

const VariablesService = () => ({
  getVariableEntries(missingVariables: TMissingVariable[], testList: TTest[]): TTestVariablesMap {
    return missingVariables.reduce((entries, {testId, variables}, index) => {
      const test = testList.find(({id}) => id === testId);

      return variables.length && !!test
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

  getDraftVariables(testVariablesMap: TTestVariablesMap): TDraftVariables {
    return {
      variables: Object.values(testVariablesMap).reduce(
        (acc, {variables}) => ({
          ...acc,
          ...variables.reduce((variablesAcc, {key, defaultValue}) => ({...variablesAcc, [key]: defaultValue}), {}),
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

  getMatchingTests(testVariables: TTestVariablesMap, key: string): TTest[] {
    const list = Object.values(testVariables).reduce<TTest[]>(
      (acc, {test, variables}) =>
        variables.find(missingVariable => missingVariable.key === key) ? acc.concat(test) : acc,
      []
    );

    return uniqBy(list, 'id');
  },
});

export default VariablesService();
