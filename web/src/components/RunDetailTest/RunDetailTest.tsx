import PanelLayout from 'components/ResizablePanels';
import {useMemo} from 'react';
import TestRun from 'models/TestRun.model';
import TestRunEvent from 'models/TestRunEvent.model';
import * as S from './RunDetailTest.styled';
import {getTestPanel} from './TestPanel';
import SetupAlert from '../SetupAlert/SetupAlert';
import {getSpanDetailsPanel} from './SpanDetailPanel';

interface IProps {
  run: TestRun;
  runEvents: TestRunEvent[];
  testId: string;
}

const RunDetailTest = ({run, runEvents, testId}: IProps) => {
  const panels = useMemo(() => [getSpanDetailsPanel(), getTestPanel(testId, run, runEvents)], [run, runEvents, testId]);
  return (
    <S.Container>
      <SetupAlert />
      <PanelLayout panels={panels} />
    </S.Container>
  );
};

export default RunDetailTest;
