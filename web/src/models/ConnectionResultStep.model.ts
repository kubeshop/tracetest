import {Model, TConfigSchemas} from 'types/Common.types';

export type TRawConnectionTestStep = TConfigSchemas['ConnectionTestStep'];
type ConnectionTestStep = Model<TRawConnectionTestStep, {}>;

const ConnectionTestStep = ({
  passed = false,
  status = 'passed',
  error = '',
  message = '',
}: TRawConnectionTestStep): ConnectionTestStep => ({
  passed,
  status,
  error,
  message,
});

export default ConnectionTestStep;
