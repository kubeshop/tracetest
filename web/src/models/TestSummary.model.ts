import {TRawTestSummary, TSummary} from 'types/Test.types';

const TestSummary = (summary: TRawTestSummary = {}): TSummary => ({
  runs: summary.runs ?? 0,
  lastRun: {
    time: summary.lastRun?.time ?? '',
    passes: summary.lastRun?.passes ?? 0,
    fails: summary.lastRun?.fails ?? 0,
  },
});

export default TestSummary;
