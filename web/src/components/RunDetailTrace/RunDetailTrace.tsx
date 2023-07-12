import ResizablePanels from 'components/ResizablePanels';
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
}

export enum VisualizationType {
  Dag,
  Timeline,
}

const RunDetailTrace = ({run, runEvents, testId}: IProps) => {
  return (
    <S.Container>
      <SetupAlert />
      <ResizablePanels>
        <SpanDetailsPanel run={run} testId={testId} />
        <TracePanel run={run} runEvents={runEvents} testId={testId} />
        <AnalyzerPanel run={run} />
      </ResizablePanels>
    </S.Container>
  );
};

export default RunDetailTrace;
