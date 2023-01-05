import {TRawTestVariables, TTestVariables} from '../types/Variables.types';
import Test from './Test.model';

const TestVariables = ({
  test = {},
  variables: {missing = [], environment = [], variables = []} = {},
}: TRawTestVariables): TTestVariables => {
  return {
    test: Test(test),
    variables: {
      missing: missing.map(({key = '', defaultValue = ''}) => ({
        key,
        defaultValue,
      })),
      environment,
      variables,
    },
  };
};

export default TestVariables;
