import {TConnectionTestStep, TRawConnectionTestStep} from '../types/Config.types';

const ConnectionTestStep = ({
  passed = false,
  error = '',
  message = '',
}: TRawConnectionTestStep): TConnectionTestStep => ({
  passed,
  error,
  message,
});

export default ConnectionTestStep;
