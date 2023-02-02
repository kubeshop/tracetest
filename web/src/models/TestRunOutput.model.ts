import {Model} from 'types/Common.types';

export type TRawTestRunOutput = {
  name?: string;
  value?: string;
  spanId?: string;
};
type TestRunOutput = Model<TRawTestRunOutput, {}>;

const TestRunOutput = ({name = '', value = '', spanId = ''}: TRawTestRunOutput): TestRunOutput => {
  return {
    name,
    value,
    spanId,
  };
};

export default TestRunOutput;
