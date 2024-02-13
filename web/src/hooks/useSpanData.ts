import Span from 'models/Span.model';
import TestRunOutput from 'models/TestRunOutput.model';
import {useTest} from 'providers/Test/Test.provider';
import {useTestRun} from 'providers/TestRun/TestRun.provider';
import {useAppSelector} from 'redux/hooks';
import {
  selectSpanById,
  selectAnalyzerErrorsBySpanId,
  selectTestSpecsBySpanId,
  selectTestOutputsBySpanId,
} from 'selectors/TestRun.selectors';
import {TAnalyzerError, TTestSpecSummary} from 'types/TestRun.types';

interface IUseSpanData {
  span: Span;
  analyzerErrors?: TAnalyzerError[];
  testSpecs?: TTestSpecSummary;
  testOutputs?: TestRunOutput[];
}

const useSpanData = (id: string): IUseSpanData => {
  const {
    test: {id: testId},
  } = useTest();
  const {
    run: {id: runId},
  } = useTestRun();

  const span = useAppSelector(state => selectSpanById(state, {testId, runId, spanId: id}));

  // TODO: should we get analyzerErrors, testSpecs and testOutputs as part of the trace struct from the BE?
  // Right now we are getting them from the testRun struct for each span by spanId
  const analyzerErrors = useAppSelector(state => selectAnalyzerErrorsBySpanId(state, {testId, runId, spanId: id}));

  const testSpecs = useAppSelector(state => selectTestSpecsBySpanId(state, {testId, runId, spanId: id}));

  const testOutputs = useAppSelector(state => selectTestOutputsBySpanId(state, {testId, runId, spanId: id}));

  return {
    span,
    analyzerErrors,
    testSpecs,
    testOutputs,
  };
};

export default useSpanData;
