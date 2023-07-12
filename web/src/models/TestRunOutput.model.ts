import {Model} from 'types/Common.types';

export type TRawTestRunOutput = {
  name?: string;
  value?: string;
  spanId?: string;
  error?: string;
};
type TestRunOutput = Model<TRawTestRunOutput, {}>;

const TestRunOutput = ({name = '', value = '', spanId = '', error = ''}: TRawTestRunOutput): TestRunOutput => {
  return {
    name,
    value,
    spanId,
    error,
  };
};

export default TestRunOutput;
