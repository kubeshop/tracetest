import {TTestSchemas} from 'types/Common.types';

export type TRawTestSummary = TTestSchemas['TestSummary'];
type Summary = {
  runs: number;
  hasRuns: boolean;
  lastRun: {
    time: string;
    passes: number;
    fails: number;
    analyzerScore: number;
  };
};

const Summary = (summary: TRawTestSummary = {}): Summary => ({
  runs: summary.runs ?? 0,
  hasRuns: !!summary.runs && summary.runs > 0,
  lastRun: {
    time: summary.lastRun?.time ?? '',
    passes: summary.lastRun?.passes ?? 0,
    fails: summary.lastRun?.fails ?? 0,
    analyzerScore: summary.lastRun?.analyzerScore ?? 0,
  },
});

export default Summary;
