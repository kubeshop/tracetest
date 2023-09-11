import {Model, TTestSuiteSchemas} from 'types/Common.types';
import VariableSet from './VariableSet.model';
import TestRun from './TestRun.model';

export type TRawTestSuiteRunResourceRun = TTestSuiteSchemas['TestSuiteRun'];
type TestSuiteRun = Model<
  TRawTestSuiteRunResourceRun,
  {
    steps: TestRun[];
    variableSet?: VariableSet;
    metadata?: {[key: string]: string};
  }
>;

const TestSuiteRun = ({
  id = 0,
  createdAt = '',
  completedAt = '',
  state = 'CREATED',
  steps = [],
  variableSet = {},
  metadata = {},
  version = 1,
  pass = 0,
  fail = 0,
  allStepsRequiredGatesPassed = false,
}: TRawTestSuiteRunResourceRun): TestSuiteRun => {
  return {
    id,
    createdAt,
    completedAt,
    state,
    steps: steps.map(step => TestRun(step)),
    variableSet: VariableSet.fromRun(variableSet),
    allStepsRequiredGatesPassed,
    metadata,
    version,
    pass,
    fail,
  };
};

export default TestSuiteRun;
