import {TTestSchemas} from 'types/Common.types';
import Date from 'utils/Date';

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

const Summary = (summary: TRawTestSummary = {}): Summary => {
  const time = summary.lastRun?.time ?? '';

  return {
    runs: summary.runs ?? 0,
    hasRuns: !!summary.runs && summary.runs > 0,
    lastRun: {
      time: Date.isDefaultDate(time) ? '' : time,
      passes: summary.lastRun?.passes ?? 0,
      fails: summary.lastRun?.fails ?? 0,
      analyzerScore: summary.lastRun?.analyzerScore ?? 0,
    },
  };
};

export default Summary;
