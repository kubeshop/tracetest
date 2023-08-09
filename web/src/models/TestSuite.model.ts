import { Model, TTestSuiteSchemas } from 'types/Common.types';
import Summary from './Summary.model';
import Test from './Test.model';

export type TRawTestSuiteResource = TTestSuiteSchemas['TestSuiteResource'];
export type TRawTestSuite = TTestSuiteSchemas['TestSuite'];
type TestSuite = Model<
  TTestSuiteSchemas['TestSuite'],
  {
    steps: string[];
    fullSteps: Test[];
    summary: Summary;
  }
>;

function TestSuite({ spec: raw = {} }: TRawTestSuiteResource): TestSuite {
  return TestSuite.FromRawTestSuite(raw);
}

TestSuite.FromRawTestSuite = ({
  id = '',
  name = '',
  description = '',
  version = 1,
  createdAt = '',
  summary = {},
  steps = [],
  fullSteps = [],
}: TRawTestSuite): TestSuite => {
  return {
    id,
    name,
    description,
    version,
    steps,
    fullSteps: fullSteps.map(step => Test.FromRawTest(step)),
    createdAt,
    summary: Summary(summary),
  };
};

export default TestSuite;
