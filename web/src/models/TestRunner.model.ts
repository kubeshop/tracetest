import {Model, TTestRunnerSchemas} from 'types/Common.types';

export type TRawTestRunnerResource = TTestRunnerSchemas['TestRunnerResource'];
export type TRawTestRunner = TTestRunnerSchemas['TestRunner'];
export type TSupportedGates = TTestRunnerSchemas['SupportedGates'];
type TestRunner = Model<TRawTestRunner, {}>;

export enum SupportedRequiredGates {
  AnalyzerScore = 'analyzer-score',
  AnalyzerRules = 'analyzer-rules',
  TestSpecs = 'test-specs',
}

export const SupportedRequiredGatesDescription = {
  [SupportedRequiredGates.AnalyzerScore]:
    'Test Runs will be marked as failed if the Analyzer Score is below the configured threshold',
  [SupportedRequiredGates.AnalyzerRules]:
    'Test Runs will be marked as failed if one the Error Level Analyzer Rules fails',
  [SupportedRequiredGates.TestSpecs]: 'Test Runs will be marked as failed if one of the defined Test Specs fails',
} as const;

const TestRunner = ({spec: rawTestRunner = {}}: TRawTestRunnerResource = {}): TestRunner =>
  TestRunner.FromRawTestRunner(rawTestRunner);

TestRunner.FromRawTestRunner = ({requiredGates = [], name = '', id = 'current'}: TRawTestRunner): TestRunner => ({
  requiredGates,
  name,
  id,
});

export default TestRunner;
