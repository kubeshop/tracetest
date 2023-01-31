import {TTestSchemas} from 'types/Common.types';

export type TRawTestOutput = TTestSchemas['TestOutput'];
type TestOutput = {
  isDeleted: boolean;
  isDraft: boolean;
  name: string;
  selector: string;
  value: string;
  valueRun: string;
  valueRunDraft: string;
  id: number;
  spanId: string;
};

function TestOutput({name = '', selector = {}, value = ''}: TRawTestOutput, id = -1): TestOutput {
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

export function toRawTestOutputs(testOutputs: TestOutput[]): TRawTestOutput[] {
  return testOutputs
    .filter(output => !output.isDeleted)
    .map<TRawTestOutput>(output => ({
      name: output.name,
      selector: {query: output.selector},
      value: output.value,
    }));
}

export default TestOutput;
