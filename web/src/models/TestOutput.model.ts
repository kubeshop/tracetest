import {TRawTestOutput, TRawTestRunOutput, TTestOutput, TTestRunOutput} from 'types/TestOutput.types';

function TestOutput({name = '', selector = {}, value = ''}: TRawTestOutput): TTestOutput {
  return {
    isDeleted: false,
    isDraft: false,
    name,
    selector: selector.query || '',
    value,
    valueRun: '',
  };
}

export function TestRunOutput({name = '', value = 'test'}: TRawTestRunOutput): TTestRunOutput {
  return {
    name,
    value,
  };
}

export function toRawTestOutputs(testOutputs: TTestOutput[]): TRawTestOutput[] {
  return testOutputs
    .filter(output => !output.isDeleted)
    .map<TRawTestOutput>(output => ({
      name: output.name,
      selector: {query: output.selector},
      value: output.value,
    }));
}

export default TestOutput;
