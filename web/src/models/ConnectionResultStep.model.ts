import {Model, TConfigSchemas} from 'types/Common.types';

export type TRawConnectionTestStep = TConfigSchemas['ConnectionTestStep'];
type ConnectionTestStep = Model<TRawConnectionTestStep, {}>;

const ConnectionTestStep = ({
  passed = false,
  error = '',
  message = '',
}: TRawConnectionTestStep): ConnectionTestStep => ({
  passed,
  error,
  message,
});

export default ConnectionTestStep;
