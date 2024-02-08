import ResizablePanels from 'components/ResizablePanels';
import {MAX_DAG_NODES} from 'constants/Visualization.constants';
import TestRun from 'models/TestRun.model';
import TestRunEvent from 'models/TestRunEvent.model';
import * as S from './RunDetailTrace.styled';
import SetupAlert from '../SetupAlert';
import AnalyzerPanel from './AnalyzerPanel';
import SpanDetailsPanel from './SpanDetailsPanel';
import TracePanel from './TracePanel';

interface IProps {
  run: TestRun;
  runEvents: TestRunEvent[];
  testId: string;
  skipTraceCollection: boolean;
}

export enum VisualizationType {
  Dag,
  Timeline,
}

export function getIsDAGDisabled(totalSpans: number = 0): boolean {
  return totalSpans > MAX_DAG_NODES;
}

const RunDetailTrace = ({run, runEvents, testId, skipTraceCollection}: IProps) => {
  return (
    <S.Container>
      <SetupAlert />
      <ResizablePanels>
        <SpanDetailsPanel run={run} testId={testId} />
        <TracePanel run={run} runEvents={runEvents} testId={testId} skipTraceCollection={skipTraceCollection} />
        <AnalyzerPanel run={run} />
      </ResizablePanels>
    </S.Container>
  );
};

export default RunDetailTrace;
