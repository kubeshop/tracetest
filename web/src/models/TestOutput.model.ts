import {TRawTestOutput, TRawTestRunOutput, TTestOutput, TTestRunOutput} from 'types/TestOutput.types';

function TestOutput({name = '', selector = {}, value = ''}: TRawTestOutput, id = -1): TTestOutput {
  return {
    id,
    isDeleted: false,
    isDraft: false,
    name,
    selector: selector.query || '',
    value,
    valueRun: '',
    valueRunDraft: '',
    spanId: '',
  };
}

export function TestRunOutput({name = '', value = '', spanId = ''}: TRawTestRunOutput): TTestRunOutput {
  return {
    name,
    value,
    spanId,
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
