import TestRun from 'models/TestRun.model';
import TestRunEvent from 'models/TestRunEvent.model';
import * as S from './RunDetailTest.styled';
import SetupAlert from '../SetupAlert/SetupAlert';
import ResizablePanels from '../ResizablePanels/ResizablePanels';
import SpanDetailsPanel from './SpanDetailsPanel';
import TestPanel from './TestPanel';

interface IProps {
  run: TestRun;
  runEvents: TestRunEvent[];
  testId: string;
}

const RunDetailTest = ({run, runEvents, testId}: IProps) => {
  return (
    <S.Container>
      <SetupAlert />
      <ResizablePanels>
        <SpanDetailsPanel />
        <TestPanel run={run} runEvents={runEvents} testId={testId} />
      </ResizablePanels>
    </S.Container>
  );
};

export default RunDetailTest;
