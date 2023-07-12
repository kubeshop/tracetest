import {Model, TTestRunnerSchemas} from 'types/Common.types';

export type TRawRequiredGatesResult = TTestRunnerSchemas['RequiredGatesResult'];
type RequiredGatesResult = Model<
  TRawRequiredGatesResult,
  {
    requiredFailedGates: string[];
  }
>;

const RequiredGatesResult = ({
  required = [],
  failed = [],
  passed = true,
}: TRawRequiredGatesResult): RequiredGatesResult => ({
  required,
  failed,
  passed,
  requiredFailedGates: required.filter(gate => failed.includes(gate)),
});

export default RequiredGatesResult;
