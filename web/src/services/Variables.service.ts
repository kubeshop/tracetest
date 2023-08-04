import {TDraftVariables, TTestVariablesMap} from 'types/Variables.types';
import {uniqBy} from 'lodash';
import MissingVariables from 'models/MissingVariables.model';
import Test from 'models/Test.model';
import {TVariableSetValue} from 'models/VariableSet.model';

const VariablesService = () => ({
  getVariableEntries(missingVariables: MissingVariables, testList: Test[]): TTestVariablesMap {
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

  getSubmitValues({variables}: TDraftVariables): TVariableSetValue[] {
    return Object.entries(variables).map(([key, value]) => ({
      key,
      value,
    }));
  },

  getMatchingTests(testVariables: TTestVariablesMap, key: string): Test[] {
    const list = Object.values(testVariables).reduce<Test[]>(
      (acc, {test, variables}) =>
        variables.find(missingVariable => missingVariable.key === key) ? acc.concat(test) : acc,
      []
    );

    return uniqBy(list, 'id');
  },
});

export default VariablesService();
